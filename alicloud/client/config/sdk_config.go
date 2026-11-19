package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
)

func (c *Config) SDKConfig(timeout time.Duration) *sdk.Config {
	// The sdk-level timeout follows the provider client_read_timeout and
	// falls back to 30s when the value is non-positive, so callers without a
	// provider configuration (e.g. the sweeper client) keep a request deadline
	// instead of dropping to the sdk default.
	if timeout <= 0 {
		timeout = time.Duration(30) * time.Second
	}
	// WithUserAgent will add a prefix Extra/ for user agent value
	return sdk.NewConfig().
		WithAutoRetry(true).
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(timeout).
		WithEnableAsync(false).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(c.getTransport()).
		WithScheme(c.Protocol).
		WithUserAgent(fmt.Sprintf("Terraform %s", c.getUserAgent()))
}

func (c *Config) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	idleConnTimeout, err := strconv.Atoi(os.Getenv("IdleConnTimeout"))
	if err != nil || idleConnTimeout <= 0 || idleConnTimeout > 86400 {
		idleConnTimeout = 5
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second
	transport.IdleConnTimeout = time.Duration(idleConnTimeout) * time.Second

	return transport
}
