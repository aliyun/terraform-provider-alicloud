package connectivity

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

// --- VCRLocalAddr ---

func TestUnitCommonVCRLocalAddr_Disabled(t *testing.T) {
	os.Unsetenv("VCR_PATH")
	assert.Equal(t, "", VCRLocalAddr())
}

// --- __vcr_host query param filtering ---

func TestUnitCommonFilterSignatureParams_VCRHost(t *testing.T) {
	q := url.Values{
		"Action":     {"DescribeVpcs"},
		"__vcr_host": {"vpc.aliyuncs.com"},
	}
	filtered := filterSignatureParams(q)
	assert.Equal(t, "DescribeVpcs", filtered.Get("Action"))
	assert.Empty(t, filtered.Get("__vcr_host"), "__vcr_host should be filtered out")
}

// --- VCRRandSeed ---

func TestUnitCommonVCRRandSeed_Disabled(t *testing.T) {
	os.Unsetenv("VCR_PATH")
	assert.Equal(t, int64(0), VCRRandSeed())
}

func TestUnitCommonVCRRandSeed_Deterministic(t *testing.T) {
	os.Setenv("VCR_PATH", "testdata/vcr/vpc")
	defer os.Unsetenv("VCR_PATH")

	seed1 := VCRRandSeed()
	seed2 := VCRRandSeed()
	assert.Equal(t, seed1, seed2, "same path should produce same seed")
	assert.NotEqual(t, int64(0), seed1)
}

func TestUnitCommonVCRRandSeed_DifferentPaths(t *testing.T) {
	os.Setenv("VCR_PATH", "path/a")
	seedA := VCRRandSeed()

	os.Setenv("VCR_PATH", "path/b")
	seedB := VCRRandSeed()
	os.Unsetenv("VCR_PATH")

	assert.NotEqual(t, seedA, seedB, "different paths should produce different seeds")
}

// --- vcrMatcher ---

func TestUnitCommonVCRMatcher_MethodMismatch(t *testing.T) {
	r := &http.Request{Method: "GET", URL: mustParseURL("https://vpc.aliyuncs.com/?Action=DescribeVpcs")}
	i := cassette.Request{Method: "POST", URL: "https://vpc.aliyuncs.com/?Action=DescribeVpcs"}
	assert.False(t, vcrMatcher(r, i))
}

func TestUnitCommonVCRMatcher_ExactMatch(t *testing.T) {
	r := &http.Request{Method: "POST", URL: mustParseURL("https://vpc.aliyuncs.com/?Action=CreateVpc&RegionId=cn-hangzhou")}
	i := cassette.Request{
		Method: "POST",
		URL:    "https://vpc.aliyuncs.com/?Action=CreateVpc&RegionId=cn-hangzhou",
		Form:   url.Values{"Action": {"CreateVpc"}, "RegionId": {"cn-hangzhou"}},
	}
	assert.True(t, vcrMatcher(r, i))
}

func TestUnitCommonVCRMatcher_IgnoresSignatureParams(t *testing.T) {
	r := &http.Request{
		Method: "POST",
		URL:    mustParseURL("https://vpc.aliyuncs.com/?Action=CreateVpc&AccessKeyId=ak1&Signature=sig1&Timestamp=t1&SignatureNonce=n1"),
	}
	i := cassette.Request{
		Method: "POST",
		URL:    "https://vpc.aliyuncs.com/?Action=CreateVpc&AccessKeyId=ak2&Signature=sig2&Timestamp=t2&SignatureNonce=n2",
		Form:   url.Values{"Action": {"CreateVpc"}, "AccessKeyId": {"ak2"}, "Signature": {"sig2"}, "Timestamp": {"t2"}, "SignatureNonce": {"n2"}},
	}
	assert.True(t, vcrMatcher(r, i), "should match ignoring signature params")
}

func TestUnitCommonVCRMatcher_DifferentAction(t *testing.T) {
	r := &http.Request{Method: "POST", URL: mustParseURL("https://vpc.aliyuncs.com/?Action=CreateVpc")}
	i := cassette.Request{Method: "POST", URL: "https://vpc.aliyuncs.com/?Action=DeleteVpc"}
	assert.False(t, vcrMatcher(r, i))
}

func TestUnitCommonVCRMatcher_DifferentHost(t *testing.T) {
	r := &http.Request{Method: "POST", URL: mustParseURL("https://vpc.aliyuncs.com/?Action=Test")}
	i := cassette.Request{Method: "POST", URL: "https://ecs.aliyuncs.com/?Action=Test"}
	assert.False(t, vcrMatcher(r, i))
}

func TestUnitCommonVCRMatcher_BodyDistinguishes(t *testing.T) {
	r := &http.Request{
		Method: "POST",
		URL:    mustParseURL("https://vpc.aliyuncs.com/?Action=DescribeVpcAttribute"),
		Body:   io.NopCloser(strings.NewReader("VpcId=vpc-001&RegionId=cn-hangzhou")),
	}
	// cassette Form = query params + body params (like Go's ParseForm)
	iMatch := cassette.Request{
		Method: "POST",
		URL:    "https://vpc.aliyuncs.com/?Action=DescribeVpcAttribute",
		Form:   url.Values{"Action": {"DescribeVpcAttribute"}, "VpcId": {"vpc-001"}, "RegionId": {"cn-hangzhou"}},
	}
	iDiff := cassette.Request{
		Method: "POST",
		URL:    "https://vpc.aliyuncs.com/?Action=DescribeVpcAttribute",
		Form:   url.Values{"Action": {"DescribeVpcAttribute"}, "VpcId": {"vpc-999"}, "RegionId": {"cn-hangzhou"}},
	}
	assert.True(t, vcrMatcher(r, iMatch), "same body should match")
	// Re-create request since body was consumed
	r.Body = io.NopCloser(strings.NewReader("VpcId=vpc-001&RegionId=cn-hangzhou"))
	assert.False(t, vcrMatcher(r, iDiff), "different body should not match")
}

// --- sanitizeCassette ---

func TestUnitCommonSanitizeCassette_Headers(t *testing.T) {
	i := &cassette.Interaction{
		Request: cassette.Request{
			Headers: http.Header{
				"Authorization":        {"Bearer secret"},
				"X-Acs-Security-Token": {"token123"},
				"Content-Type":         {"application/json"},
			},
		},
	}
	err := sanitizeCassette(i)
	assert.NoError(t, err)
	assert.Empty(t, i.Request.Headers["Authorization"])
	assert.Empty(t, i.Request.Headers["X-Acs-Security-Token"])
	assert.Equal(t, "application/json", i.Request.Headers.Get("Content-Type"))
}

func TestUnitCommonSanitizeCassette_URL(t *testing.T) {
	i := &cassette.Interaction{
		Request: cassette.Request{
			URL:     "https://vpc.aliyuncs.com/?AccessKeyId=LTAI123&Action=CreateVpc&SecurityToken=STS_TOKEN&Signature=abc123&SignatureNonce=nonce1",
			Headers: http.Header{},
		},
	}
	err := sanitizeCassette(i)
	assert.NoError(t, err)

	u, _ := url.Parse(i.Request.URL)
	q := u.Query()
	assert.Equal(t, "REDACTED", q.Get("AccessKeyId"))
	assert.Equal(t, "REDACTED", q.Get("SecurityToken"))
	assert.Equal(t, "REDACTED", q.Get("Signature"))
	assert.Equal(t, "REDACTED", q.Get("SignatureNonce"))
	assert.Equal(t, "CreateVpc", q.Get("Action"), "non-sensitive params should be preserved")
}

func TestUnitCommonSanitizeCassette_FormFields(t *testing.T) {
	i := &cassette.Interaction{
		Request: cassette.Request{
			Headers: http.Header{},
			Form: url.Values{
				"AccessKeyId":    {"LTAI123"},
				"SecurityToken":  {"STS_TOKEN"},
				"Signature":      {"abc"},
				"SignatureNonce": {"nonce"},
				"Action":         {"CreateVpc"},
				"VpcName":        {"my-vpc"},
			},
		},
	}
	err := sanitizeCassette(i)
	assert.NoError(t, err)
	assert.Equal(t, []string{"REDACTED"}, i.Request.Form["AccessKeyId"])
	assert.Equal(t, []string{"REDACTED"}, i.Request.Form["SecurityToken"])
	assert.Equal(t, []string{"REDACTED"}, i.Request.Form["Signature"])
	assert.Equal(t, []string{"REDACTED"}, i.Request.Form["SignatureNonce"])
	assert.Equal(t, []string{"CreateVpc"}, i.Request.Form["Action"])
	assert.Equal(t, []string{"my-vpc"}, i.Request.Form["VpcName"])
}

// --- redactURLParams ---

func TestUnitCommonRedactURLParams_Full(t *testing.T) {
	raw := "https://vpc.aliyuncs.com/?AccessKeyId=AK&Action=Test&SecurityToken=ST&Signature=SIG&SignatureNonce=SN"
	result := redactURLParams(raw)
	u, _ := url.Parse(result)
	q := u.Query()
	assert.Equal(t, "REDACTED", q.Get("AccessKeyId"))
	assert.Equal(t, "REDACTED", q.Get("SecurityToken"))
	assert.Equal(t, "Test", q.Get("Action"))
}

func TestUnitCommonRedactURLParams_NoSensitiveParams(t *testing.T) {
	raw := "https://vpc.aliyuncs.com/?Action=Test&RegionId=cn-hangzhou"
	result := redactURLParams(raw)
	u, _ := url.Parse(result)
	q := u.Query()
	assert.Equal(t, "Test", q.Get("Action"))
	assert.Equal(t, "cn-hangzhou", q.Get("RegionId"))
}

func TestUnitCommonRedactURLParams_InvalidURL(t *testing.T) {
	raw := "://invalid"
	assert.Equal(t, raw, redactURLParams(raw))
}

// --- filterSignatureParams ---

func TestUnitCommonFilterSignatureParams(t *testing.T) {
	q := url.Values{
		"Action":           {"CreateVpc"},
		"RegionId":         {"cn-hangzhou"},
		"AccessKeyId":      {"AK123"},
		"Signature":        {"sig"},
		"SignatureNonce":   {"nonce"},
		"Timestamp":        {"2026-01-01"},
		"SecurityToken":    {"token"},
		"SignatureMethod":  {"HMAC-SHA1"},
		"SignatureVersion": {"1.0"},
	}
	filtered := filterSignatureParams(q)
	assert.Equal(t, "CreateVpc", filtered.Get("Action"))
	assert.Equal(t, "cn-hangzhou", filtered.Get("RegionId"))
	assert.Empty(t, filtered.Get("AccessKeyId"))
	assert.Empty(t, filtered.Get("Signature"))
	assert.Empty(t, filtered.Get("Timestamp"))
}

// --- normalizeNumberedParams ---

func TestUnitCommonNormalizeNumberedParams(t *testing.T) {
	// Tag.1.Key=For, Tag.2.Key=Created should normalize the same as
	// Tag.1.Key=Created, Tag.2.Key=For (sorted by values)
	a := url.Values{
		"Action":      {"TagResources"},
		"Tag.1.Key":   {"For"},
		"Tag.1.Value": {"Test"},
		"Tag.2.Key":   {"Created"},
		"Tag.2.Value": {"TF"},
	}
	b := url.Values{
		"Action":      {"TagResources"},
		"Tag.1.Key":   {"Created"},
		"Tag.1.Value": {"TF"},
		"Tag.2.Key":   {"For"},
		"Tag.2.Value": {"Test"},
	}
	assert.Equal(t, normalizeNumberedParams(a).Encode(), normalizeNumberedParams(b).Encode(),
		"same tags in different order should normalize identically")

	// Non-numbered params should pass through unchanged
	c := url.Values{"Action": {"DescribeVpcs"}, "RegionId": {"cn-hangzhou"}}
	assert.Equal(t, c.Encode(), normalizeNumberedParams(c).Encode())
}

// --- helpers ---

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
