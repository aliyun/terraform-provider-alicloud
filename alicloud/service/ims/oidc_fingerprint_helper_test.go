package ims_test

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"strings"
	"testing"
)

// imsOidcIssuerHost is the OIDC issuer used by the ims oidc provider
// acceptance tests. The IMS CreateOIDCProvider API validates the supplied
// Fingerprints against this issuer's HTTPS CA certificate (the API contract
// describes the parameter as "HTTPS CA证书的验证指纹"), so the tests must
// supply the real fingerprint of the issuer's TLS certificate rather than a
// synthetic value copied from the API documentation example.
const imsOidcIssuerHost = "oauth.aliyun.com:443"

// imsOidcProviderFingerprint returns the SHA-1 verification fingerprint of
// the HTTPS root CA certificate presented by the OIDC issuer
// (oauth.aliyun.com). The fingerprint is read from the live TLS certificate
// at test time — never hardcoded or guessed — so it remains correct across CA
// certificate rotation. Every certificate in the presented chain is logged so
// the binding cert can be diagnosed if the IMS API expects a different one.
func imsOidcProviderFingerprint(t *testing.T) string {
	t.Helper()
	conn, err := tls.Dial("tcp", imsOidcIssuerHost, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Fatalf("ims oidc: tls.Dial %s: %v", imsOidcIssuerHost, err)
	}
	defer conn.Close()

	chain := conn.ConnectionState().PeerCertificates
	if len(chain) == 0 {
		t.Fatalf("ims oidc: no peer certificates presented by %s", imsOidcIssuerHost)
	}
	for i, c := range chain {
		sum := sha1.Sum(c.Raw)
		t.Logf("ims oidc %s chain[%d] subject=%q issuer=%q sha1=%s isCA=%v", imsOidcIssuerHost, i, c.Subject.String(), c.Issuer.String(), strings.ToUpper(hex.EncodeToString(sum[:])), c.IsCA)
	}

	// The verification fingerprint is the trust anchor: the outermost CA
	// certificate in the presented chain (root CA).
	var root *x509.Certificate
	for i := len(chain) - 1; i >= 0; i-- {
		if chain[i].IsCA {
			root = chain[i]
			break
		}
	}
	if root == nil {
		root = chain[len(chain)-1]
	}
	sum := sha1.Sum(root.Raw)
	fp := strings.ToUpper(hex.EncodeToString(sum[:]))
	t.Logf("ims oidc: selected root CA verification fingerprint %s (subject=%q)", fp, root.Subject.String())
	return fp
}
