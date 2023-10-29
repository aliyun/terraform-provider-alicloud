// This file is auto-generated, don't edit it. Thanks.
/**
 * Encode Util for Darabonba.
 */
package client

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	u "net/url"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/tjfoc/gmsm/sm3"
)

/**
 * Encode the URL
 * @param url string
 * @return encoded string
 */
func UrlEncode(url *string) (_result *string) {
	return tea.String(base64.URLEncoding.EncodeToString([]byte(tea.StringValue(url))))
}

/**
 * Special encoding for url params.
 * @param params string
 * @return encoded string
 */
func PercentEncode(raw *string) (_result *string) {
	uri := tea.StringValue(raw)
	uri = u.QueryEscape(uri)
	uri = strings.Replace(uri, "+", "%20", -1)
	uri = strings.Replace(uri, "*", "%2A", -1)
	uri = strings.Replace(uri, "%7E", "~", -1)
	return tea.String(uri)
}

/**
 * Encode the partial path of url.
 * @param path string
 * @return encoded string
 */
func PathEncode(path *string) (_result *string) {
	uri := tea.StringValue(path)
	strs := strings.Split(uri, "/")
	for i, v := range strs {
		strs[i] = u.QueryEscape(v)
	}
	uri = strings.Join(strs, "/")
	uri = strings.Replace(uri, "+", "%20", -1)
	uri = strings.Replace(uri, "*", "%2A", -1)
	uri = strings.Replace(uri, "%7E", "~", -1)
	return tea.String(uri)
}

/**
 * Hex encode for byte array.
 * @param raw byte array
 * @return encoded string
 */
func HexEncode(raw []byte) (_result *string) {
	return tea.String(hex.EncodeToString(raw))
}

/**
 * Hash the raw data with signatureAlgorithm.
 * @param raw hashing data
 * @param signatureAlgorithm the autograph method
 * @return hashed bytes
 */
func Hash(raw []byte, signatureAlgorithm *string) (_result []byte) {
	signType := tea.StringValue(signatureAlgorithm)
	if strings.Contains(signType, "HMAC-SHA256") || strings.Contains(signType, "RSA-SHA256") {
		h := sha256.New()
		h.Write(raw)
		return h.Sum(nil)
	} else if strings.Contains(signType, "HMAC-SM3") {
		h := sm3.New()
		h.Write(raw)
		return h.Sum(nil)
	}
	return nil
}

/**
 * Base64 encoder for byte array.
 * @param raw byte array
 * @return encoded string
 */
func Base64EncodeToString(raw []byte) (_result *string) {
	return tea.String(base64.StdEncoding.EncodeToString(raw))
}

/**
 * Base64 dncoder for string.
 * @param src string
 * @return dncoded byte array
 */
func Base64Decode(src *string) (_result []byte) {
	_result, err := base64.StdEncoding.DecodeString(tea.StringValue(src))
	if err != nil {
		return
	}
	return _result
}
