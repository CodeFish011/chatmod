package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	Network    string
	RemoteAddr string
	Conn       net.Conn
}

func NewClient(network, remoteAddr string) *Client {
	return &Client{
		Network:    network,
		RemoteAddr: remoteAddr,
	}
}

func (c *Client) Connect() {
	conn, err := net.Dial(c.Network, c.RemoteAddr)
	if err != nil {
		fmt.Printf("Failed to establish connection, err: %s", err.Error())
		return
	}
	c.Conn = conn

	go c.ReadMessage()
	c.SendMessage()
	select {}
}

func (c *Client) ReadMessage() {
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Client read message from server failed, err: ", err.Error())
		}
		fmt.Println(string(buf))
	}
}

func (c *Client) SendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("User input with err:%s\n", err.Error())
			continue
		}
		_, err = c.Conn.Write(input)
		if err != nil {
			fmt.Printf("Client send message to server failed, err: %s\n", err.Error())
			continue
		}
	}

}
