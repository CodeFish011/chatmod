package main

import "github.com/CodeFish011/chat-server/server/core"

func main() {
	s := core.NewServer("main", "tcp4", "127.0.0.1", "8099")
	s.Serve()
}
