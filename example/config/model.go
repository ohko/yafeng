// 数据库模型
package config

type TableUser struct {
	UserID   int    `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Account  string `gorm:"column:account;index:check;unique" json:"account"`
	Password string `gorm:"column:password;index:check" json:"-" json:"password"`

	Token string `gorm:"-" json:"token"`
}
