package esa

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestUnitErrorPagesRedirectRuleParseID exercises the composite ID parser over
// the valid case and every malformed shape that used to panic on an
// out-of-range slice index. None of these may panic; each malformed input must
// surface a diagnostic the caller returns on, never index into the slice.
func TestUnitErrorPagesRedirectRuleParseID(t *testing.T) {
	cases := []struct {
		name string
		id   string
		ok   bool
	}{
		{"valid", "site-123:cfg-456", true},
		{"empty", "", false},
		{"one part", "site", false},
		{"three parts", "a:b:c", false},
		{"empty site", ":cfg", false},
		{"empty config", "site:", false},
		{"two empty", ":", false},
		{"trailing colon", "a:b:", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			siteId, configId, diags := parseErrorPagesRedirectRuleID(c.id)
			if c.ok {
				assert.Falsef(t, diags.HasError(), "expected no diagnostics for %q, got %v", c.id, diags)
				assert.Equal(t, "site-123", siteId)
				assert.Equal(t, "cfg-456", configId)
			} else {
				assert.Truef(t, diags.HasError(), "expected a diagnostic for malformed id %q, got none", c.id)
				assert.Empty(t, siteId)
				assert.Empty(t, configId)
			}
		})
	}
}

// TestUnitErrorPagesRedirectRuleStateMapping verifies a GetErrorPagesRedirectRule
// response is mapped onto the resource model with the composite id rebuilt from
// the parsed parts, the nested error_pages_redirect list reconstructed, and JSON
// numbers coerced to Int64.
func TestUnitErrorPagesRedirectRuleStateMapping(t *testing.T) {
	ctx := context.Background()
	obj := map[string]interface{}{
		"ConfigType":  "rule",
		"Rule":        "http.host eq \"a.com\"",
		"RuleEnable":  "on",
		"RuleName":    "rule-1",
		"Sequence":    float64(10),
		"SiteVersion": float64(2),
		"ErrorPagesRedirect": []interface{}{
			map[string]interface{}{
				"TargetURL":  "https://example.com/404.html",
				"StatusCode": "404",
			},
			map[string]interface{}{
				"TargetURL":  "https://example.com/500.html",
				"StatusCode": "500",
			},
		},
	}
	d := new(diag.Diagnostics)
	m := errorPagesRedirectRuleState(ctx, "site-1", "cfg-1", obj, d)
	assert.Falsef(t, d.HasError(), "unexpected diagnostics: %v", d)
	assert.Equal(t, "site-1:cfg-1", m.Id.ValueString())
	assert.Equal(t, "site-1", m.SiteId.ValueString())
	assert.Equal(t, "cfg-1", m.ConfigId.ValueString())
	assert.Equal(t, "rule", m.ConfigType.ValueString())
	assert.Equal(t, "rule-1", m.RuleName.ValueString())
	assert.Equal(t, int64(10), m.Sequence.ValueInt64())
	assert.Equal(t, int64(2), m.SiteVersion.ValueInt64())
	assert.False(t, m.ErrorPagesRedirect.IsNull(), "list must not be null")
	assert.Equal(t, 2, len(m.ErrorPagesRedirect.Elements()))
}

// TestUnitErrorPagesRedirectRuleStateMappingOmitsEmpty checks that missing or
// empty response fields map to null state values, so optional/computed
// attributes read back absent rather than "".
func TestUnitErrorPagesRedirectRuleStateMappingOmitsEmpty(t *testing.T) {
	ctx := context.Background()
	obj := map[string]interface{}{
		"ConfigType":         "",
		"Rule":               nil,
		"RuleEnable":         "off",
		"ErrorPagesRedirect": nil,
	}
	d := new(diag.Diagnostics)
	m := errorPagesRedirectRuleState(ctx, "s", "c", obj, d)
	assert.False(t, d.HasError())
	assert.True(t, m.ConfigType.IsNull(), "empty ConfigType must be null")
	assert.True(t, m.Rule.IsNull(), "nil Rule must be null")
	assert.False(t, m.RuleEnable.IsNull())
	assert.Equal(t, "off", m.RuleEnable.ValueString())
	assert.False(t, m.ErrorPagesRedirect.IsNull(), "nil list must map to a non-null empty list")
	assert.Equal(t, 0, len(m.ErrorPagesRedirect.Elements()))
}

// TestUnitErrorPagesRedirectRuleListBuilding drives errorPagesRedirectList with
// nil, a non-array value and a real array, asserting it never produces a null
// list and never errors. A non-array value maps to an empty list, not a panic.
func TestUnitErrorPagesRedirectRuleListBuilding(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		raw  interface{}
		len  int
	}{
		{"nil", nil, 0},
		{"non-array", "not-an-array", 0},
		{"empty array", []interface{}{}, 0},
		{"one element", []interface{}{map[string]interface{}{"TargetURL": "u", "StatusCode": "404"}}, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := new(diag.Diagnostics)
			l := errorPagesRedirectList(ctx, c.raw, d)
			assert.Falsef(t, d.HasError(), "unexpected diagnostics: %v", d)
			assert.False(t, l.IsNull(), "list must never be null")
			assert.Equal(t, c.len, len(l.Elements()))
		})
	}
}

// TestUnitErrorPagesRedirectRuleRequestJSON verifies the plan's
// error_pages_redirect list is marshalled to the JSON string shape the ESA
// Create/Update API expects, and that a null/unknown list maps to "[]".
func TestUnitErrorPagesRedirectRuleRequestJSON(t *testing.T) {
	ctx := context.Background()

	t.Run("null list", func(t *testing.T) {
		d := new(diag.Diagnostics)
		s, ok := errorPagesRedirectRequestJSON(ctx, types.ListNull(errorPagesRedirectItemType), d)
		assert.True(t, ok)
		assert.False(t, d.HasError())
		assert.Equal(t, "[]", s)
	})

	t.Run("two items roundtrip", func(t *testing.T) {
		items := []errorPagesRedirectItem{
			{TargetUrl: types.StringValue("https://a.com/404.html"), StatusCode: types.StringValue("404")},
			{TargetUrl: types.StringValue("https://a.com/500.html"), StatusCode: types.StringValue("500")},
		}
		list, ldiags := types.ListValueFrom(ctx, errorPagesRedirectItemType, items)
		assert.Falsef(t, ldiags.HasError(), "building list: %v", ldiags)
		d := new(diag.Diagnostics)
		s, ok := errorPagesRedirectRequestJSON(ctx, list, d)
		assert.True(t, ok)
		assert.Falsef(t, d.HasError(), "unexpected diagnostics: %v", d)

		var got []map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(s), &got))
		assert.Equal(t, 2, len(got))
		assert.Equal(t, "404", got[0]["StatusCode"])
		assert.Equal(t, "https://a.com/404.html", got[0]["TargetURL"])
		assert.Equal(t, "500", got[1]["StatusCode"])
	})
}

// TestUnitErrorPagesRedirectRuleScalarConverters pins the edge-case behavior of
// the response→types converters: empty/missing → null, JSON numbers → Int64
// and to canonical integer strings for ID components.
func TestUnitErrorPagesRedirectRuleScalarConverters(t *testing.T) {
	// stringVal
	assert.True(t, stringVal(nil).IsNull())
	assert.True(t, stringVal("").IsNull())
	assert.Equal(t, "x", stringVal("x").ValueString())
	assert.Equal(t, "5", stringVal(float64(5)).ValueString())

	// int64Val
	assert.True(t, int64Val(nil).IsNull())
	assert.True(t, int64Val("").IsNull())
	assert.Equal(t, int64(3), int64Val(float64(3)).ValueInt64())
	assert.Equal(t, int64(9), int64Val(9).ValueInt64())
	assert.Equal(t, int64(7), int64Val("7").ValueInt64())
	assert.True(t, int64Val("not-a-number").IsNull(), "non-numeric string must map to null, not panic")

	// toString (ID component form: integer for JSON numbers)
	assert.Equal(t, "", toString(nil))
	assert.Equal(t, "abc", toString("abc"))
	assert.Equal(t, "123", toString(float64(123)))
	assert.Equal(t, "42", toString(42))
	assert.Equal(t, "42", toString(int64(42)))
}

// TestUnitErrorPagesRedirectRuleNotFoundClassifier verifies the not-found
// detector that Read (RemoveResource on nil) and Delete (idempotent ignore) rely
// on. A 404 or a NotFound-bearing message is classified true; a 400/Throttling
// error and nil are not. This is the failure path the cloudless unit layer can
// reach without driving the plugin gRPC server.
func TestUnitErrorPagesRedirectRuleNotFoundClassifier(t *testing.T) {
	assert.True(t, errs.NotFoundError(&tea.SDKError{StatusCode: tea.Int(404)}), "404 must be not-found")
	assert.True(t, errs.NotFoundError(&tea.SDKError{Message: tea.String("the specified rule is NotFound")}))

	assert.False(t, errs.NotFoundError(nil), "nil must not be not-found")
	assert.False(t, errs.NotFoundError(&tea.SDKError{Code: tea.String("Throttling"), StatusCode: tea.Int(400)}))
	assert.False(t, errs.NotFoundError(&tea.SDKError{Code: tea.String("InvalidParameter"), Message: tea.String("bad input"), StatusCode: tea.Int(400)}))
}
