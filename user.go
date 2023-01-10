package main

import (
	"net"
)

type User struct {
	Name string
	Addr string `commit:"用户地址"`
	C    chan string
	conn net.Conn `commit:"用户链接"`
}

// ListMessage 监听当前用户频道
func (this *User) ListMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// NewUser 创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListMessage()

	return user
}
