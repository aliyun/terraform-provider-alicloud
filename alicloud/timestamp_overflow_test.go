package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestUnitTimestampStringReadPathNoOverflow is the 32-bit (linux/386, linux/arm) regression
// guard for the integer-overflow panic in timestamp attributes.
//
// Background: schema.TypeInt reads raw config/state through the SDK field reader via
// strconv.ParseInt(value, 0, 0) followed by int(v) (helper/schema field_reader.go). On
// 386/arm, where int is 32-bit, that panics on 13-digit millisecond and 19-digit nanosecond
// epochs. The fix changes these timestamp attributes to schema.TypeString; the TypeString read
// path returns the raw string verbatim with no integer conversion, so it is arch-invariant
// and cannot overflow.
//
// This test feeds a 13-digit (ms) and a 19-digit (ns) value through the SDK raw-config read
// path (schema.TestResourceDataRaw -> ConfigFieldReader) for each fixed field and asserts the
// value round-trips as the exact string with no truncation and no panic. The branch selected
// (TypeString vs TypeInt) is determined by the schema, not by GOARCH, so passing on amd64
// proves the same non-ParseInt branch runs on 386/arm — the panic-prone TypeInt branch is no
// longer reachable for these fields.
func TestUnitTimestampStringReadPathNoOverflow(t *testing.T) {
	cases := []struct {
		name  string
		s     map[string]*schema.Schema
		field string
	}{
		{"ots_tunnel.create_time", resourceAlicloudOtsTunnel().Schema, "create_time"},
		{"ots_tunnel.tunnel_rpo", resourceAlicloudOtsTunnel().Schema, "tunnel_rpo"},
		{"oss_scan_config.last_modified_start_time", resourceAliCloudThreatDetectionOssScanConfig().Schema, "last_modified_start_time"},
	}
	// 13-digit millisecond epoch and 19-digit nanosecond epoch — both overflow int32, the
	// latter also overflows int64-into-int32 truncation that caused the original panic.
	bigValues := []string{"1699999999999", "1699999999999999999"}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fi, ok := c.s[c.field]
			if !ok {
				t.Fatalf("field %q not in schema", c.field)
			}
			if fi.Type != schema.TypeString {
				t.Fatalf("%s: expected schema.TypeString (the 386 overflow fix), got %v", c.name, fi.Type)
			}
			for _, raw := range bigValues {
				d := schema.TestResourceDataRaw(t, c.s, map[string]interface{}{c.field: raw})
				if got := d.Get(c.field); got != raw {
					t.Fatalf("%s: raw %q read back as %#v; the TypeString read path must return the string verbatim with no integer conversion", c.name, raw, got)
				}
			}
		})
	}
}

// TestUnitTimestampStringFieldValidateFunc guards the F7 restoration of plan-phase numeric
// validation. After TypeInt -> TypeString the original plan-phase ParseInt validation was
// lost; the TypeString timestamp fields now carry a ValidateFunc that accepts numeric input
// and rejects non-numeric input, without re-introducing the platform-width-int read path.
func TestUnitTimestampStringFieldValidateFunc(t *testing.T) {
	cases := []struct {
		name  string
		field *schema.Schema
		key   string
	}{
		{
			"oss_scan_config.last_modified_start_time",
			resourceAliCloudThreatDetectionOssScanConfig().Schema["last_modified_start_time"],
			"last_modified_start_time",
		},
		{
			"wafv3_defense_rule.time_periods.start",
			wafv3TimePeriodsSchema()["start"],
			"start",
		},
	}
	valid := []string{"0", "1699999999999", "1699999999999999999"}
	invalid := []string{"abc", "1699abc", "2024-01-01"}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.field == nil {
				t.Fatalf("%s: field not in schema", c.name)
			}
			if c.field.Type != schema.TypeString {
				t.Fatalf("%s: expected schema.TypeString, got %v", c.name, c.field.Type)
			}
			if c.field.ValidateFunc == nil {
				t.Fatalf("%s: expected a ValidateFunc restoring plan-phase numeric validation (F7)", c.name)
			}
			for _, v := range valid {
				if _, errs := c.field.ValidateFunc(v, c.key); len(errs) > 0 {
					t.Fatalf("%s: expected %q to be valid, got %v", c.name, v, errs)
				}
			}
			for _, v := range invalid {
				if _, errs := c.field.ValidateFunc(v, c.key); len(errs) == 0 {
					t.Fatalf("%s: expected %q to be rejected by the numeric ValidateFunc", c.name, v)
				}
			}
		})
	}
}

// wafv3TimePeriodsSchema navigates the nested config.time_config.time_periods schema of
// resourceAlicloudWafv3DefenseRule so tests can assert the type/ValidateFunc of its start/end
// fields without exercising the cloud read path.
func wafv3TimePeriodsSchema() map[string]*schema.Schema {
	r := resourceAliCloudWafv3DefenseRule()
	configSchema := r.Schema["config"].Elem.(*schema.Resource).Schema
	tcSchema := configSchema["time_config"].Elem.(*schema.Resource).Schema
	return tcSchema["time_periods"].Elem.(*schema.Resource).Schema
}

// TestUnitTimestampStringFlatmapStateReadNoOverflow guards the SDK *state* (InstanceState
// flatmap) read path — the actual 32-bit panic site. On 386/arm, schema.TypeInt reads a
// flatmap state value through strconv.ParseInt(value, 0, 0) followed by int(v)
// (helper/schema field_reader.go), which panics when a 13-digit millisecond or 19-digit
// nanosecond epoch overflows native int. After the fix these timestamp attributes are
// schema.TypeString; the state read path returns the raw attribute string verbatim with no
// integer conversion, so it is arch-invariant and cannot overflow.
//
// schema.TestResourceDataRaw (the config-read guard above) exercises ConfigFieldReader only;
// this test constructs a synthetic InstanceState.Attributes flatmap carrying the big epochs
// and reads via Resource.Data(state) -> StateFieldReader -> stringToPrimitive, asserting the
// TypeString branch returns the exact string. Passing on amd64 proves the same non-ParseInt
// branch runs on 386/arm — the panic-prone TypeInt branch is no longer reachable for these
// fields.
func TestUnitTimestampStringFlatmapStateReadNoOverflow(t *testing.T) {
	cases := []struct {
		name  string
		res   *schema.Resource
		field string
	}{
		{"ots_tunnel.create_time", resourceAlicloudOtsTunnel(), "create_time"},
		{"ots_tunnel.tunnel_rpo", resourceAlicloudOtsTunnel(), "tunnel_rpo"},
		{"oss_scan_config.last_modified_start_time", resourceAliCloudThreatDetectionOssScanConfig(), "last_modified_start_time"},
	}
	// 13-digit millisecond epoch and 19-digit nanosecond epoch — both overflow int32, the
	// latter also overflows the int64-into-int32 truncation that caused the original panic.
	bigValues := []string{"1699999999999", "1699999999999999999"}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fi, ok := c.res.Schema[c.field]
			if !ok {
				t.Fatalf("field %q not in schema", c.field)
			}
			if fi.Type != schema.TypeString {
				t.Fatalf("%s: expected schema.TypeString (the 386 overflow fix), got %v", c.name, fi.Type)
			}
			for _, raw := range bigValues {
				state := &sdktf.InstanceState{
					ID: "test",
					Attributes: map[string]string{
						"id":    "test",
						c.field: raw,
					},
				}
				d := c.res.Data(state)
				if got := d.Get(c.field); got != raw {
					t.Fatalf("%s: flatmap state raw %q read back as %#v; the TypeString state read path must return the string verbatim with no integer conversion", c.name, raw, got)
				}
			}
		})
	}
}
