package core

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//创建一个新用户
func NewUser(conn net.Conn, server *Server) *User {
	//获取用户地址
	userAddr := conn.RemoteAddr().String()
	//初始化结构体
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听消息
	go user.ListenMessage()
	//返回
	return user
}

//上线
func (u *User) Online() {
	//上线
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播
	u.server.BroadCast(u, "上线")
}

//下线
func (u *User) Offline() {
	//上线
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播
	u.server.BroadCast(u, "下线")
}

//发送消息
func (u *User) Send(message string) {
	u.conn.Write([]byte(message))
}

//用户处理消息业务
func (u *User) DoMessage(message string) {
	//解析指令
	if message == "who" { //查询当前在线用户列表
		u.server.mapLock.Lock()
		buf := bytes.Buffer{}
		for _, user := range u.server.OnlineMap {
			buf.WriteString(fmt.Sprintf("[%s]%s 上线\n", user.Addr, user.Name))
		}
		u.Send(buf.String())
		u.server.mapLock.Unlock()
	} else if len(message) > 7 && message[:7] == "rename|" { //更名
		//格式：rename|zhangsan
		newName := strings.Split(message, "|")[1]

		//判断是否已经存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.Send("当前用户名被使用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.Send("您已经更新用户名：" + u.Name + "\n")
		}
	} else if len(message) > 3 && message[:3] == "to|" { //消息发送指定用户
		//格式：to|zhangsan|消息内容
		//获取对方的用户名
		remoteName := strings.Split(message, "|")[1]
		if remoteName == "" {
			u.Send("消息格式不正确，请使用\"to|zhangsan|消息内容\"\n")
			return
		}
		//判断是否已经存在
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.Send("该用户名不存在\n")
			return
		}
		//消息内容
		toMessage := strings.Split(message, "|")[2]
		if toMessage == "" {
			u.Send("无消息内容，请重发\n")
			return
		}
		remoteUser.Send(fmt.Sprintf("%s对您说:%s", u.Name, toMessage))
	} else {
		//广播
		u.server.BroadCast(u, message)
	}
}

//监听当前的User channel,一旦有消息，就直接发给对应的客户端
func (u *User) ListenMessage() {
	for {
		//监听消息
		message := <-u.C
		//发送
		u.conn.Write([]byte(message + "\n"))
	}
}
