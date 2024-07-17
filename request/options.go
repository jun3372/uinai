package request

type Option func(*Request)

// WithTopP 设置请求中的 TopP 参数，用于控制生成文本的多样性
func WithTopP(top_p float32) Option {
	return func(r *Request) {
		r.TopP = top_p // 将 TopP 参数设置到请求对象中
	}
}

// WithFrequencyPenalty 设置请求中的 FrequencyPenalty 参数，用于惩罚频繁出现的词
func WithFrequencyPenalty(frequencyPenalty float32) Option {
	return func(r *Request) {
		r.FrequencyPenalty = frequencyPenalty // 将 FrequencyPenalty 参数设置到请求对象中
	}
}

// WithPresencePenalty 设置请求中的 PresencePenalty 参数，用于惩罚过于新颖的词
func WithPresencePenalty(presencePenalty float32) Option {
	return func(r *Request) {
		r.PresencePenalty = presencePenalty // 将 PresencePenalty 参数设置到请求对象中
	}
}

// WithStop 设置请求中的 Stop 参数，用于指定生成文本时需要避免的词汇列表
func WithStop(stop []string) Option {
	return func(r *Request) {
		r.Stop = stop // 将 Stop 参数设置到请求对象中
	}
}

// WithStream 是一个配置选项，用于设置请求是否以流式方式处理。
// 参数 value 为 true 时，表示启用流式处理；为 false 时，表示不启用。
func WithStream(value bool) Option {
	// 返回一个匿名函数，该函数接收一个指向 Request 结构体的指针，
	// 并根据 value 的值来设置 Request 的 Stream 字段。
	return func(r *Request) {
		r.Stream = value
	}
}

// WithMessages 设置请求中的 Messages 参数，用于提供上下文信息
func WithMessages(messages []Messages) Option {
	return func(r *Request) {
		r.Messages = messages // 将 Messages 参数设置到请求对象中
	}
}

// WithModel 设置请求中的 Model 参数，用于指定使用的模型
func WithModel(model string) Option {
	return func(r *Request) {
		r.Model = model // 将 Model 参数设置到请求对象中
	}
}

// WithEndpoint 设置请求的端点。
func WithEndpoint(endpoint string) Option {
	// 返回一个函数，该函数接收一个指向 Request 结构体的指针，并设置其 Endpoint 字段。
	return func(r *Request) {
		r.Endpoint = endpoint // 将传入的端点字符串赋值给 Request 的 Endpoint 字段。
	}
}
