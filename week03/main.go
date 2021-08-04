package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "index")
	})

	shutdown := make(chan struct{})
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "http server manual shutdown...")
		shutdown <- struct{}{}
	})

	server := http.Server{
		Addr: ":8080",
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup context exit...")
		case <-shutdown:
			log.Println("http server manual shutdown...")
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Println("shutting down server...")
		return server.Shutdown(timeoutCtx)
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-c:
			return errors.Errorf("receive system signal: %v", sig)
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup exit: %v\n", err)
	}
}
