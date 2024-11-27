// web的中间件
package web

import (
	"example/config"

	"github.com/ohko/yafeng"
)

// 数据库事务
func MiddlewareDBBegin(f func(ctx *yafeng.Context)) func(ctx *yafeng.Context) {
	return func(ctx *yafeng.Context) {
		ctx.Tx = config.DB.Begin()
		defer ctx.Tx.Rollback()

		f(ctx)

		ctx.Tx.Commit()
	}
}

// 检查登陆
func MiddlewareCheckToken(f func(ctx *yafeng.Context)) func(ctx *yafeng.Context) {
	return func(ctx *yafeng.Context) {
		if ctx.R.URL.Path != "/api/user/login" {
			checkAuth(ctx)
		}

		f(ctx)
	}
}
