package coverage

import (
	"strings"
	"testing"
)

// Framework resources with only scalar attributes have no configurable ordered
// collection. The synthesizer must read the real Schema and reach an empty
// candidate set (evidence that the absence is real, not a skip), and Check
// then reports nothing to cover.
const frameworkScalarsSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type scalarsResource struct{}

func (r *scalarsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}
`

// A framework resource with a ListNestedBlock has a configurable ordered
// collection. The synthesizer must discover it as a candidate so the
// collection-order check applies (fail-closed), instead of skipping the
// framework resource as "outside scope".
const frameworkListNestedBlockSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type ruleResource struct{}

func (r *ruleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}
`

// A Schema() method that does not declare a schema.Schema literal is a parse
// failure: the synthesizer fail-closes rather than guessing an empty schema.
const frameworkNoLiteralSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type noLiteralResource struct{}

func (r *noLiteralResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = computeSchema()
}

func computeSchema() schema.Schema {
	return schema.Schema{}
}
`

// An attribute kind the synthesizer does not model is reported, not swallowed.
const frameworkUnknownKindSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type unknownKindResource struct{}

func (r *unknownKindResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"x": schema.WeirdAttribute{},
		},
	}
}
`

// F6 positive: scalar-only framework resource has no candidate collection and
// Check reports nothing to cover.
func TestFrameworkResource_ScalarOnlyHasNoCandidates(t *testing.T) {
	r, err := FrameworkResourceSource("scalars.go", []byte(frameworkScalarsSrc))
	if err != nil {
		t.Fatalf("FrameworkResourceSource: %v", err)
	}
	if got := Candidates(r); len(got) != 0 {
		t.Fatalf("scalar-only framework resource should have no candidates, got %v", got)
	}
	if issues := Check("scalars_test.go", []byte(`package ecs_test`), "alicloud_scalars", r); len(issues) != 0 {
		t.Fatalf("scalar-only resource should have nothing to cover, got %v", issues)
	}
}

// F6 positive: a ListNestedBlock is discovered as a candidate, proving the
// synthesizer reads the real block schema and does not skip framework blocks.
func TestFrameworkResource_ListNestedBlockIsCandidate(t *testing.T) {
	r, err := FrameworkResourceSource("rule.go", []byte(frameworkListNestedBlockSrc))
	if err != nil {
		t.Fatalf("FrameworkResourceSource: %v", err)
	}
	got := Candidates(r)
	if !got["rule"] {
		t.Fatalf("ListNestedBlock 'rule' must be a candidate; got %v", got)
	}
}

// F6 negative: a framework resource with a ListNestedBlock but a test file that
// never applies a convergent member reorder is reported as uncovered
// (fail-closed), proving missing coverage fails rather than being silently
// skipped.
const ruleNoReorderTestSrc = `package ecs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{Config: testAccRuleConfig()},
		},
	})
}

func testAccRuleConfig() string {
	return "resource \"alicloud_rule\" \"default\" {}"
}
`

func TestFrameworkResource_MissingReorderCoverageFails(t *testing.T) {
	r, err := FrameworkResourceSource("rule.go", []byte(frameworkListNestedBlockSrc))
	if err != nil {
		t.Fatalf("FrameworkResourceSource: %v", err)
	}
	issues := Check("rule_test.go", []byte(ruleNoReorderTestSrc), "alicloud_rule", r)
	if len(issues) == 0 {
		t.Fatalf("a ListNestedBlock with no reorder coverage must fail; got no issues")
	}
	found := false
	for _, issue := range issues {
		if strings.Contains(issue, "rule") {
			found = true
		}
	}
	if !found {
		t.Fatalf("the uncovered candidate must be named in the issue; got %v", issues)
	}
}

// F6 negative: a Schema() method that does not declare a schema.Schema literal
// fails closed instead of synthesizing an empty schema.
func TestFrameworkResource_NoSchemaLiteralFailsClosed(t *testing.T) {
	if _, err := FrameworkResourceSource("noliteral.go", []byte(frameworkNoLiteralSrc)); err == nil {
		t.Fatalf("a Schema() method with no schema.Schema literal must fail closed")
	}
}

// F6 negative: an unrecognized attribute kind fails closed.
func TestFrameworkResource_UnknownKindFailsClosed(t *testing.T) {
	if _, err := FrameworkResourceSource("unknown.go", []byte(frameworkUnknownKindSrc)); err == nil {
		t.Fatalf("an unrecognized attribute kind must fail closed")
	}
}

// F6 regression: the scalar-only resource must be indistinguishable from a real
// schema with no ordered collection — no candidate, no suppressed paths.
func TestFrameworkResource_ScalarOnlyHasNoSuppressed(t *testing.T) {
	r, err := FrameworkResourceSource("scalars.go", []byte(frameworkScalarsSrc))
	if err != nil {
		t.Fatalf("FrameworkResourceSource: %v", err)
	}
	if got := Suppressed(r); len(got) != 0 {
		t.Fatalf("scalar-only resource should have no suppressed paths, got %v", got)
	}
}

// Real-source integration: synthesize the actual SG snapshot policy resource this
// change ships, verifying the flag-faithfulness fix (1a) end-to-end on the real
// schema rather than an in-memory fixture. Computed-only outputs (id, policy_id,
// region_id) stay Computed without a forced Optional; the Optional+Computed
// updatable fields (name, status) keep both flags. The resource is scalar-only, so
// the candidate set is empty: the collection-order gate reaches it by reading its
// real schema, not by a generic skip. The path is relative to this package dir.
func TestFrameworkResource_RealSecurityGroupSnapshotPolicy(t *testing.T) {
	r, err := FrameworkResource("../../../../alicloud/service/ecs/resource_security_group_snapshot_policy.go")
	if err != nil {
		t.Fatalf("synthesizing the real resource must succeed, not fail: %v", err)
	}
	want := map[string]struct {
		required, optional, computed bool
	}{
		"id":                                  {false, false, true},
		"security_group_snapshot_policy_id":   {false, false, true},
		"security_group_snapshot_policy_name": {false, true, true},
		"status":                              {false, true, true},
		"region_id":                           {false, false, true},
	}
	for name, w := range want {
		s, ok := r.Schema[name]
		if !ok {
			t.Fatalf("real resource missing attribute %q", name)
		}
		if s.Required != w.required || s.Optional != w.optional || s.Computed != w.computed {
			t.Fatalf("real attribute %q flags = (Required=%v Optional=%v Computed=%v); want (%v %v %v)",
				name, s.Required, s.Optional, s.Computed, w.required, w.optional, w.computed)
		}
	}
	if got := Candidates(r); len(got) != 0 {
		t.Fatalf("the scalar-only SG snapshot policy resource must have no candidates, got %v", got)
	}
}
