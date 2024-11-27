// 接口定义文件
package config

// token接口
type IToken interface {
	GetUserID() int             // 返回登陆用户ID
	EnToken() (string, error)   // 生成token
	DeToken(token string) error // 解析token
}
