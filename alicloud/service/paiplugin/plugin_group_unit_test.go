package paiplugin

import (
	"encoding/json"
	"errors"
	"testing"

	tea "github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// sdkErr builds a *tea.SDKError, the form in which the ROA path surfaces HTTP
// responses via formatError.
func sdkErr(statusCode int, code, message string) *tea.SDKError {
	e := &tea.SDKError{Code: tea.String(code), Message: tea.String(message)}
	if statusCode != 0 {
		e.StatusCode = tea.Int(statusCode)
	}
	return e
}

// TestUnitPluginGroupToStrIDValidation exercises F3: toStr accepts only string
// and integer scalars as a group id; bool/map/slice/nil and non-integer floats
// yield "" so a malformed body cannot synthesize a pseudo-id.
func TestUnitPluginGroupToStrIDValidation(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		want string
	}{
		{"string", "abc123", "abc123"},
		{"int", 123, "123"},
		{"int64", int64(456), "456"},
		{"integer float64", float64(789), "789"},
		{"json.Number integer", json.Number("321"), "321"},
		{"non-integer float64", float64(12.5), ""},
		{"json.Number non-integer", json.Number("1.5"), ""},
		{"bool", true, ""},
		{"map", map[string]interface{}{"Id": 1}, ""},
		{"slice", []interface{}{1, 2}, ""},
		{"nil", nil, ""},
		{"struct", struct{ A int }{A: 1}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, toStr(c.in))
		})
	}
}

// TestUnitPluginGroupIsNotExistContract exercises F2 (AC-R2): only an explicit
// product not-found code on a 404 removes the resource or tolerates a failed
// Delete. An empty code, an unknown code, routing/auth-layer 404s
// (InvalidApi.NotFound / InvalidAction / auth), 403, 500, nil and non-SDK
// errors are hard errors that must preserve state. *ComplexError is unwrapped
// to its SDK cause.
func TestUnitPluginGroupIsNotExistContract(t *testing.T) {
	genuineNotFound := sdkErr(404, "NotFound", "group not found")
	resourceNotFound := sdkErr(404, "ResourceNotFound", "resource not found")
	emptyCode404 := sdkErr(404, "", "")
	unknownCode404 := sdkErr(404, "SomeUnknownError", "unknown")
	invalidApiNotFound := sdkErr(404, "InvalidApi.NotFound", "specified api is not found")
	invalidAction := sdkErr(404, "InvalidAction", "action not found")
	forbidden := sdkErr(404, "Forbidden.AccessForbidden", "access forbidden")
	invalidAccessKeyId := sdkErr(404, "InvalidAccessKeyId.NotFound", "ak not found")
	missingAccessKeyId := sdkErr(404, "MissingAccessKeyId", "ak missing")
	signatureMismatch := sdkErr(404, "SignatureDoesNotMatch", "signature mismatch")
	unauthorized := sdkErr(404, "Unauthorized", "unauthorized")
	accessDenied := sdkErr(404, "AccessDenied", "access denied")
	forbidden403 := sdkErr(403, "Forbidden", "forbidden")
	serverError500 := sdkErr(500, "InternalError", "internal")

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"genuine NotFound 404", genuineNotFound, true},
		{"ResourceNotFound 404", resourceNotFound, true},
		{"empty code 404 (preserve state)", emptyCode404, false},
		{"unknown non-routing 404 (preserve state)", unknownCode404, false},
		{"InvalidApi.NotFound 404 (routing)", invalidApiNotFound, false},
		{"InvalidAction 404", invalidAction, false},
		{"Forbidden 404", forbidden, false},
		{"InvalidAccessKeyId 404", invalidAccessKeyId, false},
		{"MissingAccessKeyId 404", missingAccessKeyId, false},
		{"SignatureDoesNotMatch 404", signatureMismatch, false},
		{"Unauthorized 404", unauthorized, false},
		{"AccessDenied 404", accessDenied, false},
		{"403 Forbidden (not 404)", forbidden403, false},
		{"500 InternalError", serverError500, false},
		{"nil error", nil, false},
		{"plain non-SDK error", errors.New("something failed"), false},
		{"ComplexError wrapping genuine 404", &errs.ComplexError{Cause: genuineNotFound}, true},
		{"ComplexError wrapping InvalidApi.NotFound", &errs.ComplexError{Cause: invalidApiNotFound}, false},
		{"ComplexError wrapping 500", &errs.ComplexError{Cause: serverError500}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, isPluginGroupNotExist(c.err))
		})
	}
}

// TestUnitPluginGroupStateFromCreatePreservesID exercises F1: the persisted
// create-state carries the assigned id/plugin_group_id and the original plan
// configuration, and the caller's plan is not mutated (value semantics).
func TestUnitPluginGroupStateFromCreatePreservesID(t *testing.T) {
	plan := pluginGroupResourceModel{
		Id:            types.StringUnknown(),
		PluginGroupId: types.StringUnknown(),
		Name:          types.StringValue("g1"),
		Source:        types.Int64Value(0),
		Text:          types.StringValue("13900000000"),
		Remark:        types.StringValue("terraform-example"),
	}
	m := pluginGroupStateFromCreate(plan, "123")
	assert.Equal(t, "123", m.Id.ValueString())
	assert.Equal(t, "123", m.PluginGroupId.ValueString())
	assert.Equal(t, "g1", m.Name.ValueString())
	assert.Equal(t, int64(0), m.Source.ValueInt64())
	assert.Equal(t, "13900000000", m.Text.ValueString())
	assert.Equal(t, "terraform-example", m.Remark.ValueString())

	// Caller's plan is not mutated (id stays unknown).
	assert.True(t, plan.Id.IsUnknown(), "plan.Id must remain unknown (not mutated)")
	assert.True(t, plan.PluginGroupId.IsUnknown(), "plan.PluginGroupId must remain unknown (not mutated)")
}

// TestUnitPluginGroupStateMalformedData exercises F3: pluginGroupState reports
// ok=false for a nil/empty/invalid Id (so callers preserve prior state), and
// ok=true with a mirrored plugin_group_id for a valid integer/string Id.
func TestUnitPluginGroupStateMalformedData(t *testing.T) {
	cases := []struct {
		name string
		data map[string]interface{}
		ok   bool
		id   string
	}{
		{"nil data", nil, false, ""},
		{"empty map", map[string]interface{}{}, false, ""},
		{"missing Id", map[string]interface{}{"Name": "g1"}, false, ""},
		{"nil Id", map[string]interface{}{"Id": nil}, false, ""},
		{"bool Id", map[string]interface{}{"Id": true}, false, ""},
		{"map Id", map[string]interface{}{"Id": map[string]interface{}{"x": 1}}, false, ""},
		{"empty string Id", map[string]interface{}{"Id": ""}, false, ""},
		{"non-integer float Id", map[string]interface{}{"Id": float64(1.5)}, false, ""},
		{"integer float64 Id", map[string]interface{}{"Id": float64(123), "Name": "g1"}, true, "123"},
		{"json.Number Id", map[string]interface{}{"Id": json.Number("456")}, true, "456"},
		{"string Id", map[string]interface{}{"Id": "abc"}, true, "abc"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m, ok := pluginGroupState(c.data)
			assert.Equal(t, c.ok, ok)
			if c.ok {
				assert.Equal(t, c.id, m.Id.ValueString())
				assert.Equal(t, c.id, m.PluginGroupId.ValueString())
			}
		})
	}
}
