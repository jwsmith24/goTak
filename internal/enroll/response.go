package enroll

import (
	"fmt"
	"io"
)

const (
	// TLS configuration XML contains only a small set of subject fields.
	maxTLSConfigResponseBytes = 1 << 20
	// Signing JSON may contain a leaf certificate and several CA certificates.
	maxSignResponseBytes = 4 << 20
)

func readResponseBody(r io.Reader, limit int64, name string) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > limit {
		return nil, fmt.Errorf("enroll: %s response is too large (limit %d bytes)", name, limit)
	}
	return body, nil
}
