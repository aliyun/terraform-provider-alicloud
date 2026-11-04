package fwadapt_test

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// semanticString lives in the test package on purpose: it proves a type outside fwadapt can
// satisfy the exported semantic hooks. The unexported hooks this replaced could never be
// implemented across a Go package boundary.
type semanticString struct{ types.String }

func (v semanticString) SemanticEqual(other attr.Value) bool {
	value, ok := other.(semanticString)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false
	}
	return strings.EqualFold(v.ValueString(), value.ValueString())
}

func (v semanticString) NullValue() attr.Value {
	return semanticString{String: types.StringNull()}
}

var (
	_ fwadapt.SemanticValue = semanticString{}
	_ fwadapt.NullValue     = semanticString{}
)

func TestUnitFwadaptSemanticHooksCrossPackage(t *testing.T) {
	semanticA := semanticString{String: types.StringValue("{\"a\": 1}")}
	semanticB := semanticString{String: types.StringValue("{\"A\": 1}")}
	if fwadapt.PlannedChanged(semanticA, semanticB) {
		t.Fatal("semantic equality should treat case-only differences as unchanged")
	}
	if !fwadapt.ConfiguredMatch(semanticA, semanticB) {
		t.Fatal("configured match should use the semantic hook")
	}
	plainA := types.StringValue("{\"a\": 1}")
	plainB := types.StringValue("{\"A\": 1}")
	if !fwadapt.PlannedChanged(plainA, plainB) {
		t.Fatal("plain strings should still compare with Equal")
	}

	unknown := semanticString{String: types.StringUnknown()}
	known := fwadapt.KnownValue(unknown)
	if known == nil || !known.IsNull() {
		t.Fatalf("unknown custom value should become a typed null, got %#v", known)
	}
	if _, ok := known.(semanticString); !ok {
		t.Fatalf("typed null should keep the custom type, got %T", known)
	}
	if got := fwadapt.KnownValue(types.StringUnknown()); !got.IsNull() {
		t.Fatal("unknown built-in string should still become a null string")
	}
}

func TestUnitFwadaptSelectReadObject(t *testing.T) {
	rows := []interface{}{
		map[string]interface{}{"Namespace": "alpha"},
		map[string]interface{}{"Namespace": "beta"},
	}
	selected, err := fwadapt.SelectReadObject(rows, []fwadapt.ReadIdentity{{Path: "$.Namespace", Value: "beta"}})
	if err != nil {
		t.Fatalf("selecting an existing identity failed: %s", err)
	}
	if selected == nil || selected["Namespace"] != "beta" {
		t.Fatalf("selected the wrong row: %#v", selected)
	}

	missing, err := fwadapt.SelectReadObject(rows, []fwadapt.ReadIdentity{{Path: "$.Namespace", Value: "gamma"}})
	if err != nil || missing != nil {
		t.Fatalf("a missing identity should yield nil without error, got %#v / %v", missing, err)
	}

	duplicated := []interface{}{
		map[string]interface{}{"Namespace": "beta"},
		map[string]interface{}{"Namespace": "beta"},
	}
	if _, err := fwadapt.SelectReadObject(duplicated, []fwadapt.ReadIdentity{{Path: "$.Namespace", Value: "beta"}}); err == nil {
		t.Fatal("multiple matching identities must be an error")
	}

	unreadable := []interface{}{map[string]interface{}{"Other": "value"}}
	if _, err := fwadapt.SelectReadObject(unreadable, []fwadapt.ReadIdentity{{Path: "$.Namespace", Value: "beta"}}); err == nil {
		t.Fatal("a list item without a readable identity must be an error")
	}

	if _, err := fwadapt.SelectReadObject(rows, nil); err == nil {
		t.Fatal("reading without a declared identity must be an error")
	}
}

func TestUnitFwadaptApplyMappingsAndPreserveRequestOnlyState(t *testing.T) {
	request := map[string]interface{}{}
	values := map[string]interface{}{"$.Name": "alpha"}
	if err := fwadapt.ApplyMappings(request, values, []fwadapt.Mapping{{Source: "$.Name", Target: "Name"}}, false); err != nil {
		t.Fatalf("applying a scalar mapping failed: %s", err)
	}
	if request["Name"] != "alpha" {
		t.Fatalf("scalar mapping wrote %#v", request["Name"])
	}

	flattened := map[string]interface{}{}
	flattenValues := map[string]interface{}{"$.Deep": map[string]interface{}{"Key": "value"}}
	bindings := []fwadapt.Mapping{{Source: "$.Deep.Key", Target: "Deep.Key", Repeat: true}}
	if err := fwadapt.ApplyMappings(flattened, flattenValues, bindings, false); err != nil {
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
	fields := map[string]fwadapt.WireField{
		"secret": {Name: "secret", Read: false},
	}
	preserved, ok := fwadapt.PreserveRequestValues(observed, prior, fields).(types.Object)
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

func TestUnitFwadaptValuesEqualAndMappingPath(t *testing.T) {
	if !fwadapt.ValuesEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1}) {
		t.Fatal("equal JSON values should compare equal")
	}
	if fwadapt.ValuesEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 2}) {
		t.Fatal("different JSON values must not compare equal")
	}
	path := fwadapt.MappingPath("$.Items[*].Key")
	if len(path) != 3 || path[0] != "Items" || path[2] != "Key" {
		t.Fatalf("mapping path parsed as %#v", path)
	}
}