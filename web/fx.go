package web

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

const RouteTag = `group:"routes"`

func Provide(ctors ...any) fx.Option {
	var options []fx.Option
	for _, ctor := range ctors {
		it := fx.Provide(fx.Annotate(ctor, fx.As(new(Route)), fx.ResultTags(RouteTag)))
		if len(ctors) == 1 {
			return it
		}
		options = append(options, it)
	}
	return fx.Options(options...)
}

func Supply(values ...Route) fx.Option {
	var options []fx.Option
	for _, value := range values {
		it := fx.Supply(fx.Annotate(value, fx.As(new(Route)), fx.ResultTags(RouteTag)))
		if len(values) == 1 {
			return it
		}
		options = append(options, it)
	}
	return fx.Options(options...)
}

type Params struct {
	fx.In
	DB     *gorm.DB
	Jwt    *Jwt
	Routes []Route `group:"routes"`
}

func Run(params Params) {
	r := NewMux()
	r.Use(Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := r.Context()
			c = context.WithValue(c, DBCtxKey, params.DB)
			c = context.WithValue(c, JwtCtxKey, params.Jwt)
			c = context.WithValue(c, WriterKey, w)
			next.ServeHTTP(w, r.WithContext(c))
		})
	})

	for _, s := range params.Routes {
		Struct(r, s)
	}

	Serve(r)
}
