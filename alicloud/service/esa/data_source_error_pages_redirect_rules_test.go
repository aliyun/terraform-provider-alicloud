package esa

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestUnitErrorPagesRedirectRulesListValueFromStrings verifies the ids/names
// output list builder over a populated slice, an empty slice and a nil slice.
// nil must map to a non-null empty list so the Computed list attribute reads
// back as [] rather than null.
func TestUnitErrorPagesRedirectRulesListValueFromStrings(t *testing.T) {
	ctx := context.Background()
	t.Run("populated", func(t *testing.T) {
		d := new(diag.Diagnostics)
		l, ok := listValueFromStrings(ctx, []string{"a:1", "a:2"}, d)
		assert.True(t, ok)
		assert.False(t, d.HasError())
		assert.False(t, l.IsNull())
		assert.Equal(t, 2, len(l.Elements()))
	})
	t.Run("empty", func(t *testing.T) {
		d := new(diag.Diagnostics)
		l, ok := listValueFromStrings(ctx, []string{}, d)
		assert.True(t, ok)
		assert.False(t, l.IsNull())
		assert.Equal(t, 0, len(l.Elements()))
	})
	t.Run("nil", func(t *testing.T) {
		d := new(diag.Diagnostics)
		l, ok := listValueFromStrings(ctx, nil, d)
		assert.True(t, ok)
		assert.False(t, l.IsNull(), "nil slice must map to a non-null empty list")
		assert.Equal(t, 0, len(l.Elements()))
	})
}

// TestUnitErrorPagesRedirectRulesListValueFromObjects verifies the rules output
// list builder: a slice of tfsdk-tagged structs is turned into a types.List of
// object values whose element type carries the nested error_pages_redirect
// list. This is the framework nested-object-list-of-objects shape the data
// source pioneers, so a passing build+run here validates that ObjectType /
// ListType / ListValueFrom compose for this schema.
func TestUnitErrorPagesRedirectRulesListValueFromObjects(t *testing.T) {
	ctx := context.Background()

	// Build one rule item with a nested one-element error_pages_redirect list.
	eprList, ldiags := types.ListValueFrom(ctx, errorPagesRedirectItemType, []errorPagesRedirectItem{
		{TargetUrl: types.StringValue("https://example.com/404.html"), StatusCode: types.StringValue("404")},
	})
	assert.Falsef(t, ldiags.HasError(), "building nested list: %v", ldiags)

	items := []errorPagesRedirectRuleDSItem{
		{
			Id:                 types.StringValue("site-1:cfg-1"),
			ConfigId:           types.StringValue("cfg-1"),
			ConfigType:         types.StringValue("rule"),
			SiteVersion:        types.Int64Value(2),
			RuleEnable:         types.StringValue("on"),
			RuleName:           types.StringValue("rule-1"),
			Rule:               types.StringValue("http.host eq \"a.com\""),
			Sequence:           types.Int64Value(10),
			ErrorPagesRedirect: eprList,
		},
	}
	d := new(diag.Diagnostics)
	l, ok := listValueFromObjects(ctx, errorPagesRedirectRuleDSItemType, items, d)
	assert.True(t, ok)
	assert.Falsef(t, d.HasError(), "unexpected diagnostics: %v", d)
	assert.False(t, l.IsNull())
	assert.Equal(t, 1, len(l.Elements()))
}

// TestUnitErrorPagesRedirectRulesStringOrEmpty pins the names-list helper: nil
// and non-string values flatten to "" so the Computed names list never holds a
// null element.
func TestUnitErrorPagesRedirectRulesStringOrEmpty(t *testing.T) {
	assert.Equal(t, "", stringOrEmpty(nil))
	assert.Equal(t, "", stringOrEmpty(""))
	assert.Equal(t, "rule-1", stringOrEmpty("rule-1"))
	assert.Equal(t, "42", stringOrEmpty(42))
}
