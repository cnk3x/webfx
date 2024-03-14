package web

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/cnk3x/webfx/utils/log"
	"github.com/cnk3x/webfx/utils/strs"

	"github.com/samber/lo"
)

type Route interface{}

func Struct(r Router, s any) {
	if simple, ok := s.(Simple); ok {
		if h := ResolveHandler(simple.Handler); h != nil {
			if simple.Protected {
				r = r.With(Authenticator)
			}
			if simple.Mount {
				log.Output(0, log.DEBUG, "SIMPLE MOUNT   %-35s", simple.Pattern)
				r.Mount(simple.Pattern, h)
			} else {
				method, pattern := ResolvePattern(simple.Pattern)
				log.Output(0, log.DEBUG, "SIMPLE %-7s %-35s", method, pattern)
				r.Handle(method+" "+pattern, h)
			}
		} else {
			log.Fatalf("simple.Handler is nil: %T", simple.Handler)
		}
		return
	}

	var mount string
	if rx, ok := s.(interface{ Mount() string }); ok {
		mount = strings.TrimSuffix(rx.Mount(), "/")
	}

	var protected bool
	if rx, ok := s.(interface{ Protected() bool }); ok {
		protected = rx.Protected()
	}

	r.Group(func(r Router) {
		if protected {
			r.Use(Authenticator)
		}

		rv := reflect.ValueOf(s)
		rt := reflect.TypeOf(s)

		nameFix, _ := s.(interface{ Pattern(mName string) string })
		for i := 0; i < rt.NumMethod(); i++ {
			m := rt.Method(i)
			if !m.IsExported() {
				continue
			}

			handler := ResolveHandler(rv.Method(i).Interface())
			if handler == nil {
				continue
			}

			var (
				pattern string
				method  string
			)

			if nameFix != nil {
				if p := nameFix.Pattern(m.Name); p != "" {
					if ps := strings.Fields(p); len(ps) == 2 {
						method, pattern = ps[0], ps[1]
					} else {
						pattern = p
						method = "POST"
					}
				}
			}

			if pattern == "" {
				method, pattern = ResolvePattern(m.Name)
			}

			log.Output(0, log.DEBUG, "HANDLE %-7s %-35s => %s.%s", method, mount+pattern, rt.String(), m.Name)

			r.Handle(method+" "+mount+pattern, handler)
		}
	})
}

func ResolvePattern(in string) (method, pattern string) {
	method = "POST"
	pattern = strings.TrimSpace(in)

	func() {
		if strings.HasPrefix(pattern, "/") {
			return
		}

		if ss := strings.Fields(pattern); len(ss) == 2 {
			method = strings.ToUpper(ss[0])
			pattern = ss[1]
			return
		}

		pattern = strs.Snake(pattern)

		if ss := strings.SplitN(pattern, "_", 2); len(ss) == 2 {
			m := strings.ToUpper(ss[0])
			p := ss[1]

			if m == "GET" || m == "POST" || m == "PUT" || m == "DELETE" {
				method = m
				pattern = p
			}
		}

		if pattern == "index" {
			method = "GET"
			pattern = "/"
		} else {
			pattern = "/" + pattern
		}
	}()

	pattern = regexp.MustCompile(`\s+|/index$`).ReplaceAllString(pattern, "")
	pattern = lo.Ternary(strings.HasPrefix(pattern, "/"), pattern, "/"+pattern)

	return
}

func ResolveHandler(in any) (handler http.Handler) {
	handler, _ = in.(http.Handler)

	if handler == nil {
		if h, ok := in.(func(http.ResponseWriter, *http.Request)); ok {
			handler = http.HandlerFunc(h)
		}
	}

	if handler == nil {
		if h, ok := in.(func(Context)); ok {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h(GetContext(r))
			})
		}
	}

	if handler == nil {
		if h, ok := in.(func(Context) error); ok {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c := GetContext(r)
				if err := h(c); err != nil {
					c.Error(err)
				}
			})
		}
	}

	if handler == nil {
		if h, ok := in.(func(context.Context)); ok {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h(r.Context())
			})
		}
	}

	if handler == nil {
		if h, ok := in.(func(context.Context) error); ok {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := h(r.Context()); err != nil {
					GetContext(r).Error(err)
				}
			})
		}
	}

	if handler == nil {
		if h, ok := in.(func(context.Context) (any, error)); ok {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := h(r.Context())
				if err != nil {
					GetContext(r).Error(err)
				} else {
					GetContext(r).JSON(data)
				}
			})
		}
	}

	return
}

type Simple struct {
	Pattern   string
	Handler   any
	Mount     bool
	Protected bool
}
