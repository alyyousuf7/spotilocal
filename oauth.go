package spotilocal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	ORIGIN_URL = "https://open.spotify.com"
	TOKEN_PATH = "/token"
)

func fetchToken() (string, error) {
	return fetchTokenFromURL(ORIGIN_URL + TOKEN_PATH)
}

func fetchTokenFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tokenBody := struct {
		Token string `json:"t"`
	}{}
	if err := json.Unmarshal(buf, &tokenBody); err != nil {
		return "", err
	}

	return tokenBody.Token, nil
}
