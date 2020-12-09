package main

import (
	"context"
	"fmt"
	"time"
)

//Tracker knows how to track events for the application
type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func main() {
	tr := NewTracker()
	go tr.Run() //把并行的行为交给调用者
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	_ = tr.Event(context.Background(), "test4")
	_ = tr.Event(context.Background(), "test5")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()   //5s后关闭context
	tr.ShutDown(ctx) //执行
}

func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done(): //
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	//在这里只要把t.ch这个通道close掉就退出了
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	//当未超时的时候，通道数据接收完了
	//给Tracker实例的stop对象发消息
	//表示通道里的数据都接收完了
	t.stop <- struct{}{}
}

func (t *Tracker) ShutDown(ctx context.Context) {
	close(t.ch) //关闭ch通道
	//在此select是多路选择
	select {
	case <-t.stop: //在此是监听t.stop, Run里把通道里的数据都读取完了的时候
	case <-ctx.Done(): //这里做的超时处理,Done的意思表示context被cancel掉了
	}
}
