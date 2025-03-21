package firstKV

import (
	"encoding/json"
	"net"

	"github.com/cnlesscode/gotool"
)

// 接收消息结构体
type ReceiveMessage struct {
	Action  string
	MainKey string
	ItemKey string
	Data    Item
}

// 响应消息结构体
type ResponseMessage struct {
	ErrCode int
	Data    string
}

type Item struct {
	Data           any
	CreateTime     int64
	ExpirationTime int64
}

// 处理消息
func HandleMessage(conn net.Conn, msgByte []byte) {
	msg := ReceiveMessage{}
	err := json.Unmarshal(msgByte, &msg)
	if err != nil {
		gotool.WriteTCPResponse(
			conn, ResponseResult(200300, "消息格式错误"),
		)
		return
	}

	if msg.Action == "set" {
		receiveData := msg.Data
		SetItem(msg.MainKey, msg.ItemKey, receiveData, receiveData.ExpirationTime)
		gotool.WriteTCPResponse(
			conn, ResponseResult(0, string(ResponseResult(0, "ok"))),
		)
	} else if msg.Action == "get" {
		data, ok := Get(msg.MainKey)
		if ok {
			dataString, err := json.Marshal(data)
			if err != nil {
				gotool.WriteTCPResponse(
					conn,
					ResponseResult(200100, err.Error()),
				)
				return
			}
			gotool.WriteTCPResponse(
				conn, ResponseResult(0, string(dataString)),
			)
		} else {
			gotool.WriteTCPResponse(
				conn,
				ResponseResult(200200, msg.MainKey+" 不存在"),
			)
		}
	} else if msg.Action == "remove" {
		Remove(msg.MainKey)
		gotool.WriteTCPResponse(
			conn, ResponseResult(0, "ok"),
		)
	} else if msg.Action == "removeItem" {
		RemoveItem(msg.MainKey, msg.ItemKey)
		gotool.WriteTCPResponse(
			conn, ResponseResult(0, "ok"),
		)
	} else {
		gotool.WriteTCPResponse(
			conn, ResponseResult(400400, ""),
		)
	}
}

// 响应结果
func ResponseResult(errcode int, data string) []byte {
	responseMessage := ResponseMessage{
		ErrCode: errcode,
		Data:    data,
	}
	responseMessageByte, _ := json.Marshal(responseMessage)
	return responseMessageByte
}
