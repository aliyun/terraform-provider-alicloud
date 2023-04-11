// This file is auto-generated, don't edit it. Thanks.
package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/pkg/errors"
)

const (
	// HTTPHeaderContentMD5 key in request headers
	HTTPHeaderContentMD5 = "Content-MD5"
	// HTTPHeaderPrefix key prefix in request headers
	HTTPHeaderPrefix = "x-fc-"
	// HTTPHeaderDate key in request headers
	HTTPHeaderDate = "Date"
	// AuthQueryKeyExpires key in request headers
	AuthQueryKeyExpires = "x-fc-expires"
	// HTTPHeaderContentType key in request headers
	HTTPHeaderContentType = "Content-Type"
	// HTTPHeaderAuthorization key in request headers
	HTTPHeaderAuthorization = "Authorization"
	// HTTPHeaderSecurityToken key in request headers
	HTTPHeaderSecurityToken = "x-fc-security-token"
)

type Client struct {
	Credential credential.Credential
}

func NewClient(cred credential.Credential) (*Client, error) {
	client := new(Client)
	err := client.Init(cred)
	return client, err
}

func (client *Client) Init(cred credential.Credential) (_err error) {
	client.Credential = cred
	return nil
}

func TeeReader(r io.Reader, w io.Writer) io.Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r io.Reader
	w io.Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

// InvokeHTTPTrigger invoke a http trigger
func (client *Client) InvokeHTTPTrigger(url *string, method *string, body []byte, headers *http.Header) (_result *http.Response, _err error) {
	req, err := client.BuildHTTPRequest(url, method, body, headers)
	if err != nil {
		return nil, err
	}
	return client.SendHTTPRequestWithAuthorization(req)
}

// InvokeAnonymousHTTPTrigger invoke an anonymous http trigger
func (client *Client) InvokeAnonymousHTTPTrigger(url *string, method *string, body []byte, headers *http.Header) (_result *http.Response, _err error) {
	req, err := client.BuildHTTPRequest(url, method, body, headers)
	if err != nil {
		return nil, err
	}
	return client.SendHTTPRequest(req)
}

func (client *Client) SendHTTPRequestWithAuthorization(req *http.Request) (_result *http.Response, _err error) {
	signedRequest, err := client.SignRequest(req)
	if err != nil {
		return nil, err
	}
	return client.SendHTTPRequest(signedRequest)
}

func (client *Client) SendHTTPRequest(req *http.Request) (_result *http.Response, _err error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *Client) SignRequest(req *http.Request) (_result *http.Request, _err error) {
	var (
		contentMD5 = ""
		err        error
	)
	// CONTENT-MD5
	if req.Body != nil {
		bodyNew := bytes.NewBuffer([]byte{})
		tee := TeeReader(req.Body, bodyNew)
		contentMD5, err = md5Digest(tee)
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(bodyNew)
	}
	return client.SignRequestWithContentMD5(req, tea.String(contentMD5))
}

func (client *Client) SignRequestWithContentMD5(req *http.Request, contentMD5 *string) (_result *http.Request, _err error) {
	headerParams := make(map[string]string)
	req.Header.Set(HTTPHeaderDate, time.Now().UTC().Format(http.TimeFormat))
	if len(tea.StringValue(contentMD5)) != 0 {
		req.Header.Set(HTTPHeaderContentMD5, tea.StringValue(contentMD5))
	}
	if req.Header != nil {
		for k, _ := range req.Header {
			headerParams[k] = req.Header.Get(k)
		}
	}
	// Canonicalized
	pathWithQuery := req.URL.Path
	params := req.URL.Query()
	pathWithQuery = getSignResourceWithQueries(req.URL.Path, params)
	// Build Authorization header
	accessKeyId, _ := client.Credential.GetAccessKeyId()
	accessKeySecret, _ := client.Credential.GetAccessKeySecret()
	securityToken, _ := client.Credential.GetSecurityToken()
	if securityToken != nil && len(*securityToken) != 0 {
		req.Header.Set(HTTPHeaderSecurityToken, tea.StringValue(securityToken))
	}
	authStr := getAuthString(tea.StringValue(accessKeyId), tea.StringValue(accessKeySecret), req.Method, headerParams, pathWithQuery)
	req.Header.Set(HTTPHeaderAuthorization, authStr)
	return req, nil
}

func (client *Client) BuildHTTPRequest(url *string, method *string, body []byte, headers *http.Header) (_result *http.Request, _err error) {
	res, err := http.NewRequest(tea.StringValue(method), tea.StringValue(url), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if headers != nil && len(*headers) != 0 {
		res.Header = *headers
	}
	return res, nil
}

func md5Digest(reader io.Reader) (string, error) {
	ctx := md5.New()
	buffer := make([]byte, 8*1024)
	for {
		if c, err := reader.Read(buffer); err == nil {
			ctx.Write(buffer[:c])
		} else if err == io.EOF {
			break
		} else {
			return "", errors.Wrapf(err, "failed to read from reader")
		}
	}
	return hex.EncodeToString(ctx.Sum(nil)), nil
}

func getAuthString(accessKeyID string, accessKeySecret string, method string, header map[string]string, pathWithQuerys string) string {
	return "FC " + accessKeyID + ":" + getSignature(accessKeySecret, method, header, pathWithQuerys)
}

type headers struct {
	Keys []string
	Vals []string
}

func (h *headers) Len() int {
	return len(h.Vals)
}

func (h *headers) Less(i, j int) bool {
	return bytes.Compare([]byte(h.Keys[i]), []byte(h.Keys[j])) < 0
}

func (h *headers) Swap(i, j int) {
	h.Vals[i], h.Vals[j] = h.Vals[j], h.Vals[i]
	h.Keys[i], h.Keys[j] = h.Keys[j], h.Keys[i]
}

// GetSignature calculate user's signature
func getSignature(key string, method string, req map[string]string, path string) string {
	header := &headers{}
	lowerKeyHeaders := map[string]string{}
	for k, v := range req {
		lowerKey := strings.ToLower(k)
		if strings.HasPrefix(lowerKey, HTTPHeaderPrefix) {
			header.Keys = append(header.Keys, lowerKey)
			header.Vals = append(header.Vals, v)
		}
		lowerKeyHeaders[lowerKey] = v
	}
	sort.Sort(header)

	fcHeaders := ""
	for i := range header.Keys {
		fcHeaders += header.Keys[i] + ":" + header.Vals[i] + "\n"
	}

	date := req[HTTPHeaderDate]
	if expires, ok := getExpiresFromURLQueries(path); ok {
		date = expires
	}

	signStr := method + "\n" + lowerKeyHeaders[strings.ToLower(HTTPHeaderContentMD5)] + "\n" + lowerKeyHeaders[strings.ToLower(HTTPHeaderContentType)] + "\n" + date + "\n" + fcHeaders + path

	h := hmac.New(sha256.New, []byte(key))
	_, _ = io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signedStr
}

func getExpiresFromURLQueries(path string) (string, bool) {
	originItems := strings.Split(path, "\n")
	queriesUnescaped := map[string]string{}

	for _, item := range originItems[1:] {
		kvPair := strings.Split(item, "=")
		if len(kvPair) > 1 {
			queriesUnescaped[kvPair[0]] = kvPair[1]
		}
	}

	expires, ok := queriesUnescaped[AuthQueryKeyExpires]
	return expires, ok
}

func getSignResourceWithQueries(path string, queries map[string][]string) string {
	if len(path) == 0 {
		path = "/"
	}
	return path + "\n" + getSignQueries(queries)
}

func getSignQueries(queries map[string][]string) string {
	paramsList := []string{}
	for key, values := range queries {
		if len(values) == 0 {
			paramsList = append(paramsList, key)
			continue
		}
		for _, v := range values {
			paramsList = append(paramsList, fmt.Sprintf("%s=%s", key, v))
		}
	}
	sort.Strings(paramsList)
	return strings.Join(paramsList, "\n")
}
