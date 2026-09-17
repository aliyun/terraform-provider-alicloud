package alicloud

import (
	"reflect"
	"testing"
)

func TestUnitPolarDBGatewayAIChildID(t *testing.T) {
	id := polarDBGatewayAIChildID("pg-1", "ms-1")
	if id != "pg-1:ms-1" {
		t.Fatalf("unexpected child ID: %s", id)
	}
	gatewayID, childID, err := parsePolarDBGatewayAIChildID(id)
	if err != nil || gatewayID != "pg-1" || childID != "ms-1" {
		t.Fatalf("unexpected parse result: %q %q %v", gatewayID, childID, err)
	}
	if _, _, err = parsePolarDBGatewayAIChildID("ms-1"); err == nil {
		t.Fatal("expected invalid import ID to be rejected")
	}
}

func TestUnitNormalizePolarDBGatewayAIJSON(t *testing.T) {
	left := normalizePolarDBGatewayAIJSON(`[{"Weight":"10","Name":"primary"}]`)
	right := normalizePolarDBGatewayAIJSON(`[
      {"Name":"primary", "Weight":"10"}
    ]`)
	if left != right {
		t.Fatalf("semantically equal route rules were not normalized: %q != %q", left, right)
	}
	if warnings, errors := validatePolarDBGatewayAIRouteRules(`{"Name":"not-an-array"}`, "route_rules"); len(warnings) != 0 || len(errors) != 1 {
		t.Fatalf("expected object route_rules to be rejected, got warnings=%v errors=%v", warnings, errors)
	}
}

func TestUnitRedactPolarDBGatewayAISecrets(t *testing.T) {
	input := map[string]interface{}{
		"ApiKey": "top-secret",
		"Items":  []interface{}{map[string]interface{}{"apiKey": "nested-secret", "Name": "service"}},
	}
	want := map[string]interface{}{
		"ApiKey": "***",
		"Items":  []interface{}{map[string]interface{}{"apiKey": "***", "Name": "service"}},
	}
	got := redactPolarDBGatewayAISecrets(input)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected redacted value: %#v", got)
	}
	if input["ApiKey"] != "top-secret" {
		t.Fatal("redaction mutated the original request")
	}
}

func TestUnitPolarDBGatewayAIItems(t *testing.T) {
	raw := map[string]interface{}{"ModelService": []interface{}{
		map[string]interface{}{"ModelServiceId": "ms-1"},
	}}
	items := polarDBGatewayAIItems(raw)
	if len(items) != 1 || items[0]["ModelServiceId"] != "ms-1" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestUnitPolarDBGatewayAIResourceSchemas(t *testing.T) {
	if !resourceAlicloudPolarDBGatewayModelService().Schema["api_key"].Sensitive {
		t.Fatal("model service api_key must be sensitive")
	}
	if !resourceAlicloudPolarDBGatewayModelAPI().Schema["name"].ForceNew {
		t.Fatal("model API name must be ForceNew because ModifyModelApi cannot change it")
	}
	if resourceAlicloudPolarDBGatewayCostRule().Schema["model_name"].ForceNew {
		t.Fatal("cost rule model_name is supported by ModifyCostRule")
	}
	consumerSchema := resourceAlicloudPolarDBGatewayConsumer().Schema
	if !consumerSchema["api_key"].Sensitive || !consumerSchema["api_key"].Computed {
		t.Fatal("consumer api_key must be a sensitive computed value")
	}
	if consumerSchema["api_key_reset_token"].Sensitive {
		t.Fatal("consumer api_key_reset_token is a rotation trigger, not the secret")
	}
	budgetSchema := resourceAlicloudPolarDBGatewayBudgetPolicy().Schema
	if !budgetSchema["budget_type"].ForceNew || budgetSchema["budget_points"].ForceNew {
		t.Fatal("budget type must replace the policy while budget points must update in place")
	}
	rateLimitSchema := resourceAlicloudPolarDBGatewayRateLimitPolicy().Schema
	if !rateLimitSchema["scope_type"].ForceNew || rateLimitSchema["rate_limit_rpm"].ForceNew {
		t.Fatal("rate limit scope must replace the policy while rate limits must update in place")
	}
}
