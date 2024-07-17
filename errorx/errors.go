package errorx

import "github.com/pkg/errors"

var (
	InvalidHost    = errors.New("invalid host")
	InvalidInput   = errors.New("invalid input")
	InvalidRequest = errors.New("Invalid Request Error")
	NotFound       = errors.New("not found")
)
