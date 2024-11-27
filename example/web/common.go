// web域需要的公共函数
package web

import (
	"example/bsvs/token"
	"example/config"

	"github.com/ohko/yafeng"
)

func checkToken(ctx *yafeng.Context) config.IToken {
	tk := token.NewToken(0, "")
	if err := tk.DeToken(ctx.R.Header.Get("token")); err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}
	return tk
}
