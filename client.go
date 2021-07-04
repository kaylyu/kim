package main

import (
	"flag"
	"fmt"
	"github.com/kaylyu/kim/core"
)

var cserverIp string
var cserverPort int

func init() {
	//初始化flag
	flag.StringVar(&cserverIp, "ip", "127.0.0.1", "设置服务器IP地址，默认是127.0.0.1")
	flag.IntVar(&cserverPort, "port", 8888, "设置服务器端口，默认是8888")
}

//主入口
func main() {
	//命令行解析
	flag.Parse()
	client := core.NewClient(cserverIp, cserverPort)
	if client == nil {
		fmt.Println(">>>>>连接服务器失败<<<<<")
		return
	}
	fmt.Println(">>>>>连接服务器成功<<<<<")
	//响应处理
	go client.DealResponse()
	//业务处理
	client.Run()
}
