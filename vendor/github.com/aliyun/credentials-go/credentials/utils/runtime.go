package utils

import (
	"context"
	"net"
	"time"
)

// Runtime is for setting timeout, proxy and host
// Deprecated: it was used for internal
type Runtime struct {
	ReadTimeout    int
	ConnectTimeout int
	Proxy          string
	Host           string
	STSEndpoint    string
}

// NewRuntime returns a Runtime
// Deprecated: it was used for internal
func NewRuntime(readTimeout, connectTimeout int, proxy string, host string) *Runtime {
	return &Runtime{
		ReadTimeout:    readTimeout,
		ConnectTimeout: connectTimeout,
		Proxy:          proxy,
		Host:           host,
	}
}

// Timeout is for connect Timeout
// Deprecated: it was used for internal
func Timeout(connectTimeout time.Duration) func(cxt context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{
			Timeout:   connectTimeout,
			DualStack: true,
		}).DialContext(ctx, network, address)
	}
}
