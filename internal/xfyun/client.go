package xfyun

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alevinval/sse/pkg/decoder"
	"github.com/bytedance/sonic"

	"github.com/jun3372/uniai/client"
	"github.com/jun3372/uniai/errorx"
	"github.com/jun3372/uniai/request"
	"github.com/jun3372/uniai/response"
)

type xfyun struct{}

func NewClient() client.IClient {
	return &xfyun{}
}

// Completions 方法用于获取补全建议。
// 它根据提供的选项、上下文和请求信息，返回一个响应的channel和可能的错误。
// 参数:
//
//	opt Options - 补全请求的选项，包含了请求的具体配置。
//	ctx context.Context - 请求的上下文，用于控制请求的取消和超时等。
//	in request.Request - 补全请求的对象，包含了请求的具体内容。
//
// 返回值:
//
//	chan response.Response - 一个channel，用于接收补全响应的结果。
//	error - 如果在请求过程中出现错误，将返回错误信息。
func (h xfyun) Completions(opt client.Options, ctx context.Context, in request.Request) (chan response.Response, error) {
	// 检查提供的主机地址是否为空，如果为空则返回错误。
	if opt.Host == "" {
		return nil, errorx.InvalidHost
	}

	// 定义默认的API端点
	endpoint := "/v1/chat/completions"
	// 如果输入参数中指定了端点，则使用输入参数中的端点覆盖默认值
	if in.Endpoint != "" {
		endpoint = in.Endpoint
	}

	// 将请求信息序列化为字符串形式。
	payload, _ := in.MarshalToString()
	// 构建请求的URI。
	uri := opt.Host + endpoint
	// 创建一个POST请求到指定的URI。
	req, err := http.NewRequest(http.MethodPost, uri, in.Payload())
	if err != nil {
		return nil, errorx.InvalidInput
	}

	// 添加自定义请求头到HTTP请求。
	for k, v := range opt.Header {
		req.Header.Add(k, v[0])
	}

	// 初始化一个可取消的上下文ctx，用于后续操作中根据需要取消相应的协程或操作。
	// 这里通过context.WithCancel函数创建了一个新的可取消上下文，并获取了相应的取消函数cancel。
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	// 设置请求的上下文。
	req = req.WithContext(ctx)
	// 发送HTTP请求并获取响应。
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errorx.InvalidRequest
	}

	// 创建一个用于接收响应的通道。
	out := make(chan response.Response, 100)
	go func() {
		select {
		case <-ctx.Done():
			resp.Body.Close()
			close(out)
		}
	}()

	// 如果响应状态码不是200，则认为请求失败。
	if resp.StatusCode != http.StatusOK {
		// 尝试读取响应体的内容。
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// 记录请求失败的详细信息。
		slog.Error("openai Completions error", slog.String("uri", uri), slog.String("payload", payload), slog.String("response", string(body)))
		return nil, err
	}

	go func() {
		var err error
		defer func() {
			if err != nil {
				slog.Error("openai Completions error", slog.String("uri", uri), slog.Any("err", err))
			}
		}()

		if in.Stream {
			err = h.handlerStream(resp.Body, out, cancel)
		} else {
			err = h.handlerResponse(resp.Body, out, cancel)
		}
	}()
	return out, nil
}

// handlerStream 处理流的函数。
// 该函数负责从reader中读取流数据，并通过out通道发送处理结果。
// cancel函数用于在需要时取消处理过程。
// 参数:
//
//	reader: 一个io.ReadCloser接口，用于读取流数据。
//	out: 一个response.Response类型的通道，用于发送处理结果。
//	cancel: 一个函数，调用该函数可以取消处理过程。
//
// 返回值:
//
//	error - 解码过程中可能出现的错误。
func (xfyun) handlerStream(reader io.ReadCloser, out chan response.Response, cancel context.CancelFunc) error {
	defer cancel()

	code := decoder.New(reader)
	for {
		event, err := code.Decode()
		if err != nil {
			if err != io.EOF {
				slog.Error("openai Completions stream error", slog.String("error", err.Error()))
				return err
			}
			break
		}

		if event.Data == "[DONE]" {
			break
		}

		var resp response.Response
		if err := sonic.ConfigDefault.UnmarshalFromString(event.Data, &resp); err != nil {
			panic(err)
		}

		out <- resp
		if len(resp.Choices) > 0 && strings.ToLower(resp.Choices[0].FinishReason) == "stop" {
			break
		}
	}

	return nil
}

// handlerResponse 处理响应的函数。
// 该函数负责从reader中读取响应数据，并通过out通道发送处理结果。
// cancel函数用于在需要时取消处理过程。
// 参数:
//
//	reader: 一个io.ReadCloser接口，用于读取响应数据。
//	out: 一个response.Response类型的通道，用于发送处理结果。
//	cancel: 一个函数，调用该函数可以取消处理过程。
//
// 返回值:
//
//	error - 解码过程中可能出现的错误。
func (xfyun) handlerResponse(reader io.ReadCloser, out chan response.Response, cancel context.CancelFunc) error {
	defer cancel()
	var resp response.Response
	if err := sonic.ConfigDefault.NewDecoder(reader).Decode(&resp); err != nil {
		slog.Error("openai Completions decode error", slog.String("error", err.Error()))
		return err
	}

	out <- resp
	return nil
}
