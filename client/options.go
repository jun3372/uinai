package client

import "net/http"

// Options 结构体用于配置选项参数
type Options struct {
	// Type 字段表示选项的类型
	Type string
	// Host 字段表示选项的主机名
	Host string
	// Header 字段表示选项的HTTP头部信息
	Header http.Header
}

// Option 是一个函数类型，用于修改Options结构体
type Option func(*Options)

// NewOptions 创建一个新的 Options 结构体实例，并允许通过可变参数传入 Option 配置。
// Option 是一个函数类型，它接受一个指向 Options 的指针作为参数，用于配置 Options。
// 通过遍历传入的 Option 参数，并对每个 Option 调用，将配置应用到新创建的 Options 实例上。
func NewOptions(opts ...Option) *Options {
	o := &Options{ // 创建一个 Options 结构体指针
		Header: make(http.Header), // 初始化 header 字段为一个空的 http.Header 类型
	}

	o.Header.Add("Content-Type", "application/json")
	o.Header.Add("Accept", "*/*")
	for _, opt := range opts { // 遍历传入的 Option 参数
		opt(o) // 对每个 Option 调用，传入 Options 指针以应用配置
	}
	return o // 返回配置好的 Options 指针
}

// AddHeader 是一个函数，它接受两个字符串参数 key 和 value。
// 它返回一个 Option 类型的函数，这个函数会修改传入的 Options 结构体的 header 字段。
// 当这个返回的函数被调用时，它会将 key 和 value 作为一个键值对添加到 Options 的 header 字段中。
func AddHeader(key, value string) Option {
	return func(o *Options) {
		// 使用 Add 方法将 key 和 value 添加到 header 中
		o.Header.Add(key, value)
	}
}

func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func WithType(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}

// WithHeader 是一个函数，它接受一个 http.Header 类型的参数 header。
// 它返回一个 Option 类型的函数，这个函数会修改传入的 Options 结构体的 header 字段。
// 当这个返回的函数被调用时，它会将传入的 header 赋值给 Options 的 header 字段。
func WithHeader(header http.Header) Option {
	return func(o *Options) { o.Header = header }
}
