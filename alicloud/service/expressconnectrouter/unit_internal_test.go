package expressconnectrouter

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

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
			config := map[string]interface{}{
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_region_id": "cn-hangzhou",
				"allowed_prefixes":      tc.newValue,
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
