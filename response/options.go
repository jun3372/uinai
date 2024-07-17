package response

type Option func(*Response)

func NewResponse() *Response {
	return &Response{}
}

func WithChoicesAndMessage(choices []Choices, msg Message) Option {
	return func(r *Response) {
		r.Choices = choices
		// if len(r.Choices) > 0 {
		// 	if r.Choices[0].Message == nil {
		// 		r.Choices[0].Message = &msg
		// 	}
		// }
	}
}

func WithChoices(choices []Choices) Option {
	return func(r *Response) {
		r.Choices = choices
	}
}

func WithMessage(msg Message) Option {
	return func(r *Response) {
		if len(r.Choices) > 0 {
			if r.Choices[0].Message == nil {
				r.Choices[0].Message = &msg
			}
		}
	}
}

func WithUsage(usage *Usage) Option {
	return func(r *Response) {
		if usage == nil {
			return
		}

		r.Usage = usage
	}
}
