package yafeng

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func NewDB(models ...interface{}) (*gorm.DB, error) {
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

	dbPath := GetEnv("DSN", "./sqlite.db")
	var dsn gorm.Dialector
	if strings.HasPrefix(dbPath, "postgres://") {
		dsn = postgres.Open(dbPath)
	} else {
		dsn = sqlite.Open(dbPath)
	}
	db, err := gorm.Open(dsn, &options)
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
