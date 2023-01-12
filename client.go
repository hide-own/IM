package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

var serverIp string
var serverPort int

// 比main先执行
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器默认ip(127.0.0.1)")
	flag.IntVar(&serverPort, "port", 1234, "服务器默认ip(1234)")
}

func main() {
	// 解析命令行
	flag.Parse()

	client := NewClient("127.0.0.1", 1234)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败 >>>>>")
		return
	}

	fmt.Println(">>>>> 连接服务器成功 >>>>>")

	select {}
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//链接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}
