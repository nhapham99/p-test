package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

func ParseJWS(jws string) (header, payload, signature string, err error) {
	parts := strings.Split(jws, ".")
	if len(parts) != 3 {
		err = errors.New("invalid JWS")
		return
	}
	header = parts[0]
	payload = parts[1]
	signature = parts[2]
	return header, payload, signature, nil
}

func Base64UrlDecode(s string) ([]byte, error) {
	decodedBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(s)
	if err != nil {
		return nil, err
	}

	return decodedBytes, nil
}
