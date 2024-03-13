package web

import (
	"io"
	"net/http"

	"gorm.io/gorm"
)

type ctxKey struct{ name string }

func (k *ctxKey) String() string {
	return "context value " + k.name
}

var (
	StatusCtxKey = &ctxKey{"Status"}
	DBCtxKey     = &ctxKey{"DB"}
	JwtCtxKey    = &ctxKey{"Jwt"}
	CtxKey       = &ctxKey{"Ctx"}
	WriterKey    = &ctxKey{"Writer"}
	UIDKey       = &ctxKey{"UID"}
	TokenKey     = &ctxKey{"Token"}
)

func (c Context) GetDB() *gorm.DB  { return cValue[*gorm.DB](c, DBCtxKey) }
func (c Context) Jwt() *Jwt        { return cValue[*Jwt](c, JwtCtxKey) }
func (c Context) GetUID() (id int) { return cValue[int](c, UIDKey) }
func (c Context) Response() w      { return cValue[w](c, WriterKey) }
func (c Context) GetStatus() int   { return cValue[int](c, StatusCtxKey) }

func (c Context) Request() *http.Request         { return c.r }
func (c Context) RequestHeader() http.Header     { return c.r.Header }
func (c Context) RequestWrite(w io.Writer) error { return c.r.Write(w) }

func (c Context) WriteHeader(s int)           { c.Response().WriteHeader(s) }
func (c Context) Header() http.Header         { return c.Response().Header() }
func (c Context) Write(b []byte) (int, error) { return c.Response().Write(b) }

func (c Context) ClaimsValue(key string) any {
	if m := cValue[M](c, TokenKey); m != nil {
		return m[key]
	}
	return nil
}

func (c Context) Claims() M { return cValue[M](c, TokenKey) }

func cValue[T any](c Context, key any) T { return conv[T](c.r.Context().Value(key)) }
func conv[T any](v any) (out T)          { out, _ = v.(T); return }
