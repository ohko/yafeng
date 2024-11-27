// 接口定义文件
package config

// auth接口
type IAuth interface {
	GetUserID() int            // 返回登陆用户ID
	EnToken() (string, error)  // 生成auth
	DeToken(auth string) error // 解析auth
}
