package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/cnk3x/webfx/utils/bufpool"
	"github.com/cnk3x/webfx/utils/log"
	"github.com/cnk3x/webfx/utils/strs"
)

/* 输出方法 */
func (c Context) NoContent()                   { c.WriteHeader(http.StatusNoContent) }
func (c Context) NotFound()                    { http.Error(c.Response(), "404 page not found", http.StatusNotFound) }
func (c Context) HeaderSet(name, value string) { c.Header().Set(name, value) }

func (c Context) Redirect(url string) {
	http.Redirect(c.Response(), c.Request(), url, http.StatusFound)
}

func (c Context) Status(status int) {
	*c.r = *(c.r.WithContext(context.WithValue(c.r.Context(), StatusCtxKey, status)))
}

func (c Context) Blob(v []byte, contentType string) {
	if len(contentType) > 0 {
		c.SetContentType(contentType)
	}
	if status := c.GetStatus(); status != 0 {
		c.WriteHeader(status)
	}
	c.Write(v)
}

func (c Context) SetContentType(contentType string) {
	if contentType != "" {
		c.HeaderSet("Content-Type", contentType)
	}
}

func (c Context) Msg(errmsg string, errno ...int) {
	var s int
	if len(errno) > 0 {
		c.Status(errno[0])
		s = errno[0]
	} else {
		if s = c.GetStatus(); s < 400 {
			s = 0
		}
	}
	c.JSON(M{"errno": s, "errmsg": errmsg})
}

func (c Context) Error(err error) {
	if err != nil {
		status := c.GetStatus()
		errmsg := err.Error()

		if status < 400 {
			switch {
			case errors.Is(err, fs.ErrNotExist):
				errmsg, status = "404 page not found", http.StatusNotFound
			case errors.Is(err, fs.ErrPermission):
				errmsg, status = "403 Forbidden", http.StatusForbidden
			default:
				status = http.StatusInternalServerError
			}
			c.Status(status)
		}
		log.Warnf("[WEB] %v(%d)", errmsg, status)
		c.JSON(M{"errno": status, "errmsg": errmsg})
	}
}

func (c Context) HTML(data string) {
	c.Blob([]byte(data), ContentTypeHTML)
}

func (c Context) PlainText(data string) {
	c.Blob([]byte(data), ContentTypePlain)
}

func (c Context) JSON(v any, pretty ...bool) {
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	e := json.NewEncoder(buf)
	e.SetEscapeHTML(true)
	if len(pretty) > 0 && pretty[0] {
		e.SetIndent("", "  ")
	}
	if err := e.Encode(v); err != nil {
		log.Warnf("[WEB] encode json: %v", err)
		return
	}

	c.JSONBlob(buf.Bytes())
}

func (c Context) JSONBlob(jsonBytes []byte) {
	c.Blob(jsonBytes, ContentTypeJSON)
}

func (c Context) FlatMap(keyValue ...any) {
	m := M{}
	var k string
	for i := 1; i < len(keyValue); i += 2 {
		if k = strs.From(keyValue[i-1]); k == "" {
			continue
		}
		m[k] = keyValue[i]
	}
	c.JSON(m)
}

func (c Context) File(filename string, attachName string) {
	c.Header().Set("Content-Disposition", "attachment; filename="+attachName)
	http.ServeFile(c, c.r, filename)
}

func (c Context) FileBlob(data []byte, contentType string, attachName string) {
	c.Header().Set("Content-Disposition", "attachment; filename="+attachName)
	c.Blob(data, contentType)
}

func (c Context) Template(name string, model any) {
	tpl := GetTemplateContext(c.r, http.FS(templateFs))

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if err := tpl.Execute(name, buf, model); err != nil {
		c.Error(fmt.Errorf("templates render: %v, name=%s, model=%T", err, name, model))
		return
	}

	c.Blob(buf.Bytes(), "text/html")
}
