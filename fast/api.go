package fast

import (
	"fmt"
	"regexp"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"

	"github.com/starudream/clash-speedtest/util"
)

const (
	URL   = "https://api.fast.com/netflix/speedtest/v2"
	Count = 1
)

type Data struct {
	Client  *DataClient   `json:"client"`
	Targets []*DataTarget `json:"targets"`
}

type DataClient struct {
	ASN      string        `json:"asn"`
	IP       string        `json:"ip"`
	Location *DataLocation `json:"location"`
}

type DataTarget struct {
	URL      string        `json:"url"`
	Name     string        `json:"name"`
	Location *DataLocation `json:"location"`
}

type DataLocation struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

var (
	getURL = ""
)

func GetData() (*Data, error) {
	if getURL == "" {
		htmlBody, err := util.HTTPGet("https://fast.com", nil, nil)
		if err != nil {
			return nil, err
		}

		js := regexp.MustCompile("app-.*\\.js").FindString(htmlBody)
		if js == "" {
			return nil, fmt.Errorf("not found app-*.js")
		}

		jsBody, err := util.HTTPGet("https://fast.com/"+js, nil, nil)
		if err != nil {
			return nil, err
		}

		token := regexp.MustCompile("token:\"[a-zA-Z]*\"").FindString(jsBody)
		if len(token) <= 8 {
			return nil, fmt.Errorf("not found token")
		}
		token = token[7 : len(token)-1]

		getURL = fmt.Sprintf("%s?https=%t&token=%s&urlCount=%d", URL, true, token, Count)

		logx.Debugf("[fast.com] %s", getURL)
	}

	body, err := util.HTTPGet(getURL, nil, nil)
	if err != nil {
		return nil, err
	}

	data := &Data{}
	return data, json.UnmarshalFromString(body, data)
}
