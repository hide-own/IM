package main

func main() {
	server := NewSever("127.0.0.1", 1234)
	server.Start()
}
