// web域需要的公共函数
package web

import (
	"example/bsvs/auth"
	"example/config"

	"github.com/ohko/yafeng"
)

func checkAuth(ctx *yafeng.Context) config.IAuth {
	tk := auth.NewAuth(0, "")
	if err := tk.DeToken(ctx.GetAuthorization()); err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}
	return tk
}
