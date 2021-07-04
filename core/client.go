package core

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn

	flag int //当前客户端模式
}

//创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net.Dail error:", err)
	}
	client.conn = conn
	return client

}

//菜单
func (c *Client) Menu() bool {
	fmt.Println("1：公聊模式")
	fmt.Println("2：私聊模式")
	fmt.Println("3：更新用户名")
	fmt.Println("4：查看在线用户")
	fmt.Println("5：退出")

	var flag int
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 5 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字<<<<<")
		return false
	}
}

//响应处理
func (c *Client) DealResponse() {
	//一旦有数据，拷贝后就标准输出，永久监听
	io.Copy(os.Stdout, c.conn)
}

//业务
func (c *Client) Run() {
	defer c.conn.Close()
	for c.flag != 0 {
		for !c.Menu() {
		}
		//模式处理
		switch c.flag {
		case 1: //公聊模式
			fmt.Println(">>>>>公聊模式<<<<<")
			c.PublicChat()
			break
		case 2: //私聊模式
			fmt.Println(">>>>>私聊模式<<<<<")
			c.PrivateChat()
			break
		case 3: //更新用户名
			fmt.Println(">>>>>更新用户名<<<<<")
			c.UpdateName()
			break
		case 4: //更新用户名
			fmt.Println(">>>>>查看在线用户<<<<<")
			c.Who()
			break
		default:
			fmt.Println(">>>>>退出<<<<<")
			return //runtime.Goexit()
		}
	}
}

//公聊
func (c *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//消息不为空
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
		fmt.Scanln(&chatMsg)
	}
}

//私聊
func (c *Client) PrivateChat() {
	//查询用户列表
	c.Who()
	//指定聊天对象
	var remoteName string
	fmt.Println(">>>>>请输入聊天用户名，exit退出<<<<<")
	fmt.Scanln(&remoteName)
	//进行聊天
	for remoteName != "exit" {
		var chatMsg string
		fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			//消息不为空
			if len(chatMsg) != 0 {
				sendMsg := fmt.Sprintf("to|%s|%s\n\n", remoteName, chatMsg)
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error:", err)
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
			fmt.Scanln(&chatMsg)
		}
		//查询用户列表
		c.Who()
		//指定聊天对象
		remoteName = ""
		fmt.Println(">>>>>请输入聊天用户名，exit退出<<<<<")
		fmt.Scanln(&remoteName)
	}

}

//更名
func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名<<<<<")
	fmt.Scanln(&c.Name)

	sendMessage := fmt.Sprintf("rename|%s\n", c.Name)
	_, err := c.conn.Write([]byte(sendMessage))
	if err != nil {
		fmt.Println("conn.Write error：", err)
		return false
	}
	return true
}

//查看在线用户
func (c *Client) Who() bool {
	sendMessage := "who\n"
	_, err := c.conn.Write([]byte(sendMessage))
	if err != nil {
		fmt.Println("conn.Write error：", err)
		return false
	}
	return true
}
