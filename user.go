package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string `commit:"用户地址"`
	C    chan string
	conn net.Conn `commit:"用户链接"`

	Server *Server
}

// NewUser 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		Server: server,
	}

	go user.ListMessage()

	return user
}

// Online 用户上线
func (this *User) Online() {
	this.Server.mapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.mapLock.Unlock()

	//	广播用户上线
	this.Server.BroadCast(this, "已上线")
}

// Offline 用户下线
func (this *User) Offline() {
	this.Server.mapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.mapLock.Unlock()

	//	广播上线
	this.Server.BroadCast(this, "下线")
}

// DoMessage 用户处理消息
func (this *User) DoMessage(msg string) {
	if msg == "w" {
		//查询用户
		this.Server.mapLock.Lock()
		for _, user := range this.Server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.sendMsg(onlineMsg)
		}
		this.Server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := this.Server.OnlineMap[newName]

		if ok {
			this.sendMsg("当前用户名称已使用\n")
		} else {
			this.Server.mapLock.Lock()
			delete(this.Server.OnlineMap, this.Name)
			this.Server.OnlineMap[newName] = this
			this.Server.mapLock.Unlock()

			this.Name = newName
			this.sendMsg("您已经更新用户名:" + this.Name + "\n")
		}
	} else {
		this.Server.BroadCast(this, msg)
	}
}

// 发送消息的处理
func (this *User) sendMsg(msg string) {
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("")
		return
	}
}

// ListMessage 监听当前用户频道
func (this *User) ListMessage() {
	for {
		msg := <-this.C
		_, err := this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Printf("Error writing:%s", err)
			return
		}
	}
}
