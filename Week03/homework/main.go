package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"week03/errgroup"
)

/*
作业:
 基于 errgroup 实现一个http server的启动和关闭 ，
 以及 linux signal 信号的注册和处理，
 要保证能够一个退出，全部注销退出。
*/

func main() {

	//通知主goroutine其它server已关闭
	shutdown := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := errgroup.WithContext(ctx)

	server1 := http.Server{
		Addr: ":3000",
	}
	g.Go(func(context.Context) error {
		if err := server1.ListenAndServe(); err != nil {
			cancel()
			return err
		}
		log.Println("server1 runing")
		return nil
	})

	server2 := http.Server{
		Addr: ":3001",
	}
	g.Go(func(context.Context) error {
		if err := server2.ListenAndServe(); err != nil {
			cancel()
			return err
		}
		log.Println("server2 runing")
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			select {
			case s := <-c:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					cancel()
				}
			}
		}
	}()

	// context cancel后，关闭全部并发的 http server，全部关闭完成后通知主goroutine
	go func() {
		<-ctx.Done() //阻塞等待signal信号
		go func() {
			log.Println(ctx.Err())
			if err := server1.Shutdown(context.Background()); err != nil {
				log.Printf("server1 shutdown err: %v\n", err)
			}
			if err := server2.Shutdown(context.Background()); err != nil {
				log.Printf("server2 shutdown err: %v\n", err)
			}
			log.Println("servers graceful shutdown")
			close(shutdown)
			return
		}()
		log.Println("servers no-graceful shutdown")
		<-time.After(time.Minute * 1)
		close(shutdown)
		return
	}()

	if err := g.Wait(); err != nil {
		log.Printf("g err = %v", err)
	}

	log.Println("server stopped")
}
