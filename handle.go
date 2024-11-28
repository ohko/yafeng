package yafeng

import (
	"embed"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type MWS func(func(ctx *Context)) func(*Context)

func (o *YaFeng) HandleFuncs(path string, v any, mws ...MWS) {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fun := reflect.ValueOf(v).MethodByName(method.Name)
		name := strings.ToLower(method.Name[:1]) + method.Name[1:]

		m := fun.Interface().(func(*Context))
		if o.cross != "" {
			m = MiddlewareCross(o.cross, m)
		}
		for i := len(mws) - 1; i >= 0; i-- {
			m = mws[i](m)
		}

		o.m.HandleFunc(path+name, Middleware(m))
	}
}

func (o *YaFeng) HandleFS(path, prefix string, fsFolder embed.FS) {
	path = strings.TrimRight(path, "/")
	o.m.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + strings.TrimPrefix(r.URL.Path, path)
		if _, err := fsFolder.Open(r.URL.Path); os.IsNotExist(err) {
			r.URL.Path = prefix + "/"
		}
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "index.html")
		http.FileServer(http.FS(fsFolder)).ServeHTTP(w, r)
	})
}

func (o *YaFeng) HandleStatic(path, folder string) {
	path = strings.TrimRight(path, "/")
	folder = strings.TrimRight(folder, "/")
	o.m.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
		realUrlPath := strings.TrimPrefix(r.URL.Path, path)
		if _, err := os.Stat(folder + realUrlPath); os.IsNotExist(err) {
			http.ServeFile(w, r, folder+"/index.html")
		} else {
			http.ServeFile(w, r, folder+realUrlPath)
		}
	})
}

func (o *YaFeng) HandleProxy(path, uri string) {
	path = strings.TrimSuffix(path, "/")
	http.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
		targetUrl, _ := url.Parse(uri)
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host
			req.URL.Path = strings.TrimPrefix(req.URL.Path, path)
			if req.URL.Scheme == "wss:" {
				req.Header.Set("Connection", "upgrade")
				req.Header.Set("Upgrade", "websocket")
			}
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusBadGateway)
		}
		proxy.ServeHTTP(w, r)
	})
}

func Middleware(f func(*Context)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)

		defer func() {
			if perr := recover(); perr != nil {
				switch perr.(type) {
				case YFError:
				default:
					log.Println(perr)
					for i := 5; i >= 1; i-- {
						_, file, line, _ := runtime.Caller(i)
						log.Printf("%s | %s:%d\n", ctx.FlowID, file, line)
					}
				}
			}
		}()

		f(ctx)
	}
}

func MiddlewareCross(domain string, f func(ctx *Context)) func(ctx *Context) {
	return func(ctx *Context) {
		ctx.W.Header().Set("Access-Control-Allow-Origin", domain)
		ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.W.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")

		f(ctx)
	}
}

func TestMiddleware(f func(ctx *Context), mws ...MWS) func(w http.ResponseWriter, r *http.Request) {
	m := f
	for i := len(mws) - 1; i >= 0; i-- {
		m = mws[i](m)
	}
	return Middleware(m)
}
