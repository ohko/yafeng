package yafeng

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type YaFeng struct {
	s     *http.Server
	m     *http.ServeMux
	cross string
}

func New(s *http.ServeMux) *YaFeng {
	if s == nil {
		s = http.NewServeMux()
	}
	return &YaFeng{m: s}
}

func (o *YaFeng) Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Flags() | log.Lshortfile)

	addr := GetEnv("HTTP_HOST", ":8080")
	o.s = &http.Server{Addr: addr, Handler: o.m}
	log.Println("Listen:", addr)
	go func() {
		if err := o.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	waitShutdown(time.Second*15, o)
}

func (o *YaFeng) SetCross(domain string) { o.cross = domain }

// 等待信号，优雅的停止服务
func waitShutdown(waitTime time.Duration, ss ...*YaFeng) {
	log.Println("wait ctrl+c ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	<-c
	close(c)

	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	for _, s := range ss {
		if err := s.s.Shutdown(ctx); err != nil {
			log.Println(s.s.Addr, err)
		} else {
			log.Println(s.s.Addr, "shutdown ok.")
		}
	}
}
