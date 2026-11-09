package connectivity

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// oidcStsRecorder records what a fake STS AssumeRoleWithOIDC endpoint received.
type oidcStsRecorder struct {
	mu        sync.Mutex
	hits      int32
	tokens    []string
	durations []string
}

func (rec *oidcStsRecorder) record(r *http.Request) (int32, error) {
	if err := r.ParseForm(); err != nil {
		return 0, err
	}
	rec.mu.Lock()
	defer rec.mu.Unlock()
	n := atomic.AddInt32(&rec.hits, 1)
	rec.tokens = append(rec.tokens, r.FormValue("OIDCToken"))
	rec.durations = append(rec.durations, r.FormValue("DurationSeconds"))
	return n, nil
}

func (rec *oidcStsRecorder) snapshot() (int32, []string, []string) {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	tokens := make([]string, len(rec.tokens))
	copy(tokens, rec.tokens)
	durations := make([]string, len(rec.durations))
	copy(durations, rec.durations)
	return atomic.LoadInt32(&rec.hits), tokens, durations
}

// newOIDCStsServer starts a TLS server emulating the STS AssumeRoleWithOIDC
// endpoint. Every issued session expires after one second, so the credentials
// provider must refresh on the next GetCredential.
func newOIDCStsServer(t *testing.T, rec *oidcStsRecorder) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := rec.record(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := fmt.Sprintf(`{
			"Credentials": {
				"AccessKeyId": "STS.OIDCAccessKey%d",
				"AccessKeySecret": "oidcAccessKeySecret%d",
				"SecurityToken": "oidcSecurityToken%d",
				"Expiration": %q
			}
		}`, n, n, n, time.Now().Add(time.Second).UTC().Format("2006-01-02T15:04:05Z"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return server
}

// trustOIDCStsServer makes the credentials-go HTTP client, which is built on
// http.DefaultTransport, trust the TLS test server.
func trustOIDCStsServer(t *testing.T, server *httptest.Server) {
	t.Helper()
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected default transport %T", http.DefaultTransport)
	}
	trusted := transport.Clone()
	pool := x509.NewCertPool()
	pool.AddCert(server.Certificate())
	trusted.TLSClientConfig = &tls.Config{RootCAs: pool}
	saved := http.DefaultTransport
	http.DefaultTransport = trusted
	t.Cleanup(func() { http.DefaultTransport = saved })
}

// TestConfigSetAuthCredentialByOidcTokenFileRefresh 验证 oidc_token_file 形态装配的是
// 可刷新的 oidc_role_arn 凭证：会话临期时重新读文件并重新 AssumeRoleWithOIDC，
// 且 session_expiration 真正传到了 STS 请求（而不是被默认 3600 覆盖）。
func TestConfigSetAuthCredentialByOidcTokenFileRefresh(t *testing.T) {
	tokenFile := filepath.Join(t.TempDir(), "oidc-token")
	if err := os.WriteFile(tokenFile, []byte("jwt-token-1"), 0644); err != nil {
		t.Fatal(err)
	}

	rec := &oidcStsRecorder{}
	server := newOIDCStsServer(t, rec)
	trustOIDCStsServer(t, server)

	config := &Config{
		RegionId:    "cn-hangzhou",
		StsEndpoint: strings.TrimPrefix(server.URL, "https://"),
		AssumeRoleWithOidc: &AssumeRoleWithOidc{
			RoleARN:         "acs:ram::123456789012:role/test-role",
			OIDCProviderArn: "acs:ram::123456789012:oidc-provider/test-idp",
			RoleSessionName: "terraform",
			OIDCTokenFile:   tokenFile,
			DurationSeconds: 1800,
		},
	}

	if err := config.setAuthCredentialByOidc(); err != nil {
		t.Fatalf("setAuthCredentialByOidc failed: %v", err)
	}

	first, err := config.Credential.GetCredential()
	if err != nil {
		t.Fatalf("first GetCredential failed: %v", err)
	}
	if got := *first.Type; got != "oidc_role_arn" {
		t.Fatalf("Expected credential type oidc_role_arn, got %q", got)
	}

	// 轮换 token 文件，1 秒的会话随即触发刷新。
	if err := os.WriteFile(tokenFile, []byte("jwt-token-2"), 0644); err != nil {
		t.Fatal(err)
	}
	second, err := config.Credential.GetCredential()
	if err != nil {
		t.Fatalf("second GetCredential failed: %v", err)
	}

	if got := *second.Type; got != "oidc_role_arn" {
		t.Fatalf("Expected credential type to stay oidc_role_arn, got %q", got)
	}
	if *first.AccessKeyId == *second.AccessKeyId {
		t.Fatalf("Expected the role session to be refreshed, got %q twice", *first.AccessKeyId)
	}

	hits, tokens, durations := rec.snapshot()
	// 三次请求：setAuthCredentialByOidc 内部首次换职 + 每次测试 GetCredential 触发一次临期重换。
	if hits != 3 {
		t.Fatalf("Expected the STS endpoint to be hit 3 times, got %d", hits)
	}
	// token 依次为：首次（jwt-token-1）、文件未轮转的临期重换（jwt-token-1）、轮转后的重换（jwt-token-2）。
	if len(tokens) != 3 || tokens[0] != "jwt-token-1" || tokens[1] != "jwt-token-1" || tokens[2] != "jwt-token-2" {
		t.Fatalf("Expected the token file to be re-read on refresh, got %v", tokens)
	}
	if len(durations) != 3 || durations[0] != "1800" || durations[1] != "1800" || durations[2] != "1800" {
		t.Fatalf("Expected session_expiration 1800 to reach all STS requests, got %v", durations)
	}
}

// inline oidc_token 形态（一次性 sts 快照）没有动态测试：该分支走 tea/dara 客户端，
// 它自建 http.Transport 且 RootCAs 为空（系统信任库），进程内无法注入对自签名测试证书的
// 信任（http.DefaultTransport 交换对它无效，darwin 也不认 SSL_CERT_FILE）。该分支的行为
// 未被本次修复改动，一次性语义由文档说明。
