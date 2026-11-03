package polardb

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestUnitPolarDBClusterTDEAutomaticRotationValue covers the presence check
// behind EnableAutomaticRotation. The argument is input only, so the merged
// update state cannot tell "explicitly false" from "absent"; the raw
// configuration decides, and a removed argument is submitted through the
// planned value on a non-create update.
func TestUnitPolarDBClusterTDEAutomaticRotationValue(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]interface{}
		changes     map[string]interface{}
		rawConfig   map[string]cty.Value
		wantValue   bool
		wantPresent bool
	}{
		{
			name:        "not configured",
			attributes:  map[string]interface{}{"tde_status": "Enabled"},
			changes:     map[string]interface{}{"encrypt_new_tables": "ON"},
			rawConfig:   map[string]cty.Value{"encrypt_new_tables": cty.StringVal("ON")},
			wantPresent: false,
		},
		{
			name:        "explicit true",
			attributes:  map[string]interface{}{"tde_status": "Enabled"},
			changes:     map[string]interface{}{"enable_automatic_rotation": true},
			rawConfig:   map[string]cty.Value{"enable_automatic_rotation": cty.True},
			wantValue:   true,
			wantPresent: true,
		},
		{
			name:        "explicit false",
			attributes:  map[string]interface{}{"tde_status": "Enabled", "enable_automatic_rotation": true},
			changes:     map[string]interface{}{"enable_automatic_rotation": false},
			rawConfig:   map[string]cty.Value{"enable_automatic_rotation": cty.False},
			wantValue:   false,
			wantPresent: true,
		},
		{
			name:        "removed from configuration",
			attributes:  map[string]interface{}{"tde_status": "Enabled", "enable_automatic_rotation": true},
			changes:     map[string]interface{}{"enable_automatic_rotation": false},
			rawConfig:   map[string]cty.Value{},
			wantValue:   false,
			wantPresent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTdeUnitTestResource(t, tt.attributes, tt.changes, tt.rawConfig)
			value, present := tdeAutomaticRotationRequestValue(d)
			if value != tt.wantValue || present != tt.wantPresent {
				t.Fatalf("expected (%t, %t), got (%t, %t)", tt.wantValue, tt.wantPresent, value, present)
			}
		})
	}
}

func newTdeUnitTestResource(t *testing.T, attributes, changes map[string]interface{}, configured map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	resourceSchema := ResourceCluster().Schema
	stateData, err := schema.InternalMap(resourceSchema).Data(&sdkterraform.InstanceState{
		ID:         "pc-unit-test",
		Attributes: map[string]string{"id": "pc-unit-test"},
	}, nil)
	if err != nil {
		t.Fatalf("failed to build the unit test state: %s", err)
	}
	for key, value := range attributes {
		if err := stateData.Set(key, value); err != nil {
			t.Fatalf("failed to set %s in the unit test state: %s", key, err)
		}
	}
	state := stateData.State()
	if state == nil {
		t.Fatal("failed to build the unit test instance state")
	}
	oldData, err := schema.InternalMap(resourceSchema).Data(state, nil)
	if err != nil {
		t.Fatalf("failed to build the prior unit test resource data: %s", err)
	}
	newData, err := schema.InternalMap(resourceSchema).Data(state, nil)
	if err != nil {
		t.Fatalf("failed to build the changed unit test resource data: %s", err)
	}
	for key, value := range changes {
		if err := newData.Set(key, value); err != nil {
			t.Fatalf("failed to set %s in the changed unit test resource data: %s", key, err)
		}
	}
	diff := sdkterraform.NewInstanceDiff()
	for _, key := range helper.MapSort(newData.State().Attributes) {
		newValue := newData.State().Attributes[key]
		oldValue, ok := oldData.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: oldValue, New: newValue}
	}
	diff.RawConfig = tdeUnitTestRawConfig(t, configured)
	d, err := schema.InternalMap(resourceSchema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build the diffed unit test resource data: %s", err)
	}
	return d
}

// tdeUnitTestRawConfig builds the configuration Terraform sends: every schema
// attribute is present with its type and only the listed ones carry a value,
// so unlisted attributes are null exactly like omitted arguments.
func tdeUnitTestRawConfig(t *testing.T, configured map[string]cty.Value) cty.Value {
	t.Helper()
	impliedType := ResourceCluster().CoreConfigSchema().ImpliedType()
	values := make(map[string]cty.Value, len(impliedType.AttributeTypes()))
	for name, attributeType := range impliedType.AttributeTypes() {
		values[name] = cty.NullVal(attributeType)
	}
	for name, value := range configured {
		if _, ok := values[name]; !ok {
			t.Fatalf("raw config attribute %q is not part of the resource schema", name)
		}
		values[name] = value
	}
	return cty.ObjectVal(values)
}
