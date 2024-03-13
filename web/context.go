package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-Ip"

	HeaderXForwardedProto    = "X-Forwarded-Proto"
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"
	HeaderXForwardedSsl      = "X-Forwarded-Ssl"
	HeaderXUrlScheme         = "X-Url-Scheme"

	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"

	ContentTypeJSON  = "application/json; charset=utf-8"
	ContentTypeHTML  = "text/html; charset=utf-8"
	ContentTypePlain = "text/plain; charset=utf-8"
)

type M map[string]any

type (
	HandlerContext = func(Context)
	ContextMw      = func(next HandlerContext) HandlerContext
	Middleware     = func(next http.Handler) http.Handler
)

// func wrapMw(cmw func(next HandlerContext) HandlerContext) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		cNext := func(c Context) { next.ServeHTTP(c.Response(), c.Request()) }
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { cmw(cNext)(GetContext(r)) })
// 	}
// }

func GetContext(r *http.Request) Context {
	return Context{r: r}
}

type Context struct{ *r }

type (
	r = http.Request
	w = http.ResponseWriter
)

// 获取远程链接的JSON数据
func (c Context) GetJSON(url string) (any, error) {
	if url == "" {
		return "", fmt.Errorf("jsonApi: no URL provided")
	}

	isAbsUrl := func(url string) bool {
		return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	}

	if !isAbsUrl(url) {
		if url[0] != '/' {
			url = c.PathDir() + "/" + url
		}
		url = c.Origin() + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out any
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}
