package solve

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/opentracing/opentracing-go"

	"pow-b/internal/pkg/quotes"
)

type Service struct {
	target *big.Int
	quotes quotesService
}

type quotesService interface {
	Get(ctx context.Context) (*quotes.Response, error)
}

func New(tBits int, q quotesService) *Service {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-tBits))
	return &Service{
		target: target,
		quotes: q,
	}
}

func (s *Service) Solve(ctx context.Context, token string, ip string, hash string, nonce int) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "pkg/solve/Solve")
	defer span.Finish()

	var hashNew big.Int

	data := s.prepare(token, ip, nonce)
	h := sha256.Sum256(data)
	hashNew.SetBytes(h[:])

	equalHash := hash == hex.EncodeToString(h[:])

	if !equalHash || hashNew.Cmp(s.target) != -1 {
		return "", errors.New("hash is not valid")
	}

	res, err := s.quotes.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("couldn't get phrase: %w", err)
	}

	return res.Content, nil
}

func (s *Service) prepare(token string, ip string, nonce int) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(token),
			[]byte(ip),
			[]byte(strconv.FormatInt(int64(nonce), 16)),
		},
		[]byte{},
	)
}
