package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/cnk3x/webfx/config"
	"github.com/cnk3x/webfx/utils/fss"
	"github.com/cnk3x/webfx/utils/log"
	"github.com/cnk3x/webfx/utils/strs"

	"github.com/caddyserver/certmagic"
)

func Serve(handler http.Handler) {
	httpsEnabled := config.Get("web.https").Bool()
	if httpsEnabled {
		ServeHTTPS(handler)
	} else {
		ServeHTTP(handler)
	}
}

func ServeHTTP(handler http.Handler) {
	var (
		port = config.Select(config.Get("web.port").String(), "8080")
		host = config.Iif(config.Get("web.internal").Bool(), "127.0.0.1", "")
	)

	s := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("%s:%s", host, port),
		BaseContext: func(listener net.Listener) context.Context {
			log.Infof("web server started, listen: %s", listener.Addr())
			return context.Background()
		},
	}

	go func() {
		log.Infof("web server will start listen at %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Warnf("web server error: %s", err)
			}
		}
		log.Infof("web server stopped")
	}()
}

func ServeHTTPS(mux http.Handler) {
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "me@wen.cx"
	certmagic.DefaultACME.CA = certmagic.ZeroSSLProductionCA // lo.Ternary(config.Debug, certmagic.LetsEncryptStagingCA, certmagic.ZeroSSLProductionCA)
	certmagic.Default.KeySource = certmagic.StandardKeyGenerator{KeyType: certmagic.RSA4096}
	certmagic.Default.Storage = &certmagic.FileStorage{Path: fss.MakeDirs(config.WorkSpace("certs"))}

	domains := strs.Clean(strings.Split(config.Get("web.domain").String(), ","), strings.TrimSpace)

	log.Infof("serve web: %s", domains)
	go func() {
		if err := certmagic.HTTPS(domains, mux); err != nil {
			if err != http.ErrServerClosed {
				log.Warnf("web server error: %s", err)
			}
		}
		log.Infof("web server stopped")
	}()
}
