package core

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const ORDER = "请输入你要选择使用服务的序号:\n" + "1. 世界聊天\n" + "2. 房间聊天\n" + "3. 私人聊天\n" + "0. 返回主菜单\n"

type Server struct {
	Name            string
	Network         string
	Address         string
	Port            string
	CurrentClientID uint64

	OnlineClient map[string]*Client
	mu           sync.RWMutex
}

func NewServer(name, network, address string, port string) *Server {
	return &Server{
		Name:            network,
		Network:         network,
		Address:         address,
		Port:            port,
		CurrentClientID: 0,
		OnlineClient:    make(map[string]*Client),
		mu:              sync.RWMutex{},
	}
}

func (s *Server) Serve() {
	fmt.Println("server start serve ...")

	// 创建一个 服务器指定端口的监听器
	listener, err := net.Listen(s.Network, s.Address+":"+s.Port)
	if err != nil {
		fmt.Printf("Create listener at %s:%s failed, err:%s \n ", s.Address, s.Port, err.Error())
		return
	}

	// 监听器循环监听
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("New connection access failed, err: %s \n", err.Error())
			continue
		}
		fmt.Printf("There is a new connection %s accessed \n", conn.RemoteAddr())

		// 创建 goroutine 处理链接
		// 每一个链接对应一个 goroutine 处理
		go s.HandleConnection(conn)

	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	// 创建一个客户端
	client := NewClient(s.CurrentClientID, conn)
	defer func(c Client) {
		s.mu.Lock()
		delete(s.OnlineClient, conn.RemoteAddr().String())
		s.mu.Unlock()
		c.Conn.Close()
		fmt.Printf("Client %s has been disconnected\n", conn.RemoteAddr().String())
	}(*client)
	s.CurrentClientID++
	// 将新的链接加入到在线列表
	s.mu.Lock()
	s.OnlineClient[conn.RemoteAddr().String()] = client
	s.mu.Unlock()
	// 第一次发送菜单消息
	s.SendMessage(client, ORDER)
	// 处理消息的读写
	go s.ReadMessage(client)
	select {}

}

func (s *Server) BroadCast(message []byte) {

}

func (s *Server) SendMessage(client *Client, message string) {
	_, err := client.Conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Server send message to client %d failed, err: %s\n", client.ID, err.Error())
	}
}

func (s *Server) ReadMessage(client *Client) {
	for {
		buf := make([]byte, 512)
		n, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Printf("Server read message from client %d failed, err: %s\n", client.ID, err.Error())
			continue
		}
		if n != 0 {
			// 此处实际处理消息
			// fmt.Println(buf[0:n])
			// fmt.Println(n)
			// fmt.Println(string(buf[0:n]))
			switch string(buf[0:n]) {
			case "1":
				s.SendMessage(client, "世界聊天")
			case "2":
				s.RoomChat(client)
			case "3":
				s.P2PChat(client)
			}
		}

	}
}

func (s *Server) OnlineClientToString(c *Client) string {
	var list string
	s.mu.Lock()
	for _, client := range s.OnlineClient {
		if client == c {
			list += fmt.Sprintf("%s(本机)\n", client.Conn.RemoteAddr().String())
		} else {
			list += fmt.Sprintf("%s\n", client.Conn.RemoteAddr().String())
		}
	}
	s.mu.Unlock()
	return list
}

func (s *Server) P2PChat(client *Client) {
	s.SendMessage(client, "你正在进行私人聊天\n")

	s.SendMessage(client, "当前在线用户:\n"+s.OnlineClientToString(client))

	s.SendMessage(client, "请按照如下格式发送私聊消息:\n")
	s.SendMessage(client, "@目标用户#聊天内容\n")
	for {
		buf := make([]byte, 512)
		n, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
		}
		if n > 0 {
			content := string(buf)
			if content == "exit" {
				break
			}
			// 获取第一个符号
			if buf[0] != '@' {
				s.SendMessage(client, "请输入正确的聊天格式\n")
			}
			// 获取目标主机
			part1 := strings.TrimSpace(strings.Split(content, "#")[0])
			message := strings.Split(content, "#")[1]
			target := part1[1:]
			// 判断
			if target == client.Conn.RemoteAddr().String() {
				s.SendMessage(client, "不能将消息发送给本机\n")
			}
			if t, ok := s.OnlineClient[target]; ok {
				s.SendMessage(t, client.Conn.RemoteAddr().String() + "对您说: "+message)
			} else {
				s.SendMessage(client, "目标用户不存在，请检查输入内容\n")
			}
		}
	}
}

func (s *Server) RoomChat(client *Client) {
	s.SendMessage(client, "你进入了 1 号房间")
	for {
		buf := make([]byte, 512)
		n, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
		}
		if n > 0 {
			s.SendMessage(client, "你正在进行房间聊天\n")
		}
	}
}
