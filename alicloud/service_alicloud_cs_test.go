package alicloud

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	csSetCertsClusterAPIError = "c-api-err"
	csSetCertsClusterEmpty    = "c-empty"
	csSetCertsClusterOK       = "c-ok"
)

var (
	csSetCertsInitOnce sync.Once
	csSetCertsEndpoint string
)

// csSetCertsInit starts one shared TLS mock server and installs its certificate
// as the process root pool via x509.SetFallbackRoots. The CS SDK forces HTTPS,
// and the system root pool is already initialized (darwin platform pool) before
// tests run, so SSL_CERT_FILE has no effect inside the test process.
func csSetCertsInit(t *testing.T) string {
	t.Helper()
	csSetCertsInitOnce.Do(func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case strings.HasSuffix(r.URL.Path, "/k8s/"+csSetCertsClusterAPIError+"/user_config"):
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"Code":"InvalidCluster.NotFound","Message":"cluster not found"}`))
			case strings.HasSuffix(r.URL.Path, "/k8s/"+csSetCertsClusterEmpty+"/user_config"):
				_, _ = w.Write([]byte(`{"config":""}`))
			default:
				body, _ := json.Marshal(map[string]string{"config": csSetCertsValidKubeconfig})
				_, _ = w.Write(body)
			}
		})
		server := httptest.NewTLSServer(handler)

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}))
		t.Setenv("GODEBUG", "x509usefallbackroots=1")
		x509.SetFallbackRoots(pool)

		csSetCertsEndpoint = strings.TrimPrefix(server.URL, "https://")
	})
	return csSetCertsEndpoint
}

func csSetCertsTestClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	endpoint := csSetCertsInit(t)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoints.Store("cs", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "https",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

const csSetCertsValidKubeconfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: Y2E=
    server: https://127.0.0.1:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: admin
  name: kubernetes
current-context: kubernetes
kind: Config
users:
- name: admin
  user:
    client-certificate-data: Y2VydA==
    client-key-data: a2V5
`

func csSetCertsTestResourceData(t *testing.T, clusterId, kubeConfig, clientCert, clientKey, clusterCaCert string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"kube_config":           {Type: schema.TypeString, Optional: true},
		"client_cert":           {Type: schema.TypeString, Optional: true},
		"client_key":            {Type: schema.TypeString, Optional: true},
		"cluster_ca_cert":       {Type: schema.TypeString, Optional: true},
		"certificate_authority": {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}, map[string]interface{}{
		"kube_config":     kubeConfig,
		"client_cert":     clientCert,
		"client_key":      clientKey,
		"cluster_ca_cert": clusterCaCert,
	})
	d.SetId(clusterId)
	return d
}

func TestUnitSetCertsApiError(t *testing.T) {
	client := csSetCertsTestClient(t)
	kubeFile := filepath.Join(t.TempDir(), "kubeconfig.yaml")
	d := csSetCertsTestResourceData(t, csSetCertsClusterAPIError, kubeFile, "", "", "")

	err := setCerts(d, client, false)
	if err == nil {
		t.Fatalf("expected an error when the kubeconfig API fails, got nil")
	}
	if _, statErr := os.Stat(kubeFile); !os.IsNotExist(statErr) {
		t.Errorf("kubeconfig file must not be written on API error, stat err: %v", statErr)
	}
}

func TestUnitSetCertsEmptyConfig(t *testing.T) {
	client := csSetCertsTestClient(t)
	kubeFile := filepath.Join(t.TempDir(), "kubeconfig.yaml")
	d := csSetCertsTestResourceData(t, csSetCertsClusterEmpty, kubeFile, "", "", "")

	err := setCerts(d, client, false)
	if err == nil {
		t.Fatalf("expected an error for an empty kubeconfig response, got nil")
	}
	if _, statErr := os.Stat(kubeFile); !os.IsNotExist(statErr) {
		t.Errorf("kubeconfig file must not be written on empty config, stat err: %v", statErr)
	}
}

func TestUnitSetCertsWriteFailure(t *testing.T) {
	client := csSetCertsTestClient(t)
	kubeFile := filepath.Join(t.TempDir(), "no-such-dir", "kubeconfig.yaml")
	d := csSetCertsTestResourceData(t, csSetCertsClusterOK, kubeFile, "", "", "")

	err := setCerts(d, client, false)
	if err == nil {
		t.Fatalf("expected an error when writing kube_config fails, got nil")
	}
}

func TestUnitSetCertsSuccess(t *testing.T) {
	client := csSetCertsTestClient(t)
	dir := t.TempDir()
	kubeFile := filepath.Join(dir, "kubeconfig.yaml")
	clientCertFile := filepath.Join(dir, "client-cert.pem")
	clientKeyFile := filepath.Join(dir, "client-key.pem")
	caFile := filepath.Join(dir, "cluster-ca-cert.pem")
	d := csSetCertsTestResourceData(t, csSetCertsClusterOK, kubeFile, clientCertFile, clientKeyFile, caFile)

	if err := setCerts(d, client, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for path, want := range map[string]string{
		kubeFile:       csSetCertsValidKubeconfig,
		clientCertFile: "Y2VydA==",
		clientKeyFile:  "a2V5",
		caFile:         "Y2E=",
	} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if string(got) != want {
			t.Errorf("%s content = %q, want %q", path, got, want)
		}
	}

	authority := d.Get("certificate_authority").(map[string]interface{})
	for key, want := range map[string]string{
		"cluster_cert": "Y2E=",
		"client_cert":  "Y2VydA==",
		"client_key":   "a2V5",
	} {
		if got, ok := authority[key].(string); !ok || got != want {
			t.Errorf("certificate_authority[%s] = %#v, want %q", key, authority[key], want)
		}
	}
}

func TestUnitSetCertsSuccessSkipAuthority(t *testing.T) {
	client := csSetCertsTestClient(t)
	kubeFile := filepath.Join(t.TempDir(), "kubeconfig.yaml")
	d := csSetCertsTestResourceData(t, csSetCertsClusterOK, kubeFile, "", "", "")

	if err := setCerts(d, client, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	authority := d.Get("certificate_authority").(map[string]interface{})
	for _, key := range []string{"cluster_cert", "client_cert", "client_key"} {
		if got, ok := authority[key].(string); !ok || got != "" {
			t.Errorf("certificate_authority[%s] = %#v, want empty string", key, authority[key])
		}
	}
}
