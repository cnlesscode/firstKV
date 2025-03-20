package firstKV

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

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
		mapData := make(map[string]Item, 0)
		err = json.Unmarshal(content, &mapData)
		if err != nil {
			continue
		}
		firstKVMap.Store(v.Name[0:len(v.Name)-5], mapData)
	}
	// 间隔5分钟检查有效期
	go func() {
		for {
			time.Sleep(5 * time.Second)
			CheckExpirationTime()
		}
	}()
}

// 有效期检查
func CheckExpirationTime() {
	println("FirstKV 检查有效期")
	firstKVMap.Range(func(key, value any) bool {
		changed := false
		if value, ok := value.(map[string]Item); ok {
			for k, v := range value {
				if v.ExpirationTime < 0 {
					continue
				} else if v.ExpirationTime < time.Now().Unix() {
					RemoveItem(key.(string), k)
					changed = true
				}
			}
			if changed {
				SaveDataToLog(key.(string))
			}
		}
		return true
	})
}

func Set(k string, v any) {
	firstKVMap.Store(k, v)
	SaveDataToLog(k)
}

func SetItem(mainKey, itemKey string, item Item, expirationTime int64) {
	// 获取主库
	mainDB, ok := Get(mainKey)
	item.CreateTime = time.Now().Unix()
	if expirationTime < 0 {
		item.ExpirationTime = -1
	} else {
		item.ExpirationTime = item.CreateTime + expirationTime
	}
	// 主库为空
	if !ok {
		Set(mainKey, map[string]Item{itemKey: item})
	} else {
		// 已存在数据
		dataOld, ok := mainDB.(map[string]Item)
		if ok {
			dataOld[itemKey] = item
			Set(mainKey, dataOld)
		}
	}
}

func Get(k string) (any, bool) {
	return firstKVMap.Load(k)
}

func GetItem(mainKey, itemKey string) (Item, bool) {
	mainDB, ok := Get(mainKey)
	// 主库为空
	if !ok {
		return Item{}, false
	}
	data, ok := mainDB.(map[string]Item)
	if !ok {
		return Item{}, false
	}
	item, ok := data[itemKey]
	return item, ok
}

func Remove(mainKey string) {
	// 获取主库
	_, ok := Get(mainKey)
	if !ok {
		return
	}
	firstKVMap.Delete(mainKey)
	os.Remove(filepath.Join(FirstKVdataLogsDir, mainKey+".json"))
}

func RemoveItem(mainKey, itemKey string) {
	// 获取主库
	mainDB, ok := Get(mainKey)
	// 主库为空
	if !ok {
		return
	}
	data, ok := mainDB.(map[string]Item)
	if !ok {
		return
	}
	delete(data, itemKey)
	Set(mainKey, data)
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
