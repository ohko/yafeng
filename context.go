package yafeng

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrClose = YFError(errors.New("CLOSE"))
)

type YFError error

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Ctx    context.Context
	FlowID string
	Tx     *gorm.DB
}

type JSONData struct {
	No   int `json:"no"` // 0=无错/大于0业务错误/小于0框架内错误
	Data any `json:"data"`
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r, Ctx: context.Background(), FlowID: GenerateNonce(6)}
}

func (o Context) SetAuthorization(token string) {
	o.W.Header().Set("Authorization", "Bearer "+token)
}

func (o Context) GetAuthorization() string {
	return strings.TrimPrefix(o.R.Header.Get("Authorization"), "Bearer ")
}

func (o Context) ParsePostData(v any) error {
	if o.R.Body == nil {
		return errors.New("body empty")
	}

	bs, err := io.ReadAll(o.R.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}

	return nil
}

func (o Context) JsonOrigin(data any) {
	bs, err := json.Marshal(data)
	if err != nil {
		o.Json(-1, err.Error())
		return
	}
	o.W.Header().Set("x-flowid", o.FlowID)
	o.W.Header().Set("Content-Type", "application/json")
	o.W.Write(bs)
}

func (o Context) Json(no int, data any) {
	bs, err := json.Marshal(&JSONData{No: no, Data: data})
	if err != nil {
		o.W.Write([]byte(fmt.Sprintf(`{"no":-1, "data":"%s"}`, err.Error())))
	}
	o.W.Header().Set("x-flowid", o.FlowID)
	o.W.Header().Set("Content-Type", "application/json")
	o.W.Write(bs)
}

func (o Context) JsonSuccess(data any) {
	o.Json(0, data)
}

func (o Context) JsonFailed(data any) {
	o.Json(1, data)
}

func (o Context) HTML(data string) {
	o.W.Header().Set("x-flowid", o.FlowID)
	o.W.Header().Add("Content-Type", "text/html; charset=UTF-8")
	o.W.Header().Add("Content-Length", strconv.Itoa(len(data)))
	o.W.Write([]byte(data))
}

// IsAjax 是否是ajax请求
func (o Context) IsAjax() bool {
	if o.R.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
		strings.Contains(o.R.Header.Get("Accept"), "application/json") {
		return true
	}
	return false
}

func (o Context) GetRealIP() string {
	if o.R.Header.Get("Ali-Cdn-Real-Ip") != "" {
		return o.R.Header.Get("Ali-Cdn-Real-Ip")
	}
	if o.R.Header.Get("X-Forwarded-For") != "" {
		return o.R.Header.Get("X-Forwarded-For")
	}
	// ipv6
	if host, _, err := net.SplitHostPort(o.R.RemoteAddr); err == nil {
		return host
	}
	// ipv4
	return strings.Split(o.R.RemoteAddr, ":")[0]
}

func (o Context) Close() {
	panic(ErrClose)
}
