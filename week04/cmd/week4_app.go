package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/wire"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"go-geekbang/week04/biz"
	"go-geekbang/week04/data"
	"go-geekbang/week04/service"
)

var (
	dbUrl = "gopher:studyhard@tcp(localhost:3306)/go_training?charset=utf8"

	userHandler *service.UserHandler
)

func init() {
	userHandler_ := InitHandler()
	userHandler = &userHandler_
}

func serve(addr string, handler http.Handler, stop chan string) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		reason := <-stop
		s.Shutdown(context.Background())
		fmt.Println("server exit:", reason)
	}()
	fmt.Println("server listen:", addr)
	return s.ListenAndServe()
}

func registSignal(stop chan string) error {
	signalCh := make(chan os.Signal, 1)
	fmt.Println("regist system signal......")
	signal.Notify(signalCh, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	signal := <-signalCh
	signalMsg := signal.String()
	stop <- signalMsg
	fmt.Println("receive system signal:", signalMsg)
	return errors.Errorf("receive system signal: %s", signalMsg)
}

func main() {
	stop := make(chan string)
	var g errgroup.Group
	g.Go(func() error {
		return serve(":8080", userHandler, stop)
	})
	g.Go(func() error {
		return registSignal(stop)
	})
	if err := g.Wait(); err != nil {
		fmt.Println("server broken with error: ", err)
	}
}

func InitHandler() service.UserHandler {
	wire.Build(InitDB(), data.NewDao, biz.NewUserService, service.NewUserHandler)
	return service.UserHandler{}
}

func InitDB() *sql.DB {
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		fmt.Println("open db failed:", err)
		panic(err)
	}
	fmt.Println("connect db ok")
	return db
}
