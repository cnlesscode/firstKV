package firstKV

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/cnlesscode/gotool/gfs"
)

var firstKVMap *sync.Map = &sync.Map{}

func Init() {
	if !gfs.DirExists(FirstKVdataLogsDir) {
		err := os.Mkdir(FirstKVdataLogsDir, 0777)
		if err != nil {
			panic("FirstKV Error : 数据目录创建失败: " + err.Error() + "\n")
		}
	}
	// 加载数据到 syncMap
	res := gfs.ScanDirStruct{
		Path: FirstKVdataLogsDir,
	}
	err := gfs.ScanDir(false, &res)
	if err != nil {
		panic("FirstKV Error : 数据目录扫描失败: " + err.Error() + "\n")
	}
	for _, v := range res.Sons {
		if v.IsDir {
			continue
		}
		// 读取文件内容
		content, err := os.ReadFile(v.Path)
		if err != nil {
			continue
		}
		// 解析数据
		mapData := FirstMQAddrs{}
		err = json.Unmarshal(content, &mapData)
		if err != nil {
			continue
		}
		firstKVMap.Store(v.Name[0:len(v.Name)-5], mapData)
	}
}

func Set(k string, v any) {
	firstKVMap.Store(k, v)
	SaveDataToLog(k)
}

func Get(k string) (any, bool) {
	return firstKVMap.Load(k)
}

func Delete(k string) {
	firstKVMap.Delete(k)
	os.Remove(filepath.Join(FirstKVdataLogsDir, k+".json"))
}

func SaveDataToLog(k string) error {
	mapdata, ok := firstKVMap.Load(k)
	if !ok {
		return errors.New("FirstKV Error : 数据不存在")
	}
	str, err := json.Marshal(mapdata)
	if err != nil {
		return errors.New("FirstKV Error : JSON 格式转换失败")
	}
	filePath := filepath.Join(FirstKVdataLogsDir, k+".json")
	err = os.WriteFile(filePath, str, 0777)
	if err != nil {
		return errors.New("FirstKV Error : 数据保存失败")
	}
	return nil
}
