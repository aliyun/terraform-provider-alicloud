package agentloop

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	tfSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitAliCloudExpandAgentloopEndpointConnectorHeaders(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("empty headers block caused panic: %v", recovered)
		}
	}()

	got := expandAgentloopEndpointConnectorHeaders([]interface{}{
		nil,
		map[string]interface{}{"key": "X-Test", "value": "value"},
	})
	want := []interface{}{map[string]interface{}{"key": "X-Test", "value": "value"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestUnitAliCloudExpandAgentloopEndpointConnectorUpdateCredential(t *testing.T) {
	got, ok := expandAgentloopEndpointConnectorUpdateCredential(map[string]interface{}{
		"api_key": "secret",
	}).(map[string]interface{})
	if !ok {
		t.Fatalf("credential type = %T, want map", got)
	}
	if got["providerType"] != "plain" {
		t.Fatalf("providerType = %v, want plain", got["providerType"])
	}
	var sensitiveInfo map[string]interface{}
	if err := json.Unmarshal([]byte(got["sensitiveInfo"].(string)), &sensitiveInfo); err != nil {
		t.Fatalf("decode sensitiveInfo: %v", err)
	}
	if sensitiveInfo["apiKey"] != "secret" {
		t.Fatalf("apiKey = %v, want secret", sensitiveInfo["apiKey"])
	}

	envelope := map[string]interface{}{"providerType": "kms", "sensitiveInfo": "{}"}
	if passThrough := expandAgentloopEndpointConnectorUpdateCredential(envelope); !reflect.DeepEqual(passThrough, envelope) {
		t.Fatalf("existing envelope changed: %#v", passThrough)
	}
}

// TestUnitAliCloudAgentloopEndpointConnectorFlattenProperties pins the
// properties read back against the GetEndpointConnector shape: injected
// defaults and null keys are skipped, configured keys (including defaults
// and keys the backend does not store) are kept, JSON keeps the configured
// text, and external changes surface as drift.
func TestUnitAliCloudAgentloopEndpointConnectorFlattenProperties(t *testing.T) {
	raw := map[string]interface{}{
		"channelType": "custom",
		"maxRetries":  "3",
		"qpsLimit":    "100",
		"timeoutMs":   "300000",
		"protocol":    "openai-compatible",
		"modelList":   `[{"model":"qwen-plus"}]`,
		"vpcId":       nil,
		"natIp":       nil,
	}

	configured := map[string]interface{}{
		"protocol":  "openai-compatible",
		"qpsLimit":  "100",
		"modelList": `[ { "model": "qwen-plus" } ]`,
		"property1": "value1",
	}
	want := map[string]interface{}{
		"protocol":  "openai-compatible",
		"qpsLimit":  "100",
		"modelList": `[ { "model": "qwen-plus" } ]`,
		"property1": "value1",
	}
	if got := flattenAgentloopEndpointConnectorProperties(raw, configured, "model_service"); !reflect.DeepEqual(got, want) {
		t.Fatalf("refresh = %#v, want %#v", got, want)
	}

	imported := map[string]interface{}{
		"protocol":  "openai-compatible",
		"modelList": `[{"model":"qwen-plus"}]`,
	}
	if got := flattenAgentloopEndpointConnectorProperties(raw, nil, "model_service"); !reflect.DeepEqual(got, imported) {
		t.Fatalf("import = %#v, want %#v", got, imported)
	}

	drifted := map[string]interface{}{
		"protocol":  "anthropic",
		"qpsLimit":  "200",
		"modelList": nil,
	}
	wantDrift := map[string]interface{}{
		"protocol": "anthropic",
		"qpsLimit": "200",
	}
	if got := flattenAgentloopEndpointConnectorProperties(drifted, map[string]interface{}{
		"protocol":  "openai-compatible",
		"modelList": `[{"model":"qwen-plus"}]`,
	}, "model_service"); !reflect.DeepEqual(got, wantDrift) {
		t.Fatalf("drift = %#v, want %#v", got, wantDrift)
	}
}

// TestUnitAliCloudAgentloopEndpointConnectorImportPlan builds the state
// produced by an import Read and checks that re-applying the original
// configuration, including the tags list, plans no replacement.
func TestUnitAliCloudAgentloopEndpointConnectorImportPlan(t *testing.T) {
	r := ResourceEndpointConnector()
	data := tfSchema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"agent_space": "as",
		"type":        "model_service",
		"name":        "ec",
		"endpoint":    "https://example.com/v1",
		"tags":        []interface{}{"tag1", "tag2"},
		"headers":     []interface{}{map[string]interface{}{"key": "X-K", "value": "v"}},
	})
	data.SetId("as:cid")
	if err := data.Set("properties", flattenAgentloopEndpointConnectorProperties(map[string]interface{}{
		"channelType": "custom", "qpsLimit": "100", "protocol": "openai-compatible", "vpcId": nil,
	}, nil, "model_service")); err != nil {
		t.Fatalf("set properties: %v", err)
	}
	state := data.State()

	cfg := map[string]interface{}{
		"agent_space": "as",
		"type":        "model_service",
		"name":        "ec",
		"endpoint":    "https://example.com/v1",
		"credential":  map[string]interface{}{"api_key": "sk"},
		"tags":        []interface{}{"tag1", "tag2"},
		"headers":     []interface{}{map[string]interface{}{"key": "X-K", "value": "v"}},
		"properties":  map[string]interface{}{"protocol": "openai-compatible"},
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff != nil && diff.RequiresNew() {
		t.Fatalf("imported config must not require replacement: %#v", diff.Attributes)
	}
	if diff != nil {
		for k := range diff.Attributes {
			if k != "credential.%" && k != "credential.api_key" {
				t.Errorf("unexpected diff on %s: %#v", k, diff.Attributes[k])
			}
		}
	}
}

// TestUnitAliCloudAgentloopEndpointConnectorFlattenPropertiesAgentApp verifies
// that the qpsLimit default for agent_app (20) is correctly suppressed, while
// model_service's default (100) would surface as drift on an agent_app connector.
func TestUnitAliCloudAgentloopEndpointConnectorFlattenPropertiesAgentApp(t *testing.T) {
	// agent_app backend returns qpsLimit=20 (its default) – should be suppressed.
	raw := map[string]interface{}{
		"channelType": "custom",
		"maxRetries":  "3",
		"qpsLimit":    "20",
		"timeoutMs":   "300000",
	}
	got := flattenAgentloopEndpointConnectorProperties(raw, nil, "agent_app")
	if len(got) != 0 {
		t.Fatalf("agent_app defaults should be suppressed, got %#v", got)
	}

	// Same raw but with model_service type – qpsLimit=20 is NOT the
	// model_service default (100), so it should surface as drift.
	got2 := flattenAgentloopEndpointConnectorProperties(raw, nil, "model_service")
	if got2["qpsLimit"] != "20" {
		t.Fatalf("model_service should see qpsLimit=20 as drift, got %#v", got2)
	}

	// agent_app with explicitly configured qpsLimit=20 should keep it.
	configured := map[string]interface{}{"qpsLimit": "20"}
	got3 := flattenAgentloopEndpointConnectorProperties(raw, configured, "agent_app")
	if got3["qpsLimit"] != "20" {
		t.Fatalf("explicitly configured qpsLimit should be kept, got %#v", got3)
	}
}
