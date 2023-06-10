package plugins

import "github.com/pkg/errors"

var (
	// ErrInvalidPluginReq 错误的输入参数
	ErrInvalidPluginReq = errors.New("invalid plugin request")
)
