package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/cnk3x/webfx/utils/fss"
	"github.com/cnk3x/webfx/utils/strs"
)

var (
	debug     bool   // 调试模式
	workSpace string // 工作目录
	name      string // 程序名称
)

func init() {
	workSpace, _ = os.Getwd()
	name = getExecName()
}

func IsDebug() bool   { return debug }
func SetDebug(d bool) { debug = d }

func Display() {
	fmt.Printf("Name:        %s\n", name)
	fmt.Printf("WorkSpace:   %s\n", WorkSpace())
	fmt.Printf("Debug:       %v\n", IsDebug())
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

func SetWorkSpace(ws string) error {
	if workSpace != ws {
		workSpace = ws
		if err := os.Chdir(workSpace); err != nil {
			return err
		}
		workSpace, _ = os.Getwd()
	}
	return nil
}

// Get 获取配置
func Get(name string, force ...bool) (out Value) {
	if out = IifF(Coalesce(force...), loadDirect, loadOnce); name != "" {
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

func getExecName() string {
	name := filepath.Base(os.Args[0])

	var exts []string
	if runtime.GOOS == "windows" {
		if x := os.Getenv(`PATHEXT`); x != "" {
			for _, e := range strings.Split(strings.ToLower(x), `;`) {
				if e == "" {
					continue
				}
				if e[0] != '.' {
					e = "." + e
				}
				exts = append(exts, e)
			}
		} else {
			exts = []string{".com", ".exe", ".bat", ".cmd"}
		}
	}

	for _, ext := range exts {
		if strings.HasSuffix(name, ext) {
			name = name[:len(name)-len(ext)]
			break
		}
	}

	return strings.Trim(regexp.MustCompile(`[^\w]`).ReplaceAllString(name, "_"), "_")
}
