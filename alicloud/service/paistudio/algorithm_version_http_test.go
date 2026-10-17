package paistudio

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// These TestUnits exercise the real production CRUD wire against an in-process
// HTTP mock. They build an offline connectivity.AliyunClient (fake access key
// credentials, AccountType preset so BssOpenApi is skipped, the PaiStudio
// endpoint pre-stored in the Endpoints map so the location service is never
// called) and drive it through AliyunClient.RoaPost/RoaGet/RoaPut/RoaDelete.
// Those helpers route to roaRequest -> roa.NewClient -> (*roa.Client).DoRequest,
// the tea-roa ROA client every ROA product in this provider uses; that path
// never reaches (*openapi.Client).Execute (the nil Spi.ModifyConfiguration site
// at darabonba-openapi/v2 client.go:1271). So these tests prove the production
// request/response wire and the Framework state semantics end-to-end, against
// 2xx / 404-NoSuchObject / 5xx / malformed responses — not just the helpers.

const (
	mockAlgorithmID      = "algo-mock"
	mockAlgorithmVersion = "v1.0.0"
)

// newOfflinePaiStudioClient builds a fully offline AliyunClient whose ROA
// requests target `host` (host:port, no scheme) over HTTP. The endpoint is
// pre-stored in the Endpoints map so loadApiEndpoint never reads the env or
// calls the location service; AccountType is preset so getAccountType (the
// BssOpenApi QueryAvailableInstances cloud call) is skipped. No real
// credential or network is used; the credentials are fake and the only socket
// is the in-process httptest listener.
func newOfflinePaiStudioClient(t *testing.T, host string) *connectivity.AliyunClient {
	t.Helper()
	cred, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("mock-access-key-id").
		SetAccessKeySecret("mock-access-key-secret"))
	if err != nil {
		t.Fatalf("building fake credential: %s", err)
	}
	endpoints := &sync.Map{}
	// Pre-store the PaiStudio endpoint so loadApiEndpoint resolves from the map
	// and never reaches loadEndpoint/describeEndpointForService (location
	// service) or any ALIBABA_CLOUD_ENDPOINT_* env var.
	endpoints.Store("paistudio", host)
	conf := connectivity.Config{
		Region:      connectivity.Beijing,
		RegionId:    string(connectivity.Beijing),
		AccessKey:   "mock-access-key-id",
		SecretKey:   "mock-access-key-secret",
		Protocol:    "HTTP",
		AccountType: "Domestic",
		Credential:  cred,
		Endpoints:   endpoints,
		SignVersion: &sync.Map{},
	}
	client, err := conf.Client()
	if err != nil {
		t.Fatalf("building offline PaiStudio client: %s", err)
	}
	return client
}

// paiStudioHost trims the scheme from an httptest server URL so the roa client
// (which prepends the Protocol scheme) targets it as a bare host:port.
func paiStudioHost(serverURL string) string {
	return strings.TrimPrefix(strings.TrimPrefix(serverURL, "http://"), "https://")
}

// mockVersion is the in-process, stateful PaiStudio algorithm-version store.
type mockVersion struct {
	mu         sync.Mutex
	versions   map[string]map[string]interface{}
	lastMethod string
	lastPath   string
	// methods records every request method in arrival order. The Framework
	// Create/Update methods perform a read-back (getAlgorithmVersion -> RoaGet)
	// after the write, so lastMethod is always the trailing GET by the time the
	// RPC returns; hasMethod inspects the full history to prove the write
	// (POST/PUT) was actually issued.
	methods []string
}

// hasMethod reports whether method was issued at least once against the mock.
func (m *mockVersion) hasMethod(method string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, got := range m.methods {
		if got == method {
			return true
		}
	}
	return false
}

func newMockVersion() *mockVersion {
	return &mockVersion{versions: map[string]map[string]interface{}{}}
}

// parseVersionPath splits /api/v1/algorithms/<id>/versions/<version> into the
// id and version used to key the store and to echo back in responses.
func parseVersionPath(p string) (id, version string, ok bool) {
	const prefix = "/api/v1/algorithms/"
	if !strings.HasPrefix(p, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(p, prefix)
	parts := strings.SplitN(rest, "/versions/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func (m *mockVersion) key(id, version string) string { return id + ":" + version }

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Fatalf("encoding mock response: %s", err)
		}
	}
}

// handler returns an http.HandlerFunc implementing the PaiStudio
// algorithm-version CRUD contract over the in-process store.
func (m *mockVersion) handler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		id, version, ok := parseVersionPath(r.URL.Path)
		if !ok {
			writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
				"Code": "NoSuchObject", "Message": "unknown path", "RequestId": "r-path",
			})
			return
		}
		m.mu.Lock()
		m.lastMethod = r.Method
		m.lastPath = r.URL.Path
		m.methods = append(m.methods, r.Method)
		m.mu.Unlock()
		key := m.key(id, version)
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var req map[string]interface{}
			if r.Body != nil {
				if data, err := io.ReadAll(r.Body); err == nil && len(data) > 0 {
					_ = json.Unmarshal(data, &req)
				}
			}
			spec := req["AlgorithmSpec"]
			obj, exists := m.versions[key]
			if !exists {
				obj = map[string]interface{}{
					"AlgorithmId":       id,
					"AlgorithmVersion":  version,
					"AlgorithmName":     "echo",
					"AlgorithmProvider": "self",
					"RegionId":          "cn-beijing",
					"CreateTime":        "2026-01-01T00:00:00Z",
				}
				m.versions[key] = obj
			}
			if spec != nil {
				obj["AlgorithmSpec"] = spec
			}
			writeJSON(t, w, http.StatusOK, obj)
		case http.MethodGet:
			obj, exists := m.versions[key]
			if !exists {
				writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
					"Code": "NoSuchObject", "Message": "algorithm version not found", "RequestId": "r-get",
				})
				return
			}
			writeJSON(t, w, http.StatusOK, obj)
		case http.MethodDelete:
			if _, exists := m.versions[key]; !exists {
				writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
					"Code": "NoSuchObject", "Message": "algorithm version not found", "RequestId": "r-del",
				})
				return
			}
			delete(m.versions, key)
			writeJSON(t, w, http.StatusOK, map[string]interface{}{})
		default:
			writeJSON(t, w, http.StatusMethodNotAllowed, map[string]interface{}{
				"Code": "MethodNotAllowed", "Message": r.Method, "RequestId": "r-method",
			})
		}
	}
}

// TestUnitPaiAlgorithmVersionCRUDHTTPWire drives the production
// RoaPost -> RoaGet -> RoaPut -> RoaDelete wire against a stateful httptest
// server and asserts the full create -> read -> update -> read -> delete ->
// not-found cycle, proving the production path neither panics nor mangles
// state and that the Framework state mapping reflects the API responses.
func TestUnitPaiAlgorithmVersionCRUDHTTPWire(t *testing.T) {
	store := newMockVersion()
	server := httptest.NewServer(store.handler(t))
	defer server.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server.URL))

	path := algorithmVersionPath(mockAlgorithmID, mockAlgorithmVersion)

	// Create: RoaPost carries the AlgorithmSpec on the production write path.
	spec, err := parseAlgorithmSpec(`{"JobType":"TfJob","Command":"python main.py","Image":"reg/echo:1.0"}`)
	if err != nil {
		t.Fatalf("parseAlgorithmSpec: %s", err)
	}
	createBody := map[string]interface{}{"AlgorithmSpec": spec}
	if _, err := client.RoaPost(PaiStudioProduct, PaiStudioVersion, path, nil, nil, createBody, true); err != nil {
		t.Fatalf("RoaPost create: %s", err)
	}
	if got := store.lastMethod; got != http.MethodPost {
		t.Fatalf("create method: got %q, want POST", got)
	}
	if got := store.lastPath; got != path {
		t.Fatalf("create path: got %q, want %q", got, path)
	}

	// Read after create: getAlgorithmVersion (RoaGet) returns the stored object,
	// the response identity matches, and the Framework state maps it.
	object, err := getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err != nil {
		t.Fatalf("getAlgorithmVersion after create: %s", err)
	}
	if object == nil {
		t.Fatal("getAlgorithmVersion after create returned nil for a just-created version")
	}
	if err := validateAlgorithmVersionResponse(object, mockAlgorithmID, mockAlgorithmVersion); err != nil {
		t.Fatalf("validateAlgorithmVersionResponse after create: %s", err)
	}
	state := algorithmVersionState(object)
	if state.Id.ValueString() != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("state id: got %q", state.Id.ValueString())
	}
	if state.AlgorithmSpec.IsNull() {
		t.Fatal("state AlgorithmSpec must be populated from the create response")
	}

	// Update: RoaPut changes the spec on the production write path.
	updatedSpec, err := parseAlgorithmSpec(`{"JobType":"TfJob","Command":"python train.py","Image":"reg/echo:2.0"}`)
	if err != nil {
		t.Fatalf("parseAlgorithmSpec update: %s", err)
	}
	if _, err := client.RoaPut(PaiStudioProduct, PaiStudioVersion, path, nil, nil, map[string]interface{}{"AlgorithmSpec": updatedSpec}, true); err != nil {
		t.Fatalf("RoaPut update: %s", err)
	}
	if got := store.lastMethod; got != http.MethodPut {
		t.Fatalf("update method: got %q, want PUT", got)
	}
	object, err = getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err != nil {
		t.Fatalf("getAlgorithmVersion after update: %s", err)
	}
	if err := validateAlgorithmVersionResponse(object, mockAlgorithmID, mockAlgorithmVersion); err != nil {
		t.Fatalf("validateAlgorithmVersionResponse after update: %s", err)
	}
	state = algorithmVersionState(object)
	var specMap map[string]interface{}
	if err := json.Unmarshal([]byte(state.AlgorithmSpec.ValueString()), &specMap); err != nil {
		t.Fatalf("updated AlgorithmSpec is not valid JSON: %s", err)
	}
	if specMap["Command"] != "python train.py" {
		t.Fatalf("update did not propagate to the read-back spec: %#v", specMap)
	}

	// Delete: RoaDelete on the production write path; a second read is now not-found.
	if _, err := client.RoaDelete(PaiStudioProduct, PaiStudioVersion, path, nil, nil, nil, true); err != nil {
		t.Fatalf("RoaDelete: %s", err)
	}
	if got := store.lastMethod; got != http.MethodDelete {
		t.Fatalf("delete method: got %q, want DELETE", got)
	}
	object, err = getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err != nil {
		t.Fatalf("getAlgorithmVersion after delete: %s", err)
	}
	if object != nil {
		t.Fatalf("getAlgorithmVersion after delete must return nil (NoSuchObject), got %#v", object)
	}
}

// TestUnitPaiAlgorithmVersionReadNotFoundHTTPWire asserts that a 404 carrying
// the product NoSuchObject code on the production RoaGet path is parsed into a
// *tea.SDKError and recognized as not-found (so Read removes the resource and
// Delete is idempotent), returning (nil, nil) rather than an error.
func TestUnitPaiAlgorithmVersionReadNotFoundHTTPWire(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
			"Code": "NoSuchObject", "Message": "algorithm version not found", "RequestId": "r-nf",
		})
	}))
	defer server.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server.URL))

	object, err := getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err != nil {
		t.Fatalf("getAlgorithmVersion not-found must surface as (nil,nil), got err: %s", err)
	}
	if object != nil {
		t.Fatalf("getAlgorithmVersion not-found must return nil object, got %#v", object)
	}
}

// TestUnitPaiAlgorithmVersionReadServerErrorHTTPWire asserts that a 5xx on the
// production RoaGet path is NOT treated as not-found: getAlgorithmVersion
// returns (nil, err) so a transient server error preserves the prior state
// instead of silently dropping the resource.
func TestUnitPaiAlgorithmVersionReadServerErrorHTTPWire(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusInternalServerError, map[string]interface{}{
			"Code": "InternalError", "Message": "boom", "RequestId": "r-500",
		})
	}))
	defer server.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server.URL))

	object, err := getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err == nil {
		t.Fatal("getAlgorithmVersion on 5xx must return an error, not nil")
	}
	if object != nil {
		t.Fatalf("getAlgorithmVersion on 5xx must return nil object, got %#v", object)
	}
	if isAlgorithmVersionNotExist(err) {
		t.Fatalf("a 5xx InternalError must not be treated as not-found: %s", err)
	}
}

// TestUnitPaiAlgorithmVersionReadMalformedHTTPWire asserts that a 200 response
// missing a required identifier on the production RoaGet path is caught by the
// response validation: getAlgorithmVersion returns the (well-formed-enough) map
// with no error, but validateAlgorithmVersionResponse rejects it so it cannot
// compose a ":" id or overwrite the resource identity. This is the malformed
// path the helper-only tests could not reach.
func TestUnitPaiAlgorithmVersionReadMalformedHTTPWire(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 200, but AlgorithmVersion is missing: would compose "algo-mock:".
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"AlgorithmId": mockAlgorithmID,
		})
	}))
	defer server.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server.URL))

	object, err := getAlgorithmVersion(client, mockAlgorithmID, mockAlgorithmVersion)
	if err != nil {
		t.Fatalf("getAlgorithmVersion on a 200 malformed response must not error at read time, got: %s", err)
	}
	if object == nil {
		t.Fatal("getAlgorithmVersion on a 200 malformed response must return the parsed object, not nil")
	}
	if err := validateAlgorithmVersionResponse(object, mockAlgorithmID, mockAlgorithmVersion); err == nil {
		t.Fatal("validateAlgorithmVersionResponse must reject a response missing AlgorithmVersion")
	}
	// And the malformed object must not be flattened onto state unguarded.
	if id := algorithmVersionState(object).Id.ValueString(); id == mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("malformed response unexpectedly composed the valid id: %q", id)
	}
}
