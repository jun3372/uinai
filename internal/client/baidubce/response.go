package baidubce

type Response struct {
	ID               string          `json:"id"`                 // 本轮对话的id
	Object           string          `json:"object"`             // 回包类型 chat.completion：多轮对话返回
	Created          int             `json:"created"`            // 时间戳
	SentenceID       int             `json:"sentence_id"`        // 表示当前子句的序号。只有在流式接口模式下会返回该字段
	IsEnd            bool            `json:"is_end"`             // 表示当前子句是否是最后一句。只有在流式接口模式下会返回该字段
	IsTruncated      bool            `json:"is_truncated"`       // 当前生成的结果是否被截断
	Result           string          `json:"result"`             // 对话返回结果
	SearchInfo       []SearchResults `json:"search_info"`        //搜索数据，当请求参数enable_citation或enable_trace为true，并且触发搜索时，会返回该字段
	NeedClearHistory bool            `json:"need_clear_history"` // 表示用户输入是否存在安全风险，是否关闭当前会话，清理历史会话信息	true：是，表示用户输入存在安全风险，建议关闭当前会话，清理历史会话信息	false：否，表示用户输入无安全风险
	FinishReason     string          `json:"finish_reason"`      // 输出内容标识，说明：	· normal：输出内容完全由大模型生成，未触发截断、替换	· stop：输出结果命中入参stop中指定的字段后被截断	· length：达到了最大的token数，根据EB返回结果is_truncated来截断	· content_filter：输出内容被截断、兜底、替换为**等
	Usage            Usage           `json:"usage"`              // token统计信息
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 问题tokens数
	CompletionTokens int `json:"completion_tokens"` // 回答tokens数
	TotalTokens      int `json:"total_tokens"`      // tokens总数
}

type SearchResults []SearchResult

type SearchResult struct {
	Index int    `json:"index"` // 序号
	Url   string `json:"url"`   // 搜索结果地址
	Title string `json:"title"` // 搜索结果标题
}
