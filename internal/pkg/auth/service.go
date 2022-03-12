package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
)

type Service struct {
	stg storage
}

type storage interface {
	Add(ctx context.Context, ip string, token string) error
	Exist(ctx context.Context, ip string) bool
	Check(ctx context.Context, ip string, token string) bool
}

func New(stg storage) *Service {
	return &Service{
		stg: stg,
	}
}

func (s *Service) Auth(ctx context.Context, ip string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "pkg/auth/Auth")
	defer span.Finish()

	data := bytes.Join([][]byte{[]byte(ip), []byte(strconv.FormatInt(time.Now().Unix(), 10))}, []byte{})
	tokenBytes := sha256.Sum256(data)
	token := hex.EncodeToString(tokenBytes[:])

	err := s.stg.Add(ctx, ip, token)
	if err != nil {
		return "", fmt.Errorf("couldn't auth: %w", err)
	}

	return token, nil
}

func (s *Service) Check(ctx context.Context, ip string, token string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "pkg/auth/Check")
	defer span.Finish()

	return s.stg.Check(ctx, ip, token)
}
