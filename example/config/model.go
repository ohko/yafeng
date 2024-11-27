// 数据库模型
package config

type TableUser struct {
	ID       int    `gorm:"user_id;primary_key"`
	Account  string `gorm:"account;index:check;unique"`
	Password string `gorm:"password;index:check" json:"-"`

	Token string `gorm:"-"`
}
