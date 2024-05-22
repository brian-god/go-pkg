package ratelimit

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var limiter *rate.Limiter

// WithTimeoutHandler 超时时间限定
func WithTimeoutHandler(r int) app.HandlerFunc {
	if limiter == nil && r > 0 {
		limiter = rate.NewLimiter(rate.Limit(r), r)
	}
	return func(ctx context.Context, c *app.RequestContext) {
		if limiter != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			if err := limiter.Wait(ctx); err != nil {
				c.JSON(http.StatusForbidden, utils.H{"code": http.StatusForbidden, "msg": hertzI18n.MustGetMessage(ctx, "ServerBusy")})
				c.Abort()
				return
			}
		}
		c.Next(ctx)
	}
}
