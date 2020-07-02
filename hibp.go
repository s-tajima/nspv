package nspv

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	HIBP_API_BASE_URL = "https://api.pwnedpasswords.com"
)

type hibpClient struct {
	baseUrl string
	client  *http.Client
	ctx     context.Context
}

func newHibpClient() *hibpClient {
	hc := hibpClient{}
	hc.baseUrl = HIBP_API_BASE_URL
	hc.client = &http.Client{}

	return &hc
}

func (hc *hibpClient) pwnedCount(password string) (int, error) {
	h := sha1.New()
	h.Write([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	body, err := hc.callRange(hash[:5])
	if err != nil {
		return -1, err
	}

	hashes, err := hc.parseRange(body)
	if err != nil {
		return -1, err
	}

	if c, ok := hashes[hash[5:]]; ok {
		return c, nil
	}

	return 0, nil
}

func (hc *hibpClient) callRange(hash5 string) (string, error) {
	url := fmt.Sprintf("%s/range/%s", HIBP_API_BASE_URL, hash5)

	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)
	req.Header.Add("Add-Padding", "true")

	if hc.ctx != nil {
		req = req.WithContext(hc.ctx)
	}

	response, err := hc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Request failed (%d)", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (hc *hibpClient) parseRange(body string) (map[string]int, error) {
	strs := strings.Split(body, "\r\n")

	hashes := map[string]int{}
	for _, str := range strs {
		s := strings.Split(str, ":")

		n, _ := strconv.Atoi(s[1])
		if n < 0 {
			continue
		}

		hashes[s[0]] = n
	}

	return hashes, nil
}
