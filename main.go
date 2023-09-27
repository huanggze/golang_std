package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/huanggze/src/net"
	"github.com/huanggze/src/net/http"
)

func main() {
	method := "PUT"
	println(method[:1] + strings.ToLower(method[1:]))

	c := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Proxy:        http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   150,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.baidu.com", bytes.NewBufferString(`{"data": "123"}`))
	if err != nil {
		panic(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
