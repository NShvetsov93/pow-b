package middlewares

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type request struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
}

type auth interface {
	Check(ctx context.Context, ip string, token string) bool
}

func WithAuth(authService auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			req := &request{}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(b, &req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if ip != req.Ip {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if !authService.Check(r.Context(), req.Ip, req.Token) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			r.Body = ioutil.NopCloser(strings.NewReader(string(b)))

			next.ServeHTTP(w, r)
		})
	}
}
