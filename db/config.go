package db

import (
	"io"

	"github.com/valyala/fasttemplate"
)

const (
	MYSQL_DSN  = "{user}:{password}@tcp({server}:{port})/{name}?charset={charset}&parseTime=True&loc=Local&timeout=10s"
	PQ_DSN     = "host={server} port={port} user={user} password={password} dbname={name} sslmode=disable TimeZone=Asia/Shanghai"
	SQLITE_DSN = "file:{path}?_pragma=busy_timeout(3000)&_pragma=journal_mode(WAL)"
)

type DatabaseOptions struct {
	Type     string `json:"type,omitempty"`     // 数据库类型
	User     string `json:"user,omitempty"`     // 用户
	Password string `json:"password,omitempty"` // 密码
	Server   string `json:"server,omitempty"`   // 服务器
	Port     string `json:"port,omitempty"`     // 端口
	Name     string `json:"name,omitempty"`     // 数据库名
	Charset  string `json:"charset,omitempty"`  // 字符集
	Path     string `json:"path,omitempty"`     // 文件路径
}

func (o DatabaseOptions) DSN() string {
	var tpl string
	switch o.Type {
	case "mysql":
		tpl = MYSQL_DSN
	case "postgres":
		tpl = PQ_DSN
	case "sqlite":
		tpl = SQLITE_DSN
	}

	return fasttemplate.New(tpl, "{", "}").ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "user":
			return w.Write([]byte(o.User))
		case "password":
			return w.Write([]byte(o.Password))
		case "server":
			return w.Write([]byte(o.Server))
		case "port":
			return w.Write([]byte(o.Port))
		case "name":
			return w.Write([]byte(o.Name))
		case "charset":
			return w.Write([]byte(o.Charset))
		case "path":
			return w.Write([]byte(o.Path))
		default:
			return 0, nil
		}
	})
}
