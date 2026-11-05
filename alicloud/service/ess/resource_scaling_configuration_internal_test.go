package ess

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// scalingConfigurationUnitTestServer serves the ESS actions the resource issues
// during a create/update unit test. Captured forms are grouped by action so a
// test can assert the exact parameters the provider serialized.
type scalingConfigurationUnitTestServer struct {
	mutex    sync.Mutex
	requests map[string][]map[string]string
}

func newScalingConfigurationUnitTestServer(t *testing.T) *scalingConfigurationUnitTestServer {
	t.Helper()
	return &scalingConfigurationUnitTestServer{requests: map[string][]map[string]string{}}
}

func (s *scalingConfigurationUnitTestServer) handler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("cannot parse the provider request: %s", err)
		}
		form := make(map[string]string, len(r.Form))
		for key, values := range r.Form {
			form[key] = values[0]
		}
		action := r.Form.Get("Action")
		s.mutex.Lock()
		s.requests[action] = append(s.requests[action], form)
		s.mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "CreateScalingConfiguration":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"ScalingConfigurationId": "asc-unit-test"})
		case "DescribeScalingConfigurations":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ScalingConfigurations": map[string]interface{}{
					"ScalingConfiguration": []interface{}{map[string]interface{}{
						"ScalingConfigurationId": "asc-unit-test",
						"ScalingGroupId":         "asg-unit-test",
						"LifecycleState":         "Active",
						"InstanceType":           "ecs.g6.large",
						"SecurityGroupId":        "sg-unit-test",
					}},
				},
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "unit-test"})
		}
	}
}

func (s *scalingConfigurationUnitTestServer) captured(action string) []map[string]string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return append([]map[string]string(nil), s.requests[action]...)
}

// applyScalingConfigurationUnitTest drives the public resource Diff/Apply entry
// and reproduces the raw configuration the gRPC plan path attaches to the diff
// (helper/schema/grpc_provider.go sets diff.RawConfig from the protocol config).
// The legacy Resource.Diff helper only copies RawConfig from the prior state, so
// without this step GetRawConfig inside Create/Update would always be null and
// the test would not exercise the real Optional+Computed behavior.
func applyScalingConfigurationUnitTest(t *testing.T, server *scalingConfigurationUnitTestServer, previous *sdktf.InstanceState, raw map[string]interface{}, rawConfig cty.Value) {
	t.Helper()
	client := essScalingConfigurationUnitTestClient(t, server.handler(t))
	scalingConfigurationResource := resourceAlicloudEssScalingConfiguration()
	config := sdktf.NewResourceConfigRaw(raw)
	diff, err := scalingConfigurationResource.Diff(context.Background(), previous, config, client)
	if err != nil {
		t.Fatalf("diff failed: %s", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for the unit test configuration")
	}
	diff.RawConfig = rawConfig
	if _, diags := scalingConfigurationResource.Apply(context.Background(), previous, diff, client); diags.HasError() {
		t.Fatalf("apply reported errors: %v", diags)
	}
}

func assertRequestKeyAbsent(t *testing.T, action string, request map[string]string, keys ...string) {
	t.Helper()
	for _, key := range keys {
		if value, ok := request[key]; ok {
			t.Fatalf("%s must not serialize %s, got %q", action, key, value)
		}
	}
}

// TestUnitEssScalingConfigurationCreateRequestMatrix proves through the public
// resource entry (Diff + Apply + real client serialization) that omitted
// Optional+Computed attributes contribute no request key, while explicitly
// configured values, including a legal explicit false, are serialized.
func TestUnitEssScalingConfigurationCreateRequestMatrix(t *testing.T) {
	t.Run("omitted attributes serialize no keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
		}, cty.EmptyObjectVal)
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		assertRequestKeyAbsent(t, "CreateScalingConfiguration", creates[0],
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.MatchCriteria",
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("explicit values including false are serialized", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"security_group_id":            "sg-unit-test",
			"cpu":                          2,
			"memory":                       4,
			"load_balancer_weight":         50,
			"system_disk_bursting_enabled": false,
			"private_pool_options":         []interface{}{map[string]interface{}{}},
			"scheduler_options":            []interface{}{map[string]interface{}{}},
		}, cty.ObjectVal(map[string]cty.Value{
			"cpu":                          cty.NumberIntVal(2),
			"memory":                       cty.NumberIntVal(4),
			"load_balancer_weight":         cty.NumberIntVal(50),
			"system_disk_bursting_enabled": cty.False,
			"private_pool_options":         cty.ListValEmpty(cty.EmptyObject),
			"scheduler_options":            cty.ListValEmpty(cty.EmptyObject),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		request := creates[0]
		for key, want := range map[string]string{
			"Cpu":                        "2",
			"Memory":                     "4",
			"LoadBalancerWeight":         "50",
			"SystemDisk.BurstingEnabled": "false",
		} {
			if got := request[key]; got != want {
				t.Fatalf("%s = %q, expected %q", key, got, want)
			}
		}
		assertRequestKeyAbsent(t, "CreateScalingConfiguration", request,
			"PrivatePoolOptions.MatchCriteria", "PrivatePoolOptions.Id", "SchedulerOptions")
	})
}

// TestUnitEssScalingConfigurationUpdateRequestMatrix proves that an unrelated
// update and an attribute removal resend nothing for the removed
// Optional+Computed attributes, while an explicit false is still sent.
func TestUnitEssScalingConfigurationUpdateRequestMatrix(t *testing.T) {
	previousState := func() *sdktf.InstanceState {
		return &sdktf.InstanceState{
			ID: "asc-unit-test",
			Attributes: map[string]string{
				"id":                           "asc-unit-test",
				"scaling_group_id":             "asg-unit-test",
				"instance_type":                "ecs.g6.large",
				"security_group_id":            "sg-unit-test",
				"cpu":                          "2",
				"memory":                       "4",
				"load_balancer_weight":         "50",
				"system_disk_bursting_enabled": "true",
				"scaling_configuration_name":   "old-name",
			},
		}
	}

	t.Run("removal and unrelated update serialize no stale keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, previousState(), map[string]interface{}{
			"scaling_group_id":           "asg-unit-test",
			"instance_type":              "ecs.g6.large",
			"security_group_id":          "sg-unit-test",
			"scaling_configuration_name": "new-name",
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name": cty.StringVal("new-name"),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		request := modifies[0]
		if got := request["ScalingConfigurationName"]; got != "new-name" {
			t.Fatalf("ScalingConfigurationName = %q, expected new-name", got)
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", request,
			"Cpu", "Memory", "LoadBalancerWeight", "SystemDisk.BurstingEnabled")
	})

	t.Run("explicit false is sent on update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, previousState(), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"system_disk_bursting_enabled": false,
		}, cty.ObjectVal(map[string]cty.Value{
			"system_disk_bursting_enabled": cty.False,
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := modifies[0]["SystemDisk.BurstingEnabled"]; got != "false" {
			t.Fatalf("SystemDisk.BurstingEnabled = %q, expected false", got)
		}
	})
}
