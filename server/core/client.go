package core

import (
	"fmt"
	"net"
)

// 作为服务端操作的客户端对象
type Client struct {
	ID        uint64
	Conn      net.Conn
	OrderChan chan uint
	IsAlive   chan bool
}

func NewClient(id uint64, conn net.Conn) *Client {
	return &Client{
		Conn:      conn,
		ID:        id,
		OrderChan: make(chan uint),
		IsAlive:   make(chan bool),
	}
}

func (c *Client) ReadMessage() {
	buf := make([]byte, 512)
	_, err := c.Conn.Read(buf)
	if err != nil {
		fmt.Printf("Client read from conn failed, err: %s", err.Error())
		return
	}

	fmt.Println(string(buf))
}
