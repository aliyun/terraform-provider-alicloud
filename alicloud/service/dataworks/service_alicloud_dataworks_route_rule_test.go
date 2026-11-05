package dataworks

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// TestUnitRouteRuleIdFromResponse covers the route rule id conversion for every
// shape the tea RPC decoder can return: string, json.Number (the lossless path
// for large ids), float64 (with its lossless-range guard), int, int64, uint64
// and the rejected containers/NaN/fractional/beyond-exact values.
func TestUnitRouteRuleIdFromResponse(t *testing.T) {
	cases := map[string]struct {
		raw     interface{}
		want    string
		wantErr bool
	}{
		"string":             {raw: "123", want: "123"},
		"emptyString":        {raw: "", wantErr: true},
		"jsonNumber":         {raw: json.Number("1183029123456789012"), want: "1183029123456789012"},
		"jsonNumberFraction": {raw: json.Number("1.5"), wantErr: true},
		"jsonNumberExponent": {raw: json.Number("1e3"), wantErr: true},
		"jsonNumberZero":     {raw: json.Number("0"), wantErr: true},
		"float":              {raw: float64(123), want: "123"},
		"floatExactLimit":    {raw: float64(1 << 53), want: "9007199254740992"},
		"floatBeyondExact":   {raw: float64(1<<53 + 2), wantErr: true},
		"floatFraction":      {raw: 1.5, wantErr: true},
		"floatNaN":           {raw: math.NaN(), wantErr: true},
		"floatInf":           {raw: math.Inf(1), wantErr: true},
		"floatZero":          {raw: float64(0), wantErr: true},
		"floatNegative":      {raw: float64(-1), wantErr: true},
		"int":                {raw: 123, want: "123"},
		"intNegative":        {raw: -1, wantErr: true},
		"int64":              {raw: int64(123), want: "123"},
		"uint64":             {raw: uint64(123), want: "123"},
		"uint64Zero":         {raw: uint64(0), wantErr: true},
		"map":                {raw: map[string]interface{}{"id": "123"}, wantErr: true},
		"slice":              {raw: []interface{}{"123"}, wantErr: true},
		"bool":               {raw: true, wantErr: true},
		"nil":                {raw: nil, wantErr: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := routeRuleIdFromResponse(c.raw)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("routeRuleIdFromResponse(%v) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}

// TestUnitRouteRuleParseCompositeId covers the composite id contract: the
// middle segment (resource_id) may itself contain a colon, a bare route rule
// id is rejected so Read surfaces the missing identity, and empty segments fail.
func TestUnitRouteRuleParseCompositeId(t *testing.T) {
	cases := map[string]struct {
		id           string
		wantGroup    string
		wantResource string
		wantRule     string
		wantErr      bool
	}{
		"numericTriple":       {id: "123:456:789", wantGroup: "123", wantResource: "456", wantRule: "789"},
		"resourceIdWithColon": {id: "123:ns-12:34:789", wantGroup: "123", wantResource: "ns-12:34", wantRule: "789"},
		"bareRuleIdRejected":  {id: "789", wantErr: true},
		"twoSegmentsRejected": {id: "123:789", wantErr: true},
		"emptyGroup":          {id: ":456:789", wantErr: true},
		"emptyResource":       {id: "123::789", wantErr: true},
		"emptyRule":           {id: "123:456:", wantErr: true},
		"empty":               {id: "", wantErr: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			group, resource, rule, err := parseRouteRuleId(c.id)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %q, got %s:%s:%s", c.id, group, resource, rule)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", c.id, err)
			}
			if group != c.wantGroup || resource != c.wantResource || rule != c.wantRule {
				t.Fatalf("parseRouteRuleId(%q) = (%q,%q,%q), want (%q,%q,%q)", c.id, group, resource, rule, c.wantGroup, c.wantResource, c.wantRule)
			}
		})
	}
}

func TestUnitRouteRuleCompositeIdRoundTrip(t *testing.T) {
	pairs := []struct{ group, resource, rule string }{
		{"123", "456", "789"},
		{"123", "ns-12:34", "789"},
	}
	for _, p := range pairs {
		id := routeRuleCompositeId(p.group, p.resource, p.rule)
		group, resource, rule, err := parseRouteRuleId(id)
		if err != nil {
			t.Fatalf("round-trip parse failed for %q: %v", id, err)
		}
		if group != p.group || resource != p.resource || rule != p.rule {
			t.Fatalf("round-trip mismatch: got (%q,%q,%q) from %q", group, resource, rule, id)
		}
	}
}

// TestUnitRouteRuleSuccess covers the OpenApiResultModel success flag in both
// its bool and string spellings, plus camelCase/PascalCase fallback.
func TestUnitRouteRuleSuccess(t *testing.T) {
	cases := map[string]struct {
		response map[string]interface{}
		want     bool
	}{
		"boolTrue":        {response: map[string]interface{}{"success": true}, want: true},
		"boolFalse":       {response: map[string]interface{}{"success": false}, want: false},
		"stringTrue":      {response: map[string]interface{}{"success": "true"}, want: true},
		"stringTrueUpper": {response: map[string]interface{}{"Success": "TRUE"}, want: true},
		"pascalBool":      {response: map[string]interface{}{"Success": true}, want: true},
		"missing":         {response: map[string]interface{}{"data": "x"}, want: false},
		"nil":             {response: nil, want: false},
		"wrongType":       {response: map[string]interface{}{"success": 1}, want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := routeRuleSuccess(c.response); got != c.want {
				t.Fatalf("routeRuleSuccess(%v) = %v, want %v", c.response, got, c.want)
			}
		})
	}
}

// TestUnitRouteRuleBusinessError covers rebuilding the error from the
// OpenApiResultModel errorCode/errorMessage fields (and the dynamic fallbacks),
// including the empty-but-success=false body that must still surface a failure.
func TestUnitRouteRuleBusinessError(t *testing.T) {
	cases := map[string]struct {
		response map[string]interface{}
		wantErr  bool
		wantMsg  string
	}{
		"successTrueReturnsNil": {response: map[string]interface{}{"success": true, "errorCode": "X"}, wantErr: false},
		"nilReturnsNil":         {response: nil, wantErr: false},
		"errorCodeAndMessage":   {response: map[string]interface{}{"success": false, "errorCode": "InvalidParam", "errorMessage": "bad cidr"}, wantErr: true, wantMsg: "InvalidParam"},
		"pascalCaseFields":      {response: map[string]interface{}{"Success": false, "ErrorCode": "NotFound.RouteRule", "ErrorMessage": "gone"}, wantErr: true, wantMsg: "NotFound.RouteRule"},
		"dynamicFallback":       {response: map[string]interface{}{"success": false, "dynamicErrorCode": "Dyn.Code", "dynamicErrorMessage": "dyn msg"}, wantErr: true, wantMsg: "Dyn.Code"},
		"codeOnly":              {response: map[string]interface{}{"success": false, "errorCode": "OnlyCode"}, wantErr: true, wantMsg: "OnlyCode"},
		"messageOnly":           {response: map[string]interface{}{"success": false, "errorMessage": "only message"}, wantErr: true, wantMsg: "only message"},
		"emptyBody":             {response: map[string]interface{}{"success": false}, wantErr: true, wantMsg: "success=false"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := routeRuleBusinessError(c.response, "CreateRouteRule")
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %s", name)
				}
				if c.wantMsg != "" && !strings.Contains(err.Error(), c.wantMsg) {
					t.Fatalf("error %q does not mention %q", err.Error(), c.wantMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error for %s, got %v", name, err)
			}
		})
	}
}

// TestUnitRouteRuleNotFoundClassification guards the not-found classifier: only
// an explicit not-found token (or a typed NotFound) counts as gone, so a
// permission error whose message happens to mention "NotFound" still surfaces as
// a real failure during delete/read instead of being swallowed as absence.
func TestUnitRouteRuleNotFoundClassification(t *testing.T) {
	notFound := errs.WrapErrorf(errs.NotFoundErr("RouteRule", "1"), errs.NotFoundMsg, map[string]interface{}{"RequestId": "r"})
	cases := map[string]struct {
		err      error
		wantGone bool
	}{
		"nil":               {err: nil, wantGone: false},
		"typedNotFound":     {err: notFound, wantGone: true},
		"notFoundInMessage": {err: errors.New("the route rule was NotFound"), wantGone: true},
		"doesNotExist":      {err: errors.New("route rule does not exist"), wantGone: true},
		"permissionDenied":  {err: errors.New("APINotAuthorize: denied while reading NotFound path"), wantGone: false},
		"internalError":     {err: errors.New("InternalError: try again"), wantGone: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := routeRuleNotFound(c.err, "1"); got != c.wantGone {
				t.Fatalf("routeRuleNotFound(%v) = %v, want %v", c.err, got, c.wantGone)
			}
		})
	}
}

// TestUnitRouteRuleRetryable guards the retry classifier: the shared transient
// classifier and explicit internal-error/system-busy tokens retry, while a
// documented terminal answer does not.
func TestUnitRouteRuleRetryable(t *testing.T) {
	cases := map[string]struct {
		err  error
		want bool
	}{
		"nil":                {err: nil, want: false},
		"internalError":      {err: errors.New("InternalError: boom"), want: true},
		"systemBusy":         {err: errors.New("SystemBusy"), want: true},
		"serviceUnavailable": {err: errors.New("ServiceUnavailable"), want: true},
		"rateLimited":        {err: errors.New("Request was rate limited"), want: true},
		"invalidParam":       {err: errors.New("InvalidParam: bad cidr"), want: false},
		"notFound":           {err: errors.New("NotFound.RouteRule: gone"), want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := routeRuleRetryable(c.err); got != c.want {
				t.Fatalf("routeRuleRetryable(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

// TestUnitRouteRuleNormalize covers the ListRouteRules entry mapping: the
// composite id is filled from the request identity, and the destination fields
// are read by their camelCase backend name with the PascalCase fallback.
func TestUnitRouteRuleNormalize(t *testing.T) {
	t.Run("camelCase", func(t *testing.T) {
		entry := map[string]interface{}{
			"destinationCidr":           "192.168.10.0/24",
			"destinationType":           "VPC",
			"destinationAddressingType": "CIDR",
		}
		got := routeRuleNormalize(entry, "123", "456", "789")
		if got["Id"] != "123:456:789" {
			t.Fatalf("Id = %v, want 123:456:789", got["Id"])
		}
		if got["RouteRuleId"] != "789" {
			t.Fatalf("RouteRuleId = %v", got["RouteRuleId"])
		}
		if got["DestinationCidr"] != "192.168.10.0/24" {
			t.Fatalf("DestinationCidr = %v", got["DestinationCidr"])
		}
		if got["DestinationType"] != "VPC" {
			t.Fatalf("DestinationType = %v", got["DestinationType"])
		}
		if got["DestinationAddressingType"] != "CIDR" {
			t.Fatalf("DestinationAddressingType = %v", got["DestinationAddressingType"])
		}
	})
	t.Run("pascalCaseFallback", func(t *testing.T) {
		entry := map[string]interface{}{
			"DestinationCidr":           "10.0.0.0/8",
			"DestinationType":           "IDC",
			"DestinationAddressingType": "IP",
		}
		got := routeRuleNormalize(entry, "1", "2", "3")
		if got["DestinationCidr"] != "10.0.0.0/8" || got["DestinationType"] != "IDC" || got["DestinationAddressingType"] != "IP" {
			t.Fatalf("PascalCase fallback failed: %v", got)
		}
	})
	t.Run("missingFieldsRenderEmpty", func(t *testing.T) {
		got := routeRuleNormalize(map[string]interface{}{}, "1", "2", "3")
		if got["DestinationCidr"] != "" || got["DestinationType"] != "" {
			t.Fatalf("missing fields must render as empty, got %v", got)
		}
	})
}
