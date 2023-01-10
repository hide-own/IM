package main

import (
	"fmt"
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
	//fmt.Println("链接建立成功")

	//	用户上线
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	//user.ListMessage()
	this.mapLock.Unlock()

	//	广播当前用户上线消息
	this.BroadCast(user, "已上线")

	//	当前header阻塞
	select {}
}
