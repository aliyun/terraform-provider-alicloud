package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	pop "github.com/alibabacloud-go/alibabacloud-gateway-pop/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	"github.com/alibabacloud-go/tea/tea"
	fc "github.com/aliyun/fc-go-sdk"
)

// fcV4ClientOption installs the V4 signing transport on the owning alias's Resty client.
func fcV4ClientOption(client *AliyunClient, apiVersion, fcHost string, base *http.Transport) fc.ClientOption {
	return func(*fc.Client) {
		client.fcHTTP.SetTransport(&fcV4Transport{client: client, apiVersion: apiVersion, fcHost: fcHost, base: base})
	}
}

// fcV4Transport adapts requests produced by the legacy fc-go-sdk client (FC
// 1.0, 2016-08-15) to the shared provider-level V4 signature layer. The SDK
// signs every request with its own legacy HMAC scheme before dispatch, and it
// mutates the process-wide resty client instead of accepting an injectable
// transport, so the V4 signature is applied on the final wire request inside
// the transport and the legacy Authorization header is dropped there.
//
// Only requests addressed to the FC endpoint this transport was built for are
// signed; anything else routed through the shared resty client passes through
// unchanged.
type fcV4Transport struct {
	client     *AliyunClient
	apiVersion string
	fcHost     string
	base       *http.Transport
}

func (t *fcV4Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !strings.EqualFold(requestHost(req), t.fcHost) {
		return t.base.RoundTrip(req)
	}
	if err := t.sign(req); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

func requestHost(req *http.Request) string {
	if req.Host != "" {
		return req.Host
	}
	return req.URL.Host
}

// sign replaces the legacy FC signature on the outgoing request with an
// ACS4-HMAC-SHA256 one computed by the shared gateway signer. The FC 1.0 API
// is a REST surface without an action header, so the signed request keeps the
// SDK's pathname, query, payload and business headers (content-type, MD5)
// untouched; only the credential headers change.
func (t *fcV4Transport) sign(req *http.Request) error {
	gateway, err := pop.NewClient()
	if err != nil {
		return err
	}
	var payload []byte
	if req.Body != nil {
		payload, err = io.ReadAll(req.Body)
		if err != nil {
			return err
		}
	}

	headers := make(map[string]*string)
	for name, values := range req.Header {
		if len(values) == 0 {
			continue
		}
		// The pop signer lowercases header names and only signs host,
		// content-type and x-acs-* values; the legacy credential headers are
		// superseded by the V4 ones and must not reach the wire.
		switch lower := strings.ToLower(name); lower {
		case "authorization", "date", "host", "x-fc-security-token":
			continue
		default:
			headers[lower] = tea.String(values[0])
		}
	}
	query := make(map[string]*string)
	for name, values := range req.URL.Query() {
		if len(values) > 0 {
			query[name] = tea.String(values[0])
		}
	}

	signed := &spi.InterceptorContextRequest{
		ProductId: tea.String("fc"), Version: tea.String(t.apiVersion),
		Credential: t.client.config.Credential, AuthType: tea.String("AK"),
		UserAgent: tea.String(req.Header.Get("User-Agent")),
		Pathname:  tea.String(req.URL.EscapedPath()), Method: tea.String(req.Method),
		Query: query, Headers: headers, ReqBodyType: tea.String("byte"),
	}
	if len(payload) > 0 {
		signed.Body = payload
	}
	err = gateway.ModifyRequest(&spi.InterceptorContext{Request: signed, Configuration: &spi.InterceptorContextConfiguration{
		RegionId: tea.String(t.client.config.RegionId), Endpoint: tea.String(requestHost(req)),
	}}, nil)
	if err != nil {
		return err
	}

	req.Header = make(http.Header)
	for name, value := range signed.Headers {
		if name == "host" || value == nil {
			continue
		}
		req.Header.Set(name, tea.StringValue(value))
	}
	if len(payload) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(payload))
		req.ContentLength = int64(len(payload))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(payload)), nil
		}
	}
	return nil
}
