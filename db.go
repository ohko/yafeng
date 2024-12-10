package yafeng

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// "./sqlite.db"
// postgres://postgres:postgres@postgres/gorm?sslmode=disable&TimeZone=Asia/Shanghai
// mysql://gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local
func NewDB(dsn string, models ...interface{}) (*gorm.DB, error) {
	sqlLevel, _ := strconv.Atoi(GetEnv("SQL_LEVEL", "0"))
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.LogLevel(sqlLevel), // 0=Silent/1=Error/2=Warn/3=Info
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	options := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",                                // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false,                             // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	}

	var dial gorm.Dialector
	if strings.HasPrefix(dsn, "postgres://") {
		dial = postgres.Open(dsn)
	} else if strings.HasPrefix(dsn, "mysql://") {
		dial = mysql.New(mysql.Config{
			// DSN:                       "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
			DSN:                       strings.TrimPrefix(dsn, "mysql://"), // data source name
			DefaultStringSize:         256,                                 // default size for string fields
			DisableDatetimePrecision:  true,                                // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,                                // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,                                // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,                               // auto configure based on currently MySQL version
		})
	} else {
		dial = sqlite.Open(dsn)
	}
	db, err := gorm.Open(dial, &options)
	if err != nil {
		return nil, err
	}

	db.Use(dbresolver.Register(dbresolver.Config{ /* xxx */ }).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200))

	// 初始化数据库
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("DEBUG") != "" {
		return db.Debug(), nil
	}
	return db, nil
}
