package requestchallenge

import (
	"context"
	"fmt"
	"net"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
)

type Response struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
}

type Implemetation struct {
	service auth
}

type auth interface {
	Auth(ctx context.Context, ip string) (string, error)
}

func New(service auth) *Implemetation {
	return &Implemetation{
		service: service,
	}
}

func (i *Implemetation) Gen(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "app/Gen")
	defer span.Finish()

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	token, err := i.service.Auth(ctx, ip)
	if err != nil {
		err = fmt.Errorf("couldn't generate token: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &Response{
		Token: token,
		Ip:    ip,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := jsoniter.NewEncoder(w).Encode(res); err != nil {
		err = fmt.Errorf("error encoding response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
