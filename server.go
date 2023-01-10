package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//	当前用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//	消息广播
	Message chan string
}

// NewSever  创建server接口
func NewSever(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// BroadCast 广播
func (this *Server) BroadCast(user *User, msg string) {
	this.Message <- "[" + user.Addr + "]" + user.Name + ":" + msg
}

// Start 启动服务器的方法
func (this *Server) Start() {
	//	监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listener error:", err)
		return
	}

	//	防止回调
	defer listener.Close()

	go this.ListenMessage()

	// 循环监听下一个链接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener Accept error:", err)
			continue
		}
		go this.Header(conn)
	}
}

// ListenMessage 监听Message进行广播
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		//消息全局发送
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Header(conn net.Conn) {
	//....当前链接啊的业务

	//	用户上线
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	//user.ListMessage()
	this.mapLock.Unlock()

	//	广播当前用户上线消息
	this.BroadCast(user, "已上线")

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			length, err := conn.Read(buf)
			if length == 0 {
				this.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read error: ", err)
				return
			}

			//	除去消息的"\n"全局广播
			msg := string(buf[:length])
			fmt.Println(msg)
			this.BroadCast(user, msg)
		}
	}()

	//	当前header阻塞
	select {}
}
