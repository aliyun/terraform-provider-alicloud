package ess

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// scalingConfigurationPublicResource resolves the resource through the public
// service package, i.e. the same factory the provider registers, so the request
// matrix exercises the shipped entry instead of the package-private function.
// This file stays in package ess only because it reuses the package-level
// fake-client helper declared in scaling_configuration_internal_test.go; the
// provider entry itself is public.
func scalingConfigurationPublicResource(t *testing.T) *schema.Resource {
	t.Helper()
	for _, resource := range ServicePackage().SDKResources {
		if resource.TypeName == "alicloud_ess_scaling_configuration" {
			return resource.Factory()
		}
	}
	t.Fatal("alicloud_ess_scaling_configuration is not registered in the public Ess service package")
	return nil
}

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
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
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
	scalingConfigurationResource := scalingConfigurationPublicResource(t)
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

// requestSerializedValue returns the serialized form of an aggregated request
// parameter, such as a JSON array, that the RPC client may emit under a single
// key or under indexed dotted keys. It joins every matching "key=value" pair so
// a test can assert on the document and the key names instead of guessing the
// transport shape.
func requestSerializedValue(request map[string]string, key string) string {
	var parts []string
	for name, value := range request {
		if name == key || strings.HasPrefix(name, key+".") {
			parts = append(parts, name+"="+value)
		}
	}
	return strings.Join(parts, ",")
}

// scalingConfigurationOptionalBaseState is the prior state used by the update
// matrix: a created configuration with the Optional+Computed attributes and both
// blocks configured, all values stored as state strings.
func scalingConfigurationOptionalBaseState() *sdktf.InstanceState {
	return scalingConfigurationOptionalBaseStateWithBursting(true)
}

func scalingConfigurationOptionalBaseStateWithBursting(bursting bool) *sdktf.InstanceState {
	return &sdktf.InstanceState{
		ID: "asc-unit-test",
		Attributes: map[string]string{
			"id":                                    "asc-unit-test",
			"scaling_group_id":                      "asg-unit-test",
			"instance_type":                         "ecs.g6.large",
			"security_group_id":                     "sg-unit-test",
			"cpu":                                   "2",
			"memory":                                "4",
			"load_balancer_weight":                  "50",
			"ipv6_address_count":                    "1",
			"system_disk_bursting_enabled":          fmt.Sprint(bursting),
			"scaling_configuration_name":            "old-name",
			"private_pool_options.#":                "1",
			"private_pool_options.0.match_criteria": "Open",
			"private_pool_options.0.pool_id":        "pp-unit-test",
			"scheduler_options.#":                   "1",
			"scheduler_options.0.managed_private_space_id": "mps-unit-test",
		},
	}
}

// TestUnitScalingConfigurationCreateRequestMatrix proves through the public
// resource entry (Diff + Apply + real client serialization) that omitted
// Optional+Computed attributes contribute no request key, while explicitly
// configured values, including a legal explicit zero and false, are serialized.
func TestUnitScalingConfigurationCreateRequestMatrix(t *testing.T) {
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

	t.Run("a legal zero ipv6 count is serialized", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":   "asg-unit-test",
			"instance_type":      "ecs.g6.large",
			"security_group_id":  "sg-unit-test",
			"ipv6_address_count": 0,
		}, cty.ObjectVal(map[string]cty.Value{
			"ipv6_address_count": cty.NumberIntVal(0),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := creates[0]["Ipv6AddressCount"]; got != "0" {
			t.Fatalf("Ipv6AddressCount = %q, expected the configured zero", got)
		}
	})

	t.Run("a configured block serializes its keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"private_pool_options": []interface{}{map[string]interface{}{
				"match_criteria": "Target",
				"pool_id":        "pp-unit-test",
			}},
			"scheduler_options": []interface{}{map[string]interface{}{
				"managed_private_space_id": "mps-unit-test",
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"private_pool_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"match_criteria": cty.StringVal("Target"),
				"pool_id":        cty.StringVal("pp-unit-test"),
			})}),
			"scheduler_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"managed_private_space_id": cty.StringVal("mps-unit-test"),
			})}),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		request := creates[0]
		if got := request["PrivatePoolOptions.MatchCriteria"]; got != "Target" {
			t.Fatalf("PrivatePoolOptions.MatchCriteria = %q, expected Target", got)
		}
		if got := request["PrivatePoolOptions.Id"]; got != "pp-unit-test" {
			t.Fatalf("PrivatePoolOptions.Id = %q, expected pp-unit-test", got)
		}
		if got := requestSerializedValue(request, "SchedulerOptions"); !strings.Contains(got, "mps-unit-test") {
			t.Fatalf("SchedulerOptions = %q, expected the managed private space id", got)
		}
	})

	t.Run("data disk bursting false is serialized and omitted when unset", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"data_disk": []interface{}{map[string]interface{}{
				"size":             20,
				"bursting_enabled": false,
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"data_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"size":             cty.NumberIntVal(20),
				"bursting_enabled": cty.False,
			})}),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := requestSerializedValue(creates[0], "DataDisk"); !strings.Contains(got, `"BurstingEnabled":false`) && !strings.Contains(got, "BurstingEnabled=false") {
			t.Fatalf("DataDisk = %q, expected the explicit false BurstingEnabled", got)
		}

		omitted := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, omitted, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"data_disk": []interface{}{map[string]interface{}{
				"size": 20,
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"data_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"size": cty.NumberIntVal(20),
			})}),
		}))
		creates = omitted.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := requestSerializedValue(creates[0], "DataDisk"); strings.Contains(got, "BurstingEnabled") {
			t.Fatalf("DataDisk = %q, expected no BurstingEnabled for an omitted value", got)
		}
	})
}

// TestUnitScalingConfigurationUpdateRequestMatrix proves that removals, pure
// unrelated updates, explicit values and cleared blocks serialize exactly the
// intended keys on ModifyScalingConfiguration.
func TestUnitScalingConfigurationUpdateRequestMatrix(t *testing.T) {
	t.Run("removal and unrelated update serialize no stale keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
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
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.MatchCriteria",
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("pure unrelated update resends unchanged configured attributes", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"scaling_configuration_name":   "new-name",
			"cpu":                          2,
			"memory":                       4,
			"load_balancer_weight":         50,
			"ipv6_address_count":           1,
			"system_disk_bursting_enabled": true,
			"private_pool_options": []interface{}{map[string]interface{}{
				"match_criteria": "Open",
				"pool_id":        "pp-unit-test",
			}},
			"scheduler_options": []interface{}{map[string]interface{}{
				"managed_private_space_id": "mps-unit-test",
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name":   cty.StringVal("new-name"),
			"cpu":                          cty.NumberIntVal(2),
			"memory":                       cty.NumberIntVal(4),
			"load_balancer_weight":         cty.NumberIntVal(50),
			"ipv6_address_count":           cty.NumberIntVal(1),
			"system_disk_bursting_enabled": cty.True,
			"private_pool_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"match_criteria": cty.StringVal("Open"),
				"pool_id":        cty.StringVal("pp-unit-test"),
			})}),
			"scheduler_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"managed_private_space_id": cty.StringVal("mps-unit-test"),
			})}),
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
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.MatchCriteria",
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("explicit false is sent on update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
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

	t.Run("non-zero and true values are sent on update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseStateWithBursting(false), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"load_balancer_weight":         60,
			"ipv6_address_count":           2,
			"system_disk_bursting_enabled": true,
		}, cty.ObjectVal(map[string]cty.Value{
			"load_balancer_weight":         cty.NumberIntVal(60),
			"ipv6_address_count":           cty.NumberIntVal(2),
			"system_disk_bursting_enabled": cty.True,
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		request := modifies[0]
		for key, want := range map[string]string{
			"LoadBalancerWeight":         "60",
			"Ipv6AddressCount":           "2",
			"SystemDisk.BurstingEnabled": "true",
		} {
			if got := request[key]; got != want {
				t.Fatalf("%s = %q, expected %q", key, got, want)
			}
		}
	})

	t.Run("configured cpu and memory are sent without an instance type", func(t *testing.T) {
		previous := &sdktf.InstanceState{
			ID: "asc-unit-test",
			Attributes: map[string]string{
				"id":                "asc-unit-test",
				"scaling_group_id":  "asg-unit-test",
				"security_group_id": "sg-unit-test",
			},
		}
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, previous, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"security_group_id": "sg-unit-test",
			"cpu":               4,
			"memory":            8,
		}, cty.ObjectVal(map[string]cty.Value{
			"cpu":    cty.NumberIntVal(4),
			"memory": cty.NumberIntVal(8),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := modifies[0]["Cpu"]; got != "4" {
			t.Fatalf("Cpu = %q, expected 4", got)
		}
		if got := modifies[0]["Memory"]; got != "8" {
			t.Fatalf("Memory = %q, expected 8", got)
		}
	})

	t.Run("clearing a block serializes no block keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
		}, cty.EmptyObjectVal)
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", modifies[0],
			"PrivatePoolOptions.MatchCriteria", "PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("an unchanged configured value is not resent on an unrelated update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":           "asg-unit-test",
			"instance_type":              "ecs.g6.large",
			"security_group_id":          "sg-unit-test",
			"scaling_configuration_name": "new-name",
			"ipv6_address_count":         1,
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name": cty.StringVal("new-name"),
			"ipv6_address_count":         cty.NumberIntVal(1),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", modifies[0], "Ipv6AddressCount")
	})
}

// TestUnitScalingConfigurationRequestSerializedValue keeps the local helper
// honest: requestSerializedValue must see indexed dotted keys as well as one
// aggregated JSON value.
func TestUnitScalingConfigurationRequestSerializedValue(t *testing.T) {
	request := map[string]string{
		"DataDisk.1.Size":         "20",
		"DataDisk.1.Categories.1": "cloud",
		"Tags":                    "{}",
	}
	if got := requestSerializedValue(request, "DataDisk"); !strings.Contains(got, "DataDisk.1.Size=20") {
		t.Fatalf("requestSerializedValue = %q, expected the indexed keys", got)
	}
	if got := requestSerializedValue(request, "Missing"); got != "" {
		t.Fatalf("requestSerializedValue = %q, expected empty for a missing key", got)
	}
}
