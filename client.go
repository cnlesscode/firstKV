package firstKV

import (
	"encoding/json"
	"net"

	"github.com/cnlesscode/gotool"
)

// 此函数用于其他工具调用 FirstKV 时使用
func Send(conn net.Conn, msg ReceiveMessage, close bool) (ResponseMessage, error) {
	if close {
		defer conn.Close()
	}
	response := ResponseMessage{}
	msgByte, _ := json.Marshal(msg)

	// 写消息
	err := gotool.WriteTCPResponse(conn, msgByte)
	if err != nil {
		return response, err
	}

	// 读取消息
	buf, err := gotool.ReadTCPResponse(conn)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return response, err
	}

	// 返回消息
	return response, nil
}
