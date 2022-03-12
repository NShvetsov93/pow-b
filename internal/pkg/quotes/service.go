package quotes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	url    string
	client *http.Client
}

type Response struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func New(url string, timeout time.Duration) *Service {
	return &Service{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *Service) Get(ctx context.Context) (*Response, error) {
	res := &Response{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return res, fmt.Errorf("couldn't create request fot external service: %w", err)
	}
	r, err := s.client.Do(req)
	if err != nil {
		return res, fmt.Errorf("couldn't get response from external service: %w", err)
	}

	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil && err != io.EOF {
		return res, fmt.Errorf("couldn't decode response from external service: %w", err)
	}

	return res, nil
}
