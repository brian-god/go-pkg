package i18n

import (
	"context"
	"github.com/brian-god/go-pkg/configs"
	"github.com/cloudwego/hertz/pkg/app"
	hertzI18n "github.com/hertz-contrib/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func Handler() app.HandlerFunc {
	return hertzI18n.Localize(
		hertzI18n.WithBundle(&hertzI18n.BundleCfg{
			RootPath:         configs.GwtI18RootPath(),
			AcceptLanguage:   []language.Tag{language.Chinese, language.English, language.TraditionalChinese},
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "yaml",
			UnmarshalFunc:    yaml.Unmarshal,
		}),
		hertzI18n.WithGetLangHandle(func(c context.Context, ctx *app.RequestContext, defaultLang string) string {
			lang := ctx.GetHeader("lang")
			if len(lang) == 0 {
				return defaultLang
			}
			return string(lang)
		}),
	)
}
