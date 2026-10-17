package paillmtrace

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	expectedPath   = "/api/v1/PAILLMTrace/eval/templates"
	expectedRegion = "cn-hangzhou"
)

// newMockClient mirrors the established HTTP-mock recipe (see the ESA lists
// unit test): test credentials, a localhost httptest server registered as the
// product endpoint, and NO_PROXY so the SDK dial reaches the mock directly.
func newMockClient(t *testing.T, server *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: expectedRegion, AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client: %v", err)
	}
	// Register the mock endpoint under the lowercased product code so
	// loadApiEndpoint returns it without consulting the location service.
	endpoints.Store("paillmtrace", endpoint)
	return client
}

// readEvaluationTemplates injects the mock client into a fresh data source and
// drives its read path. read builds the state model from the GetEvaluationTemplates
// response without touching Framework state plumbing, so it can be exercised
// directly here.
func readEvaluationTemplates(t *testing.T, client *connectivity.AliyunClient) (evaluationTemplatesModel, diag.Diagnostics) {
	t.Helper()
	ds := NewEvaluationTemplatesDataSource().(*evaluationTemplatesDataSource)

	var cfgResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, &cfgResp)
	if cfgResp.Diagnostics.HasError() {
		t.Fatalf("configure data source: %v", cfgResp.Diagnostics.Errors())
	}

	return ds.read(context.Background())
}

func TestUnitEvaluationTemplatesDataSource_ReadSuccessPreservesOrder(t *testing.T) {
	templates := []map[string]interface{}{
		{"TemplateId": "tpl-001", "Name": "alpha", "Metrics": []interface{}{"accuracy", "f1"}},
		{"TemplateId": "tpl-002", "Name": "beta", "Kind": "regression"},
		{"TemplateId": "tpl-003", "Name": "gamma", "Enabled": true},
	}

	var sawRequest bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawRequest = true
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"EvaluationTemplates": templates,
		})
	}))
	t.Cleanup(server.Close)

	client := newMockClient(t, server)
	state, diags := readEvaluationTemplates(t, client)

	if diags.HasError() {
		t.Fatalf("read returned diagnostics: %v", diags.Errors())
	}
	if !sawRequest {
		t.Fatal("the mock server never received a request")
	}
	if state.Id.ValueString() != expectedRegion {
		t.Errorf("id = %q, want %q", state.Id.ValueString(), expectedRegion)
	}
	if len(state.Templates) != len(templates) {
		t.Fatalf("templates count = %d, want %d", len(state.Templates), len(templates))
	}
	for i, item := range templates {
		want, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshal expected item %d: %v", i, err)
		}
		if got := state.Templates[i].ValueString(); got != string(want) {
			t.Errorf("templates[%d] = %s, want %s (order must be preserved)", i, got, string(want))
		}
	}
}

func TestUnitEvaluationTemplatesDataSource_ReadEmptySet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"EvaluationTemplates": []interface{}{},
		})
	}))
	t.Cleanup(server.Close)

	client := newMockClient(t, server)
	state, diags := readEvaluationTemplates(t, client)

	if diags.HasError() {
		t.Fatalf("empty-set read returned diagnostics: %v", diags.Errors())
	}
	if len(state.Templates) != 0 {
		t.Fatalf("templates count = %d, want 0 for empty response", len(state.Templates))
	}
	if state.Id.ValueString() != expectedRegion {
		t.Errorf("id = %q, want %q", state.Id.ValueString(), expectedRegion)
	}
}

func TestUnitEvaluationTemplatesDataSource_ReadFailurePropagatesError(t *testing.T) {
	// A 200 response carrying a non-success Code/Message is the fast,
	// non-retried failure path through formatError; it asserts the data
	// source surfaces API errors as diagnostics rather than swallowing them.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":      "InternalError",
			"Message":   "fixture failure",
			"RequestId": "req-fixture",
		})
	}))
	t.Cleanup(server.Close)

	client := newMockClient(t, server)
	state, diags := readEvaluationTemplates(t, client)

	if !diags.HasError() {
		t.Fatalf("read with an error response must surface diagnostics, got none")
	}
	if len(state.Templates) != 0 {
		t.Errorf("templates count = %d, want 0 on failure", len(state.Templates))
	}
}
