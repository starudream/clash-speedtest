package speedtest

import (
	"encoding/xml"
	"fmt"
)

type Config struct {
	XMLName xml.Name `xml:"settings"`
	Client  struct {
		Ip        string `xml:"ip,attr"`
		Country   string `xml:"country,attr"`
		Lat       string `xml:"lat,attr"`
		Lon       string `xml:"lon,attr"`
		Isp       string `xml:"isp,attr"`
		IspRating string `xml:"isprating,attr"`
		Rating    string `xml:"rating,attr"`
		IspDlAvg  string `xml:"ispdlavg,attr"`
		IspUlAvg  string `xml:"ispulavg,attr"`
		LoggedIn  string `xml:"loggedin,attr"`
	} `xml:"client"`
	ServerConfig struct {
		ThreadCount       string `xml:"threadcount,attr"`
		IgnoreIds         string `xml:"ignoreids,attr"`
		NotOnMap          string `xml:"notonmap,attr"`
		ForcePingId       string `xml:"forcepingid,attr"`
		PreferredServerId string `xml:"preferredserverid,attr"`
	} `xml:"server-config"`
	LicenseKey string `xml:"licensekey"`
	Customer   string `xml:"customer"`
	Odometer   struct {
		Start string `xml:"start,attr"`
		Rate  string `xml:"rate,attr"`
	} `xml:"odometer"`
	Times struct {
		Dl1 string `xml:"dl1,attr"`
		Dl2 string `xml:"dl2,attr"`
		Dl3 string `xml:"dl3,attr"`
		Ul1 string `xml:"ul1,attr"`
		Ul2 string `xml:"ul2,attr"`
		Ul3 string `xml:"ul3,attr"`
	} `xml:"times"`
	Download struct {
		TestLength    string `xml:"testlength,attr"`
		InitialTest   string `xml:"initialtest,attr"`
		MinTestSize   string `xml:"mintestsize,attr"`
		ThreadsPerUrl string `xml:"threadsperurl,attr"`
	} `xml:"download"`
	Upload struct {
		TestLength    string `xml:"testlength,attr"`
		Ratio         string `xml:"ratio,attr"`
		InitialTest   string `xml:"initialtest,attr"`
		MinTestSize   string `xml:"mintestsize,attr"`
		Threads       string `xml:"threads,attr"`
		MaxChunkSize  string `xml:"maxchunksize,attr"`
		MaxChunkCount string `xml:"maxchunkcount,attr"`
		ThreadsPerUrl string `xml:"threadsperurl,attr"`
	} `xml:"upload"`
	Latency struct {
		TestLength string `xml:"testlength,attr"`
		WaitTime   string `xml:"waittime,attr"`
		Timeout    string `xml:"timeout,attr"`
	} `xml:"latency"`
	SocketDownload struct {
		TestLength      string `xml:"testlength,attr"`
		InitialThreads  string `xml:"initialthreads,attr"`
		MinThreads      string `xml:"minthreads,attr"`
		MaxThreads      string `xml:"maxthreads,attr"`
		ThreadRatio     string `xml:"threadratio,attr"`
		MaxSampleSize   string `xml:"maxsamplesize,attr"`
		MinSampleSize   string `xml:"minsamplesize,attr"`
		StartSampleSize string `xml:"startsamplesize,attr"`
		StartBufferSize string `xml:"startbuffersize,attr"`
		BufferLength    string `xml:"bufferlength,attr"`
		PacketLength    string `xml:"packetlength,attr"`
		ReadBuffer      string `xml:"readbuffer,attr"`
	} `xml:"socket-download"`
	SocketUpload struct {
		TestLength      string `xml:"testlength,attr"`
		InitialThreads  string `xml:"initialthreads,attr"`
		MinThreads      string `xml:"minthreads,attr"`
		MaxThreads      string `xml:"maxthreads,attr"`
		ThreadRatio     string `xml:"threadratio,attr"`
		MaxSampleSize   string `xml:"maxsamplesize,attr"`
		MinSampleSize   string `xml:"minsamplesize,attr"`
		StartSampleSize string `xml:"startsamplesize,attr"`
		StartBufferSize string `xml:"startbuffersize,attr"`
		BufferLength    string `xml:"bufferlength,attr"`
		PacketLength    string `xml:"packetlength,attr"`
		Disabled        string `xml:"disabled,attr"`
	} `xml:"socket-upload"`
	SocketLatency struct {
		TestLength string `xml:"testlength,attr"`
		WaitTime   string `xml:"waittime,attr"`
		Timeout    string `xml:"timeout,attr"`
	} `xml:"socket-latency"`
	Translation struct {
		Lang string `xml:"lang,attr"`
	} `xml:"translation"`
}

func (c *Client) GetConfig() (*Config, error) {
	resp, err := c.c.R().Get("https://www.speedtest.net/speedtest-config.php")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	config := &Config{}
	err = xml.Unmarshal(resp.Body(), config)
	return config, err
}
