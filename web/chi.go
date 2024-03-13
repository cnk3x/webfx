package web

import (
	"strings"

	"github.com/cnk3x/webfx/utils/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	Router = chi.Router
	Routes = chi.Routes
	Mux    = chi.Mux
)

var (
	Recoverer = middleware.Recoverer
	Logger    = middleware.Logger
	NewMux    = chi.NewMux
)

func printRoute(r Routes) {
	var printRoute func(r Routes, prefix string)

	printRoute = func(r Routes, prefix string) {
		prefix = strings.TrimSuffix(prefix, "/*")
		for _, route := range r.Routes() {
			var methods []string
			for method := range route.Handlers {
				if method == "*" {
					methods = []string{method}
					break
				}
				methods = append(methods, method)
			}

			if strings.HasSuffix(route.Pattern, "/*") {
				if route.SubRoutes != nil {
					printRoute(route.SubRoutes, route.Pattern)
				} else {
					for _, method := range methods {
						log.Output(0, log.DEBUG, "ROUTE %-7s %s%s", method, prefix, route.Pattern)
					}
				}
			} else {
				for _, method := range methods {
					log.Output(0, log.DEBUG, "ROUTE %-7s %s%s", method, prefix, route.Pattern)
				}
			}
		}
	}

	printRoute(r, "")
}
