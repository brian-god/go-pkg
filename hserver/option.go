package hserver

import (
	"github.com/brian-god/go-pkg/token"
)

// Option 定义一个函数类型，用于修改Server配置
type Option func(*Service)

// WithTokenizer 设置token
func WithTokenizer(token token.IToken) Option {
	return func(a *Service) {
		a.Tokenizer = token
	}
}

//// WithConfigs 设置config
//func WithConfigs(config *configs.Bootstrap) Option {
//	return func(a *Service) {
//		a.config = config
//	}
//}

// WithEnv 设置环境变量
func WithEnv(env string) Option {
	return func(a *Service) {
		a.Env = env
	}
}
