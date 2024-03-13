package web

import (
	"context"
	"crypto/ed25519"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cnk3x/webfx/config"
	"github.com/cnk3x/webfx/utils/eddsa"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenMalformed = jwt.ErrTokenMalformed

func NewJwt() (s *Jwt, err error) {
	s = &Jwt{method: jwt.SigningMethodEdDSA}

	priFile := config.WorkSpace("auth.pem")
	if s.priKey, err = eddsa.ReadPrivateFile(priFile); err != nil {
		if os.IsNotExist(err) {
			if s.pubKey, s.priKey, err = eddsa.GenerateKey(); err != nil {
				return
			}
			if err = eddsa.WritePrivateFile(priFile, s.priKey); err != nil {
				return
			}
		}
		return
	}

	s.pubKey = s.priKey.Public().(ed25519.PublicKey)
	return
}

type Jwt struct {
	method jwt.SigningMethod
	pubKey ed25519.PublicKey
	priKey ed25519.PrivateKey
}

func (s *Jwt) PublicKey(*jwt.Token) (any, error) {
	return s.pubKey, nil
}

func (s *Jwt) Validate(bearer string) (claims M, err error) {
	var token *jwt.Token
	if token, err = jwt.Parse(bearer, s.PublicKey); err != nil {
		return
	}

	if m, ok := token.Claims.(jwt.MapClaims); ok {
		claims = M(m)
		return
	}

	err = ErrTokenMalformed
	return
}

func (s *Jwt) Create(claims M) (bearer string, err error) {
	if _, find := claims["exp"]; !find {
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}
	bearer, _ = jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims(claims)).SignedString(s.priKey)
	return
}

func (s *Jwt) MustCreate(claims M) (bearer string) {
	var err error
	if bearer, err = s.Create(claims); err != nil {
		panic(err)
	}
	return
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := GetContext(r)
		bearer := c.HeaderGet("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			bearer = bearer[7:]
		}

		if bearer == "" {
			bearer = c.Cookie("token")
		}

		if bearer == "" {
			c.Msg("unauthorized", 401)
			return
		}

		claims, err := c.Jwt().Validate(bearer)
		if err != nil {
			c.Msg(err.Error(), 401)
			return
		}

		r = r.WithContext(context.WithValue(c.r.Context(), TokenKey, claims))
		next.ServeHTTP(w, r)
	})
}

func Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := GetContext(r)
		bearer := c.HeaderGet("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			bearer = bearer[7:]
		}

		if bearer == "" {
			bearer = c.Cookie("token")
		}

		if bearer == "" {
			c.Msg("unauthorized", 401)
			return
		}

		claims, err := c.Jwt().Validate(bearer)
		if err != nil {
			c.Msg(err.Error(), 401)
			return
		}

		r = r.WithContext(context.WithValue(c.r.Context(), TokenKey, claims))
		next.ServeHTTP(w, r)
	})
}
