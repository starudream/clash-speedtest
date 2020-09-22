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
	AllBytes   float64
	TotalBytes float64
	TotalTime  time.Duration

	finish chan bool
	err    chan error
}

const (
	Timeout = 30 * time.Second
)

func speedtest(url string, timeout time.Duration) (*Result, error) {
	fmt.Println()
	defer fmt.Println()

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

	result := &Result{}

	i, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	result.AllBytes = float64(i)

	go func() {
		for {
			bs := make([]byte, 100*1024)
			n, err := resp.Body.Read(bs)
			result.TotalBytes += float64(n)
			fmt.Printf("\r  %.02f%% %.03fkb/%.03fkb", result.TotalBytes/result.AllBytes*100, result.TotalBytes/1024, result.AllBytes/1024)
			if err != nil {
				if err != io.EOF {
					result.err <- err
					return
				}
			}
			result.finish <- true
		}
	}()

	if timeout < Timeout {
		timeout = Timeout
	}

	select {
	case v := <-result.err:
		fmt.Println()
		return nil, v
	case <-result.finish:
		fmt.Println()
		result.TotalTime = time.Now().Sub(begin)
		return result, nil
	case <-time.After(timeout):
		fmt.Println()
		result.TotalTime = Timeout
		return result, nil
	}
}
