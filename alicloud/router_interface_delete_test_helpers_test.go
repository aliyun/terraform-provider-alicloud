package alicloud

// Root-package compatibility copies of the router-interface delete test helpers.
// These helpers were originally declared in resource_alicloud_router_interface_test.go
// and shared by resource_alicloud_vpc_test.go. The router_interface resource (and
// its test file) migrated to alicloud/service/expressconnect; the VPC delete unit
// test in resource_alicloud_vpc_test.go is not part of the expressconnect migration
// scope and still lives in the root package, so the helpers it depends on are
// re-declared here verbatim to keep the root test package compiling unchanged.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func routerInterfaceDeleteTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("RegionId") != "cn-hangzhou" {
			t.Errorf("unexpected region: %q", r.Form.Get("RegionId"))
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("vpc", endpoint)
	return client
}

func routerInterfaceDeleteTestResponse(t *testing.T, w http.ResponseWriter, response interface{}) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Error(err)
	}
}

func routerInterfaceDeleteTestError(t *testing.T, w http.ResponseWriter, code string) {
	t.Helper()
	if strings.HasSuffix(code, ".NotFound") {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{"Code": code, "Message": "test operation failed"})
}
