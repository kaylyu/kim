package core

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//线上用户列表
	OnlineMap map[string]*User
	//读写锁
	mapLock sync.RWMutex
	//广播队列
	Message chan string
}

//创建一个Server的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

//处理业务
func (s *Server) handler(conn net.Conn) {
	user := NewUser(conn, s)
	//上线
	user.Online()
	fmt.Printf("%s连接建立成功\n", user.Name)

	//监听用户活跃度
	isLive := make(chan bool)

	//接收客户端消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			//提取用户信息,去除'\n'
			message := string(buf[:n-1])

			//针对用户message处理
			user.DoMessage(message)

			//用户任意消息代表用户还在活跃
			isLive <- true
		}
	}()

	//当前协程阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃，重置定时器
			//不做任何事情，只为了激活select，更新下面的定时器
			break
		case <-time.After(600 * time.Second):
			//已经超时
			user.Send("您被踢下线\n")
			//销毁资源
			close(user.C)
			//关闭连接
			conn.Close()
			//退出handler
			return //runtime.Goexit()
		}
	}
}

//广播消息
func (s *Server) BroadCast(user *User, message string) {
	//
	sendMessage := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, message)

	//发送
	s.Message <- sendMessage
}

//监听message channel，发送给用户列表的channel
func (s *Server) ListenMessage() {
	for {
		message := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- message
		}
		s.mapLock.Unlock()
	}
}

//启动服务
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
	}
	//close listen socket
	defer listener.Close()

	fmt.Printf("服务端口%d启动成功...\n", s.Port)

	//监听Message channel
	go s.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}
		//do handler
		go s.handler(conn)
	}
}
