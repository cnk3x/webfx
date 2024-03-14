package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/cnk3x/webfx/utils/fss"
	"github.com/cnk3x/webfx/utils/strs"
)

var workSpace string // 工作目录

func init() {
	workSpace, _ = os.Getwd()
}

func Display() {
	fmt.Printf("Version:   %s\n", Version)
	fmt.Printf("BuildHash: %s\n", BuildHash)
	fmt.Printf("BuildTime: %s\n", BuildTime)
	fmt.Printf("WorkSpace: %s\n", WorkSpace())
	fmt.Println("--------------------------")
}

// WorkSpace 解析工作路径
func WorkSpace(names ...string) string {
	base := fss.MakeDirs(workSpace)
	if len(names) > 1 {
		return filepath.Join(slices.Insert(names, 0, base)...)
	} else if len(names) > 0 {
		return filepath.Join(base, names[0])
	} else {
		return base
	}
}

func SetWorkSpace(ws string) {
	if ws, _ = filepath.Abs(ws); ws != "" {
		workSpace = ws
	}
}

// Get 获取配置
func Get(name string, force ...bool) (out Value) {
	if out = IifF(Select(force...), loadDirect, loadOnce); name != "" {
		out = out.Get(name)
	}

	if !out.Exists() {
		out = Value{Type: String, Str: os.Getenv(strs.Snake(filepath.Base(workSpace)+"_"+name, true))}
	}

	return
}

var loadOnce = sync.OnceValue(loadDirect)

func loadDirect() Value {
	var dataAny any
	_ = ReadJSONFile(WorkSpace("config.json"), &dataAny)
	_ = ReadYAMLFile(WorkSpace("config.yaml"), &dataAny)
	data, _ := json.Marshal(dataAny)
	return ParseJSON(data)
}
