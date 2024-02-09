package connect

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func init() {
	proxy.RegisterDialerType("http", newConnectDialer)
	proxy.RegisterDialerType("https", newConnectDialer)
}

type ConnectDialer struct {
	proxyURL *url.URL
	forward  proxy.Dialer
	client   *http.Client
}

func newConnectDialer(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	s := new(ConnectDialer)
	s.proxyURL = uri
	s.forward = forward

	// set proxy to nil in the http client
	s.client = &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return s, nil
}

func (cd *ConnectDialer) Dial(network, addr string) (c net.Conn, err error) {
	c, err = cd.forward.Dial("tcp", cd.proxyURL.Host)
	if err != nil {
		return nil, err
	}

	reqURL, err := url.Parse("http://" + addr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodConnect, reqURL.String(), nil)
	if err != nil {
		c.Close()
		return nil, err
	}
	req.Close = false

	if cd.proxyURL.User != nil {
		passwd, _ := cd.proxyURL.User.Password()
		req.SetBasicAuth(cd.proxyURL.User.Username(), passwd)
	}
	req.Header.Set("User-Agent", "GoLang Connect Proxy Dialer")

	err = req.Write(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		c.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		c.Close()
		return nil, fmt.Errorf("connection to proxy server StatusCode: %d", resp.StatusCode)
	}

	return c, nil
}
