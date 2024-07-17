package client

import (
	"context"

	"github.com/jun3372/uinai/request"
	"github.com/jun3372/uinai/response"
)

// IClient 接口定义了客户端应提供的完成请求的方法。
// 它允许客户端根据一定的选项和上下文，发起请求并获取响应。
type IClient interface {
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
	Completions(opt Options, ctx context.Context, in request.Request) (chan response.Response, error)
}
