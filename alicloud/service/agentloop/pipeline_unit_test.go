package agentloop

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUnitAliCloudAgentloopPipelineJSONDecodeError(t *testing.T) {
	var decoded map[string]interface{}
	decodeErr := json.Unmarshal([]byte("Internal Server Error"), &decoded)
	if !isAgentloopPipelineJSONDecodeError(decodeErr) {
		t.Fatalf("expected direct JSON syntax error to be detected: %v", decodeErr)
	}
	if !isAgentloopPipelineJSONDecodeError(fmt.Errorf("wrapped decode failure: %w", decodeErr)) {
		t.Fatalf("expected wrapped JSON syntax error to be detected: %v", decodeErr)
	}
	if isAgentloopPipelineJSONDecodeError(fmt.Errorf("request rejected")) {
		t.Fatal("non-JSON error must not be treated as an ambiguous create response")
	}
}

func TestUnitAliCloudExpandAgentloopPipelineNodeParameters(t *testing.T) {
	raw := map[string]interface{}{
		"plain":     "input IS NOT NULL",
		"array":     `["a","b"]`,
		"object":    `{"x":1}`,
		"number":    "10",
		"numberPad": " 10 ",
		"bigNumber": "9007199254740993",
		"boolean":   "true",
		"nullValue": "null",
		"empty":     "",
		"blank":     "  ",
		"broken":    `[1,2`,
		"multiDoc":  `{"a":1} {"b":2}`,
		"nonString": 42,
	}
	got := expandAgentloopPipelineNodeParameters(raw)

	if got["plain"] != "input IS NOT NULL" {
		t.Errorf("plain = %#v, want verbatim string", got["plain"])
	}
	if want := []interface{}{"a", "b"}; !reflect.DeepEqual(got["array"], want) {
		t.Errorf("array = %#v, want %#v", got["array"], want)
	}
	if want := map[string]interface{}{"x": json.Number("1")}; !reflect.DeepEqual(got["object"], want) {
		t.Errorf("object = %#v, want %#v", got["object"], want)
	}
	if got["number"] != json.Number("10") {
		t.Errorf("number = %#v, want json.Number(10)", got["number"])
	}
	if got["numberPad"] != json.Number("10") {
		t.Errorf("numberPad = %#v, want json.Number(10)", got["numberPad"])
	}
	// Integers beyond the float64 precision range must reach the API unrounded.
	if got["bigNumber"] != json.Number("9007199254740993") {
		t.Errorf("bigNumber = %#v, want exact json.Number(9007199254740993)", got["bigNumber"])
	}
	if got["boolean"] != true {
		t.Errorf("boolean = %#v, want true", got["boolean"])
	}
	if v, ok := got["nullValue"]; !ok || v != nil {
		t.Errorf("nullValue = %#v (present=%v), want present nil", v, ok)
	}
	if got["empty"] != "" {
		t.Errorf("empty = %#v, want verbatim empty string", got["empty"])
	}
	if got["blank"] != "  " {
		t.Errorf("blank = %#v, want verbatim blank string", got["blank"])
	}
	if got["broken"] != `[1,2` {
		t.Errorf("broken = %#v, want verbatim invalid JSON", got["broken"])
	}
	// A second JSON value after the first must not be silently dropped: the
	// whole string stays verbatim instead of sending only the first document.
	if got["multiDoc"] != `{"a":1} {"b":2}` {
		t.Errorf("multiDoc = %#v, want verbatim multi-document string", got["multiDoc"])
	}
	if got["nonString"] != 42 {
		t.Errorf("nonString = %#v, want 42", got["nonString"])
	}
}

func TestUnitAliCloudFlattenAgentloopPipelineNodeParameters(t *testing.T) {
	t.Run("keeps configured text when semantically equal", func(t *testing.T) {
		raw := map[string]interface{}{
			"str": "plain",
			"arr": []interface{}{"a", "b"},
			"num": float64(10),
			"obj": map[string]interface{}{"chn": true, "type": "text"},
		}
		configured := map[string]interface{}{
			"str": "plain",
			"arr": `[ "a", "b" ]`,
			"num": "10",
			"obj": `{"type":"text","chn":true}`,
		}
		got := flattenAgentloopPipelineNodeParameters(raw, configured)
		if got["str"] != "plain" {
			t.Errorf("str = %q, want verbatim API string", got["str"])
		}
		if got["arr"] != `[ "a", "b" ]` {
			t.Errorf("arr = %q, want configured text kept", got["arr"])
		}
		if got["num"] != "10" {
			t.Errorf("num = %q, want configured text kept", got["num"])
		}
		if got["obj"] != `{"type":"text","chn":true}` {
			t.Errorf("obj = %q, want configured text kept", got["obj"])
		}
	})

	t.Run("restores configured json string literal", func(t *testing.T) {
		// A JSON string literal such as `"qwen-plus"` is decoded to the bare
		// string qwen-plus by expand before reaching the API. Flatten must
		// restore the configured literal so this mutable attribute does not
		// keep diffing and trigger an update on every apply.
		raw := map[string]interface{}{
			"model": "qwen-plus",
		}
		configured := map[string]interface{}{
			"model": `"qwen-plus"`,
		}
		got := flattenAgentloopPipelineNodeParameters(raw, configured)
		if got["model"] != `"qwen-plus"` {
			t.Errorf("model = %q, want configured JSON string literal kept", got["model"])
		}
	})

	t.Run("plain string round-trips unchanged", func(t *testing.T) {
		raw := map[string]interface{}{
			"model": "qwen-plus",
		}
		configured := map[string]interface{}{
			"model": "qwen-plus",
		}
		got := flattenAgentloopPipelineNodeParameters(raw, configured)
		if got["model"] != "qwen-plus" {
			t.Errorf("model = %q, want plain string unchanged", got["model"])
		}
	})

	t.Run("falls back to serialized value when different", func(t *testing.T) {
		raw := map[string]interface{}{
			"arr": []interface{}{"x"},
		}
		configured := map[string]interface{}{
			"arr": `["y"]`,
		}
		got := flattenAgentloopPipelineNodeParameters(raw, configured)
		if got["arr"] != `["x"]` {
			t.Errorf("arr = %q, want serialized API value", got["arr"])
		}
	})

	t.Run("serializes without configured (import)", func(t *testing.T) {
		raw := map[string]interface{}{
			"arr": []interface{}{"a"},
			"obj": map[string]interface{}{"chn": true},
		}
		got := flattenAgentloopPipelineNodeParameters(raw, nil)
		if got["arr"] != `["a"]` {
			t.Errorf("arr = %q, want serialized value", got["arr"])
		}
		if got["obj"] != `{"chn":true}` {
			t.Errorf("obj = %q, want serialized value", got["obj"])
		}
	})

	t.Run("json-looking strings keep their quotes without configured (import)", func(t *testing.T) {
		// jsonencode("true") / jsonencode("10") reach the API as the strings
		// true / 10. Re-expanding the bare text would send a bool / number,
		// so the imported value must be the JSON string literal.
		raw := map[string]interface{}{
			"flag":  "true",
			"limit": "10",
			"text":  "hello",
		}
		got := flattenAgentloopPipelineNodeParameters(raw, nil)
		if got["flag"] != `"true"` {
			t.Errorf("flag = %q, want JSON string literal", got["flag"])
		}
		if got["limit"] != `"10"` {
			t.Errorf("limit = %q, want JSON string literal", got["limit"])
		}
		if got["text"] != "hello" {
			t.Errorf("text = %q, want plain string unchanged", got["text"])
		}
		for k, v := range got {
			expanded := expandAgentloopPipelineNodeParameters(map[string]interface{}{k: v})
			if expanded[k] != raw[k] {
				t.Errorf("%s: re-expanded value %#v, want API value %#v", k, expanded[k], raw[k])
			}
		}
	})

	t.Run("json-encoded string literal keeps configured text", func(t *testing.T) {
		raw := map[string]interface{}{
			"flag": "true",
		}
		configured := map[string]interface{}{
			"flag": `"true"`,
		}
		got := flattenAgentloopPipelineNodeParameters(raw, configured)
		if got["flag"] != `"true"` {
			t.Errorf("flag = %q, want configured JSON string literal kept", got["flag"])
		}
	})

	t.Run("non map input returns empty map", func(t *testing.T) {
		if got := flattenAgentloopPipelineNodeParameters("not-a-map", nil); len(got) != 0 {
			t.Errorf("got %v entries, want empty map", len(got))
		}
	})

	t.Run("unmarshalable values are skipped", func(t *testing.T) {
		raw := map[string]interface{}{
			"bad": make(chan int),
		}
		got := flattenAgentloopPipelineNodeParameters(raw, nil)
		if _, ok := got["bad"]; ok {
			t.Errorf("bad should have been skipped, got %v", got["bad"])
		}
	})
}

func TestUnitAliCloudFlattenAgentloopPipelinePipeline(t *testing.T) {
	t.Run("non map returns empty slice", func(t *testing.T) {
		if got := flattenAgentloopPipelinePipeline("x", nil); len(got) != 0 {
			t.Errorf("got %v entries, want empty slice", len(got))
		}
	})

	t.Run("empty map returns empty slice", func(t *testing.T) {
		if got := flattenAgentloopPipelinePipeline(map[string]interface{}{}, nil); len(got) != 0 {
			t.Errorf("got %v entries, want empty slice", len(got))
		}
	})

	t.Run("nodes keep configured parameter text by index", func(t *testing.T) {
		raw := map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{
					"id":   "n1",
					"type": "filter",
					"parameters": map[string]interface{}{
						"filters": []interface{}{"a", "b"},
					},
				},
				"not-a-map",
			},
		}
		configured := []interface{}{
			map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"id":   "n1",
						"type": "filter",
						"parameters": map[string]interface{}{
							"filters": `[ "a", "b" ]`,
						},
					},
				},
			},
		}
		got := flattenAgentloopPipelinePipeline(raw, configured)
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		nodes, ok := got[0]["nodes"].([]map[string]interface{})
		if !ok || len(nodes) != 1 {
			t.Fatalf("nodes = %#v, want 1 node", got[0]["nodes"])
		}
		if nodes[0]["id"] != "n1" || nodes[0]["type"] != "filter" {
			t.Errorf("node = %#v, want id/type echoed", nodes[0])
		}
		params, ok := nodes[0]["parameters"].(map[string]interface{})
		if !ok {
			t.Fatalf("parameters = %#v, want map", nodes[0]["parameters"])
		}
		if params["filters"] != `[ "a", "b" ]` {
			t.Errorf("filters = %q, want configured text kept", params["filters"])
		}
	})

	t.Run("import without configured falls back to serialized", func(t *testing.T) {
		raw := map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{
					"id":   "n1",
					"type": "filter",
					"parameters": map[string]interface{}{
						"filters": []interface{}{"a"},
					},
				},
			},
		}
		got := flattenAgentloopPipelinePipeline(raw, nil)
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		nodes := got[0]["nodes"].([]map[string]interface{})
		params := nodes[0]["parameters"].(map[string]interface{})
		if params["filters"] != `["a"]` {
			t.Errorf("filters = %q, want serialized value", params["filters"])
		}
	})
}
