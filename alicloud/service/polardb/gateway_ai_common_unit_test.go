package polardb

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitPolarDBGatewayAIUnsetOptionalRequests(t *testing.T) {
	serviceData := schema.TestResourceDataRaw(t, ResourceGatewayModelService().Schema, map[string]interface{}{
		"gateway_id": "gw-1", "name": "service", "model_category": "text", "protocol": "openai",
		"base_url": "https://example.com/v1", "api_key": "test-only", "vendor": "custom",
	})
	serviceRequest := polarDBGatewayModelServiceRequest(serviceData, "cn-beijing")
	for _, key := range []string{"InputCostPointsPerMillion", "OutputCostPointsPerMillion", "RequestCostPoints"} {
		if _, exists := serviceRequest[key]; exists {
			t.Fatalf("request unexpectedly contains unset %s", key)
		}
	}

	apiData := schema.TestResourceDataRaw(t, ResourceGatewayModelAPI().Schema, map[string]interface{}{
		"gateway_id": "gw-1", "name": "api", "model_category": "text", "path_prefix": "/api",
		"protocol": "openai", "route_rules": `[{"RuleName":"primary"}]`,
	})
	apiCreate := polarDBGatewayModelAPIRequest(apiData, "cn-beijing", true)
	apiModify := polarDBGatewayModelAPIRequest(apiData, "cn-beijing", false)
	for _, key := range []string{"RecordInput", "RecordOutput"} {
		if _, exists := apiCreate[key]; exists {
			t.Fatalf("create request unexpectedly contains unset %s", key)
		}
		if _, exists := apiModify[key]; exists {
			t.Fatalf("modify request unexpectedly contains unset %s", key)
		}
	}
	if !suppressPolarDBGatewayAIUnsetOptional("record_input", "2", "", apiData) {
		t.Fatal("removing a server-retained optional value must suppress the diff")
	}
	if suppressPolarDBGatewayAIUnsetOptional("record_input", "2", "3", apiData) {
		t.Fatal("changing an optional value must not suppress the diff")
	}
}
func TestUnitPolarDBGatewayAIFilterMatching(t *testing.T) {
	data := schema.TestResourceDataRaw(t, DataSourceGatewayModelAPIs().Schema, map[string]interface{}{
		"gateway_id": "pg-1", "name": "api-one", "model_category": "text",
	})
	fields := map[string][]string{
		"name":           []string{"Name"},
		"model_category": []string{"Category", "ModelCategory"},
	}
	if !matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-one", "Category": "text"}, fields) {
		t.Fatal("expected matching filters to be accepted")
	}
	if !matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-one", "ModelCategory": "text"}, fields) {
		t.Fatal("expected ModelCategory response shape to be accepted")
	}
	if matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-two", "Category": "text"}, fields) {
		t.Fatal("expected a mismatched name to be rejected")
	}
}
func TestUnitPolarDBGatewayAIChildID(t *testing.T) {
	id := polarDBGatewayAIChildID("pg-1", "ms-1")
	if id != "pg-1:ms-1" {
		t.Fatalf("unexpected child ID: %s", id)
	}
	gatewayID, childID, err := ParsePolarDBGatewayAIChildID(id)
	if err != nil || gatewayID != "pg-1" || childID != "ms-1" {
		t.Fatalf("unexpected parse result: %q %q %v", gatewayID, childID, err)
	}
	if _, _, err = ParsePolarDBGatewayAIChildID("ms-1"); err == nil {
		t.Fatal("expected invalid import ID to be rejected")
	}
}
func TestUnitPolarDBGatewayAIInt(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{input: " 28 ", want: 28},
		{input: float64(80), want: 80},
		{input: "unknown", want: "unknown"},
	}
	for _, test := range tests {
		if got := polarDBGatewayAIInt(test.input); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("polarDBGatewayAIInt(%#v) = %#v, want %#v", test.input, got, test.want)
		}
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
	modelServiceSchema := ResourceGatewayModelService().Schema
	if !modelServiceSchema["api_key"].Sensitive {
		t.Fatal("model service api_key must be sensitive")
	}
	validateBaseURL := modelServiceSchema["base_url"].ValidateFunc
	if _, errors := validateBaseURL("https://example.com/v1", "base_url"); len(errors) != 0 {
		t.Fatalf("expected HTTPS base_url to be accepted: %v", errors)
	}
	if _, errors := validateBaseURL("ftp://example.com/v1", "base_url"); len(errors) == 0 {
		t.Fatal("expected non-HTTP base_url to be rejected")
	}
	if !ResourceGatewayModelAPI().Schema["name"].ForceNew {
		t.Fatal("model API name must be ForceNew because ModifyModelApi cannot change it")
	}
	if !ResourceGatewayModelService().Schema["name"].ForceNew {
		t.Fatal("model service name must be ForceNew because ModifyModelService cannot change it")
	}
	if !ResourceGatewayCostRule().Schema["model_name"].ForceNew {
		t.Fatal("cost rule model_name must be ForceNew because ModifyCostRule does not change it")
	}
	if !ResourceGatewayCostRule().Schema["model_service_id"].ForceNew {
		t.Fatal("cost rule model_service_id must be ForceNew because ModifyCostRule retains it")
	}
	consumerSchema := ResourceGatewayConsumer().Schema
	for _, key := range []string{"name", "consumer_group_id", "nickname", "is_default"} {
		if !consumerSchema[key].ForceNew {
			t.Fatalf("consumer %s must be ForceNew because live updates did not converge", key)
		}
	}
	if consumerSchema["api_key_reset_token"].ForceNew {
		t.Fatal("consumer api_key_reset_token must update in place to reset the API key")
	}
	if !ResourceGatewayConsumerGroup().Schema["is_default"].Computed {
		t.Fatal("consumer group is_default must be read-only because live creates and updates did not persist it")
	}
	if !consumerSchema["api_key"].Sensitive || !consumerSchema["api_key"].Computed {
		t.Fatal("consumer api_key must be a sensitive computed value")
	}
	if consumerSchema["api_key_reset_token"].Sensitive {
		t.Fatal("consumer api_key_reset_token is a rotation trigger, not the secret")
	}
	budgetSchema := ResourceGatewayBudgetPolicy().Schema
	if !budgetSchema["budget_type"].ForceNew || budgetSchema["budget_points"].ForceNew {
		t.Fatal("budget type must replace the policy while budget points must update in place")
	}
	rateLimitSchema := ResourceGatewayRateLimitPolicy().Schema
	if !rateLimitSchema["scope_type"].ForceNew || rateLimitSchema["rate_limit_rpm"].ForceNew {
		t.Fatal("rate limit scope must replace the policy while rate limits must update in place")
	}
}
