package client

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	pop "github.com/alibabacloud-go/alibabacloud-gateway-pop/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	"github.com/alibabacloud-go/tea/tea"
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/valyala/fasthttp"
)

// The MNS data plane API version, kept identical to the legacy SDK wire value
// (the ali_mns package keeps its copy unexported).
const mnsV4APIVersion = "2015-06-06"

// mnsV4Client lets the MNS REST data plane sign with ACS4-HMAC-SHA256 when
// sign_version selects v4. The legacy ali_mns.MNSClient interface carries the
// unexported getAccountID/getRegion methods the SDK managers rely on, so the
// wrapper embeds the interface to promote them instead of re-implementing the
// service layer. The embedded client keeps its URL for account/region parsing;
// this wrapper owns the actual transport and signature.
type mnsV4Client struct {
	ali_mns.MNSClient
	// sendURL is the base URL every request is sent to. It matches the embedded
	// client's URL in production; wire-capture tests point it at a local server.
	sendURL string
	// sendHost is the host:port of sendURL, pre-parsed so the value handed to
	// the signature gateway is always the host actually on the wire.
	sendHost string
	// gateway signs each request with the shared POP ACS4 implementation.
	gateway *pop.Client
	// httpclient applies the provider proxy policy to the signed request.
	httpclient *http.Client
	// owner supplies credentials, the physical signing region and the user
	// agent at request time, so refreshed temporary credentials are honored.
	owner *AliyunClient
}

// newMnsV4Client wraps a configured legacy MNS client for v4 signing.
func newMnsV4Client(wrapped ali_mns.MNSClient, sendURL string, owner *AliyunClient) (*mnsV4Client, error) {
	gateway, err := pop.NewClient()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the MNS v4 signing gateway: %#v", err)
	}
	parsed, err := url.Parse(sendURL)
	if err != nil || parsed.Host == "" {
		return nil, fmt.Errorf("invalid MNS endpoint url %q", sendURL)
	}
	timeout := time.Second * time.Duration(ali_mns.DefaultTimeout)
	proxy, err := owner.getHttpProxy()
	if err != nil {
		return nil, err
	}
	transport := owner.getSdkConfig().HttpTransport
	if proxy != nil {
		skip, err := owner.skipProxy(parsed.Host)
		if err != nil {
			return nil, err
		}
		if !skip {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}
	return &mnsV4Client{
		MNSClient: wrapped,
		sendURL:   sendURL,
		sendHost:  parsed.Host,
		gateway:   gateway,
		httpclient: &http.Client{
			Transport: transport, Timeout: timeout,
			// Credentials and signatures are scoped to the original endpoint.
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		owner: owner,
	}, nil
}

// Send signs and sends one MNS REST request with the v4 signature. Body
// serialization, the Content-MD5/Date/x-mns-version headers and the no-retry
// semantics mirror the legacy client so responses stay byte-compatible; the
// Authorization/x-acs-* headers carry the ACS4 signature instead of the legacy
// MNS HMAC-SHA1 one.
func (c *mnsV4Client) Send(method ali_mns.Method, headers map[string]string, message interface{}, resource string) (*fasthttp.Response, error) {
	// Serialize the body exactly like the legacy client: nil -> empty body,
	// []byte verbatim, everything else XML-marshaled.
	var xmlContent []byte
	if message == nil {
		xmlContent = []byte{}
	} else {
		switch m := message.(type) {
		case []byte:
			xmlContent = m
		default:
			bXml, err := xml.Marshal(message)
			if err != nil {
				return nil, fmt.Errorf("mns v4: marshal message failed: %#v", err)
			}
			xmlContent = bXml
		}
	}

	// Keep the legacy wire headers, with the same Content-MD5 digest algorithm
	// (base64 of the hex md5). The v4 signature additionally protects the body
	// through x-acs-content-sha256.
	if headers == nil {
		headers = make(map[string]string)
	}
	xmlMD5 := md5.Sum(xmlContent)
	headers[ali_mns.CONTENT_MD5] = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", xmlMD5)))
	headers[ali_mns.CONTENT_TYPE] = "application/xml"
	headers[ali_mns.DATE] = time.Now().UTC().Format(http.TimeFormat)
	headers[ali_mns.MQ_VERSION] = mnsV4APIVersion

	// Split the resource into path and decoded query pairs. The callers embed
	// already-escaped query values (e.g. ReceiptHandle) in the resource; the
	// canonical v4 signature and the sent URL are both rebuilt from the decoded
	// values so they cannot diverge.
	path := resource
	query := make(url.Values)
	if idx := strings.Index(resource, "?"); idx >= 0 {
		path = resource[:idx]
		var err error
		query, err = url.ParseQuery(resource[idx+1:])
		if err != nil {
			return nil, fmt.Errorf("mns v4: parse resource query failed: %#v", err)
		}
	}
	flat := make(map[string]*string, len(query))
	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		flat[key] = &value
	}

	// The gateway signs the x-acs-*, host and content-type headers, so the
	// header keys handed to it are lowercased.
	signedHeaders := make(map[string]*string, len(headers)+8)
	for key, value := range headers {
		lowerKey, lowerValue := strings.ToLower(key), value
		signedHeaders[lowerKey] = &lowerValue
	}

	if c.owner.config.RegionId == "" {
		return nil, fmt.Errorf("sign_version.default = v4 requires a signing region")
	}
	signed := &spi.InterceptorContextRequest{
		ProductId:   tea.String("mns"),
		Version:     tea.String(mnsV4APIVersion),
		Credential:  c.owner.config.Credential,
		AuthType:    tea.String("AK"),
		UserAgent:   tea.String(c.owner.config.UserAgent()),
		Pathname:    tea.String("/" + path),
		Method:      tea.String(string(method)),
		Query:       flat,
		Headers:     signedHeaders,
		Body:        xmlContent,
		ReqBodyType: tea.String("byte"),
	}
	err := c.gateway.ModifyRequest(&spi.InterceptorContext{Request: signed, Configuration: &spi.InterceptorContextConfiguration{
		RegionId: tea.String(c.owner.config.RegionId), Endpoint: tea.String(c.sendHost),
	}}, nil)
	if err != nil {
		return nil, fmt.Errorf("mns v4: sign request failed: %#v", err)
	}

	// The sent query is rebuilt from the signed flat map, so a resource that
	// ever repeated a key cannot put a value on the wire the signature does
	// not cover.
	target := c.sendURL + "/" + path
	if len(flat) > 0 {
		keys := make([]string, 0, len(flat))
		for key := range flat {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		encoded := make([]string, 0, len(keys))
		for _, key := range keys {
			encoded = append(encoded, url.QueryEscape(key)+"="+url.QueryEscape(*flat[key]))
		}
		target += "?" + strings.Join(encoded, "&")
	}
	req, err := http.NewRequest(string(method), target, bytes.NewReader(xmlContent))
	if err != nil {
		return nil, err
	}
	for key, value := range signed.Headers {
		if key != "host" && value != nil {
			req.Header.Set(key, *value)
		}
	}
	response, err := c.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mns v4: send request: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// The MNS managers consume the legacy response type.
	resp := fasthttp.AcquireResponse()
	resp.SetStatusCode(response.StatusCode)
	for key, values := range response.Header {
		for _, value := range values {
			resp.Header.Add(key, value)
		}
	}
	resp.SetBodyRaw(body)
	return resp, nil
}
