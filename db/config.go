package db

import (
	"io"
	"strconv"
	"strings"

	"github.com/cnk3x/webfx/config"
	"github.com/valyala/fasttemplate"
)

// var (
// 	ErrDatabaseName     = errors.New("database name is required")
// 	ErrDatabaseUser     = errors.New("database user is required")
// 	ErrDatabasePassword = errors.New("database password is required")
// )

const (
	MYSQL_DSN  = "{user}:{password}@tcp({server}:{port})/{name}?{args}"                           //?charset={charset}&parseTime=True&loc=Local&timeout=10s
	PQ_DSN     = "user={user} password={password} host={server} port={port} dbname={name} {args}" // sslmode=disable TimeZone=Asia/Shanghai
	SQLITE_DSN = "{path}?{args}"                                                                  //?_pragma=busy_timeout(3000)&_pragma=journal_mode(WAL)
)

type Config struct {
	Type     string `json:"type,omitempty"`     // 数据库类型
	User     string `json:"user,omitempty"`     // 用户
	Password string `json:"password,omitempty"` // 密码
	Host     string `json:"host,omitempty"`     // 服务器
	Port     int    `json:"port,omitempty"`     // 端口
	Name     string `json:"name,omitempty"`     // 数据库名
	Path     string `json:"path,omitempty"`     // 文件路径
	Args     string `json:"args,omitempty"`     // 扩展参数
	Debug    bool   `json:"debug,omitempty"`    // 调试模式
}

func DefineConfig(c Config) Config {
	if c.Type == "" {
		c.Type = "sqlite"
	}

	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		switch c.Type {
		case "mysql":
			c.Port = 3306
		case "postgres", "pg", "postgresql":
			c.Port = 5432
		}
	}

	if c.User == "" {
		c.User = "root"
	}

	if c.Password == "" {
		c.Password = "root"
	}

	if c.Name == "" {
		c.Name = "app"
	}

	if c.Path == "" {
		c.Path = config.WorkSpace("app.db")
	}

	if c.Args == "" {
		switch c.Type {
		case "mysql":
			c.Args = "charset=utf8mb4&parseTime=True&loc=Asia/Shanghai&timeout=10s"
		case "postgres", "pg", "postgresql":
			c.Args = "sslmode=disable TimeZone=Asia/Shanghai"
		case "sqlite", "sqlite3":
			c.Args = "_pragma=busy_timeout(3000)&_pragma=journal_mode(WAL)"
		}
	}

	return c
}

func (c Config) DSN() string {
	var tpl string
	switch c.Type {
	case "mysql":
		tpl = MYSQL_DSN
	case "postgres", "pg", "postgresql":
		tpl = PQ_DSN
	case "sqlite", "sqlite3":
		tpl = SQLITE_DSN
	}

	s := fasttemplate.New(tpl, "{", "}").ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "user":
			return w.Write([]byte(c.User))
		case "password":
			return w.Write([]byte(c.Password))
		case "server":
			return w.Write([]byte(c.Host))
		case "port":
			return w.Write([]byte(strconv.Itoa(c.Port)))
		case "name":
			return w.Write([]byte(c.Name))
		case "path":
			return w.Write([]byte(c.Path))
		case "args":
			return w.Write([]byte(c.Args))
		default:
			return 0, nil
		}
	})

	s = strings.TrimSuffix(strings.TrimSpace(s), "?")
	return s
}
