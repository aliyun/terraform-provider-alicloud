package alicloud

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// sasProbeNotReadyBody mocks the GetHoneypotProbe response when a probe is in
// install_failed state: the API returns HoneypotProbeNotReady instead of data.
const sasProbeNotReadyBody = `{"Code":"HoneypotProbeNotReady","Message":"probe is in install_failed state","Success":false,"RequestId":"r-not-ready"}`

// sasProbeSuccessNoDataBody mocks a successful GetHoneypotProbe response that
// carries no Data field, which Describe maps to NotFound so Read clears state.
const sasProbeSuccessNoDataBody = `{"Success":true,"RequestId":"r-no-data"}`

// sasProbeTestClient stands up an in-process HTTP server that always replies
// with body and routes the Sas RPC client to it via endpoints.Store("sas").
// The returned counter tracks how many RPC calls reached the server, so tests
// can assert that install_failed is not retried for 5 minutes.
func sasProbeTestClient(t *testing.T, body string) (*connectivity.AliyunClient, *int32) {
	t.Helper()
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("new credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("config client: %v", err)
	}
	endpoints.Store("sas", endpoint)
	return client, &requests
}

// honeypotProbeResourceData builds an empty ResourceData for the honeypot probe
// resource so tests can set its id and exercise Read/StateRefresh in isolation.
func honeypotProbeResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d, err := schema.InternalMap(resourceAlicloudThreatDetectionHoneypotProbe().Schema).Data(nil, nil)
	if err != nil {
		t.Fatalf("build resource data: %v", err)
	}
	return d
}

// TestUnitHoneypotProbeDescribeNoRetryOnNotReady verifies that when
// GetHoneypotProbe returns HoneypotProbeNotReady, DescribeThreatDetectionHoneypotProbe
// fails fast (no 5-minute retry) and exposes the original error code so callers
// can map it to install_failed.
func TestUnitHoneypotProbeDescribeNoRetryOnNotReady(t *testing.T) {
	client, requests := sasProbeTestClient(t, sasProbeNotReadyBody)
	sas := SasService{client}

	_, err := sas.DescribeThreatDetectionHoneypotProbe("probe-1")
	if err == nil {
		t.Fatal("expected error from DescribeThreatDetectionHoneypotProbe when probe is install_failed, got nil")
	}
	if !IsExpectedErrors(err, []string{"HoneypotProbeNotReady"}) {
		t.Fatalf("expected HoneypotProbeNotReady error code, got %#v", err)
	}
	if got := atomic.LoadInt32(requests); got != 1 {
		t.Fatalf("expected exactly 1 request (install_failed must not be retried), got %d", got)
	}
}

// TestUnitHoneypotProbeReadInstallFailedKeepsState verifies AC1: after
// install_failed, Read keeps the resource in state (no d.SetId("")), marks
// status=install_failed, and returns nil so plan has no diff and destroy can
// still reach DeleteHoneypotProbe.
func TestUnitHoneypotProbeReadInstallFailedKeepsState(t *testing.T) {
	client, requests := sasProbeTestClient(t, sasProbeNotReadyBody)
	d := honeypotProbeResourceData(t)
	d.SetId("probe-1")

	if err := resourceAlicloudThreatDetectionHoneypotProbeRead(d, client); err != nil {
		t.Fatalf("expected Read to keep state on install_failed (AC1), got error: %v", err)
	}
	if d.Id() != "probe-1" {
		t.Fatalf("expected state kept (id=probe-1), got id=%q", d.Id())
	}
	if got := d.Get("status"); got != "install_failed" {
		t.Fatalf("expected status=install_failed, got %v", got)
	}
	if got := atomic.LoadInt32(requests); got != 1 {
		t.Fatalf("expected exactly 1 request, got %d", got)
	}
}

// TestUnitHoneypotProbeStateRefreshInstallFailed verifies that the create/update
// state refresh maps HoneypotProbeNotReady to the install_failed terminal state,
// so BuildStateConf fails fast (AC3) instead of looping until timeout.
func TestUnitHoneypotProbeStateRefreshInstallFailed(t *testing.T) {
	client, requests := sasProbeTestClient(t, sasProbeNotReadyBody)
	d := honeypotProbeResourceData(t)
	d.SetId("probe-1")
	sas := SasService{client}

	refresh := sas.ThreatDetectionHoneypotProbeStateRefreshFunc(d, []string{"install_failed"})
	obj, state, err := refresh()
	if obj != nil {
		t.Fatalf("expected nil object on install_failed, got %#v", obj)
	}
	if state != "install_failed" {
		t.Fatalf("expected state=install_failed, got %q", state)
	}
	if err == nil {
		t.Fatal("expected error describing install_failed terminal state, got nil")
	}
	if got := atomic.LoadInt32(requests); got != 1 {
		t.Fatalf("expected exactly 1 request, got %d", got)
	}
}

// TestUnitHoneypotProbeReadNotFoundClearsState verifies that when the probe no
// longer exists (Success=true, no Data), Read clears state instead of erroring,
// so a stale install_failed resource can be removed by destroy.
func TestUnitHoneypotProbeReadNotFoundClearsState(t *testing.T) {
	client, requests := sasProbeTestClient(t, sasProbeSuccessNoDataBody)
	d := honeypotProbeResourceData(t)
	d.SetId("probe-1")

	if err := resourceAlicloudThreatDetectionHoneypotProbeRead(d, client); err != nil {
		t.Fatalf("expected Read to clear state on NotFound, got error: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected id cleared on NotFound, got id=%q", d.Id())
	}
	if got := atomic.LoadInt32(requests); got != 1 {
		t.Fatalf("expected exactly 1 request, got %d", got)
	}
}
