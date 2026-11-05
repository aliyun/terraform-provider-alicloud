package acctest

import (
	"strings"
	"testing"

	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNewInstanceDiff(t *testing.T) {
	t.Run("changed attribute is diffed and unchanged attribute is skipped", func(t *testing.T) {
		attrs := map[string]interface{}{"vpc_name": "tf-test-vpc", "cidr_block": "10.0.0.0/8"}
		diff, err := NewInstanceDiff("alicloud_vpc", attrs,
			map[string]interface{}{"vpc_name": "tf-test-vpc-updated"}, &sdktf.InstanceState{ID: "vpc-test"})
		if err != nil {
			t.Fatalf("NewInstanceDiff error: %s", err)
		}
		d, ok := diff.Attributes["vpc_name"]
		if !ok {
			t.Fatalf("vpc_name entry missing: %#v", diff.Attributes)
		}
		if d.Old != "tf-test-vpc" || d.New != "tf-test-vpc-updated" {
			t.Fatalf("vpc_name = Old %q New %q, want Old %q New %q", d.Old, d.New, "tf-test-vpc", "tf-test-vpc-updated")
		}
		if _, ok := diff.Attributes["cidr_block"]; ok {
			t.Fatalf("unchanged cidr_block must not appear in the diff: %#v", diff.Attributes)
		}
		// attributesDiff is merged into attributes by design.
		if attrs["vpc_name"] != "tf-test-vpc-updated" {
			t.Fatalf("attributesDiff must merge into attributes, got %v", attrs["vpc_name"])
		}
	})

	t.Run("identical old and new attributes produce an empty diff", func(t *testing.T) {
		diff, err := NewInstanceDiff("alicloud_vpc",
			map[string]interface{}{"vpc_name": "tf-test-vpc", "cidr_block": "10.0.0.0/8"},
			nil, &sdktf.InstanceState{ID: "vpc-test"})
		if err != nil {
			t.Fatalf("NewInstanceDiff error: %s", err)
		}
		if len(diff.Attributes) != 0 {
			t.Fatalf("expected empty diff, got: %#v", diff.Attributes)
		}
	})

	t.Run("shrunken collection keeps count and marks dropped members", func(t *testing.T) {
		diff, err := NewInstanceDiff("alicloud_instance",
			map[string]interface{}{"security_groups": []interface{}{"sg-1"}},
			map[string]interface{}{"security_groups": []interface{}{}},
			&sdktf.InstanceState{ID: "i-test"})
		if err != nil {
			t.Fatalf("NewInstanceDiff error: %s", err)
		}
		count, ok := diff.Attributes["security_groups.#"]
		if !ok {
			t.Fatalf("security_groups.# entry missing: %#v", diff.Attributes)
		}
		if count.Old != "1" || count.New != "0" {
			t.Fatalf("security_groups.# = Old %q New %q, want Old \"1\" New \"0\"", count.Old, count.New)
		}
		var dropped string
		for key, d := range diff.Attributes {
			if key != "security_groups.#" && strings.HasPrefix(key, "security_groups.") && d.New == "" {
				dropped = key
				if d.Old != "sg-1" {
					t.Fatalf("%s = Old %q, want \"sg-1\"", key, d.Old)
				}
			}
		}
		if dropped == "" {
			t.Fatalf("dropped security_groups member not diffed: %#v", diff.Attributes)
		}
	})

	t.Run("unknown attribute returns a wrapped setting error", func(t *testing.T) {
		_, err := NewInstanceDiff("alicloud_vpc",
			map[string]interface{}{"not_a_field": "x"}, nil, &sdktf.InstanceState{ID: "vpc-test"})
		if err == nil {
			t.Fatalf("expected an error for unknown attribute")
		}
		if !strings.Contains(err.Error(), "the field not_a_field setting error") {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}
