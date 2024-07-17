package request

import (
	"bytes"
	"io"

	"github.com/bytedance/sonic"
)

// Request 结构体定义了一个请求的参数
type Request struct {
	TopP             float32    `json:"top_p,omitempty"`             // TopP 是一个浮点数，表示选择的概率阈值
	FrequencyPenalty float32    `json:"frequency_penalty,omitempty"` // FrequencyPenalty 是一个整数，用于对高频词汇进行惩罚
	PresencePenalty  float32    `json:"presence_penalty,omitempty"`  // PresencePenalty 是一个整数，用于对存在的词汇进行惩罚
	Temperature      float32    `json:"temperature,omitempty"`       // 使用什么采样温度，介于 0 和 2 之间。较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使输出更加集中和确定。 我们通常建议改变这个或top_p但不是两者。
	MaxTokens        int        `json:"max_tokens,omitempty"`        // 最大生成标记数
	Model            string     `json:"model"`                       // Model 是一个字符串，指定了要使用的模型
	Stop             []string   `json:"stop,omitempty"`              // Stop 是一个字符串切片，包含需要过滤停止的词汇
	Messages         []Messages `json:"messages,omitempty"`          // Messages 是一个消息切片，包含了请求中的消息内容
	Stream           bool       `json:"stream,omitempty"`            // 默认为 false 如果设置,则像在 ChatGPT 中一样会发送部分消息增量。标记将以仅数据的服务器发送事件的形式发送,这些事件在可用时,并在 data: [DONE] 消息终止流。Python 代码示例。
	Endpoint         string     `json:"-"`                           // EndPoint 是一个字符串，表示请求的端点
}

// Messages 结构体用于表示一个消息，包含内容和角色信息
type Messages struct {
	Role    string `json:"role"`    // Role 字段表示消息的角色，例如发送者、接收者等
	Content string `json:"content"` // Content 字段表示消息的内容
}

// NewRequest 创建并返回一个新的Request实例。
// 使用者可以通过传递一系列选项函数来定制Request的属性。
// 这种方法允许灵活的配置，而不必直接在结构体初始化时指定所有参数。
func NewRequest(opts ...func(*Request)) *Request {
	// 初始化Request结构体，默认设置了一些基本参数。
	resp := &Request{
		TopP:             0.8,          // 设置默认的 TopP 值为 0.8
		FrequencyPenalty: 0,            // 设置默认的 FrequencyPenalty 值为 0
		PresencePenalty:  0,            // 设置默认的 PresencePenalty 值为 0
		Temperature:      0,            // 设置默认的 Temperature 值为 0
		Stop:             []string{},   // 设置默认的 Stop 切片为空
		Messages:         []Messages{}, // 设置默认的 Messages 切片为空
		Model:            "",           // 设置默认的 Model 值为空字符串
	}

	// 遍历所有传入的选项函数，并依次应用它们到Request实例上。
	// 这些选项函数可能会修改Request的属性，以定制其行为。
	for _, fn := range opts {
		fn(resp)
	}

	// 返回配置好的Request实例。
	return resp
}

// MarshalToString 将请求对象序列化为字符串。
//
// 本方法利用了全局默认配置的sonic配置实例，对当前请求对象进行序列化。
// 它主要适用于需要将请求数据以字符串形式传输或存储的场景。
//
// 返回值:
//
//	string - 序列化后的字符串
//	error - 如果序列化过程中发生错误，则返回该错误
func (r *Request) MarshalToString() (string, error) {
	return sonic.ConfigDefault.MarshalToString(r) // 序列化 Request
}

// Marshal 将Request对象序列化为字节切片。
//
// 本方法利用全局默认配置的sonic实例来执行序列化操作，
// 它重用了配置好的序列化器，避免了每次序列化都创建新实例的开销。
//
// 返回值:
//
//	[]byte - 序列化后的字节切片
//	error - 如果序列化过程中发生错误，则返回该错误
func (r *Request) Marshal() ([]byte, error) {
	return sonic.ConfigDefault.Marshal(r) // 序列化 Request 结构体为字符串
}

// Payload 返回请求的Payload。
//
// 该方法将请求对象序列化为字节流，以便于后续操作如发送网络请求或存储。
// 不处理Marshal方法的错误，因为在此上下文中，错误处理可能不是必要的。
// 如果序列化失败，可能导致返回的io.Reader为空，调用方应对此进行处理。
func (r *Request) Payload() io.Reader {
	body, _ := r.Marshal()       // 序列化 Request 结构体为字符串
	return bytes.NewReader(body) // 将字符串转换为 io.Reader 接口
}

// NewMessage 创建一个新消息对象。
// 参数:
//
//	role: 消息发送者的角色。
//	content: 消息的内容。
//
// 返回值:
//
//	Messages: 一个包含角色和内容的消息对象。
func NewMessage(role, content string) Messages {
	return Messages{Role: role, Content: content}
}

// NewSystemMessage 创建一个系统消息对象。
// 这个函数用于封装一个具有系统角色的消息，方便后续处理和分发。
// 参数:
//
//	content: 消息的内容，它是系统消息的具体信息。
//
// 返回值:
//
//	Messages: 返回一个 Messages 结构体实例，其中 Role 字段设置为 MessageRoleSystem，Content 字段设置为传入的 content。
func NewSystemMessage(content string) Messages {
	// 返回一个 Messages 实例，表示系统消息。
	return Messages{Role: MessageRoleSystem, Content: content}
}

// NewUserMessage 创建一个表示用户消息的 Messages 实例。
// 这个函数是为了简化消息对象的创建，特别用于表示用户发出的消息。
// 参数:
//
//	content - 消息的内容，这是用户实际发送的文字信息。
//
// 返回值:
//
//	Messages - 一个已经设置了消息角色和内容的 Messages 实例。
func NewUserMessage(content string) Messages {
	return Messages{Role: MessageRoleUser, Content: content}
}

// NewAssistantMessage 创建一个表示助手回复的消息对象。
// 这个函数用于封装一个消息，明确表示这是助手的回复消息。
// 参数:
//
//	content - 消息的内容。
//
// 返回值:
//
//	Messages - 包含消息角色和内容的结构体。
func NewAssistantMessage(content string) Messages {
	return Messages{Role: MessageRoleAssistant, Content: content}
}
