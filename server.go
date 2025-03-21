package firstKV

import (
	"log"
	"net"

	"github.com/cnlesscode/gotool"
)

var FirstKVdataLogsDir string = ""

// TCPServer TCP服务器结构
type TCPServer struct {
	listener net.Listener
}

// 创建TCP服务器
func NewTCPServer(addr string) *TCPServer {
	// 创建 Socket 端口监听
	// listener 是一个用于面向流的网络协议的公用网络监听器接口，
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("✔ FirstKV : 服务启动成功, 端口" + addr)
	// 返回实例
	return &TCPServer{listener: listener}
}

// Accept 等待客户端连接
func (t *TCPServer) Accept() {
	// 关闭接口解除阻塞的 Accept 操作并返回错误
	defer t.listener.Close()
	// 循环等待客户端连接
	for {
		// 等待客户端连接
		conn, err := t.listener.Accept()
		if err == nil {
			// 处理客户端连接
			go t.Handle(conn)
		}
	}
}

// Handle 处理客户端连接
func (t *TCPServer) Handle(conn net.Conn) {
	for {
		// 创建字节切片
		buf, err := gotool.ReadTCPResponse(conn)
		if err != nil {
			// 退出协程
			conn.Close()
			break
		}
		// 处理消息
		HandleMessage(conn, buf)
	}
}

// 开启 TCP 服务
func StartServer(port string, dataLogsDir string) {
	FirstKVdataLogsDir = dataLogsDir + gotool.SystemSeparator
	Init()
	tcpServer := NewTCPServer(":" + port)
	tcpServer.Accept()
}
