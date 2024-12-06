package firstKV

import (
	"fmt"
	"testing"
)

// 测试命令 :
// go test -v -run=TestRun
func TestRun(t *testing.T) {
	FirstKVdataLogsDir = "./data"
	Init()
	Set("key1", "value1")
	res, ok := Get("key1")
	if !ok {
		println("key1 not found")
	}
	fmt.Printf("res: %v\n", res)
}
