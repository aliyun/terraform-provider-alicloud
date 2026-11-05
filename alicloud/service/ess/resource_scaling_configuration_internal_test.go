package ess

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// TestUnitOptionalComputedRequestMatrix asserts the serialized request keys for
// the Optional+Computed attributes: an omitted value contributes no key, an
// explicit (including zero) value contributes the key, and a removed value
// contributes nothing so the server keeps its own state.
func TestUnitOptionalComputedRequestMatrix(t *testing.T) {
	zero := 0
	two := 2
	hundred := 100
	on := true
	off := false

	omitted := map[string]interface{}{}
	addOptionalComputedRequestFields(omitted, optionalComputedRequestFields{})
	if len(omitted) != 0 {
		t.Fatalf("omitted attributes must serialize nothing, got %#v", omitted)
	}

	explicitZero := map[string]interface{}{}
	addOptionalComputedRequestFields(explicitZero, optionalComputedRequestFields{
		LoadBalancerWeight: &zero,
		Ipv6AddressCount:   &zero,
		Cpu:                &zero,
		Memory:             &zero,
	})
	for _, key := range []string{"LoadBalancerWeight", "Ipv6AddressCount", "Cpu", "Memory"} {
		value, ok := explicitZero[key]
		if !ok {
			t.Fatalf("%s must be serialized when explicitly configured as zero", key)
		}
		if string(value.(requests.Integer)) != "0" {
			t.Fatalf("%s must serialize 0, got %s", key, value)
		}
	}

	explicitValues := map[string]interface{}{}
	addOptionalComputedRequestFields(explicitValues, optionalComputedRequestFields{
		LoadBalancerWeight: &hundred,
		Cpu:                &two,
		Memory:             &two,
	})
	if string(explicitValues["LoadBalancerWeight"].(requests.Integer)) != "100" {
		t.Fatalf("LoadBalancerWeight must serialize 100, got %v", explicitValues["LoadBalancerWeight"])
	}
	if string(explicitValues["Cpu"].(requests.Integer)) != "2" {
		t.Fatalf("Cpu must serialize 2, got %v", explicitValues["Cpu"])
	}

	// Removal yields nil pointers (the attribute is no longer configured), so the
	// builder must resend nothing and leave the server value untouched.
	removed := map[string]interface{}{}
	addOptionalComputedRequestFields(removed, optionalComputedRequestFields{})
	if len(removed) != 0 {
		t.Fatalf("removed attributes must resend nothing, got %#v", removed)
	}

	burstingOn := map[string]interface{}{}
	addOptionalComputedRequestFields(burstingOn, optionalComputedRequestFields{SystemDiskBursting: &on})
	if v, ok := burstingOn["SystemDisk.BurstingEnabled"]; !ok || v != true {
		t.Fatalf("explicit true bursting must be serialized, got %#v", burstingOn)
	}
	burstingOff := map[string]interface{}{}
	addOptionalComputedRequestFields(burstingOff, optionalComputedRequestFields{SystemDiskBursting: &off})
	if v, ok := burstingOff["SystemDisk.BurstingEnabled"]; !ok || v != false {
		t.Fatalf("explicit false bursting must be serialized, got %#v", burstingOff)
	}
}

// TestUnitBurstingRequestMatrix asserts that an explicitly configured false is
// sent (so the remote value can be turned off) while an omitted value is not.
func TestUnitBurstingRequestMatrix(t *testing.T) {
	omitted := map[string]interface{}{}
	putBurstingRequest(omitted, false, false)
	if _, ok := omitted["BurstingEnabled"]; ok {
		t.Fatalf("omitted bursting must not be serialized, got %#v", omitted)
	}

	disabled := map[string]interface{}{}
	putBurstingRequest(disabled, false, true)
	if v, ok := disabled["BurstingEnabled"]; !ok || v != false {
		t.Fatalf("explicit false bursting must be serialized, got %#v", disabled)
	}

	enabled := map[string]interface{}{}
	putBurstingRequest(enabled, true, false)
	if v, ok := enabled["BurstingEnabled"]; !ok || v != true {
		t.Fatalf("true bursting must be serialized, got %#v", enabled)
	}
}

// TestUnitPrivatePoolAndSchedulerRequest asserts the block request shapes,
// including empty and nil blocks that must contribute no keys.
func TestUnitPrivatePoolAndSchedulerRequest(t *testing.T) {
	empty := map[string]interface{}{}
	buildPrivatePoolOptionsRequest(empty, nil)
	buildPrivatePoolOptionsRequest(empty, []interface{}{})
	buildPrivatePoolOptionsRequest(empty, []interface{}{nil})
	buildPrivatePoolOptionsRequest(empty, []interface{}{map[string]interface{}{}})
	if len(empty) != 0 {
		t.Fatalf("empty private pool blocks must serialize nothing, got %#v", empty)
	}

	filled := map[string]interface{}{}
	buildPrivatePoolOptionsRequest(filled, []interface{}{map[string]interface{}{
		"match_criteria": "Target",
		"pool_id":        "pp-123",
	}})
	if filled["PrivatePoolOptions.MatchCriteria"] != "Target" || filled["PrivatePoolOptions.Id"] != "pp-123" {
		t.Fatalf("private pool must serialize dotted keys, got %#v", filled)
	}

	emptyScheduler := map[string]interface{}{}
	buildSchedulerOptionsRequest(emptyScheduler, []interface{}{map[string]interface{}{}})
	if len(emptyScheduler) != 0 {
		t.Fatalf("empty scheduler block must serialize nothing, got %#v", emptyScheduler)
	}

	scheduler := map[string]interface{}{}
	buildSchedulerOptionsRequest(scheduler, []interface{}{map[string]interface{}{
		"managed_private_space_id": "mps-123",
	}})
	object, ok := scheduler["SchedulerOptions"].(map[string]interface{})
	if !ok || object["ManagedPrivateSpaceId"] != "mps-123" {
		t.Fatalf("scheduler options must serialize as an object, got %#v", scheduler)
	}
}
