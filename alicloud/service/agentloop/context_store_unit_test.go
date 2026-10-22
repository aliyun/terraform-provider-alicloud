package agentloop

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	tfSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// 删除返回 Deleting 后仍须等到真实 404；未知状态也不能提前当作已删除。
func TestUnitAliCloudAgentloopContextStoreDeleteCompletion(t *testing.T) {
	for _, tc := range []struct {
		name       string
		deleteCode int
		deleteBody string
		readFails  bool
		wantReads  int
		wantError  string
	}{
		{name: "wait_until_404", wantReads: 4},
		{name: "already_absent", deleteCode: http.StatusNotFound, deleteBody: `{"code":"STORE_NOT_FOUND","message":"store absent"}`},
		{name: "delete_error", deleteCode: http.StatusForbidden, deleteBody: `{"code":"Forbidden","message":"delete denied"}`, wantError: "delete denied"},
		{name: "binding_conflict_is_not_retried", deleteCode: http.StatusConflict, deleteBody: `{"code":"TaskBindingConflict","message":"binding conflict"}`, wantError: "binding conflict"},
		{name: "read_error", readFails: true, wantReads: 1, wantError: "read denied"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deletes, reads := 0, 0
			client := agentloopDeleteTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path != "/agentspace/as/contextstore/cs" {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				switch req.Method {
				case "DELETE":
					deletes++
					if tc.deleteCode != 0 {
						w.WriteHeader(tc.deleteCode)
						_, _ = w.Write([]byte(tc.deleteBody))
						return
					}
					_, _ = w.Write([]byte(`{"status":"Deleting"}`))
				case "GET":
					reads++
					if tc.readFails {
						w.WriteHeader(http.StatusForbidden)
						_, _ = w.Write([]byte(`{"code":"Forbidden","message":"read denied"}`))
						return
					}
					if reads <= 3 {
						_, _ = w.Write([]byte(`{"status":"` + []string{"Running", "Deleting", "NewBackendState"}[reads-1] + `"}`))
						return
					}
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte(`{"code":"STORE_NOT_FOUND","message":"store absent"}`))
				default:
					t.Errorf("unexpected method: %s", req.Method)
				}
			})
			d := tfSchema.TestResourceDataRaw(t, ResourceContextStore().Schema, nil)
			d.SetId("as:cs")
			err := resourceContextStoreDelete(d, client)
			if tc.wantError == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("error = %v, want %q", err, tc.wantError)
			}
			if deletes != 1 || reads != tc.wantReads {
				t.Fatalf("DELETE=%d GET=%d, want DELETE=1 GET=%d", deletes, reads, tc.wantReads)
			}
		})
	}
}

func TestUnitAliCloudAgentloopContextStoreFlattenConfig(t *testing.T) {
	apiConfig := map[string]interface{}{
		"metadataField": map[string]interface{}{
			"defaultRecallPolicyId":     "p-1",
			"experienceMetadataVersion": "v1",
		},
		"miningInterval": "1d",
		"serviceNames":   []interface{}{"svc-b", "svc-a"},
		"source": map[string]interface{}{
			"agentSpace": "as",
			"startTime":  "2026-01-01T00:00:00Z",
		},
	}

	// Import: no prior state, metadata_field cannot be recovered.
	got := flattenAgentloopContextStoreConfig(apiConfig, []interface{}{})
	want := []map[string]interface{}{{
		"service_names": []interface{}{"svc-b", "svc-a"},
		"source": []interface{}{map[string]interface{}{
			"agent_space": "as",
			"start_time":  "2026-01-01T00:00:00Z",
		}},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("import flatten = %#v, want %#v", got, want)
	}

	// Refresh: the configured mapping and an equivalent time spelling are kept.
	prior := []interface{}{map[string]interface{}{
		"metadata_field": map[string]interface{}{"userId": "user_id"},
		"service_names":  []interface{}{"svc-b", "svc-a"},
		"source": []interface{}{map[string]interface{}{
			"agent_space": "as",
			"start_time":  "2026-01-01T08:00:00+08:00",
		}},
	}}
	got = flattenAgentloopContextStoreConfig(apiConfig, prior)
	if m := got[0]["metadata_field"]; !reflect.DeepEqual(m, map[string]interface{}{"userId": "user_id"}) {
		t.Errorf("metadata_field = %#v, want prior mapping", m)
	}
	source := got[0]["source"].([]interface{})[0].(map[string]interface{})
	if source["start_time"] != "2026-01-01T08:00:00+08:00" {
		t.Errorf("start_time = %#v, want configured spelling", source["start_time"])
	}

	// Drift: a different instant is surfaced.
	prior[0].(map[string]interface{})["source"] = []interface{}{map[string]interface{}{
		"agent_space": "as",
		"start_time":  "2026-02-01T00:00:00Z",
	}}
	got = flattenAgentloopContextStoreConfig(apiConfig, prior)
	source = got[0]["source"].([]interface{})[0].(map[string]interface{})
	if source["start_time"] != "2026-01-01T00:00:00Z" {
		t.Errorf("start_time = %#v, want API value", source["start_time"])
	}

	if got := flattenAgentloopContextStoreConfig(map[string]interface{}{"miningInterval": "1d"}, nil); len(got) != 0 {
		t.Errorf("server-only config = %#v, want empty", got)
	}
	if got := flattenAgentloopContextStoreConfig(nil, nil); len(got) != 0 {
		t.Errorf("nil config = %#v, want empty", got)
	}
}

func TestUnitAliCloudAgentloopContextStoreExpandUpdateConfig(t *testing.T) {
	got := expandAgentloopContextStoreUpdateConfig([]interface{}{map[string]interface{}{
		"metadata_field": map[string]interface{}{"userId": "user_id"},
		"service_names":  []interface{}{"svc-b", "svc-a"},
		"source": []interface{}{map[string]interface{}{
			"agent_space": "as",
			"start_time":  "2026-01-01T00:00:00Z",
		}},
	}})
	want := map[string]interface{}{
		"metadataField": map[string]interface{}{"userId": "user_id"},
		"serviceNames":  []interface{}{"svc-b", "svc-a"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expand = %#v, want %#v", got, want)
	}

	// A removed mapping is sent as an empty map so the stored user keys are cleared.
	got = expandAgentloopContextStoreUpdateConfig([]interface{}{map[string]interface{}{
		"service_names": []interface{}{"svc"},
	}})
	want = map[string]interface{}{
		"metadataField": map[string]interface{}{},
		"serviceNames":  []interface{}{"svc"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expand without mapping = %#v, want %#v", got, want)
	}
}

// TestUnitAliCloudAgentloopContextStoreImportPlan builds the state produced by
// an import Read and checks which config changes are planned in place and which
// require a new Context Store.
func TestUnitAliCloudAgentloopContextStoreImportPlan(t *testing.T) {
	r := ResourceContextStore()
	imported := map[string]interface{}{
		"agent_space":        "as",
		"context_store_name": "cs",
		"context_type":       "experience",
	}
	data := tfSchema.TestResourceDataRaw(t, r.Schema, imported)
	data.SetId("as:cs")
	if err := data.Set("config", flattenAgentloopContextStoreConfig(map[string]interface{}{
		"serviceNames": []interface{}{"svc"},
		"source":       map[string]interface{}{"agentSpace": "as", "startTime": "2026-01-01T00:00:00Z"},
	}, nil)); err != nil {
		t.Fatalf("set config: %v", err)
	}
	state := data.State()

	cfg := map[string]interface{}{
		"agent_space":        "as",
		"context_store_name": "cs",
		"context_type":       "experience",
		"config": []interface{}{map[string]interface{}{
			"service_names": []interface{}{"svc"},
			"source": []interface{}{map[string]interface{}{
				"agent_space": "as",
				"start_time":  "2026-01-01T00:00:00Z",
			}},
		}},
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff != nil && diff.RequiresNew() {
		t.Fatalf("imported config must not require replacement: %#v", diff.Attributes)
	}

	cases := []struct {
		name        string
		config      []interface{}
		requiresNew bool
	}{
		{
			name: "service_names and metadata_field update in place",
			config: []interface{}{map[string]interface{}{
				"metadata_field": map[string]interface{}{"userId": "user_id"},
				"service_names":  []interface{}{"svc-other", "svc"},
			}},
		},
		{
			name: "source change requires replacement",
			config: []interface{}{map[string]interface{}{
				"service_names": []interface{}{"svc"},
				"source": []interface{}{map[string]interface{}{
					"agent_space": "as",
					"start_time":  "2026-02-01T00:00:00Z",
				}},
			}},
			requiresNew: true,
		},
		{
			name: "clearing service_names requires replacement",
			config: []interface{}{map[string]interface{}{
				"metadata_field": map[string]interface{}{"userId": "user_id"},
			}},
			requiresNew: true,
		},
		{
			name:        "removing config requires replacement",
			requiresNew: true,
		},
	}
	for _, tc := range cases {
		cfg := map[string]interface{}{
			"agent_space":        "as",
			"context_store_name": "cs",
			"context_type":       "experience",
		}
		if tc.config != nil {
			cfg["config"] = tc.config
		}
		diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
		if err != nil {
			t.Fatalf("%s: diff error: %v", tc.name, err)
		}
		if diff == nil || diff.Empty() {
			t.Fatalf("%s: want a non-empty diff", tc.name)
		}
		if diff.RequiresNew() != tc.requiresNew {
			t.Errorf("%s: RequiresNew = %v, want %v: %#v", tc.name, diff.RequiresNew(), tc.requiresNew, diff.Attributes)
		}
	}
}
