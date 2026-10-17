package fwschema

import (
	"strings"
	"testing"
)

// 1a: the synthesizer must read Required/Optional/Computed faithfully instead of
// forcing Optional+Computed on every attribute kind. A Required attribute stays
// Required, a Computed-only (output) attribute stays Computed without Optional,
// and a plain Optional attribute is not reported Computed. Forcing Optional+
// Computed would hide a Required or a read-only field from the attribute-coverage
// and collection-order checkers.
const flagsSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type flagsResource struct{}

func (r *flagsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true},
			"name":  schema.StringAttribute{Required: true},
			"opt":   schema.StringAttribute{Optional: true, Computed: true},
			"plain": schema.StringAttribute{Optional: true},
		},
	}
}
`

func TestUnitFrameworkResource_FlagsAreFaithful(t *testing.T) {
	r, err := FrameworkResourceSource("flags.go", []byte(flagsSrc))
	if err != nil {
		t.Fatalf("FrameworkResourceSource: %v", err)
	}
	cases := []struct {
		name                         string
		required, optional, computed bool
	}{
		{"id", false, false, true},    // Computed-only: must not gain Optional.
		{"name", true, false, false},  // Required: must not become Optional+Computed.
		{"opt", false, true, true},    // Optional+Computed: both preserved.
		{"plain", false, true, false}, // Optional-only: must not gain Computed.
	}
	for _, c := range cases {
		s, ok := r.Schema[c.name]
		if !ok {
			t.Fatalf("missing attribute %q", c.name)
		}
		if s.Required != c.required || s.Optional != c.optional || s.Computed != c.computed {
			t.Fatalf("attribute %q flags = (Required=%v Optional=%v Computed=%v); want (%v %v %v)",
				c.name, s.Required, s.Optional, s.Computed, c.required, c.optional, c.computed)
		}
	}
}

// 1b: SingleNestedAttribute has Attributes directly (no NestedObject field). A
// valid one must synthesize without error rather than failing "has no NestedObject
// composite literal", and its Optional flag must be preserved.
const singleNestedAttrSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type singleNestedResource struct{}

func (r *singleNestedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{Required: true},
				},
			},
		},
	}
}
`

func TestUnitFrameworkResource_SingleNestedAttributeIsNotRejected(t *testing.T) {
	r, err := FrameworkResourceSource("single.go", []byte(singleNestedAttrSrc))
	if err != nil {
		t.Fatalf("a valid SingleNestedAttribute must synthesize, not fail: %v", err)
	}
	s, ok := r.Schema["config"]
	if !ok {
		t.Fatalf("missing attribute config")
	}
	if !s.Optional || s.Required || s.Computed {
		t.Fatalf("SingleNestedAttribute Optional flag must be preserved, got (Required=%v Optional=%v Computed=%v)", s.Required, s.Optional, s.Computed)
	}
	// A SingleNestedAttribute is a single object, not a re-orderable collection:
	// MaxItems==1 keeps it out of the candidate set (collect skips TypeList with
	// MaxItems==1), so it is reached by reading its real schema, not skipped.
	if s.MaxItems != 1 {
		t.Fatalf("SingleNestedAttribute must synthesize MaxItems=1, got %d", s.MaxItems)
	}
}

// 1b negative: ListNestedAttribute still requires its NestedObject. A
// ListNestedAttribute written without NestedObject must fail closed rather than
// synthesizing an empty nested schema.
const listNestedAttrNoObjectSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type listNestedNoObjectResource struct{}

func (r *listNestedNoObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rules": schema.ListNestedAttribute{
				Optional: true,
			},
		},
	}
}
`

func TestUnitFrameworkResource_ListNestedAttributeWithoutNestedObjectFailsClosed(t *testing.T) {
	if _, err := FrameworkResourceSource("listnoobj.go", []byte(listNestedAttrNoObjectSrc)); err == nil {
		t.Fatalf("a ListNestedAttribute without NestedObject must fail closed")
	}
}

// #12: a bool flag written as a non-literal expression — a call like
// configurableFromAnotherFile() — is unresolved. The synthesizer must fail
// closed (return an error) rather than defaulting the flag to false, because a
// false Optional on a ListAttribute would hide a possibly-configurable ordered
// collection from Candidates. This mirrors the orchestrator's
// TestIndependentUnresolvedConfigurabilityCannotHideList real FAIL case.
const unresolvedConfigurableSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type unresolvedResource struct{}

func configurableFromAnotherFile() bool { return true }

func (r *unresolvedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"names": schema.ListAttribute{Optional: configurableFromAnotherFile()},
		},
	}
}
`

func TestUnitFrameworkResource_UnresolvedConfigurabilityFailsClosed(t *testing.T) {
	_, err := FrameworkResourceSource("unresolved.go", []byte(unresolvedConfigurableSrc))
	if err == nil {
		t.Fatalf("an unresolved bool flag (a call expression) must fail closed, not default to false and hide the collection")
	}
	if !strings.Contains(err.Error(), "Optional") {
		t.Fatalf("the error must name the unresolved flag; got %v", err)
	}
}

// #12 positive: a missing bool field defaults to false (the framework zero
// value) and an explicit false ident is read literally. Neither errors; both
// keep all flags false. The true-literal positive is covered by
// TestFrameworkResource_FlagsAreFaithful.
const falseDefaultSrc = `package ecs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type falseDefaultResource struct{}

func (r *falseDefaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"plain": schema.StringAttribute{},
			"off":   schema.StringAttribute{Optional: false, Computed: false},
		},
	}
}
`

func TestUnitFrameworkResource_MissingAndExplicitFalseDefaultFalse(t *testing.T) {
	r, err := FrameworkResourceSource("falsedefault.go", []byte(falseDefaultSrc))
	if err != nil {
		t.Fatalf("missing/explicit-false flags must not error: %v", err)
	}
	for _, name := range []string{"plain", "off"} {
		s, ok := r.Schema[name]
		if !ok {
			t.Fatalf("missing attribute %q", name)
		}
		if s.Required || s.Optional || s.Computed {
			t.Fatalf("attribute %q flags must all be false, got (Required=%v Optional=%v Computed=%v)", name, s.Required, s.Optional, s.Computed)
		}
	}
}
