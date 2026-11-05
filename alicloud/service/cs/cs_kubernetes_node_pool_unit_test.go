package cs

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitFlattenContainerdConfig covers the API-response-to-schema conversion of the
// containerd_config block: only keys explicitly returned by the API are written into the
// element map (missing key = "not set"), numeric fields accept both json.Number and the
// float64 fallback and are rendered as decimal strings for the nullable TypeString fields,
// and an empty/nil response yields an empty list.
func TestUnitFlattenContainerdConfig(t *testing.T) {
	tests := []struct {
		name string
		raw  interface{}
		want []map[string]interface{}
	}{
		{
			name: "nil input returns empty list",
			raw:  nil,
			want: []map[string]interface{}{},
		},
		{
			name: "empty map returns empty list",
			raw:  map[string]interface{}{},
			want: []map[string]interface{}{},
		},
		{
			name: "non-map input returns empty list",
			raw:  "not-a-map",
			want: []map[string]interface{}{},
		},
		{
			name: "partial response only writes returned keys",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": json.Number("3"),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 3,
			}},
		},
		{
			name: "json.Number values render as decimal strings",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": json.Number("5"),
				"limitCore":              json.Number("1"),
				"limitNoFile":            json.Number("65536"),
				"limitMemLock":           json.Number("128849018880"),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 5,
				"limit_core":               "1",
				"limit_no_file":            "65536",
				"limit_mem_lock":           "128849018880",
			}},
		},
		{
			name: "float64 fallback values render as decimal strings",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": float64(10),
				"limitCore":              float64(2),
				"limitNoFile":            float64(2048),
				"limitMemLock":           float64(65536),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 10,
				"limit_core":               "2",
				"limit_no_file":            "2048",
				"limit_mem_lock":           "65536",
			}},
		},
		{
			name: "bool true renders as string true",
			raw: map[string]interface{}{
				"ignoreImageDefinedVolume": true,
			},
			want: []map[string]interface{}{{
				"ignore_image_defined_volume": "true",
			}},
		},
		{
			name: "bool false renders as string false",
			raw: map[string]interface{}{
				"ignoreImageDefinedVolume": false,
			},
			want: []map[string]interface{}{{
				"ignore_image_defined_volume": "false",
			}},
		},
		{
			name: "large json.Number is not truncated",
			raw: map[string]interface{}{
				"limitNoFile": json.Number("9007199254740991"),
			},
			want: []map[string]interface{}{{
				"limit_no_file": "9007199254740991",
			}},
		},
		{
			name: "non-empty registry lists are written",
			raw: map[string]interface{}{
				"registryMirrors":    []interface{}{"https://mirror.example.com"},
				"insecureRegistries": []interface{}{"registry.example.com:5000"},
			},
			want: []map[string]interface{}{{
				"registry_mirrors":    []interface{}{"https://mirror.example.com"},
				"insecure_registries": []interface{}{"registry.example.com:5000"},
			}},
		},
		{
			name: "empty or nil registry lists are skipped",
			raw: map[string]interface{}{
				"limitCore":          json.Number("1"),
				"registryMirrors":    []interface{}{},
				"insecureRegistries": nil,
			},
			want: []map[string]interface{}{{
				"limit_core": "1",
			}},
		},
		{
			name: "explicit zero values stay distinct from missing keys",
			raw: map[string]interface{}{
				"limitCore":                json.Number("0"),
				"ignoreImageDefinedVolume": false,
			},
			want: []map[string]interface{}{{
				"limit_core":                  "0",
				"ignore_image_defined_volume": "false",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FlattenContainerdConfig(tt.raw))
		})
	}
}

// TestUnitExpandContainerdConfig covers the schema-to-API-request conversion of the
// containerd_config block using the resource's real schema via schema.TestResourceDataRaw.
// The nullable string fields (limit_core/limit_no_file/limit_mem_lock/
// ignore_image_defined_volume) skip "" ("not set") and send explicit values, including
// "0"/"false"; removing the whole block yields an empty map ({} clears the cloud-side
// configuration under the API's full-replacement semantics). Invalid values bypassing the
// schema ValidateFunc layer surface as errors from strconv.ParseInt/ParseBool.
func TestUnitExpandContainerdConfig(t *testing.T) {
	resourceSchema := ResourceAckNodepool().Schema

	tests := []struct {
		name    string
		raw     map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "block absent returns empty map",
			raw:  map[string]interface{}{},
			want: map[string]interface{}{},
		},
		{
			name: "block removed (empty list) returns empty map",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{},
			},
			want: map[string]interface{}{},
		},
		{
			name: "all fields unset sends no keys",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"max_concurrent_downloads":    0,
					"limit_core":                  "",
					"limit_no_file":               "",
					"limit_mem_lock":              "",
					"ignore_image_defined_volume": "",
					"registry_mirrors":            []interface{}{},
					"insecure_registries":         []interface{}{},
				}},
			},
			want: map[string]interface{}{},
		},
		{
			name: "explicit zero and false are sent",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core":                  "0",
					"ignore_image_defined_volume": "false",
				}},
			},
			want: map[string]interface{}{
				"limitCore":                int64(0),
				"ignoreImageDefinedVolume": false,
			},
		},
		{
			name: "normal values convert to int64, bool and lists",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"max_concurrent_downloads":    5,
					"limit_core":                  "1",
					"limit_no_file":               "65536",
					"limit_mem_lock":              "128849018880",
					"ignore_image_defined_volume": "true",
					"registry_mirrors":            []interface{}{"https://mirror.example.com"},
					"insecure_registries":         []interface{}{"registry.example.com:5000"},
				}},
			},
			want: map[string]interface{}{
				"maxConcurrentDownloads":   int64(5),
				"limitCore":                int64(1),
				"limitNoFile":              int64(65536),
				"limitMemLock":             int64(128849018880),
				"ignoreImageDefinedVolume": true,
				"registryMirrors":          []interface{}{"https://mirror.example.com"},
				"insecureRegistries":       []interface{}{"registry.example.com:5000"},
			},
		},
		{
			name: "empty registry_mirrors list omits registryMirrors key",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core":       "1",
					"registry_mirrors": []interface{}{},
				}},
			},
			want: map[string]interface{}{
				"limitCore": int64(1),
			},
		},
		{
			name: "invalid limit_core returns error",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core": "abc",
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid ignore_image_defined_volume returns error",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"ignore_image_defined_volume": "maybe",
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, tt.raw)
			got, err := ExpandContainerdConfig(d)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
