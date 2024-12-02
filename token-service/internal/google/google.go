package google

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	certEndpoint = "https://www.googleapis.com/oauth2/v1/certs"
)

func GetPublicKey(keyID string) (string, error) {
	resp, err := http.Get(certEndpoint)
	if err != nil {
		return "", err
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := map[string]string{}
	if err = json.Unmarshal(raw, &data); err != nil {
		return "", err
	}
	key, exists := data[keyID]
	if !exists {
		return "", errors.New("Key does not exists")
	}

	return key, nil
}
