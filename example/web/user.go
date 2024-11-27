package web

import (
	"example/biz"
	"example/bsvs/token"

	"github.com/ohko/yafeng"
)

type User struct{}

// Login 用户登陆
// curl 'http://127.0.0.1:8080/user/login' -d '{"Account":"demo", "Password":"demo"}'
// curl 'http://127.0.0.1:8080/user/login' -d '{"Account":"admin", "Password":"123456"}'
func (User) Login(ctx *yafeng.Context) {

	// 用户输入数据结构
	var data struct {
		Account  string
		Password string
	}

	// 解析用户提交数据
	if err := ctx.ParsePostData(&data); err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 用户登陆
	info, err := biz.UserLogin(ctx, data.Account, data.Password)
	if err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 设置token
	tk := token.NewToken(info.ID, ctx.GetRealIP())
	info.Token, err = tk.EnToken()
	if err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 返回登陆信息
	ctx.W.Header().Set("token", info.Token)
	ctx.JsonSuccess(info)
}

// Change 修改密码
func (User) Change(ctx *yafeng.Context) {
	// 获取token
	token := checkToken(ctx)

	// 用户输入数据结构
	var data struct {
		Password string
	}

	// 解析用户提交数据
	if err := ctx.ParsePostData(&data); err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 获取信息
	if err := biz.UserChange(ctx, token.GetUserID(), data.Password); err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 返回登陆信息
	ctx.JsonSuccess("ok")
}

// Info 获取用户信息
func (User) Info(ctx *yafeng.Context) {
	// 获取token
	token := checkToken(ctx)

	// 获取信息
	info, err := biz.UserInfo(ctx, token.GetUserID())
	if err != nil {
		ctx.JsonFailed(err.Error())
		ctx.Close()
	}

	// 返回登陆信息
	ctx.JsonSuccess(info)
}
