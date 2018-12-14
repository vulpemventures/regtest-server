package router

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func httpPOST(url string, bodyString string, header map[string]string) (int, string, error) {

	body := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, "", err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	rs, err := client.Do(req)
	if err != nil {
		return 0, "", errors.New("Failed to create named key request: " + err.Error())
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return 0, "", errors.New("Failed to parse response body: " + err.Error())
	}

	return rs.StatusCode, string(bodyBytes), nil
}
