package firstKV

import (
	"encoding/json"
	"net"

	"github.com/cnlesscode/gotool"
)

// 接收消息结构体
type ReceiveMessage struct {
	Action string
	Key    string
	Data   FirstMQAddr
}

// 响应消息结构体
type ResponseMessage struct {
	ErrCode int
	Data    string
}

type FirstMQAddr struct {
	Host     string
	Port     string
	Addr     string
	JoinTime string
}

type FirstMQAddrs map[string]FirstMQAddr

// 处理消息
func HandleMessage(conn net.Conn, msgByte []byte) {
	msg := ReceiveMessage{}
	err := json.Unmarshal(msgByte, &msg)
	if err != nil {
		conn.Write(ResponseResult(200300, "消息格式错误"))
		return
	}
	// 根据消息类型，调用不同的处理函数
	// map add
	if msg.Action == "add mqServer" {
		receiveData := msg.Data
		data, ok := Get(msg.Key)
		mapKey := receiveData.Addr
		if !ok {
			Set(msg.Key, FirstMQAddrs{mapKey: receiveData})
		} else {
			dataOld, ok := data.(FirstMQAddrs)
			if ok {
				dataOld[mapKey] = receiveData
				Set(msg.Key, dataOld)
			}
		}
		conn.Write(ResponseResult(0, string(ResponseResult(0, "ok"))))
	} else if msg.Action == "get mqServers" {
		data, ok := Get(msg.Key)
		if ok {
			dataString, err := json.Marshal(data)
			if err != nil {
				conn.Write(ResponseResult(200100, err.Error()))
				return
			}
			conn.Write(ResponseResult(0, string(dataString)))
		} else {
			conn.Write(ResponseResult(200200, msg.Key+" 不存在"))
		}
	} else if msg.Action == "remove mqServer" {
		data, ok := Get(msg.Key)
		if !ok {
			conn.Write(ResponseResult(0, string(ResponseResult(0, "ok"))))
			return
		} else {
			mapData, ok := data.(FirstMQAddrs)
			if !ok {
				conn.Write(ResponseResult(200300, "服务端数据格式错误"))
				return
			}
			delete(mapData, msg.Data.Addr)
			Set(msg.Key, mapData)
			conn.Write(ResponseResult(0, "ok"))
		}
	} else {
		conn.Write(ResponseResult(400400, ""))
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

func Send(conn net.Conn, msg ReceiveMessage) (ResponseMessage, error) {
	defer conn.Close()
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
