package main

import (
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

//获得当前接入的用户名
// func (this *User) GetUserName() {
//     hostname, err := os.Hostname()
//     if err != nil {
//         fmt.Println("Error:", err)
//         return 
//     }
// 	newHostname := hostname + "[" + this.Addr + "]"

// 	this.server.mapLock.Lock()
// 	delete(this.server.OnlineMap, this.Name)
// 	this.server.OnlineMap[newHostname] = this
// 	this.server.mapLock.Unlock()
// 	this.Name = newHostname
// }


//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	
	user := &User{
		Name: "用户" + "[" + userAddr + "]",
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server: server,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

//给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

// 更新用户名
func (this *User) ReName(msg string) {
	if len(msg) > 7 && msg[:7] == "rename:" {
		// 更改用户名
		newName := strings.Split(msg, ":")[1] + "[" + this.Addr + "]"
		// fmt.Println(newName)
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.OnlineMap[newName] = this
		this.server.mapLock.Unlock()
		this.Name = newName
		this.SendMsg("您已经更新用户名:" + this.Name)

	} 
}

// 用户的上线业务
func (this *User) Online(conn net.Conn) {
	//
	bufName := make([]byte, 4096)
	n, _ := conn.Read(bufName)
	msg := string(bufName[:n-1])
	this.ReName(msg)
	//

	// 用户上线,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 更新用户名！！！！！
	onlineSend := "-* 🙋" + this.Name + " 已上线 *-"
	this.server.Message <- onlineSend
}

//用户的下线业务
func (this *User) Offline() {

	//用户下线,将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//广播当前用户上线消息
	// this.server.BroadCast(this, "下线")
	onlineSend := "-* " + this.Name + " 已下线 *-"
	this.server.Message <- onlineSend
}



//用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户都有哪些
		// this.SendMsg("当前在线用户：")
		this.SendMsg("---------- 当前在线用户 -----------")
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := user.Name + ":" + "在线..."
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
		this.SendMsg("---------------------------------")
	} else {
		this.server.BroadCast(this, msg)
	}
}

//监听当前User channel的方法, 一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.SendMsg(msg)
	}
}
