package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/cnk3x/webfx/utils/bufpool"
	"github.com/cnk3x/webfx/utils/fss"
	"github.com/cnk3x/webfx/utils/log"

	"github.com/samber/lo"
)

/* 输入方法 */

func (c Context) IsTLS() bool                  { return c.TLS != nil }                                    // IsTLS returns true if HTTP connection is TLS otherwise false.
func (c Context) Origin() string               { return c.Scheme() + "://" + c.r.Host }                   // Origin returns the scheme and host from the request URL.
func (c Context) Path() string                 { return c.r.URL.Path }                                    // Path returns the registered path for the handler.
func (c Context) PathDir() string              { return path.Dir(strings.TrimSuffix(c.r.URL.Path, "/")) } // PathDir returns the registered path for the handler.
func (c Context) Query(name string) string     { return c.r.URL.Query().Get(name) }                       // Query returns the query param for the provided name.
func (c Context) Querys(name string) []string  { return c.r.URL.Query()[name] }                           // Querys returns the query params for the provided name.
func (c Context) HeaderGet(name string) string { return c.r.Header.Get(name) }                            // get request header

func (c Context) ReadJSON(v any) error {
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	io.Copy(buf, c.r.Body)
	if buf.Len() > 0 {
		if err := json.NewDecoder(buf).Decode(v); err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return err
		}
	}
	return nil
}

func (c Context) ReadBody() (body []byte, err error) {
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	if _, err = io.Copy(buf, c.r.Body); err == nil {
		body = buf.Bytes()
	} else {
		c.Status(http.StatusUnprocessableEntity)
	}
	return
}

// Cookie gets the value of a cookie with name name.
func (c Context) Cookie(name string) string {
	cookie, _ := c.r.Cookie(name)
	if cookie != nil {
		return cookie.Value
	}
	return ""
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c Context) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.r.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.r.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.r.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.r.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

// Host 返回请求的不带端口号的主机名
//
//	从请求返回过来的，所以忽略格式错误的状态
func (c Context) Host() (host string) {
	if host = c.r.Host; host != "" {
		if i := strings.Index(host, "]"); i != -1 {
			host = host[:i+1]
		} else if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
	}
	return
}

// Port 返回请求的端口号
func (c Context) Port() (port int) {
	if host := c.r.Host; host != "" {
		if i := strings.Index(host, ":"); i != -1 {
			port, _ = strconv.Atoi(host[i+1:])
		}
	}
	port = lo.If(port > 0, port).ElseIf(c.IsTLS(), 443).Else(80)
	return
}

// RealIp returns the client's network address based on `X-Forwarded-For` or `X-Real-IP` request header.
func (c Context) RealIp() string {
	if ip := c.r.Header.Get(HeaderXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			xffIp := strings.TrimSpace(ip[:i])
			xffIp = strings.TrimPrefix(xffIp, "[")
			xffIp = strings.TrimSuffix(xffIp, "]")
			return xffIp
		}
		return ip
	}
	if ip := c.r.Header.Get(HeaderXRealIP); ip != "" {
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.RemoteAddr)
	return ra
}

// 获取所有文件
func (c Context) FormFiles(fields ...string) (hs []*multipart.FileHeader) {
	const defaultMaxMemory = 32 << 20 // 32 MB

	if c.r.MultipartForm == nil {
		if err := c.r.ParseMultipartForm(defaultMaxMemory); err != nil {
			log.Warnf("[WEB] parse multipart form: %v", err)
			return
		}
	}

	if c.r.MultipartForm != nil && c.r.MultipartForm.File != nil {
		if len(fields) == 0 {
			for _, files := range c.r.MultipartForm.File {
				hs = append(hs, files...)
			} // append all files
		} else {
			for _, field := range fields {
				hs = append(hs, c.r.MultipartForm.File[field]...)
			}
		}
	}

	return
}

// 获取第一个文件
func (c Context) FormFile(fields ...string) (h *multipart.FileHeader) {
	if hs := c.FormFiles(); len(hs) > 0 {
		h = hs[0]
	}
	return
}

// 保存上传文件，可以指定上传文件的字段名，如果没有指定，则保存第一个文件
func (c Context) SaveUpload(dstPath string, fields ...string) (h *multipart.FileHeader, hash string, err error) {
	if h = c.FormFile(fields...); h == nil {
		c.Status(422)
		err = fmt.Errorf("no file uploaded")
		return
	}

	var src io.ReadCloser
	if src, err = h.Open(); err != nil {
		return
	}
	defer src.Close()
	hash, err = fss.SaveFile(src, dstPath, true)
	return
}

// // 获取当前请求的授权信息
// func (c Context) GetAuthorization() (r Authorization, ok bool) {
// 	if r.Raw = strings.TrimSpace(c.r.Header.Get(HeaderAuthorization)); r.Raw == "" {
// 		return
// 	}

// 	t, a, found := strings.Cut(r.Raw, " ")
// 	if found {
// 		t, a = strings.TrimSpace(t), strings.TrimSpace(a)
// 		if strings.EqualFold(t, "basic") {
// 			c, err := base64.StdEncoding.DecodeString(a)
// 			if err != nil {
// 				return
// 			}
// 			r.Key, r.Token, ok = strings.Cut(string(c), ":")
// 		} else {
// 			r.Type = t
// 			r.Key, r.Token, ok = strings.Cut(a, ":")
// 			if !ok {
// 				r.Token = a
// 				ok = true
// 			}
// 		}
// 	} else {
// 		r.Token = r.Raw
// 		ok = true
// 	}

// 	return
// }
// type Authorization struct {
// 	Type  string `json:"type,omitempty"` // bearer, sso-key, sso-token, basic
// 	Key   string `json:"key,omitempty"`
// 	Token string `json:"token,omitempty"`
// 	Raw   string `json:"-"`
// }
