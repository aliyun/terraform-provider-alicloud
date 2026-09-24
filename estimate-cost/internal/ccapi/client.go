// Package ccapi is a client for the CloudControl GetApiPrice endpoint.
//
// Requests are signed with POP v3 (ACS3-HMAC-SHA256). The AK/SK pair is read
// from the user's environment — see NewFromEnv for the variable name
// precedence.
package ccapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	// Default CloudControl endpoint. The HTTP Host header is the same
	// hostname. To target a different endpoint, set ALICLOUD_CC_API_ENDPOINT
	// (and optionally ALICLOUD_CC_API_HOST) — see NewFromEnv.
	defaultEndpoint = "cloudcontrol.aliyuncs.com"
	defaultHost     = "cloudcontrol.aliyuncs.com"
	apiVersion      = "2022-08-30"
	apiAction       = "GetApiPrice"
	apiPath         = "/api/v1/price/quote"
)

// Request is the GetApiPrice request body: the (popCode, popVersion, apiName,
// parameters) four-tuple that names an Alibaba Cloud OpenAPI operation and
// the parameters to price.
type Request struct {
	PopCode    string                 `json:"popCode"`
	PopVersion string                 `json:"popVersion"`
	APIName    string                 `json:"apiName"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Response decodes a GetApiPrice response. The endpoint returns one of three
// shapes: success, upstream business failure, or a system-level error.
type Response struct {
	RequestID string `json:"requestId"`
	Price     *Price `json:"price"`

	// System errors (HTTP 4xx) use CamelCase keys at the top level.
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	HostID  string `json:"HostId,omitempty"`
}

// Price mirrors the response's `price.*` section.
type Price struct {
	PricingMode       string                  `json:"pricingMode"` // single / composite / delta
	Success           bool                    `json:"success"`
	Currency          string                  `json:"currency"`
	CalculatedAmount  *float64                `json:"calculatedAmount"` // composite/delta total
	PriceSummary      *PriceSummary           `json:"priceSummary"`     // single-mode total
	Components        map[string]PriceSummary `json:"components"`       // composite/delta per-component breakdown
	UpstreamRequestID string                  `json:"upstreamRequestId"`
	ErrorCode         string                  `json:"errorCode"`
	ErrorMessage      string                  `json:"errorMessage"`
}

// PriceSummary is the core amount structure. priceSummary and each entry of
// components[*] use the same shape.
type PriceSummary struct {
	Currency           string  `json:"currency"`
	TradePrice         float64 `json:"tradePrice"`
	OriginalPrice      float64 `json:"originalPrice"`
	ModuleSum          float64 `json:"moduleSum"`          // per-unit subtotal
	Quantity           float64 `json:"quantity"`
	EffectiveModuleSum float64 `json:"effectiveModuleSum"` // ★ final user-paid amount (includes quantity)
	Modules            []struct {
		ModuleCode        string  `json:"moduleCode"`
		CostAfterDiscount float64 `json:"costAfterDiscount"`
		OriginalCost      float64 `json:"originalCost"`
	} `json:"modules"`
}

// FinalAmount flattens the single / composite / delta variants and returns
// the unified final amount. Per the GetApiPrice spec, single mode reads
// priceSummary.effectiveModuleSum and composite/delta read calculatedAmount.
func (p *Price) FinalAmount() (float64, string) {
	if p == nil {
		return 0, ""
	}
	if p.PricingMode == "composite" || p.PricingMode == "delta" {
		if p.CalculatedAmount != nil {
			return *p.CalculatedAmount, p.Currency
		}
	}
	if p.PriceSummary != nil {
		return p.PriceSummary.EffectiveModuleSum, p.PriceSummary.Currency
	}
	return 0, p.Currency
}

// Client is the GetApiPrice HTTP client.
type Client struct {
	AK       string
	SK       string
	Endpoint string // default: cloudcontrol.aliyuncs.com (production)
	Host     string // default: cloudcontrol.aliyuncs.com (production)
	HTTP     *http.Client
}

// NewFromEnv builds a Client from environment variables.
//
// Credentials. The ALIBABA_CLOUD_ACCESS_KEY_ID / ALIBABA_CLOUD_ACCESS_KEY_SECRET
// pair (matching the aliyun CLI) takes precedence; it falls back to ALICLOUD_*
// (matching the Terraform alicloud provider) when the canonical pair is unset.
//
// Endpoint. By default the client targets the CloudControl endpoint
// cloudcontrol.aliyuncs.com. To target a different endpoint, override with:
//
//	ALICLOUD_CC_API_ENDPOINT  alternate request hostname
//	ALICLOUD_CC_API_HOST      alternate HTTP Host header (rarely needed)
//
// Setting ENDPOINT alone is enough for most cases; HOST is only needed when
// a custom gateway expects a Host header different from the endpoint
// hostname.
func NewFromEnv() (*Client, error) {
	ak := firstEnv("ALIBABA_CLOUD_ACCESS_KEY_ID", "ALICLOUD_ACCESS_KEY")
	sk := firstEnv("ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALICLOUD_SECRET_KEY")
	if ak == "" || sk == "" {
		return nil, fmt.Errorf("env vars ALIBABA_CLOUD_ACCESS_KEY_ID/SECRET (or ALICLOUD_*) not set")
	}
	endpoint := defaultEndpoint
	if v := os.Getenv("ALICLOUD_CC_API_ENDPOINT"); v != "" {
		endpoint = v
	}
	host := defaultHost
	if v := os.Getenv("ALICLOUD_CC_API_HOST"); v != "" {
		host = v
	}
	return &Client{
		AK:       ak,
		SK:       sk,
		Endpoint: endpoint,
		Host:     host,
		HTTP:     &http.Client{Timeout: 15 * time.Second},
	}, nil
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// Quote performs a single GetApiPrice call.
func (c *Client) Quote(req *Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := c.signedRequest(body)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("CC API HTTP call failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		// System-level error (PricingNotSupported / InvalidParameter / ...).
		var sysErr Response
		_ = json.Unmarshal(raw, &sysErr)
		return &sysErr, fmt.Errorf("CC API system error HTTP %d %s: %s", resp.StatusCode, sysErr.Code, sysErr.Message)
	}
	var out Response
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("parse response failed: %w raw: %s", err, raw)
	}
	return &out, nil
}

// signedRequest builds the HTTP request signed with POP v3 (ACS3-HMAC-SHA256).
//
// Signing rules in brief:
//   - signed headers include only: host, content-type, and every x-acs-* header
//   - canonical headers are sorted by lowercase name and joined as "name:value\n"
//   - StringToSign = "ACS3-HMAC-SHA256\n" + sha256(canonicalRequest) hex
//   - signature = hmacsha256(SK, StringToSign) hex
func (c *Client) signedRequest(body []byte) (*http.Request, error) {
	bodyHash := sha256Hex(body)
	nonce, _ := randomNonce()
	date := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	headers := map[string]string{
		"host":                  c.Host,
		"content-type":          "application/json",
		"x-acs-action":          apiAction,
		"x-acs-version":         apiVersion,
		"x-acs-date":            date,
		"x-acs-signature-nonce": nonce,
		"x-acs-content-sha256":  bodyHash,
	}

	signedHeaderNames := make([]string, 0, len(headers))
	for k := range headers {
		signedHeaderNames = append(signedHeaderNames, k)
	}
	sort.Strings(signedHeaderNames)

	var canonicalHeaders strings.Builder
	for _, k := range signedHeaderNames {
		canonicalHeaders.WriteString(k)
		canonicalHeaders.WriteByte(':')
		canonicalHeaders.WriteString(strings.TrimSpace(headers[k]))
		canonicalHeaders.WriteByte('\n')
	}
	signedHeaders := strings.Join(signedHeaderNames, ";")

	canonicalRequest := strings.Join([]string{
		"POST",
		apiPath,
		"", // no query string
		canonicalHeaders.String(),
		signedHeaders,
		bodyHash,
	}, "\n")

	stringToSign := "ACS3-HMAC-SHA256\n" + sha256Hex([]byte(canonicalRequest))
	signature := hmacSha256Hex(c.SK, stringToSign)

	auth := fmt.Sprintf("ACS3-HMAC-SHA256 Credential=%s,SignedHeaders=%s,Signature=%s",
		c.AK, signedHeaders, signature)

	u := url.URL{Scheme: "https", Host: c.Endpoint, Path: apiPath}
	httpReq, err := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	// Host needs special handling: in Go the value lives on Request.Host,
	// not in the Header map.
	httpReq.Host = c.Host
	httpReq.Header.Set("Authorization", auth)
	return httpReq, nil
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func hmacSha256Hex(key, msg string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func randomNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
