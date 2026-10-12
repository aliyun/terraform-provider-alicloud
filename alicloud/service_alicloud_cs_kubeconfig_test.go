package alicloud

import (
	"crypto/x509"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

var (
	csKubeconfigServerOnce sync.Once
	csKubeconfigServer     *httptest.Server
	csKubeconfigHandlerMu  sync.RWMutex
	csKubeconfigHandlerFn  http.HandlerFunc
)

// csKubeconfigFixture starts one shared local TLS test server per test process
// and replaces the process-wide system root pool with one that trusts the
// server certificate. The ROA CS SDK always dials HTTPS and verifies against
// the system root pool, and both that pool and SetFallbackRoots are
// once-per-process, hence the single shared server instead of one per test.
func csKubeconfigFixture(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	// SetFallbackRoots only replaces an already-loaded system pool when
	// GODEBUG=x509usefallbackroots=1 is set before the call.
	if g := os.Getenv("GODEBUG"); !strings.Contains(g, "x509usefallbackroots") {
		ng := "x509usefallbackroots=1"
		if g != "" {
			ng = g + "," + ng
		}
		t.Setenv("GODEBUG", ng)
	}
	csKubeconfigServerOnce.Do(func() {
		csKubeconfigServer = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			csKubeconfigHandlerMu.RLock()
			fn := csKubeconfigHandlerFn
			csKubeconfigHandlerMu.RUnlock()
			if fn != nil {
				fn(w, r)
				return
			}
			http.NotFound(w, r)
		}))
		pool := x509.NewCertPool()
		pool.AddCert(csKubeconfigServer.Certificate())
		x509.SetFallbackRoots(pool)
	})
	csKubeconfigHandlerMu.Lock()
	csKubeconfigHandlerFn = handler
	csKubeconfigHandlerMu.Unlock()
}

// csKubeconfigTestClient returns an AliyunClient whose "cs" endpoint points at
// the shared local TLS test server.
func csKubeconfigTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	csKubeconfigFixture(t, handler)
	endpoint := strings.TrimPrefix(csKubeconfigServer.URL, "https://")
	t.Setenv("NO_PROXY", endpoint)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("cs", endpoint)
	return client
}

func kubeConfigJSONResponse(t *testing.T, config string) []byte {
	t.Helper()
	body, err := json.Marshal(map[string]string{"config": config, "expiration": "2099-01-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func csKubeconfigHandler(t *testing.T, config string) http.HandlerFunc {
	t.Helper()
	body := kubeConfigJSONResponse(t, config)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/k8s/c1/user_config" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(body)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"test api error"}`))
	}
}

func csKubeconfigAPIErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"test api error"}`))
}

func TestUnitSetCertsKubeconfigSuccess(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigHandler(t, testKubeConfigYaml))
	d := schema.TestResourceDataRaw(t, resourceAlicloudCSKubernetes().Schema, map[string]interface{}{})
	d.SetId("c1")
	if err := setCerts(d, client, false); err != nil {
		t.Fatalf("setCerts returned error: %s", err)
	}
	want := map[string]string{
		"cluster_cert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==",
		"client_cert":  "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==",
		"client_key":   "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQo=",
	}
	got, ok := d.Get("certificate_authority").(map[string]interface{})
	if !ok {
		t.Fatalf("certificate_authority = %#v, want a map", d.Get("certificate_authority"))
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("certificate_authority[%s] = %q, want %q", k, got[k], v)
		}
	}
}

func TestUnitSetCertsKubeconfigSkipSetCertificateAuthority(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigHandler(t, testKubeConfigYaml))
	d := schema.TestResourceDataRaw(t, resourceAlicloudCSKubernetes().Schema, map[string]interface{}{})
	d.SetId("c1")
	if err := setCerts(d, client, true); err != nil {
		t.Fatalf("setCerts returned error: %s", err)
	}
	got, ok := d.Get("certificate_authority").(map[string]interface{})
	if !ok {
		t.Fatalf("certificate_authority = %#v, want a map", d.Get("certificate_authority"))
	}
	for _, k := range []string{"cluster_cert", "client_cert", "client_key"} {
		if v, exists := got[k]; !exists || v != "" {
			t.Fatalf("certificate_authority[%s] = %q, want an empty string", k, v)
		}
	}
}

func TestUnitSetCertsKubeconfigMalformed(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigHandler(t, "clusters: []\n"))
	d := schema.TestResourceDataRaw(t, resourceAlicloudCSKubernetes().Schema, map[string]interface{}{})
	d.SetId("c1")
	err := setCerts(d, client, false)
	if err == nil || !strings.Contains(err.Error(), "failed to parse kubeconfig") {
		t.Fatalf("setCerts error = %v, want it to contain 'failed to parse kubeconfig'", err)
	}
}

func TestUnitSetCertsKubeconfigAPIError(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigAPIErrorHandler)
	d := schema.TestResourceDataRaw(t, resourceAlicloudCSKubernetes().Schema, map[string]interface{}{})
	d.SetId("c1")
	err := setCerts(d, client, false)
	if err == nil || !strings.Contains(err.Error(), "failed to get kubeconfig") {
		t.Fatalf("setCerts error = %v, want it to contain 'failed to get kubeconfig'", err)
	}
}

func csClusterCredentialTestData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, dataSourceAlicloudCSClusterCredential().Schema, map[string]interface{}{
		"cluster_id": "c1",
	})
}

func TestUnitCSClusterCredentialReadSuccess(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigHandler(t, testKubeConfigYaml))
	d := csClusterCredentialTestData(t)
	cluster := &roacs.DescribeClusterDetailResponseBody{
		ClusterId: tea.String("c1"),
		Name:      tea.String("tf-test"),
	}
	if err := csClusterAuthDescriptionAttributes(d, client, cluster); err != nil {
		t.Fatalf("read returned error: %s", err)
	}
	if d.Id() == "" {
		t.Fatal("data source id is empty")
	}
	if got := d.Get("cluster_name"); got != "tf-test" {
		t.Fatalf("cluster_name = %q, want %q", got, "tf-test")
	}
	got, ok := d.Get("certificate_authority").(map[string]interface{})
	if !ok {
		t.Fatalf("certificate_authority = %#v, want a map", d.Get("certificate_authority"))
	}
	if got["cluster_cert"] != "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==" {
		t.Fatalf("certificate_authority.cluster_cert = %q", got["cluster_cert"])
	}
}

func TestUnitCSClusterCredentialReadMalformed(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigHandler(t, "clusters: []\n"))
	d := csClusterCredentialTestData(t)
	cluster := &roacs.DescribeClusterDetailResponseBody{ClusterId: tea.String("c1")}
	err := csClusterAuthDescriptionAttributes(d, client, cluster)
	if err == nil || !strings.Contains(err.Error(), "failed to parse kubeconfig") {
		t.Fatalf("read error = %v, want it to contain 'failed to parse kubeconfig'", err)
	}
}

func TestUnitCSClusterCredentialReadAPIError(t *testing.T) {
	client := csKubeconfigTestClient(t, csKubeconfigAPIErrorHandler)
	d := csClusterCredentialTestData(t)
	cluster := &roacs.DescribeClusterDetailResponseBody{ClusterId: tea.String("c1")}
	err := csClusterAuthDescriptionAttributes(d, client, cluster)
	if err == nil || !strings.Contains(err.Error(), "DescribeClusterKubeConfigWithExpiration") {
		t.Fatalf("read error = %v, want it to contain 'DescribeClusterKubeConfigWithExpiration'", err)
	}
}
