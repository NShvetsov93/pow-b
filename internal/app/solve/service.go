package solve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
)

type Request struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
	Hash  string `json:"hash"`
	Nonce int    `json:"nonce"`
}

type Response struct {
	Phrase string `json:"phrase"`
}

type Implemetation struct {
	service solve
}

type solve interface {
	Solve(ctx context.Context, token string, ip string, hash string, nonce int) (string, error)
}

func New(service solve) *Implemetation {
	return &Implemetation{
		service: service,
	}
}

func (i *Implemetation) Solve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "app/Solve")
	defer span.Finish()

	req := &Request{}

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &req)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := i.service.Solve(ctx, req.Token, req.Ip, req.Hash, req.Nonce)
	if err != nil {
		err = fmt.Errorf("solve error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &Response{
		Phrase: result,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := jsoniter.NewEncoder(w).Encode(res); err != nil {
		err = fmt.Errorf("error encoding response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
