package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
)

// SignerTripper Подписывает запросы алгоритмом sha256.
type SignerTripper struct {
	Proxied http.RoundTripper
	key     string
}

func New(cfg *config.Config) *SignerTripper {
	return &SignerTripper{
		key: cfg.Key,
	}
}

// SignBody подписывает запрос алгоритмом sha256.
// Возвращает *http.Request, в котором проставлен
// HTTP заголовок (header) HashSHA256, значение которого - hash от тела запроса.
// Хэш считается с учетом ключа, передаваемом через config.
func (st *SignerTripper) SignBody(req *http.Request) (*http.Request, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	defer func() {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}()

	h := hmac.New(sha256.New, []byte(st.key))
	h.Write(body)
	dst := h.Sum(nil)
	v := base64.StdEncoding.EncodeToString(dst)
	req.Header.Set("HashSHA256", v)

	return req, nil
}

func (st *SignerTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req, err := st.SignBody(req)
	if err != nil {
		return nil, errors.New("failed to sign body")
	}

	res, err := st.Proxied.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to round trip from signer tripper: %w", err)
	}
	return res, nil
}
