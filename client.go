package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"flag"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	// flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		Name : hostname,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial err:", err)
		return nil
	}
	client.conn = conn

	return client
}

// 处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

// 更新用户名
func (client *Client) ReName() {
	// fmt.Println("测试名称:", client.Name)
	sendName := "rename:" + client.Name +"\n"
	// fmt.Println(sendName)
	_, err := client.conn.Write([]byte(sendName))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) Run() {

	var chatMsg string
	fmt.Println("<- 直接输入聊天内容 | 输入exit退出 | 输入who查询在线用户 ->")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg)!=0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		// fmt.Println(">>>>>>>>请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
	}
}


var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "10.170.46.245", "设置服务器IP地址")
	flag.IntVar(&serverPort, "port", 7890, "设置服务器端口")
}


func main() {
	// 命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("🔗 >>> 服务器链接失败")
		return
	}

	// 开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	fmt.Println("🔗 >>> 服务器链接成功")
	client.ReName()
	// 启动客户端的业务
	client.Run()
	
}