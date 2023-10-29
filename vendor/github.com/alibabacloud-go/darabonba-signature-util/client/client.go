// This file is auto-generated, don't edit it. Thanks.
/**
 * Signature Util for Darabonba.
 */
package client

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/tjfoc/gmsm/sm3"
)

const (
	PEM_BEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END   = "\n-----END RSA PRIVATE KEY-----"
)

/**
 * HmacSHA1 Signature
 * @param stringToSign string
 * @param secret string
 * @return signed bytes
 */
func HmacSHA1Sign(stringToSign *string, secret *string) (_result []byte) {
	return HmacSHA1SignByBytes(stringToSign, []byte(tea.StringValue(secret)))
}

/**
 * HmacSHA1 Signature
 * @param stringToSign string
 * @param secret bytes
 * @return signed bytes
 */
func HmacSHA1SignByBytes(stringToSign *string, secret []byte) (_result []byte) {
	hmac := hmac.New(sha1.New, secret)
	hmac.Write([]byte(tea.StringValue(stringToSign)))
	return hmac.Sum(nil)
}

/**
 * HmacSHA256 Signature
 * @param stringToSign string
 * @param secret string
 * @return signed bytes
 */
func HmacSHA256Sign(stringToSign *string, secret *string) (_result []byte) {
	return HmacSHA256SignByBytes(stringToSign, []byte(tea.StringValue(secret)))
}

/**
 * HmacSHA256 Signature
 * @param stringToSign string
 * @param secret bytes
 * @return signed bytes
 */
func HmacSHA256SignByBytes(stringToSign *string, secret []byte) (_result []byte) {
	h := hmac.New(sha256.New, secret)
	_, err := h.Write([]byte(tea.StringValue(stringToSign)))
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}

/**
 * HmacSM3 Signature
 * @param stringToSign string
 * @param secret string
 * @return signed bytes
 */
func HmacSM3Sign(stringToSign *string, secret *string) (_result []byte) {
	return HmacSM3SignByBytes(stringToSign, []byte(tea.StringValue(secret)))
}

/**
 * HmacSM3 Signature
 * @param stringToSign string
 * @param secret bytes
 * @return signed bytes
 */
func HmacSM3SignByBytes(stringToSign *string, secret []byte) (_result []byte) {
	h := hmac.New(sm3.New, secret)
	_, err := h.Write([]byte(tea.StringValue(stringToSign)))
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}

/**
 * SHA256withRSA Signature
 * @param stringToSign string
 * @param secret string
 * @return signed bytes
 */
func SHA256withRSASign(stringToSign *string, secret *string) (_result []byte) {
	byt := rsaSign(tea.StringValue(stringToSign), tea.StringValue(secret))
	return byt
}

func rsaSign(content, secret string) []byte {
	h := crypto.SHA256.New()
	h.Write([]byte(content))
	hashed := h.Sum(nil)
	priv, err := parsePrivateKey(secret)
	if err != nil {
		return nil
	}
	sign, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		return nil
	}
	return sign
}

func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = formatPrivateKey(privateKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("PrivateKey is invalid")
	}
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch priKey.(type) {
	case *rsa.PrivateKey:
		return priKey.(*rsa.PrivateKey), nil
	default:
		return nil, nil
	}
}

func formatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PEM_BEGIN) {
		privateKey = PEM_BEGIN + privateKey
	}

	if !strings.HasSuffix(privateKey, PEM_END) {
		privateKey += PEM_END
	}
	return privateKey
}

func MD5Sign(stringToSign *string) (_result []byte) {
	return MD5SignForBytes([]byte(tea.StringValue(stringToSign)))
}

func MD5SignForBytes(bytesToSign []byte) (_result []byte) {
	h := md5.New()
	h.Write(bytesToSign)
	return h.Sum(nil)
}
