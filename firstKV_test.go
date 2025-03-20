package firstKV

import (
	"fmt"
	"testing"
	"time"
)

// 测试命令 :
// go test -v -run=TestRun
func TestRun(t *testing.T) {
	FirstKVdataLogsDir = "./data"
	Init()
	SetItem("key1", "skey1", Item{Data: "value2"}, 500)
	SetItem("key1", "skey2", Item{Data: "value2"}, 5)
	SetItem("key2", "skey1", Item{Data: "value3"}, -1)
	res, ok := GetItem("key1", "skey1")
	if ok {
		fmt.Printf("res: %v\n", res)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
