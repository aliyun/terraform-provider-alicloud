package cs

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitCSKubernetesPolicyInstanceParametersForState(t *testing.T) {
	input := map[string]interface{}{
		"string": "value",
		"bool":   false,
		"int":    30,
		"object": map[string]interface{}{
			"z": 3,
			"a": true,
		},
		"list": []interface{}{3, true, "value"},
	}
	expected := map[string]interface{}{
		"string": "value",
		"bool":   "false",
		"int":    "30",
		"object": `{"a":true,"z":3}`,
		"list":   `[3,true,"value"]`,
	}

	actual, err := policyParametersForState(input)
	if err != nil {
		t.Fatalf("policyParametersForState returned an error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("policyParametersForState() = %#v, want %#v", actual, expected)
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceCSKubernetesPolicyInstance().Schema, nil)
	if err := resourceData.Set("parameters", actual); err != nil {
		t.Fatalf("setting normalized policy parameters in state returned an error: %v", err)
	}

	roundTrip := convert.NormalizeMap(actual)
	roundTripJSON, err := json.Marshal(roundTrip)
	if err != nil {
		t.Fatalf("marshalling round-trip policy parameters returned an error: %v", err)
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshalling input policy parameters returned an error: %v", err)
	}
	if string(roundTripJSON) != string(inputJSON) {
		t.Fatalf("NormalizeMap(policyParametersForState()) = %s, want %s", roundTripJSON, inputJSON)
	}
}

func TestUnitCSKubernetesPolicyInstanceParametersForStateUnsupportedType(t *testing.T) {
	_, err := policyParametersForState(map[string]interface{}{"invalid": make(chan int)})
	if err == nil {
		t.Fatal("policyParametersForState should reject unsupported value types")
	}
}
