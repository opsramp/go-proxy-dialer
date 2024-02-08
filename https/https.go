package https

import (
	"crypto/tls"
	"net"
)

var HTTPSDialer = HTTPSProxy{}

type HTTPSProxy struct{}

func (h *HTTPSProxy) Dial(network, addr string) (net.Conn, error) {
	c, err := tls.Dial("tcp", addr, &tls.Config{})
	if err != nil {
		return nil, err
	}
	return c, err
}
