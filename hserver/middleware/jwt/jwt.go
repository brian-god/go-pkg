package jwt

import (
	"context"
	"github.com/brian-god/go-pkg/constant"
	"github.com/brian-god/go-pkg/hctx"
	"github.com/brian-god/go-pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"net/http"
	"strings"
)

// Handler 校验的处理器
func Handler(tokenizer token.IToken) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenEmpty)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenEmpty, constant.RespData: utils.H{}})
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenEmpty)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenEmpty, constant.RespData: utils.H{}})
			c.Abort()
			return
		}

		var accessToken token.AccessToken
		if err := tokenizer.Verify(parts[1], &accessToken); err != nil {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenVerifyFail)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenVerifyFail, constant.RespData: utils.H{}})
			c.Abort()
			return
		}
		accessToken.AccessToken = parts[1]
		ctx = hctx.WithUserId(ctx, accessToken.UserId)
		ctx = hctx.WithPlatform(ctx, accessToken.Platform)
		ctx = hctx.WithToken(ctx, accessToken.AccessToken)
		ctx = hctx.WithRole(ctx, accessToken.Role)
		ctx = hctx.WithTenantId(ctx, accessToken.TenantId)
		ctx = hctx.WithOpUserID(ctx, accessToken.UserId)
		// 将身份信息缓存到Context
		c.Set(constant.KeyAccessToken, accessToken)
		c.Next(ctx)
	}
}
