package wire_test

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper/wire"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// semanticString lives in the test package on purpose: it proves a type outside wire can
// satisfy the exported semantic hooks. The unexported hooks this replaced could never be
// implemented across a Go package boundary.
type semanticString struct{ basetypes.StringValue }

func semanticValueOf(value basetypes.StringValue) semanticString {
	return semanticString{StringValue: value}
}

func (v semanticString) SemanticEqual(other attr.Value) bool {
	value, ok := other.(semanticString)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false
	}
	return strings.EqualFold(v.ValueString(), value.ValueString())
}

func (v semanticString) NullValue() attr.Value {
	return semanticValueOf(basetypes.NewStringNull())
}

var (
	_ wire.SemanticValue = semanticString{}
	_ wire.NullValue     = semanticString{}
)

func TestUnitWireSemanticHooksCrossPackage(t *testing.T) {
	semanticA := semanticValueOf(basetypes.NewStringValue("{\"a\": 1}"))
	semanticB := semanticValueOf(basetypes.NewStringValue("{\"A\": 1}"))
	if wire.PlannedChanged(semanticA, semanticB) {
		t.Fatal("semantic equality should treat case-only differences as unchanged")
	}
	if !wire.ConfiguredMatch(semanticA, semanticB) {
		t.Fatal("configured match should use the semantic hook")
	}
	plainA := types.StringValue("{\"a\": 1}")
	plainB := types.StringValue("{\"A\": 1}")
	if !wire.PlannedChanged(plainA, plainB) {
		t.Fatal("plain strings should still compare with Equal")
	}

	unknown := semanticValueOf(basetypes.NewStringUnknown())
	known := wire.KnownValue(unknown)
	if known == nil || !known.IsNull() {
		t.Fatalf("unknown custom value should become a typed null, got %#v", known)
	}
	if _, ok := known.(semanticString); !ok {
		t.Fatalf("typed null should keep the custom type, got %T", known)
	}
	if got := wire.KnownValue(types.StringUnknown()); !got.IsNull() {
		t.Fatal("unknown built-in string should still become a null string")
	}
}

func TestUnitWireSelectReadObject(t *testing.T) {
	rows := []interface{}{
		map[string]interface{}{"Namespace": "alpha"},
		map[string]interface{}{"Namespace": "beta"},
	}
	selected, err := wire.SelectReadObject(rows, []wire.ReadIdentity{{Path: "$.Namespace", Value: "beta"}})
	if err != nil {
		t.Fatalf("selecting an existing identity failed: %s", err)
	}
	if selected == nil || selected["Namespace"] != "beta" {
		t.Fatalf("selected the wrong row: %#v", selected)
	}

	missing, err := wire.SelectReadObject(rows, []wire.ReadIdentity{{Path: "$.Namespace", Value: "gamma"}})
	if err != nil || missing != nil {
		t.Fatalf("a missing identity should yield nil without error, got %#v / %v", missing, err)
	}

	empty, err := wire.SelectReadObject([]interface{}{}, []wire.ReadIdentity{{Path: "$.Namespace", Value: "beta"}})
	if err != nil || empty != nil {
		t.Fatalf("an empty collection should yield nil without error, got %#v / %v", empty, err)
	}

	if _, err := wire.SelectReadObject([]interface{}{"unexpected"}, []wire.ReadIdentity{{Path: "$.Namespace", Value: "beta"}}); err == nil || !strings.Contains(err.Error(), "not an object") {
		t.Fatalf("a non-object collection item must report its shape, got %v", err)
	}

	duplicated := []interface{}{
		map[string]interface{}{"Namespace": "beta"},
		map[string]interface{}{"Namespace": "beta"},
	}
	if _, err := wire.SelectReadObject(duplicated, []wire.ReadIdentity{{Path: "$.Namespace", Value: "beta"}}); err == nil {
		t.Fatal("multiple matching identities must be an error")
	}

	unreadable := []interface{}{map[string]interface{}{"Other": "value"}}
	if _, err := wire.SelectReadObject(unreadable, []wire.ReadIdentity{{Path: "$.Namespace", Value: "beta"}}); err == nil {
		t.Fatal("a list item without a readable identity must be an error")
	}

	if _, err := wire.SelectReadObject(rows, nil); err == nil {
		t.Fatal("reading without a declared identity must be an error")
	}
}

func TestUnitWireApplyMappingsAndPreserveRequestOnlyState(t *testing.T) {
	request := map[string]interface{}{}
	values := map[string]interface{}{"Name": "alpha"}
	if err := wire.ApplyMappings(request, values, []wire.Mapping{{Source: "$.Name", Target: "Name"}}, false); err != nil {
		t.Fatalf("applying a scalar mapping failed: %s", err)
	}
	if request["Name"] != "alpha" {
		t.Fatalf("scalar mapping wrote %#v", request["Name"])
	}

	flattened := map[string]interface{}{}
	flattenValues := map[string]interface{}{"Deep": map[string]interface{}{"Key": "value"}}
	bindings := []wire.Mapping{{Source: "$.Deep.Key", Target: "Deep.Key", Repeat: true}}
	if err := wire.ApplyMappings(flattened, flattenValues, bindings, false); err != nil {
		t.Fatalf("applying a repeated nested mapping failed: %s", err)
	}
	if flattened["Deep.Key"] != "value" {
		t.Fatalf("repeated nested mapping should flatten to a dotted key, got %#v", flattened)
	}

	attributeTypes := map[string]attr.Type{"name": types.StringType, "secret": types.StringType}
	observed := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"name":   types.StringValue("alpha"),
		"secret": types.StringValue(""),
	})
	prior := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"name":   types.StringValue("alpha"),
		"secret": types.StringValue("request-only"),
	})
	fields := map[string]wire.WireField{
		"secret": {Name: "secret", Read: false},
	}
	preserved, ok := wire.PreserveRequestValues(observed, prior, fields).(types.Object)
	if !ok {
		t.Fatal("preserved state should stay an object")
	}
	if !preserved.Attributes()["secret"].Equal(types.StringValue("request-only")) {
		t.Fatalf("request-only field was not preserved: %#v", preserved.Attributes()["secret"])
	}
	if !preserved.Attributes()["name"].Equal(types.StringValue("alpha")) {
		t.Fatal("readable fields should keep the observed value")
	}
}

func TestUnitWireValuesEqualAndMappingPath(t *testing.T) {
	if !wire.ValuesEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1}) {
		t.Fatal("equal JSON values should compare equal")
	}
	if wire.ValuesEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 2}) {
		t.Fatal("different JSON values must not compare equal")
	}
	path := wire.MappingPath("$.Items[*].Key")
	if len(path) != 3 || path[0] != "Items" || path[2] != "Key" {
		t.Fatalf("mapping path parsed as %#v", path)
	}
}
