package biz

import (
	"example/config"
	"os"
	"testing"

	"github.com/ohko/yafeng"
)

func init() {
	os.Setenv("DSN", "../sqlite.db")
	config.DBInit()
}

// go test -timeout 15s -run ^TestUserLogin$ example/biz -v -count=1
func TestUserLogin(t *testing.T) {
	ctx := yafeng.NewContext(nil, nil)
	ctx.Tx = config.DB.Begin()
	defer ctx.Tx.Rollback()

	{
		info, err := UserLogin(ctx, "demo", "1")
		if err == nil {
			t.Fatal(info)
		}
	}

	{
		info, err := UserLogin(ctx, "demo", "demo")
		if err != nil {
			t.Fatal(err)
		}
		if info.Account != "demo" {
			t.Fatal(info)
		}
	}
}

// go test -timeout 15s -run ^TestUserInfo$ example/biz -v -count=1
func TestUserInfo(t *testing.T) {
	ctx := yafeng.NewContext(nil, nil)
	ctx.Tx = config.DB.Begin()
	defer ctx.Tx.Rollback()

	{
		info, err := UserInfo(ctx, 1)
		if err != nil {
			t.Fatal(info)
		}

		if info.Account != "demo" {
			t.Fatal(info)
		}
	}
}
