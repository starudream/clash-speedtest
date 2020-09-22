package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
)

const DefaultHTTPTimeout = 30 * time.Second

func HTTPGet(url string, headers map[string]string, data interface{}) (str string, err error) {
	return HTTP(http.MethodGet, url, headers, nil, data)
}

func HTTPPut(url string, headers map[string]string, body, data interface{}) (str string, err error) {
	return HTTP(http.MethodPut, url, headers, body, data)
}

func HTTPPatch(url string, headers map[string]string, body, data interface{}) (str string, err error) {
	return HTTP(http.MethodPatch, url, headers, body, data)
}

func HTTP(method, url string, headers map[string]string, body, data interface{}) (str string, err error) {
	var code int
	var reqBytes, respBytes []byte

	if body != nil {
		reqBytes, err = json.Marshal(body)
		if err != nil {
			return "", err
		}
	}

	defer func(begin time.Time) {
		logx.
			WithField("code", code).
			WithField("took", int64(time.Now().Sub(begin)/time.Millisecond)).
			Debugf("[%s] %s, request: %s, headers: %s, response: %s", method, url, json.ReMarshal(reqBytes), json.MustMarshal(headers), json.ReMarshal(respBytes))
	}(time.Now())

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBytes))
	if err != nil {
		return "", err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := (&http.Client{Timeout: DefaultHTTPTimeout}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	code = resp.StatusCode

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	str = string(respBytes)

	if code < 200 || code >= 300 {
		return str, fmt.Errorf("http: response status code is not successful")
	}

	if data != nil {
		return str, json.Unmarshal(respBytes, data)
	}

	return str, nil
}
