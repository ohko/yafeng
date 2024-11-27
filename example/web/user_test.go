package web

import (
	"encoding/json"
	"example/config"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ohko/yafeng"
	"github.com/valyala/fastjson"
)

func init() {
	os.Setenv("DSN", "../sqlite.db")
	config.DBInit()
}

func parseResponse(bs []byte) *yafeng.JSONData {
	var rs yafeng.JSONData
	if err := json.Unmarshal(bs, &rs); err != nil {
		log.Fatal(err)
	}
	return &rs
}

// go test -timeout 15s -run ^TestUser_Login$ example/web -v -count=1
func TestUser_Login(t *testing.T) {
	user := &User{}

	for i, v := range []struct {
		Input string
		Want  string
	}{
		{Input: ``, Want: "unexpected end of JSON input"},
		{Input: `{}`, Want: "Account/Password is empty"},
		{Input: `{"Account":"demo", "Password":"1"}`, Want: "登陆失败"},
		{Input: `{"Account":"demo", "Password":"demo"}`},
	} {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/api/user/login", strings.NewReader(v.Input))
		if err != nil {
			t.Fatal(i, err)
		}

		http.HandlerFunc(yafeng.TestMiddleware(user.Login, MiddlewareDBBegin, MiddlewareCheckToken)).ServeHTTP(w, r)
		rs := parseResponse(w.Body.Bytes())
		if rs.No != 0 {
			if rs.Data != v.Want {
				t.Fatal(i, rs)
			}
		} else {
			token, ok := rs.Data.(map[string]any)["Token"]
			if !ok || token == "" {
				t.Fatal(i, rs)
			}
			// t.Log("Token:", token)
		}
	}
}

// go test -timeout 15s -run ^TestUser_Info$ example/web -v -count=1
func TestUser_Info(t *testing.T) {
	user := &User{}

	for i, v := range []struct {
		Token string
		Want  string
	}{
		{Token: ``, Want: "need token"},
		{Token: `123`, Want: "token error"},
		{Token: `5d7d96e4e7c526bcb4b1fcc6dc0de84fadf399bc6f51564efef7651977031b50710c898f`, Want: "请先登陆"},
	} {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/api/user/info", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Token", v.Token)

		http.HandlerFunc(yafeng.TestMiddleware(user.Info, MiddlewareDBBegin, MiddlewareCheckToken)).ServeHTTP(w, r)
		p, err := (&fastjson.Parser{}).ParseBytes(w.Body.Bytes())
		if err != nil {
			t.Fatal(i, err)
		}
		if p.GetInt("no") != 0 {
			if string(p.GetStringBytes("data")) != v.Want {
				t.Fatal(i, w.Body.String())
			}
		} else {
			if p.GetInt("data", "ID") != 1 {
				t.Fatal(i, w.Body.String())
			}
			// t.Log("ID:", p.GetInt("data", "ID"))
		}
	}
}
