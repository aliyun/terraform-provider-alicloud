package paiplugin_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// regexCreateFails matches the diagnostic summary every resource emits when a
// Create ROA call fails (e.g. the gateway returns HTTP 400), so the failure-path
// steps assert the apply aborted on creation rather than on an unrelated cause.
var regexCreateFails = regexp.MustCompile(`(?i)Creating alicloud_pai_plugin`)

// paiPluginEnvelope wraps a PAIPlugin gateway response. The ROA client
// unwraps the tea "body" nesting, so the handler writes this shape directly;
// roaRequest returns it verbatim and the service helpers read Data/ErrorCode.
func paiPluginEnvelope(data interface{}) []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"Data":         data,
		"ErrorCode":    0,
		"ErrorMessage": "OK",
		"RequestId":    "test-request-id",
		"Success":      true,
	})
	return b
}

// paiPluginNotFound writes the not-found signal PAIPlugin uses: HTTP 400 with a
// body carrying ErrorCode 400. isNotFound treats both the HTTP 400 (surfaced as
// a Go error) and the body ErrorCode 400 as the gone signal, so either shape is
// accepted; the handler emits both to mirror the real gateway.
func paiPluginNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"Data":null,"ErrorCode":400,"ErrorMessage":"NotFound","RequestId":"test-not-found","Success":false}`))
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, payload []byte) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

// paiPluginTestServer stands up the in-process httptest gateway and points the
// provider at it via the ALIBABA_CLOUD_ENDPOINT_PAIPLUGIN env override, the same
// mechanism loadEndpoint reads. The provider block sets protocol=HTTP so the ROA
// tea client speaks plain http to the server (the cloudfw/rds unit tests prove
// ROA honors Protocol:http), and skip_region_validation suppresses the only
// network call Configure would otherwise make. No real credentials are used.
func paiPluginTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ROA GET carries params on the query string; POST/PUT carry a JSON body.
		// Drain the body so the connection stays reusable regardless of method.
		if r.Body != nil {
			_, _ = io.ReadAll(r.Body)
			r.Body.Close()
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret")
	t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "true")
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_PAIPLUGIN", endpoint)
	t.Setenv("NO_PROXY", endpoint)
	return server
}

// paiPluginProviderConfig is the no-cloud provider block shared by every mock
// test: HTTP protocol so the ROA transport matches the httptest server, dummy
// credentials, and region validation skipped so Configure builds offline.
func paiPluginProviderConfig() string {
	return `
provider "alicloud" {
  access_key            = "test-key"
  secret_key           = "test-secret"
  region               = "cn-hangzhou"
  protocol             = "HTTP"
  skip_region_validation = true
}
`
}

// factories returns the real muxed provider (SDKv2 + Framework behind one
// proto5 server), so the Framework-native PAIPlugin resources exercise their
// real Create/Read/Update/Delete/Import against the mock gateway.
func factories() map[string]func() (tfprotov5.ProviderServer, error) {
	return acctest.ProtoV5ProviderFactories()
}
