package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/starudream/clash-speedtest/util"
)

type Result struct {
	AllBytes   int64
	TotalBytes int64
	TotalTime  time.Duration

	err error

	finish chan bool
}

const (
	SpeedTestTimeout = 30 * time.Second
)

func SpeedTest(url string, timeout time.Duration, process bool) (*Result, error) {
	if process {
		fmt.Printf("\n")
		defer fmt.Printf("\n\n")
	}

	begin := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{Transport: util.Transport()}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &Result{finish: make(chan bool)}

	result.AllBytes, _ = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	go func() {
		for {
			bs := make([]byte, 1024)
			n, err := resp.Body.Read(bs)
			result.TotalBytes += int64(n)
			if process {
				fmt.Printf("\r  %.02f%% %dkb/%dkb", float64(result.TotalBytes)/float64(result.AllBytes)*100, result.TotalBytes/1024, result.AllBytes/1024)
			}
			if err != nil {
				if err != io.EOF {
					result.err = err
				}
				result.finish <- true
				return
			}
		}
	}()

	if timeout < SpeedTestTimeout {
		timeout = SpeedTestTimeout
	}

	select {
	case <-result.finish:
		result.TotalTime = time.Now().Sub(begin)
		return result, err
	case <-time.After(timeout):
		result.TotalTime = SpeedTestTimeout
		return result, nil
	}
}
