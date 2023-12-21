package speedtest

import (
	"encoding/xml"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/starudream/go-lib/core/v2/slog"
)

type Servers struct {
	XMLName xml.Name `xml:"settings"`
	Servers struct {
		Server []*Server `xml:"server"`
	} `xml:"servers"`
}

type Server struct {
	Id      string `xml:"id,attr"`
	Country string `xml:"cc,attr"`
	City    string `xml:"name,attr"`
	Name    string `xml:"country,attr"`
	Lat     string `xml:"lat,attr"`
	Lon     string `xml:"lon,attr"`
	Host    string `xml:"host,attr"`
	URL     string `xml:"url,attr"`
	Sponsor string `xml:"sponsor,attr"`
}

func (s Server) BaseURL(f string, a ...any) string {
	base, _ := path.Split(s.URL)
	suffix := strings.TrimPrefix(f, "/")
	if len(a) == 0 {
		return base + suffix
	}
	return base + fmt.Sprintf(suffix, a...)
}

func (c *Client) GetServers() ([]*Server, error) {
	resp, err := c.c.R().Get("https://www.speedtest.net/speedtest-servers-static.php")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	var servers Servers
	err = xml.Unmarshal(resp.Body(), &servers)
	if err != nil {
		return nil, err
	}

	return servers.Servers.Server, nil
}

func (c *Client) LatencyServer(s *Server) (time.Duration, error) {
	resp, err := c.c.R().Get(s.BaseURL("latency.txt?x=%d", time.Now().UnixMilli()))
	if err != nil {
		return 0, err
	}

	if !resp.IsSuccess() {
		return 0, fmt.Errorf("response status %s", resp.Status())
	}

	if strings.TrimSuffix(string(resp.Body()), "\n") != "test=test" {
		slog.Warn("unrecognized response body: %s", resp.Body())
	}

	return resp.Time(), nil
}
