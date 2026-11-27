package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"

	roa "github.com/alibabacloud-go/tea-roa/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	otsCommon "github.com/aliyun/aliyun-tablestore-go-sdk/common"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// newLocalOtsClient points the four OTS request paths of a client at a local
// test server. The ots endpoint is stored under the control-plane product key
// and the table/tunnel SDKs read OtsEndpoint directly; the REST data plane
// deliberately does not get a "tablestore" endpoints entry so the OtsEndpoint
// fallback can be exercised.
func newLocalOtsClient(t *testing.T, endpoint string, region string, v4 bool) *AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("sts").
		SetAccessKeyId("test-ak").
		SetAccessKeySecret("test-secret").
		SetSecurityToken("test-token"))
	if err != nil {
		t.Fatalf("building the credential: %s", err)
	}
	endpoints := &sync.Map{}
	endpoints.Store("ots", endpoint)
	signVersion := &sync.Map{}
	if v4 {
		signVersion.Store("default", "v4")
	}
	return &AliyunClient{
		RegionId:                     region,
		tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		otsTunnelConnByInstanceName:  make(map[string]otsTunnel.TunnelClient),
		config: &config.Config{
			RegionId:    region,
			Protocol:    "HTTP",
			Credential:  credential,
			OtsEndpoint: "http://" + endpoint,
			Endpoints:   endpoints,
			SignVersion: signVersion,
		},
		teaSdkConfig:    rpc.Config{Protocol: tea.String("HTTP")},
		teaRoaSdkConfig: roa.Config{Protocol: tea.String("HTTP")},
	}
}

// verifyOtsSdkWireV4 recomputes the x-ots-signaturev4 header from what the
// HTTP server received: the string-to-sign is the uri, the POST method, an
// empty query, and the sorted non-empty x-ots request headers, signed with
// the daily key derived by the SDK. An algorithm name alone does not prove
// the request was actually signed.
func verifyOtsSdkWireV4(t *testing.T, r *http.Request, uri, region string) {
	t.Helper()
	if r.Header.Get("x-ots-signature") != "" {
		t.Error("V4 table/tunnel request must not carry the legacy x-ots-signature header")
	}
	if r.Header.Get("x-ots-signregion") != region {
		t.Errorf("x-ots-signregion = %q, want the physical signing region %q", r.Header.Get("x-ots-signregion"), region)
	}
	if r.Header.Get("x-ots-ststoken") != "test-token" {
		t.Error("temporary credential token missing from the wire")
	}
	signDate := r.Header.Get("x-ots-signdate")
	if len(signDate) != 8 {
		t.Errorf("x-ots-signdate = %q, want a YYYYMMDD signing date", signDate)
		return
	}
	stringToSign := uri + "\n" + "POST" + "\n" + "\n"
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, "x-ots-") && lower != "x-ots-signature" && lower != "x-ots-signaturev4" {
			names = append(names, lower)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		stringToSign += name + ":" + strings.TrimSpace(r.Header.Get(name)) + "\n"
	}
	key := otsCommon.FinalSigningKeyString("test-secret", signDate, region, "ots", "HmacSHA256")
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(stringToSign + otsCommon.V4_SIGNATURE_SALT))
	if want := base64.StdEncoding.EncodeToString(mac.Sum(nil)); want != r.Header.Get("x-ots-signaturev4") {
		t.Error("wire x-ots-signaturev4 does not match the signature recomputed from the request headers")
	}
}

func respondOtsJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"RequestId":"request-id"}`)
}

func TestUnitOtsV4SignVersionOnWire(t *testing.T) {
	clearProxyEnv(t)
	const region = "cn-hangzhou"
	const instance = "test-instance"

	t.Run("control plane/v4", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("InstanceName"); got != instance {
				t.Errorf("InstanceName query = %q, want %q", got, instance)
			}
			verifyWireV4(t, r, "ots", region)
			respondOtsJson(w)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, true)
		response, err := c.OtsPopV4("GET", "ListVpcInfoByInstance", map[string]interface{}{
			"InstanceName": instance,
		})
		if err != nil || response["RequestId"] != "request-id" {
			t.Fatalf("unexpected API response: %v, %v", response, err)
		}
	})

	t.Run("control plane/legacy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("SignatureMethod") != "HMAC-SHA1" || q.Get("SignatureVersion") != "1.0" || q.Get("Signature") == "" {
				t.Error("expected legacy RPC signature when the default is not v4")
			}
			respondOtsJson(w)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, false)
		response, err := c.OtsPopV4("GET", "ListVpcInfoByInstance", map[string]interface{}{
			"InstanceName": instance,
		})
		if err != nil || response["RequestId"] != "request-id" {
			t.Fatalf("unexpected API response: %v, %v", response, err)
		}
	})

	t.Run("REST/v4", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("InstanceName"); got != instance {
				t.Errorf("InstanceName query = %q, want %q", got, instance)
			}
			verifyWireV4(t, r, "tablestore", region)
			respondOtsJson(w)
		}))
		defer server.Close()
		// Only OtsEndpoint is configured: the REST data plane must fall back to
		// it for EO stations where no public cloud endpoint rule applies.
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, true)
		response, err := c.OtsRestV4("tablestore", "2020-12-09", "/v2/openapi/getinstance", "GET",
			map[string]*string{"InstanceName": tea.String(instance)}, nil, nil)
		if err != nil || response["RequestId"] != "request-id" {
			t.Fatalf("unexpected API response: %v, %v", response, err)
		}
	})

	t.Run("REST/legacy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-acs-signature-method") != "HMAC-SHA1" || !strings.HasPrefix(r.Header.Get("Authorization"), "acs ") {
				t.Error("expected legacy ROA signature when the default is not v4")
			}
			respondOtsJson(w)
		}))
		defer server.Close()
		endpoint := strings.TrimPrefix(server.URL, "http://")
		c := newLocalOtsClient(t, endpoint, region, false)
		c.config.Endpoints.Store("tablestore", endpoint)
		conn, err := c.NewOtsRoaClient("tablestore")
		if err != nil {
			t.Fatal(err)
		}
		response, err := conn.DoRequest(tea.String("2020-12-09"), nil, tea.String("GET"), tea.String("AK"),
			tea.String("/v2/openapi/getinstance"), map[string]*string{"InstanceName": tea.String(instance)}, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if body, ok := response["body"].(map[string]interface{}); ok {
			response = body
		}
		if response["RequestId"] != "request-id" {
			t.Fatalf("unexpected API response: %v", response)
		}
	})

	t.Run("REST/public endpoint rule", func(t *testing.T) {
		c := newLocalOtsClient(t, "unused.example", region, false)
		c.config.OtsEndpoint = ""
		// No endpoints entry for "tablestore": NewOtsRoaClient must resolve the
		// public cloud rule rather than an EO endpoint.
		if _, err := c.NewOtsRoaClient("tablestore"); err != nil {
			t.Fatal(err)
		}
		if resolved, ok := c.config.Endpoints.Load("tablestore"); !ok || resolved.(string) != "tablestore."+region+".aliyuncs.com" {
			t.Errorf("resolved tablestore endpoint = %v, want the public cloud rule tablestore.%s.aliyuncs.com", resolved, region)
		}
	})

	t.Run("table/v4", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ListTable" {
				t.Errorf("table request path = %q, want /ListTable", r.URL.Path)
			}
			if got := r.Header.Get("x-ots-instancename"); got != instance {
				t.Errorf("x-ots-instancename = %q, want %q", got, instance)
			}
			if got := r.Header.Get("x-ots-sourceip"); got != "10.0.0.1" {
				t.Errorf("x-ots-sourceip = %q, want the preserved provider SourceIp", got)
			}
			verifyOtsSdkWireV4(t, r, "/ListTable", region)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, true)
		c.config.SourceIp = "10.0.0.1"
		_, err := c.WithTableStoreClient(instance, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			_, listErr := tableStoreClient.ListTable()
			return nil, listErr
		})
		if err != nil {
			t.Fatalf("list tables through the V4 table client: %v", err)
		}
	})

	t.Run("table/legacy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-ots-signature") == "" || r.Header.Get("x-ots-signaturev4") != "" {
				t.Error("expected the legacy x-ots-signature header when the default is not v4")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, false)
		_, err := c.WithTableStoreClient(instance, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			_, listErr := tableStoreClient.ListTable()
			return nil, listErr
		})
		if err != nil {
			t.Fatalf("list tables through the legacy table client: %v", err)
		}
	})

	t.Run("table/v4 with ots override", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-ots-signature") == "" || r.Header.Get("x-ots-signaturev4") != "" {
				t.Error("an explicit ots = v1 override must restore the legacy table signature")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, true)
		c.config.SignVersion.Store("ots", "v1")
		if c.UseOtsV4() {
			t.Fatal("the ots override did not take precedence over the v4 default")
		}
		_, err := c.WithTableStoreClient(instance, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			_, listErr := tableStoreClient.ListTable()
			return nil, listErr
		})
		if err != nil {
			t.Fatalf("list tables through the overridden table client: %v", err)
		}
	})

	t.Run("tunnel/v4", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/tunnel/list" {
				t.Errorf("tunnel request path = %q, want /tunnel/list", r.URL.Path)
			}
			if got := r.Header.Get("x-ots-instancename"); got != instance {
				t.Errorf("x-ots-instancename = %q, want %q", got, instance)
			}
			verifyOtsSdkWireV4(t, r, "/tunnel/list", region)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, true)
		_, err := c.WithTableStoreTunnelClient(instance, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			_, listErr := tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{})
			return nil, listErr
		})
		if err != nil {
			t.Fatalf("list tunnels through the V4 tunnel client: %v", err)
		}
	})

	t.Run("tunnel/legacy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-ots-signature") == "" || r.Header.Get("x-ots-signaturev4") != "" {
				t.Error("expected the legacy x-ots-signature header when the default is not v4")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		c := newLocalOtsClient(t, strings.TrimPrefix(server.URL, "http://"), region, false)
		_, err := c.WithTableStoreTunnelClient(instance, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			_, listErr := tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{})
			return nil, listErr
		})
		if err != nil {
			t.Fatalf("list tunnels through the legacy tunnel client: %v", err)
		}
	})
}

func TestUnitOtsV4Preconditions(t *testing.T) {
	clearProxyEnv(t)
	const region = "cn-hangzhou"

	t.Run("missing signing region", func(t *testing.T) {
		c := newLocalOtsClient(t, "unused.example", "", true)
		if _, err := c.WithTableStoreClient("instance", func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			_, listErr := tableStoreClient.ListTable()
			return nil, listErr
		}); err == nil || !strings.Contains(err.Error(), "signing region") {
			t.Fatalf("expected the missing signing region diagnostic, got %v", err)
		}
		if _, err := c.WithTableStoreTunnelClient("instance", func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			_, listErr := tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{})
			return nil, listErr
		}); err == nil || !strings.Contains(err.Error(), "signing region") {
			t.Fatalf("expected the missing signing region diagnostic, got %v", err)
		}
	})

	t.Run("tunnel workers report the V4 gap explicitly", func(t *testing.T) {
		c := newLocalOtsClient(t, "unused.example", region, true)
		_, err := c.WithTableStoreTunnelClient("instance", func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			_, workerErr := tunnelClient.NewTunnelWorker("tunnel-id", &otsTunnel.TunnelWorkerConfig{})
			return nil, workerErr
		})
		if err == nil || !strings.Contains(err.Error(), "tunnel workers are not supported") {
			t.Fatalf("expected the explicit tunnel worker diagnostic, got %v", err)
		}
	})
}
