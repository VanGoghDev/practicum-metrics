package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignerTripper_SignBody(t *testing.T) {
	type fields struct {
		Proxied http.RoundTripper
		key     string
	}
	type args struct {
		reqBody string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid sign",
			fields: fields{
				Proxied: nil,
				key:     "secret",
			},
			args: args{
				reqBody: "ddd",
			},
			want: "ddd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &SignerTripper{
				Proxied: tt.fields.Proxied,
				key:     tt.fields.key,
			}
			buf := bytes.NewBufferString(tt.args.reqBody)
			req, err := http.NewRequest(http.MethodPost, "/test", buf)
			if err != nil {
				t.Error("failed to create request")
			}

			got, err := st.SignBody(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignerTripper.SignBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sign := getSignature(got, tt.fields.key)
			hV := got.Header.Get("HashSHA256")
			assert.NotEmpty(t, hV)
			assert.True(t, hmac.Equal(sign, []byte(hV)))
		})
	}
}

func getSignature(req *http.Request, secretkey string) []byte {
	body, _ := io.ReadAll(req.Body)
	h := hmac.New(sha256.New, []byte(secretkey))
	h.Write(body)
	sign := h.Sum(nil)
	return sign
}
