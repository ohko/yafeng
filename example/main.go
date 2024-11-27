package main

import (
	"embed"
	"example/config"
	"example/web"
	"log"
	"os"
	"runtime"

	"github.com/ohko/yafeng"
)

//go:embed public
var fsFolder embed.FS

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Flags() | log.Lshortfile)

	// 解析.env
	if _, err := yafeng.ReadDotEnv(); err != nil {
		log.Fatal(err)
	}

	// 初始化数据库
	config.DBInit()

	// YaFeng实例
	app := yafeng.New(nil)

	// 设置跨域
	app.SetCross("*")

	// 中间件
	mws := []yafeng.MWS{web.MiddlewareCheckToken, web.MiddlewareDBBegin}

	// user 路由
	app.HandleFuncs("/api/user/", &web.User{}, mws...)

	if os.Getenv("DEBUG") != "" {
		// 静态文件
		app.HandleStatic("/public/", "./public/")
	} else {
		// 嵌入文件
		app.HandleFS("/public/", "public", fsFolder)
	}

	// 启动服务
	app.Start()
}
