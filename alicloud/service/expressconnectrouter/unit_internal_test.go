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

func TestUnitVpcAssociationPrefixOrderSuppress(t *testing.T) {
	for _, tc := range []struct {
		name       string
		oldValue   []string
		newValue   []string
		suppressed bool
	}{
		{name: "ReorderOnly", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.2.0/24", "172.16.1.0/24"}, suppressed: true},
		{name: "LongerReorder", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"}, newValue: []string{"172.16.3.0/24", "172.16.1.0/24", "172.16.2.0/24"}, suppressed: true},
		{name: "MemberAdded", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"}},
		{name: "MemberRemoved", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24"}},
		{name: "MemberReplaced", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.9.0/24"}},
		{name: "MultiplicityChanged", oldValue: []string{"172.16.1.0/24", "172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.2.0/24"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attributes := map[string]string{
				"id":                    "ecr-1:assoc-1:vpc-1",
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_id":        "assoc-1",
				"association_region_id": "cn-hangzhou",
				"allowed_prefixes.#":    strconv.Itoa(len(tc.oldValue)),
				"allowed_prefixes_mode": "MatchMode",
			}
			for i, prefix := range tc.oldValue {
				attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = prefix
			}
			state := &terraform.InstanceState{ID: "ecr-1:assoc-1:vpc-1", Attributes: attributes}
			newPrefixes := make([]interface{}, 0, len(tc.newValue))
			for _, prefix := range tc.newValue {
				newPrefixes = append(newPrefixes, prefix)
			}
			config := map[string]interface{}{
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_region_id": "cn-hangzhou",
				"allowed_prefixes":      newPrefixes,
			}
			diff, err := ResourceVpcAssociation().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("build prefix diff: %v", err)
			}
			changed := diff != nil && !diff.Empty()
			if changed == tc.suppressed {
				t.Fatalf("changed = %t, want %t for old=%v new=%v", changed, !tc.suppressed, tc.oldValue, tc.newValue)
			}
		})
	}

	for _, tc := range []struct {
		name     string
		oldValue string
		newValue string
		expected bool
	}{
		{"ReorderOnly", "172.16.1.0/24,172.16.2.0/24", "172.16.2.0/24,172.16.1.0/24", true},
		{"LongerReorder", "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", "172.16.3.0/24,172.16.1.0/24,172.16.2.0/24", true},
		{"MemberAdded", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", false},
		{"MemberRemoved", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24", false},
		{"MemberReplaced", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.9.0/24", false},
		{"MultiplicityChanged", "172.16.1.0/24,172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.2.0/24,172.16.2.0/24", false},
	} {
		t.Run("direct_"+tc.name, func(t *testing.T) {
			d := newVpcAssociationPrefixDiffData(t, tc.oldValue, tc.newValue)
			if got := suppressVpcAssociationPrefixReorder("allowed_prefixes.#", tc.oldValue, tc.newValue, d); got != tc.expected {
				t.Fatalf("suppress(%q, %q) = %t, want %t", tc.oldValue, tc.newValue, got, tc.expected)
			}
		})
	}

	if suppressVpcAssociationPrefixReorder("allowed_prefixes.#", "172.16.1.0/24", "172.16.1.0/24", nil) {
		t.Fatal("a nil ResourceData carries no old/new pair and must not suppress")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{"172.16.1.0/24", 42}); ok {
		t.Fatal("a list with a non-string member must not be treated as a known prefix list")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{}); ok {
		t.Fatal("an empty list must not be treated as a known prefix list")
	}
}

func newVpcAssociationPrefixDiffData(t *testing.T, oldValue, newValue string) *schema.ResourceData {
	t.Helper()
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"allowed_prefixes": {
			Type:             schema.TypeList,
			Optional:         true,
			Elem:             &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: suppressVpcAssociationPrefixReorder,
		},
	}}
	oldMembers := strings.Split(oldValue, ",")
	state := &terraform.InstanceState{ID: "ecr-1:vpc-1", Attributes: map[string]string{
		"allowed_prefixes.#": strconv.Itoa(len(oldMembers)),
	}}
	for i, member := range oldMembers {
		state.Attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = member
	}
	newMembers := make([]interface{}, 0, len(strings.Split(newValue, ",")))
	for _, member := range strings.Split(newValue, ",") {
		newMembers = append(newMembers, member)
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{
		"allowed_prefixes": newMembers,
	}), nil)
	if err != nil {
		t.Fatalf("build prefix diff: %v", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build prefix resource data: %v", err)
	}
	return d
}
