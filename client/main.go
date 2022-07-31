package main

import "github.com/CodeFish011/chat-server/client/core"

func main() {
	c := core.NewClient("tcp4", "127.0.0.1:8099")
	c.Connect()

}
