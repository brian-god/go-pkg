package jwt

import (
	"errors"
	"fmt"
	"github.com/brian-god/go-pkg/constant"
	"github.com/brian-god/go-pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	NotFoundError = errors.New("not found from context")
)

// ParseToken 解析TOKEN
func ParseToken(c *app.RequestContext) (*token.AccessToken, error) {
	val, ok := c.Get(constant.KeyAccessToken)
	if !ok {
		return nil, NotFoundError
	}

	accessToken, ok := val.(token.AccessToken)
	if !ok {
		return nil, fmt.Errorf("parse error")
	}

	return &accessToken, nil
}
