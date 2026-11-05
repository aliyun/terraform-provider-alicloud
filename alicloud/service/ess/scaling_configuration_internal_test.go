package ess

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
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// essScalingConfigurationUnitTestClient returns a provider client whose ESS
// endpoint is a local HTTP server, so a unit test can inspect the requests the
// provider itself serializes without reaching the cloud.
func essScalingConfigurationUnitTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	endpoints.Store("ess", endpoint)
	return client
}

// TestUnitEssScalingConfigurationProviderClearRequest captures the request the
// provider serializes when tags become an explicit empty map and user_data
// becomes an empty string: the empty tag document is sent, while no UserData
// parameter is built at all.
func TestUnitEssScalingConfigurationProviderClearRequest(t *testing.T) {
	var mutex sync.Mutex
	modifyRequests := make([]map[string]string, 0)
	described := 0
	client := essScalingConfigurationUnitTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("cannot parse the provider request: %s", err)
		}
		form := make(map[string]string, len(r.Form))
		for key, values := range r.Form {
			form[key] = strings.Join(values, ",")
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "ModifyScalingConfiguration":
			mutex.Lock()
			modifyRequests = append(modifyRequests, form)
			mutex.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"ScalingConfigurationId": "asc-unit-test"})
		case "DescribeScalingConfigurations":
			mutex.Lock()
			described++
			mutex.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ScalingConfigurations": map[string]interface{}{
					"ScalingConfiguration": []interface{}{map[string]interface{}{
						"ScalingConfigurationId": "asc-unit-test",
						"ScalingGroupId":         "asg-unit-test",
						"LifecycleState":         "Active",
					}},
				},
			})
		default:
			t.Errorf("unexpected Action %q", r.Form.Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	scalingConfigurationResource := resourceAlicloudEssScalingConfiguration()
	previous := &sdktf.InstanceState{
		ID: "asc-unit-test",
		Attributes: map[string]string{
			"id":               "asc-unit-test",
			"scaling_group_id": "asg-unit-test",
			"force_delete":     "true",
			"tags.%":           "3",
			"tags.external":    "ack",
			"tags.managed-by":  "autoscaler",
			"tags.created":     "tf",
			"user_data":        "IyEvYmluL2Jhc2gKZWNobyB0ZXJyYWZvcm0tdXBkYXRlZA==",
		},
	}
	config := sdktf.NewResourceConfigRaw(map[string]interface{}{
		"scaling_group_id": "asg-unit-test",
		"force_delete":     true,
		"tags":             map[string]interface{}{},
		"user_data":        "",
	})
	diff, err := scalingConfigurationResource.Diff(context.Background(), previous, config, client)
	if err != nil {
		t.Fatalf("diff failed: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("expected the explicit empty values to produce a diff")
	}
	newState, diags := scalingConfigurationResource.Apply(context.Background(), previous, diff, client)
	if diags.HasError() {
		t.Fatalf("apply reported errors: %v", diags)
	}
	if newState == nil {
		t.Fatal("apply returned no state")
	}

	mutex.Lock()
	defer mutex.Unlock()
	if len(modifyRequests) != 1 {
		t.Fatalf("expected exactly one ModifyScalingConfiguration request, got %d: %v", len(modifyRequests), modifyRequests)
	}
	request := modifyRequests[0]
	if request["ScalingConfigurationId"] != "asc-unit-test" {
		t.Fatalf("unexpected ScalingConfigurationId %q", request["ScalingConfigurationId"])
	}
	if request["Tags"] != "{}" {
		t.Fatalf("expected the provider to send Tags=\"{}\", got %q", request["Tags"])
	}
	if value, ok := request["UserData"]; ok {
		t.Fatalf("expected no UserData parameter, got %q", value)
	}
	if described == 0 {
		t.Fatal("the update did not run the resource Read")
	}
	t.Logf("captured ModifyScalingConfiguration parameters: Tags=%q, UserData=<absent>, Read described %d times", request["Tags"], described)
}
