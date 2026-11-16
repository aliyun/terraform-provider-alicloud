package expressconnectrouter

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
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

func TestUnitVpcAssociationPrefixOrderPrivateBoundaries(t *testing.T) {
	if suppressVpcAssociationPrefixReorder("allowed_prefixes.#", "172.16.1.0/24", "172.16.1.0/24", nil) {
		t.Fatal("a nil ResourceData carries no old/new pair and must not suppress")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{"172.16.1.0/24", 42}); ok {
		t.Fatal("a list with a non-string member must not be treated as a known prefix list")
	}
}

// The order-coverage checker requires a direct named-suppressor net in addition to the public Diff regression.
func TestUnitVpcAssociationPrefixOrderSuppressNet(t *testing.T) {
	for _, tc := range []struct {
		name     string
		oldValue string
		newValue string
		expected bool
	}{
		{"ReorderOnly", "172.16.1.0/24,172.16.2.0/24", "172.16.2.0/24,172.16.1.0/24", true},
		{"MemberAdded", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := newVpcAssociationPrefixDiffData(t, tc.oldValue, tc.newValue)
			if got := suppressVpcAssociationPrefixReorder("allowed_prefixes.#", tc.oldValue, tc.newValue, d); got != tc.expected {
				t.Fatalf("suppress(%q, %q) = %t, want %t", tc.oldValue, tc.newValue, got, tc.expected)
			}
		})
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
	newMembers := strings.Split(newValue, ",")
	state := &terraform.InstanceState{ID: "ecr-1:vpc-1", Attributes: map[string]string{
		"allowed_prefixes.#": strconv.Itoa(len(oldMembers)),
	}}
	for i, member := range oldMembers {
		state.Attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = member
	}
	rawMembers := make([]cty.Value, 0, len(newMembers))
	for _, member := range newMembers {
		rawMembers = append(rawMembers, cty.StringVal(member))
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{"allowed_prefixes": cty.ListVal(rawMembers)})
	state.RawConfig = rawConfig
	diff := &terraform.InstanceDiff{
		RawConfig: rawConfig,
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"allowed_prefixes.#": {Old: strconv.Itoa(len(oldMembers)), New: strconv.Itoa(len(newMembers))},
		},
	}
	for i := 0; i < len(oldMembers) || i < len(newMembers); i++ {
		attribute := &terraform.ResourceAttrDiff{}
		if i < len(oldMembers) {
			attribute.Old = oldMembers[i]
		}
		if i < len(newMembers) {
			attribute.New = newMembers[i]
		} else {
			attribute.NewRemoved = true
		}
		diff.Attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = attribute
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build prefix resource data: %v", err)
	}
	return d
}
