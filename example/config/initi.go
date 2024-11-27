// 系统初始化
package config

import (
	"log"

	"github.com/ohko/yafeng"
)

func DBInit() {
	var err error

	// 连接、初始化数据库
	if DB, err = yafeng.NewDB(&TableUser{}); err != nil {
		log.Fatal(err)
	}

	// 初始化预制数据
	// if err := DB.Save(&TableUser{ID: 1, Account: "demo", Password: yafeng.Hash("demo")}).Error; err != nil {
	// 	log.Fatal(err)
	// }
}
