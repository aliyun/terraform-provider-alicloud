package polardb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestUnitPolarDBEndpointSccModeWire asserts the scc_mode request parameters sent by
// ModifyDBClusterEndpoint and the case normalization applied when Describe returns the value.
func TestUnitPolarDBEndpointSccModeWire(t *testing.T) {
	var mu sync.Mutex
	var modifyForm url.Values
	var describeCalls atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ModifyDBClusterEndpoint":
			mu.Lock()
			modifyForm = r.Form
			mu.Unlock()
			_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
		case "DescribeDBClusterEndpoints":
			describeCalls.Add(1)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Items":[{"DBEndpointId":"pe-unit","EndpointType":"Custom","SccMode":"ON","EndpointConfig":"{\"ConsistLevel\":\"0\"}","AutoAddNewNodes":"Disable","ReadWriteMode":"ReadOnly","AddressItems":[{"NetType":"Private","Port":"3306","ConnectionString":"pe-unit.rwlb.rds.aliyuncs.com"}]}]}`))
		case "DescribeDBClusterSSL":
			_, _ = w.Write([]byte(`{"RequestId":"mock","Items":[]}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
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
	endpoints.Store("polardb", endpoint)

	d := schema.TestResourceDataRaw(t, ResourceEndpoint().Schema, map[string]interface{}{
		"db_cluster_id": "pc-unit",
		"endpoint_type": "Custom",
		"endpoint_config": map[string]interface{}{
			"ConsistLevel": "0",
		},
		"scc_mode": "on",
	})
	d.MarkNewResource()
	d.SetId("pc-unit:pe-unit")

	if err := resourceEndpointUpdate(d, client); err != nil {
		t.Fatalf("update failed: %s", err)
	}

	mu.Lock()
	form := modifyForm
	mu.Unlock()
	if form == nil {
		t.Fatal("ModifyDBClusterEndpoint was not called")
	}
	if got := form.Get("DBClusterId"); got != "pc-unit" {
		t.Errorf("DBClusterId = %q, want %q", got, "pc-unit")
	}
	if got := form.Get("DBEndpointId"); got != "pe-unit" {
		t.Errorf("DBEndpointId = %q, want %q", got, "pe-unit")
	}
	if got := form.Get("SccMode"); got != "on" {
		t.Errorf("SccMode request parameter = %q, want %q", got, "on")
	}
	var endpointConfig map[string]string
	if err := json.Unmarshal([]byte(form.Get("EndpointConfig")), &endpointConfig); err != nil {
		t.Fatalf("decode EndpointConfig request parameter %q: %s", form.Get("EndpointConfig"), err)
	}
	if endpointConfig["ConsistLevel"] != "0" {
		t.Errorf("EndpointConfig.ConsistLevel = %q, want %q", endpointConfig["ConsistLevel"], "0")
	}
	if got := d.Get("scc_mode"); got != "on" {
		t.Errorf("scc_mode after read = %q, want %q for Describe value ON", got, "on")
	}
	if describeCalls.Load() == 0 {
		t.Error("DescribeDBClusterEndpoints was not called")
	}
}
