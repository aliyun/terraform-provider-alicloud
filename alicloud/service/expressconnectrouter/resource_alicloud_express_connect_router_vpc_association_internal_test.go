package expressconnectrouter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zclconf/go-cty/cty"
)

func TestUnitVpcAssociationUpdateBudgetIsShared(t *testing.T) {
	deadline := time.Unix(1_000, 0)
	for _, tc := range []struct {
		name string
		now  time.Time
		want time.Duration
	}{
		{name: "first_step_gets_full_window", now: time.Unix(990, 0), want: 10 * time.Second},
		{name: "second_step_only_gets_remainder", now: time.Unix(997, 0), want: 3 * time.Second},
		{name: "exhausted_budget_is_zero", now: time.Unix(1_005, 0), want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := remainingVpcAssociationTimeout(deadline, tc.now); got != tc.want {
				t.Fatalf("remaining timeout = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUnitVpcAssociationPrefixOrderPrivateBoundaries(t *testing.T) {
	if suppressVpcAssociationPrefixReorder("allowed_prefixes.#", "172.16.1.0/24", "172.16.1.0/24", nil) {
		t.Fatal("a nil ResourceData carries no old/new pair and must not suppress")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{"172.16.1.0/24", 42}); ok {
		t.Fatal("a list with a non-string member must not be treated as a known prefix list")
	}
}

func TestUnitVpcAssociationPrefixOrderSuppressNet(t *testing.T) {
	testCases := []struct {
		name     string
		oldValue string
		newValue string
		expected bool
	}{
		{name: "ReorderOnly", oldValue: "172.16.1.0/24,172.16.2.0/24", newValue: "172.16.2.0/24,172.16.1.0/24", expected: true},
		{name: "MemberAdded", oldValue: "172.16.1.0/24,172.16.2.0/24", newValue: "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", expected: false},
		{name: "MemberReplaced", oldValue: "172.16.1.0/24,172.16.2.0/24", newValue: "172.16.1.0/24,172.16.9.0/24", expected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := vpcAssociationPrefixSuppressData(t, tc.oldValue, tc.newValue)
			if got := suppressVpcAssociationPrefixReorder("allowed_prefixes", tc.oldValue, tc.newValue, data); got != tc.expected {
				t.Fatalf("case %s: suppressor = %v, want %v", tc.name, got, tc.expected)
			}
		})
	}
}

func vpcAssociationPrefixSuppressData(t *testing.T, oldValue, newValue string) *schema.ResourceData {
	t.Helper()
	oldPrefixes := strings.Split(oldValue, ",")
	attributes := map[string]string{
		"id":                    "ecr-1:assoc-1:vpc-1",
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_id":        "assoc-1",
		"association_region_id": "cn-hangzhou",
		"allowed_prefixes.#":    strconv.Itoa(len(oldPrefixes)),
		"allowed_prefixes_mode": "MatchMode",
	}
	for i, prefix := range oldPrefixes {
		attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = prefix
	}
	newPrefixes := make([]interface{}, 0, len(strings.Split(newValue, ",")))
	for _, prefix := range strings.Split(newValue, ",") {
		newPrefixes = append(newPrefixes, prefix)
	}
	config := map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
		"allowed_prefixes":      newPrefixes,
	}
	resource := ResourceVpcAssociation()
	suppress := resource.Schema["allowed_prefixes"].DiffSuppressFunc
	var data *schema.ResourceData
	resource.Schema["allowed_prefixes"].DiffSuppressFunc = func(key, old, new string, d *schema.ResourceData) bool {
		data = d
		return suppress(key, old, new, d)
	}
	rawPrefixes := make([]cty.Value, 0, len(newPrefixes))
	for _, prefix := range newPrefixes {
		rawPrefixes = append(rawPrefixes, cty.StringVal(prefix.(string)))
	}
	state := &terraform.InstanceState{ID: "ecr-1:assoc-1:vpc-1", Attributes: attributes}
	state.RawConfig = cty.ObjectVal(map[string]cty.Value{
		"allowed_prefixes":      cty.ListVal(rawPrefixes),
		"allowed_prefixes_mode": cty.NullVal(cty.String),
		"association_id":        cty.NullVal(cty.String),
		"association_region_id": cty.StringVal("cn-hangzhou"),
		"create_time":           cty.NullVal(cty.String),
		"description":           cty.NullVal(cty.String),
		"ecr_id":                cty.StringVal("ecr-1"),
		"status":                cty.NullVal(cty.String),
		"vpc_id":                cty.StringVal("vpc-1"),
		"vpc_owner_id":          cty.NullVal(cty.Number),
	})
	if _, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil); err != nil {
		t.Fatalf("build suppressor resource data: %v", err)
	}
	if data == nil {
		t.Fatal("prefix diff did not invoke the suppressor")
	}
	return data
}
