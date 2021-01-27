package main

import (
	"bufio"
	"net"
	"strconv"
	"sync"
	"time"

	log "github.com/lanewebchina/Go-000/Week09/log"
)

/*
  作业
  用 Go 实现一个 tcp server ，用两个 goroutine 读写 conn，
  两个 goroutine 通过 chan 可以传递 message，能够正确退出
*/
var (
	ID     int
	lockID sync.Mutex
)

func generateID() int {
	lockID.Lock()
	defer lockID.Unlock()
	ID++
	return ID
}

type User struct {
	id        int
	addr      string
	createAt  time.Time
	messageCh chan string
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		log.Infof("sendMessage = %v, msg = %v", conn, msg)
	}
}

func readMessage(conn net.Conn) {
	defer conn.Close()
	user := &User{
		id:        generateID(),
		addr:      conn.RemoteAddr().String(),
		createAt:  time.Now(),
		messageCh: make(chan string, 8),
	}

	//run write conn goroutine
	go sendMessage(conn, user.messageCh)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		log.Infof("readMessage userID= %v, input = %v", strconv.Itoa(user.id), input.Text())
		user.messageCh <- input.Text()
	}

	if err := input.Err(); err != nil {
		log.Errorf("Read goroutine error=%v", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("err = %v", err)
		}
		//run read conn goroutine
		go readMessage(conn)
	}
}
