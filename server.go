package main

import (
	"flag"
	"github.com/kaylyu/kim/core"
)

var serverIp string
var serverPort int

func init() {
	//初始化flag
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址，默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口，默认是8888")
}

//主入口
func main() {

	//命令行解析
	flag.Parse()

	//启动服务
	server := core.NewServer(serverIp, serverPort)
	server.Start()
}
