package fc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"sort"
	"strings"
)

const (
	// AuthQueryKeyExpires : keys in request queries
	AuthQueryKeyExpires = "x-fc-expires"

	// AuthQueryKeyAccessKeyID : keys in request queries
	AuthQueryKeyAccessKeyID = "x-fc-access-key-id"

	// AuthQueryKeySignature : keys in request queries
	AuthQueryKeySignature = "x-fc-signature"

	// AuthQueryKeySecurityToken : keys in requres queries
	AuthQueryKeySecurityToken = "x-fc-security-token"
)

type headers struct {
	Keys []string
	Vals []string
}

// GetAuthStr get signature strings
func GetAuthStr(accessKeyID string, accessKeySecret string, method string, header map[string]string, resource string) string {
	return "FC " + accessKeyID + ":" + GetSignature(accessKeySecret, method, header, resource)
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

// GetSignResourceWithQueries get signature resource with queries
func GetSignResourceWithQueries(path string, queries map[string][]string) string {
	return path + "\n" + getSignQueries(queries)
}

// getExpiresFromURLQueries returns expires and true if AuthQueryKeyExpires in queries
func getExpiresFromURLQueries(fcResource string) (string, bool) {
	originItems := strings.Split(fcResource, "\n")
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

// GetSignature calculate user's signature
func GetSignature(key string, method string, req map[string]string, fcResource string) string {
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
	if expires, ok := getExpiresFromURLQueries(fcResource); ok {
		date = expires
	}

	signStr := method + "\n" + lowerKeyHeaders[strings.ToLower(HTTPHeaderContentMD5)] + "\n" + lowerKeyHeaders[strings.ToLower(HTTPHeaderContentType)] + "\n" + date + "\n" + fcHeaders + fcResource

	h := hmac.New(func() hash.Hash { return sha256.New() }, []byte(key))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
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
