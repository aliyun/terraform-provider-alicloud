package agentloop

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnitAliCloudExpandAgentloopEvaluatorConfig(t *testing.T) {
	raw := map[string]interface{}{
		"prompt":          `{"json":"example"}`,
		"variables":       `[{"name":"input"}]`,
		"variablesPadded": ` [1] `,
		"outputSchema":    `{"type":"object"}`,
		"faasConfig":      "not-json",
		"parameters":      `[broken`,
		"codeType":        "python",
		"number":          "10",
		"nonString":       1,
	}
	got := expandAgentloopEvaluatorConfig(raw)

	if got["prompt"] != `{"json":"example"}` {
		t.Errorf("prompt = %#v, want verbatim string (not a structured key)", got["prompt"])
	}
	if want := []interface{}{map[string]interface{}{"name": "input"}}; !reflect.DeepEqual(got["variables"], want) {
		t.Errorf("variables = %#v, want %#v", got["variables"], want)
	}
	if want := map[string]interface{}{"type": "object"}; !reflect.DeepEqual(got["outputSchema"], want) {
		t.Errorf("outputSchema = %#v, want %#v", got["outputSchema"], want)
	}
	if got["faasConfig"] != "not-json" {
		t.Errorf("faasConfig = %#v, want verbatim string", got["faasConfig"])
	}
	if got["parameters"] != `[broken` {
		t.Errorf("parameters = %#v, want verbatim invalid JSON", got["parameters"])
	}
	if got["codeType"] != "python" {
		t.Errorf("codeType = %#v, want verbatim string", got["codeType"])
	}
	if got["number"] != "10" {
		t.Errorf("number = %#v, want verbatim string (scalar JSON is not expanded)", got["number"])
	}
	if got["nonString"] != 1 {
		t.Errorf("nonString = %#v, want 1", got["nonString"])
	}
	// "variablesPadded" is not in the structured-keys whitelist, so it must be
	// sent verbatim even though it parses as JSON.
	if got["variablesPadded"] != ` [1] ` {
		t.Errorf("variablesPadded = %#v, want verbatim string (not whitelisted)", got["variablesPadded"])
	}

	if got := expandAgentloopEvaluatorConfig("not-a-map"); len(got) != 0 {
		t.Errorf("non map input: got %v entries, want empty map", len(got))
	}

	// A structured key carrying two JSON values must stay verbatim instead of
	// being silently truncated to the first document.
	multi := expandAgentloopEvaluatorConfig(map[string]interface{}{
		"variables": `[{"name":"a"}] [{"name":"b"}]`,
	})
	if multi["variables"] != `[{"name":"a"}] [{"name":"b"}]` {
		t.Errorf("variables multi-doc = %#v, want verbatim multi-document string", multi["variables"])
	}
}

func TestUnitAliCloudFlattenAgentloopEvaluatorConfig(t *testing.T) {
	raw := map[string]interface{}{
		"prompt":       "text",
		"variables":    []interface{}{"a"},
		"outputSchema": map[string]interface{}{"type": "object"},
		"bigNumber":    []interface{}{json.Number("9007199254740993")},
	}
	got := flattenAgentloopEvaluatorConfig(raw, nil)
	if got["prompt"] != "text" {
		t.Errorf("prompt = %q, want verbatim string", got["prompt"])
	}
	if got["variables"] != `["a"]` {
		t.Errorf("variables = %q, want serialized array", got["variables"])
	}
	if got["outputSchema"] != `{"type":"object"}` {
		t.Errorf("outputSchema = %q, want serialized object", got["outputSchema"])
	}
	// json.Number marshals back to its exact literal without float64 rounding.
	if got["bigNumber"] != `[9007199254740993]` {
		t.Errorf("bigNumber = %q, want exact large-integer literal", got["bigNumber"])
	}

	if got := flattenAgentloopEvaluatorConfig("not-a-map", nil); len(got) != 0 {
		t.Errorf("non map input: got %v entries, want empty map", len(got))
	}

	// Semantically equal JSON keeps the configured text verbatim, so a
	// differently formatted configuration does not become a perpetual diff;
	// semantically different JSON falls back to the canonical API form.
	configured := map[string]interface{}{
		"variables":    `[ "a" ]`,
		"outputSchema": `{"type":"object","extra":true}`,
	}
	gotConfigured := flattenAgentloopEvaluatorConfig(raw, configured)
	if gotConfigured["variables"] != `[ "a" ]` {
		t.Errorf("variables (configured) = %q, want configured text kept verbatim", gotConfigured["variables"])
	}
	if gotConfigured["outputSchema"] != `{"type":"object"}` {
		t.Errorf("outputSchema (configured) = %q, want canonical API serialization", gotConfigured["outputSchema"])
	}
}
