package uniai

import (
	"context"
	"strings"
	"sync"

	"github.com/jun3372/uinai/internal/client"
	"github.com/jun3372/uinai/internal/client/openai"
	"github.com/jun3372/uinai/internal/client/xfyun"
	"github.com/jun3372/uinai/request"
	"github.com/jun3372/uinai/response"
)

// IUinai 接口定义了UI nai需要实现的方法
type IUinai interface {
	Completions(ctx context.Context, in request.Request) (chan response.Response, error) // 该方法用于处理请求并返回补全结果
}

// uinai 结构体实现了IUinai接口
type uinai struct {
	opts   *client.Options // 客户端选项配置
	client client.IClient  // 客户端接口实例
	onces  sync.Once       // 用于确保某些操作只执行一次
}

// New 函数用于创建一个新的 uinai 实例
// 它接受一个可变参数列表 opts，每个元素都是一个 Option 类型的函数
// 这些函数会依次应用到新创建的 uinai 实例的 opts 字段上
func New(opts ...client.Option) IUinai {
	// 创建一个新的 uinai 实例，使用结构体 uinai 和指针 opts
	resp := &uinai{
		opts: client.NewOptions(opts...), // 初始化 opts 为空 Options
	}

	return resp // 返回配置好的 uinai 实例
}

func (u *uinai) Completions(ctx context.Context, in request.Request) (chan response.Response, error) {
	if u.client == nil {
		u.onces.Do(func() {
			switch strings.ToLower(u.opts.Type) {
			case client.Xfyun:
				u.client = xfyun.NewClient()
			case client.OpenAI, "":
				u.client = openai.NewClient()
			default:
				panic("invalid client type")
			}
		})
	}

	return u.client.Completions(*u.opts, ctx, in)
}
