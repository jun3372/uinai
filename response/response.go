package response

// Response 结构体定义了API响应的数据结构
type Response struct {
	ID                string    `json:"id"`                 // 请求的唯一标识符
	Object            string    `json:"object"`             // 响应对象的类型
	Created           int       `json:"created"`            // 响应创建的时间戳
	Model             string    `json:"model"`              // 使用的模型名称
	SystemFingerprint any       `json:"system_fingerprint"` // 系统指纹，用于识别请求来源
	Choices           []Choices `json:"choices"`            // 选项列表，包含用户的选择信息
	Usage             *Usage    `json:"usage"`              // 使用情况统计，包括使用的token数量
}

// Usage 结构体定义了API的使用统计信息
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示部分使用的token数量
	CompletionTokens int `json:"completion_tokens"` // 完成部分使用的token数量
	TotalTokens      int `json:"total_tokens"`      // 总共使用的token数量
}

// Choices 结构体定义了用户选择的详细信息
type Choices struct {
	Index        int      `json:"index"`             // 选择的索引位置
	FinishReason string   `json:"finish_reason"`     // 完成选择的原因
	Message      *Message `json:"message,omitempty"` // 与选择相关的消息内容
	Delta        *Delta   `json:"delta,omitempty"`
}

// Message 结构体定义了消息的内容和角色
type Message struct {
	Role    string `json:"role"`    // 消息的角色，如提示、回答等
	Content string `json:"content"` // 消息的实际内容
}

// Delta 表示一个内容差异的结构体。
// 它用于存储两个版本之间内容的差异，通常用于版本控制系统或编辑器中。
// Content 字段存储了具体的差异内容，以文本形式表示。
type Delta struct {
	Content string `json:"content"`
}
