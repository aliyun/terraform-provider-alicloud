package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// polardbVpcTestClient builds a hermetic client whose vpc endpoint is a local
// test server, so the VSwitch fallback lookup can run without cloud access.
func polardbVpcTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build test credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-beijing", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build test client: %v", err)
	}
	endpoints.Store("vpc", endpoint)
	return client
}

// When zone_id is not configured, CreateApplication resolves it from the
// VSwitch; the fallback must still send the VSwitchId and take the returned
// ZoneId.
func TestUnitPolarDBApplicationVSwitchZoneFallback(t *testing.T) {
	// The handler runs in the test server goroutine; hand the value back over a
	// channel so the assertion below does not read it unsynchronized.
	queried := make(chan string, 1)
	client := polardbVpcTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.Form.Get("Action"); action != "DescribeVSwitchAttributes" {
			t.Errorf("unexpected action %q", action)
		}
		select {
		case queried <- r.Form.Get("VSwitchId"):
		default:
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"VSwitchId":"vsw-1","ZoneId":"cn-beijing-h"}`)
	})

	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"description":      "test",
		"application_type": "polarclaw",
		"architecture":     "x86",
		"pay_type":         "PostPaid",
		"region_id":        "cn-beijing",
		"vswitch_id":       "vsw-1",
		"vpc_id":           "vpc-1",
		"model_from":       "bailian",
		"model_base_url":   "https://example.invalid/compatible-mode/v1",
		"model_name":       "qwen",
	})

	request, err := NewService(client).buildApplicationCreateRequest(d)
	if err != nil {
		t.Fatalf("buildApplicationCreateRequest returned error: %v", err)
	}
	select {
	case got := <-queried:
		if got != "vsw-1" {
			t.Errorf("expected VSwitch lookup for vsw-1, got %q", got)
		}
	default:
		t.Error("expected a DescribeVSwitchAttributes call for the configured VSwitch")
	}
	if got := request["ZoneId"]; got != "cn-beijing-h" {
		t.Errorf("expected ZoneId resolved from the VSwitch, got %#v", got)
	}
}

// The rewritten lookup keeps the identity check of the original VSwitch call
// and rejects a malformed ZoneId instead of returning an empty zone.
func TestUnitPolarDBApplicationDescribeVSwitchValidation(t *testing.T) {
	for _, tc := range []struct {
		name         string
		responseBody string
		wantZone     string
		wantErr      bool
		wantNotFound bool
	}{
		{
			name:         "matching identity returns zone",
			responseBody: `{"VSwitchId":"vsw-1","ZoneId":"cn-beijing-h"}`,
			wantZone:     "cn-beijing-h",
		},
		{
			name:         "identity mismatch is not found",
			responseBody: `{"VSwitchId":"vsw-other","ZoneId":"cn-beijing-h"}`,
			wantErr:      true,
			wantNotFound: true,
		},
		{
			name:         "missing zone is an error",
			responseBody: `{"VSwitchId":"vsw-1"}`,
			wantErr:      true,
		},
		{
			name:         "non string zone is an error",
			responseBody: `{"VSwitchId":"vsw-1","ZoneId":3}`,
			wantErr:      true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := polardbVpcTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tc.responseBody)
			})

			zone, err := NewService(client).describeVSwitch("vsw-1")
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got zone %q", zone)
				}
				if got := errs.NotFoundError(err); got != tc.wantNotFound {
					t.Errorf("NotFoundError = %v, want %v (err: %v)", got, tc.wantNotFound, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if zone != tc.wantZone {
				t.Errorf("zone = %q, want %q", zone, tc.wantZone)
			}
		})
	}
}

// polardbAPITestClient builds a hermetic client whose polardb endpoint is a
// local test server, so request construction and Read can be exercised without
// cloud access.
func polardbAPITestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build test credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-beijing", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build test client: %v", err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

func decodeRequestObject(t *testing.T, request map[string]interface{}, key string) map[string]interface{} {
	t.Helper()
	raw, ok := request[key].(string)
	if !ok {
		t.Fatalf("request field %s is %#v, want a JSON string", key, request[key])
	}
	object := map[string]interface{}{}
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		t.Fatalf("decode request field %s (%s): %v", key, raw, err)
	}
	return object
}

func decodeRequestArray(t *testing.T, request map[string]interface{}, key string) []map[string]interface{} {
	t.Helper()
	raw, ok := request[key].(string)
	if !ok {
		t.Fatalf("request field %s is %#v, want a JSON string", key, request[key])
	}
	items := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("decode request field %s (%s): %v", key, raw, err)
	}
	return items
}

// The CreateApplication request carries every property_completeness parameter
// with its configured value, including Parameters, which previously reached
// only ModifyApplicationParameter. This is the API-boundary proof for the
// fields whose cloud prerequisites (NAT gateway, PolarFS storage, mem0/
// knowledge application types, agentic database) are not part of the ACC
// fixture.
func TestUnitPolarDBApplicationCreateRequestSerialization(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"description":              "tf-test",
		"application_type":         "polarclaw",
		"architecture":             "x86",
		"region_id":                "cn-beijing",
		"zone_id":                  "cn-beijing-k",
		"vswitch_id":               "vsw-1",
		"vpc_id":                   "vpc-1",
		"pay_type":                 "PostPaid",
		"db_cluster_id":            "pc-1",
		"resource_group_id":        "rg-1",
		"model_from":               "bailian",
		"ai_db_cluster_id":         "aidb-1",
		"model_api_key":            "sk-test",
		"model_base_url":           "https://dashscope.aliyuncs.com/compatible-mode/v1",
		"model_api":                "openai-completions",
		"model_name":               "qwen3.6-plus",
		"security_ip_list":         []interface{}{"127.0.0.1", "10.0.0.2"},
		"security_ip_array_name":   "wl",
		"target_version":           "v2.0.0",
		"skill_template_id":        "skill-1",
		"dnat_ip_address":          "10.0.0.8",
		"vpc_nat_gateway_id":       "ngw-1",
		"auth_provider":            "oidc",
		"auth_provider_config":     `{"issuer":"https://example.invalid"}`,
		"security_ip_type":         "ipv4",
		"auto_allocate_public_eip": true,
		"agentic_db_cluster_id":    "adb-1",
		"storages": []interface{}{map[string]interface{}{
			"storage_type":              "oss",
			"storage_instance_id":       "pfs-1",
			"storage_capacity":          "100",
			"storage_performance_level": "PL1",
			"storage_endpoint":          "oss.example.invalid",
			"endpoint_id":               "pe-1",
			"mount_path":                "/data/source",
			"container_mount_path":      "/data/container",
		}},
		"mem_application_spec": []interface{}{map[string]interface{}{
			"db_name":                         "mydb",
			"db_password":                     "secret",
			"db_user":                         "myuser",
			"embedder_model":                  "text-embedding-v4",
			"embedder_model_dimension":        1024,
			"graph_llm_model":                 "qwen-plus",
			"llm_model":                       "qwen3-max",
			"memory_type":                     "MULTIMODAL_TEXT",
			"project_name":                    "proj",
			"reranker_model":                  "qwen3-rerank",
			"shard":                           1,
			"visual_embedder_model":           "qwen3.7-text-embedding",
			"visual_embedder_model_dimension": 2560,
			"visual_llm_modalities":           []interface{}{"image"},
			"visual_llm_model":                "qwen3.6-flash",
		}},
		"agentic_db_branch_spec": []interface{}{map[string]interface{}{
			"branch_id":                "branch-1",
			"db_cluster_id":            "pc-2",
			"fork_from_application_id": "pa-parent",
			"fork_from_branch":         true,
			"project_id":               "project-1",
			"tenant_id":                "tenant-1",
		}},
		"endpoints": []interface{}{map[string]interface{}{
			"endpoint_type": "Primary",
			"description":   "my_endpoint",
		}},
		"knowledge_application_spec": []interface{}{map[string]interface{}{
			"dashboard_password": "dashboard-secret",
			"db_password":        "db-secret",
			"llm_model":          "qwen3-max",
		}},
		"dnat_entries": []interface{}{map[string]interface{}{
			"port_name":  "webui",
			"front_port": 10001,
		}},
		"parameters": []interface{}{map[string]interface{}{
			"parameter_name":  "secret.dashscope.apiKey",
			"parameter_value": "abc",
		}},
		"tags": map[string]interface{}{"Created": "allin", "Env": "test"},
	})

	request, err := NewService(nil).buildApplicationCreateRequest(d)
	if err != nil {
		t.Fatalf("buildApplicationCreateRequest returned error: %v", err)
	}

	for key, want := range map[string]interface{}{
		"RegionId":              "cn-beijing",
		"Description":           "tf-test",
		"ApplicationType":       "polarclaw",
		"Architecture":          "x86",
		"DBClusterId":           "pc-1",
		"ResourceGroupId":       "rg-1",
		"ZoneId":                "cn-beijing-k",
		"VSwitchId":             "vsw-1",
		"VpcId":                 "vpc-1",
		"PayType":               "Postpaid",
		"SecurityIPList":        "127.0.0.1,10.0.0.2",
		"SecurityIPArrayName":   "wl",
		"ModelFrom":             "bailian",
		"AIDBClusterId":         "aidb-1",
		"ModelApiKey":           "sk-test",
		"ModelBaseUrl":          "https://dashscope.aliyuncs.com/compatible-mode/v1",
		"ModelApi":              "openai-completions",
		"ModelName":             "qwen3.6-plus",
		"TargetVersion":         "v2.0.0",
		"SkillTemplateId":       "skill-1",
		"DnatIpAddress":         "10.0.0.8",
		"VpcNatGatewayId":       "ngw-1",
		"AuthProvider":          "oidc",
		"AuthProviderConfig":    `{"issuer":"https://example.invalid"}`,
		"SecurityIPType":        "ipv4",
		"AutoAllocatePublicEip": true,
		"AgenticDBClusterId":    "adb-1",
	} {
		if got := request[key]; got != want {
			t.Errorf("request[%s] = %#v, want %#v", key, got, want)
		}
	}

	parameters := decodeRequestArray(t, request, "Parameters")
	if len(parameters) != 1 || parameters[0]["ParameterName"] != "secret.dashscope.apiKey" || parameters[0]["ParameterValue"] != "abc" {
		t.Errorf("Parameters = %#v, want the configured parameter at create time", parameters)
	}

	storages := decodeRequestArray(t, request, "Storages")
	if len(storages) != 1 {
		t.Fatalf("Storages = %#v, want one member", storages)
	}
	for key, want := range map[string]interface{}{
		"StorageType":             "oss",
		"StorageInstanceId":       "pfs-1",
		"StorageCapacity":         "100",
		"StoragePerformanceLevel": "PL1",
		"StorageEndpoint":         "oss.example.invalid",
		"EndpointId":              "pe-1",
		"MountPath":               "/data/source",
		"ContainerMountPath":      "/data/container",
	} {
		if got := storages[0][key]; got != want {
			t.Errorf("Storages[0][%s] = %#v, want %#v", key, got, want)
		}
	}

	memSpec := decodeRequestObject(t, request, "MemApplicationSpec")
	for key, want := range map[string]interface{}{
		"DbName":                       "mydb",
		"DbPassword":                   "secret",
		"DbUser":                       "myuser",
		"EmbedderModel":                "text-embedding-v4",
		"EmbedderModelDimension":       float64(1024),
		"GraphLlmModel":                "qwen-plus",
		"LlmModel":                     "qwen3-max",
		"MemoryType":                   "MULTIMODAL_TEXT",
		"ProjectName":                  "proj",
		"RerankerModel":                "qwen3-rerank",
		"Shard":                        float64(1),
		"VisualEmbedderModel":          "qwen3.7-text-embedding",
		"VisualEmbedderModelDimension": float64(2560),
		"VisualLlmModel":               "qwen3.6-flash",
	} {
		if got := memSpec[key]; got != want {
			t.Errorf("MemApplicationSpec[%s] = %#v, want %#v", key, got, want)
		}
	}
	if modalities, ok := memSpec["VisualLlmModalities"].([]interface{}); !ok || len(modalities) != 1 || modalities[0] != "image" {
		t.Errorf("MemApplicationSpec[VisualLlmModalities] = %#v, want [image]", memSpec["VisualLlmModalities"])
	}

	branchSpec := decodeRequestObject(t, request, "AgenticDBBranchSpec")
	for key, want := range map[string]interface{}{
		"BranchId":              "branch-1",
		"DBClusterId":           "pc-2",
		"ForkFromApplicationId": "pa-parent",
		"ForkFromBranch":        true,
		"ProjectId":             "project-1",
		"TenantId":              "tenant-1",
	} {
		if got := branchSpec[key]; got != want {
			t.Errorf("AgenticDBBranchSpec[%s] = %#v, want %#v", key, got, want)
		}
	}

	endpoints := decodeRequestArray(t, request, "Endpoints")
	if len(endpoints) != 1 || endpoints[0]["EndpointType"] != "Primary" || endpoints[0]["Description"] != "my_endpoint" {
		t.Errorf("Endpoints = %#v, want the configured endpoint", endpoints)
	}

	knowledgeSpec := decodeRequestObject(t, request, "KnowledgeApplicationSpec")
	for key, want := range map[string]interface{}{
		"DashboardPassword": "dashboard-secret",
		"DbPassword":        "db-secret",
		"LlmModel":          "qwen3-max",
	} {
		if got := knowledgeSpec[key]; got != want {
			t.Errorf("KnowledgeApplicationSpec[%s] = %#v, want %#v", key, got, want)
		}
	}

	dnatEntries := decodeRequestArray(t, request, "DnatEntries")
	if len(dnatEntries) != 1 || dnatEntries[0]["PortName"] != "webui" || dnatEntries[0]["FrontPort"] != float64(10001) {
		t.Errorf("DnatEntries = %#v, want the configured entry", dnatEntries)
	}

	tagValues := map[string]string{}
	for i := 1; i <= 2; i++ {
		key, _ := request[fmt.Sprintf("Tag.%d.Key", i)].(string)
		value, _ := request[fmt.Sprintf("Tag.%d.Value", i)].(string)
		tagValues[key] = value
	}
	if tagValues["Created"] != "allin" || tagValues["Env"] != "test" {
		t.Errorf("Tag.N = %#v, want Created=allin and Env=test", tagValues)
	}
}

// CreateApplication requires a mem0 application to carry MemApplicationSpec
// with its model names, and a knowledge application to carry
// KnowledgeApplicationSpec.LlmModel. Both are checked before the request is
// built so the provider reports the missing item instead of forwarding an
// incomplete request.
func TestUnitPolarDBApplicationConditionalSpecValidation(t *testing.T) {
	memSpec := func(overrides map[string]interface{}) []interface{} {
		item := map[string]interface{}{
			"embedder_model": "text-embedding-v4",
			"llm_model":      "qwen3-max",
			"reranker_model": "qwen3-rerank",
		}
		for key, value := range overrides {
			item[key] = value
		}
		return []interface{}{item}
	}
	knowledgeSpec := func(llmModel string) []interface{} {
		// dashboard_password keeps the block non-empty; an all-empty block is
		// dropped by the SDK diff and reads back as "no block".
		return []interface{}{map[string]interface{}{"llm_model": llmModel, "dashboard_password": "dashboard-secret"}}
	}
	base := func(applicationType string) map[string]interface{} {
		return map[string]interface{}{
			"application_type": applicationType,
			"architecture":     "x86",
			"region_id":        "cn-beijing",
			"zone_id":          "cn-beijing-k",
		}
	}

	for _, tc := range []struct {
		name    string
		raw     map[string]interface{}
		wantErr string
	}{
		{
			name:    "mem0 without mem_application_spec",
			raw:     base("mem0"),
			wantErr: "mem_application_spec",
		},
		{
			name: "mem0 without reranker_model",
			raw: func() map[string]interface{} {
				raw := base("mem0")
				raw["mem_application_spec"] = memSpec(map[string]interface{}{"reranker_model": ""})
				return raw
			}(),
			wantErr: "reranker_model",
		},
		{
			name: "mem0 with every required model",
			raw: func() map[string]interface{} {
				raw := base("mem0")
				raw["mem_application_spec"] = memSpec(nil)
				return raw
			}(),
		},
		{
			name:    "knowledge without knowledge_application_spec",
			raw:     base("knowledge"),
			wantErr: "knowledge_application_spec",
		},
		{
			name: "knowledge without llm_model",
			raw: func() map[string]interface{} {
				raw := base("knowledge")
				raw["knowledge_application_spec"] = knowledgeSpec("")
				return raw
			}(),
			wantErr: "llm_model",
		},
		{
			name: "knowledge with llm_model",
			raw: func() map[string]interface{} {
				raw := base("knowledge")
				raw["knowledge_application_spec"] = knowledgeSpec("qwen3-max")
				return raw
			}(),
		},
		{
			name: "polarclaw needs neither block",
			raw:  base("polarclaw"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, tc.raw)
			_, err := NewService(nil).buildApplicationCreateRequest(d)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error mentioning %q, got none", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not mention %q", err.Error(), tc.wantErr)
			}
		})
	}
}

// Read restores the create-only fields DescribeApplicationAttribute returns, so
// an imported application no longer plans a replacement against an empty state,
// while write-only members keep their configured value.
func TestUnitPolarDBApplicationReadCreateOnlyAttributes(t *testing.T) {
	const describeResponse = `{"Status":"Activated","VSwitchId":"vsw-test",
		"VpcNatGatewayId":"ngw-remote",
		"Storages":[{"StorageType":"oss","StorageInstanceId":"pfs-1","StorageCapacity":"100Gi","StoragePerformanceLevel":"PL1"}],
		"MemApplicationAttribute":{"DbName":"mydb","UserName":"myuser","LlmModelName":"qwen3-max\n",
			"RerankerModelName":"qwen3-rerank","ProjectName":" proj \n","MemoryType":"MULTIMODAL_TEXT",
			"VisualEmbedderDimensions":2560,"VisualLlmModalities":["image"],"VisualLlmModelName":"qwen3.6-flash"}}`

	t.Run("configured members keep write-only fields and take readable values", func(t *testing.T) {
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				fmt.Fprint(w, describeResponse)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
			"application_type":   "polarclaw",
			"architecture":       "x86",
			"region_id":          "cn-beijing",
			"vpc_nat_gateway_id": "ngw-configured",
			"storages": []interface{}{map[string]interface{}{
				"storage_type":         "oss",
				"storage_instance_id":  "pfs-1",
				"storage_capacity":     "100",
				"storage_endpoint":     "oss.example.invalid",
				"endpoint_id":          "pe-1",
				"mount_path":           "/data/source",
				"container_mount_path": "/data/container",
			}},
			"mem_application_spec": []interface{}{map[string]interface{}{
				"db_password":              "secret",
				"embedder_model_dimension": 1024,
				"shard":                    2,
				"llm_model":                "qwen3-max",
			}},
		})
		d.SetId("app-test")

		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}

		if got := d.Get("vpc_nat_gateway_id"); got != "ngw-remote" {
			t.Errorf("vpc_nat_gateway_id = %#v, want ngw-remote", got)
		}

		storages, _ := d.Get("storages").(*schema.Set)
		if storages == nil || storages.Len() != 1 {
			t.Fatalf("storages = %#v, want one member", d.Get("storages"))
		}
		storage := storages.List()[0].(map[string]interface{})
		for key, want := range map[string]interface{}{
			"storage_type":        "oss",
			"storage_instance_id": "pfs-1",
			// The configured capacity spelling is kept and the omitted
			// performance level is filled from the response.
			"storage_capacity":          "100",
			"storage_performance_level": "PL1",
			"storage_endpoint":          "oss.example.invalid",
			"endpoint_id":               "pe-1",
			"mount_path":                "/data/source",
			"container_mount_path":      "/data/container",
		} {
			if got := storage[key]; got != want {
				t.Errorf("storages[%s] = %#v, want %#v", key, got, want)
			}
		}

		memItems, _ := d.Get("mem_application_spec").([]interface{})
		if len(memItems) != 1 {
			t.Fatalf("mem_application_spec = %#v, want one member", d.Get("mem_application_spec"))
		}
		mem := memItems[0].(map[string]interface{})
		for key, want := range map[string]interface{}{
			"db_password":                     "secret",
			"embedder_model_dimension":        1024,
			"shard":                           2,
			"db_name":                         "mydb",
			"db_user":                         "myuser",
			"llm_model":                       "qwen3-max",
			"reranker_model":                  "qwen3-rerank",
			"project_name":                    "proj",
			"memory_type":                     "MULTIMODAL_TEXT",
			"visual_embedder_model_dimension": 2560,
			"visual_llm_model":                "qwen3.6-flash",
		} {
			if got := mem[key]; got != want {
				t.Errorf("mem_application_spec[%s] = %#v, want %#v", key, got, want)
			}
		}
		modalities, _ := mem["visual_llm_modalities"].(*schema.Set)
		if modalities == nil || modalities.Len() != 1 || !modalities.Contains("image") {
			t.Errorf("mem_application_spec[visual_llm_modalities] = %#v, want [image]", mem["visual_llm_modalities"])
		}
	})

	t.Run("import restores readable members from an empty state", func(t *testing.T) {
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				fmt.Fprint(w, describeResponse)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
			"application_type": "polarclaw",
			"architecture":     "x86",
			"region_id":        "cn-beijing",
		})
		d.SetId("app-imported")

		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}

		if got := d.Get("vpc_nat_gateway_id"); got != "ngw-remote" {
			t.Errorf("vpc_nat_gateway_id = %#v, want ngw-remote", got)
		}
		storages, _ := d.Get("storages").(*schema.Set)
		if storages == nil || storages.Len() != 1 {
			t.Fatalf("storages = %#v, want one restored member", d.Get("storages"))
		}
		storage := storages.List()[0].(map[string]interface{})
		if storage["storage_instance_id"] != "pfs-1" || storage["storage_capacity"] != "100Gi" || storage["storage_performance_level"] != "PL1" {
			t.Errorf("restored storage = %#v, want the readable response fields", storage)
		}
		memItems, _ := d.Get("mem_application_spec").([]interface{})
		if len(memItems) != 1 {
			t.Fatalf("mem_application_spec = %#v, want one restored member", d.Get("mem_application_spec"))
		}
		mem := memItems[0].(map[string]interface{})
		if mem["db_name"] != "mydb" || mem["llm_model"] != "qwen3-max" || mem["visual_embedder_model_dimension"] != 2560 {
			t.Errorf("restored mem_application_spec = %#v, want the readable response fields", mem)
		}
	})

	t.Run("an explicitly empty gateway clears the state field", func(t *testing.T) {
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test","VpcNatGatewayId":""}`)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
			"application_type":   "polarclaw",
			"architecture":       "x86",
			"region_id":          "cn-beijing",
			"vpc_nat_gateway_id": "ngw-configured",
		})
		d.SetId("app-detached")

		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}
		if got := d.Get("vpc_nat_gateway_id"); got != "" {
			t.Errorf("vpc_nat_gateway_id = %#v, want it cleared", got)
		}
	})

	t.Run("an omitted gateway key keeps the configured value", func(t *testing.T) {
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test"}`)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
			"application_type":   "polarclaw",
			"architecture":       "x86",
			"region_id":          "cn-beijing",
			"vpc_nat_gateway_id": "ngw-configured",
		})
		d.SetId("app-fresh")

		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}
		if got := d.Get("vpc_nat_gateway_id"); got != "ngw-configured" {
			t.Errorf("vpc_nat_gateway_id = %#v, want the configured value kept", got)
		}
	})
}

// The storages hash covers every member attribute. SDKv2 skips the per-member
// diff when the old and new set hashes match, so an attribute left out of the
// hash would make its configuration change invisible and would let two members
// that differ only in it collapse into one. Both readable ForceNew attributes
// below therefore have to change the identity.
func TestUnitPolarDBApplicationStorageHashCoversEveryAttribute(t *testing.T) {
	base := map[string]interface{}{
		"mount_path": "/data/source", "container_mount_path": "/data/container",
		"storage_endpoint": "oss.example.invalid", "endpoint_id": "pe-1",
		"storage_type": "essd", "storage_instance_id": "pfs-1",
		"storage_capacity": "100Gi", "storage_performance_level": "PL1",
	}
	withCapacity := map[string]interface{}{}
	for key, value := range base {
		withCapacity[key] = value
	}
	withCapacity["storage_capacity"] = "200Gi"
	if storageMemberHash(base) == storageMemberHash(withCapacity) {
		t.Error("members that differ only in storage_capacity must hash differently")
	}
	withPerformance := map[string]interface{}{}
	for key, value := range base {
		withPerformance[key] = value
	}
	withPerformance["storage_performance_level"] = "PL2"
	if storageMemberHash(base) == storageMemberHash(withPerformance) {
		t.Error("members that differ only in storage_performance_level must hash differently")
	}
	other := map[string]interface{}{}
	for key, value := range base {
		other[key] = value
	}
	other["mount_path"] = "/data/other"
	if storageMemberHash(base) == storageMemberHash(other) {
		t.Error("members with different mount paths must hash differently")
	}
	// The same attributes Read back-fills make the hash differ from a
	// configuration that omits them; that difference is what
	// resourceApplicationCustomizeDiff removes at plan time instead of being
	// hidden by the hash.
	backFilled := map[string]interface{}{"mount_path": "/data/source"}
	if storageMemberHash(backFilled) == storageMemberHash(base) {
		t.Error("a member carrying reported values must hash differently from a configuration that omits them")
	}
}

// polardbStorageMember builds a storages member with the given overrides.
func polardbStorageMember(overrides map[string]interface{}) map[string]interface{} {
	member := map[string]interface{}{}
	for key, value := range overrides {
		member[key] = value
	}
	return member
}

// polardbStorageBackFilledMember is a member as Read leaves it in state: the
// configured write-only attributes plus the values DescribeApplicationAttribute
// reports.
func polardbStorageBackFilledMember() map[string]interface{} {
	return polardbStorageMember(map[string]interface{}{
		"mount_path":                "/data/source",
		"container_mount_path":      "/data/container",
		"storage_endpoint":          "oss.example.invalid",
		"endpoint_id":               "pe-1",
		"storage_type":              "essd",
		"storage_instance_id":       "pfs-1",
		"storage_capacity":          "100Gi",
		"storage_performance_level": "PL1",
	})
}

// polardbPlanRawConfig converts a raw configuration map into the cty value
// Terraform hands to the provider as RawConfig during plan.
func polardbPlanRawConfig(t *testing.T, raw map[string]interface{}) cty.Value {
	t.Helper()
	implied := ResourceApplication().CoreConfigSchema().ImpliedType()
	values := make(map[string]cty.Value, len(implied.AttributeTypes()))
	for name, attributeType := range implied.AttributeTypes() {
		values[name] = cty.NullVal(attributeType)
	}
	for name, value := range raw {
		attributeType := implied.AttributeType(name)
		if attributeType == cty.String {
			text, ok := value.(string)
			if !ok {
				t.Fatalf("raw config %s: want string, got %T", name, value)
			}
			values[name] = cty.StringVal(text)
			continue
		}
		if name != "storages" {
			t.Fatalf("raw config attribute %q is not supported by this helper", name)
		}
		elementType := attributeType.ElementType()
		members, ok := value.([]interface{})
		if !ok {
			t.Fatalf("raw config %s: want []interface{}, got %T", name, value)
		}
		elements := make([]cty.Value, 0, len(members))
		for _, member := range members {
			memberMap, ok := member.(map[string]interface{})
			if !ok {
				t.Fatalf("raw config %s member: want map, got %T", name, member)
			}
			attributes := make(map[string]cty.Value, len(elementType.AttributeTypes()))
			for attribute, attributeType := range elementType.AttributeTypes() {
				attributes[attribute] = cty.NullVal(attributeType)
			}
			for attribute, attributeValue := range memberMap {
				text, ok := attributeValue.(string)
				if !ok {
					t.Fatalf("raw config %s.%s: want string, got %T", name, attribute, attributeValue)
				}
				attributes[attribute] = cty.StringVal(text)
			}
			elements = append(elements, cty.ObjectVal(attributes))
		}
		values[name] = cty.SetVal(elements)
	}
	return cty.ObjectVal(values)
}

// polardbPlanDiff runs Resource.Diff — the plan path Terraform's
// PlanResourceChange uses, including the set hash early return and the
// CustomizeDiff hook — for a configuration against a synthetic prior state.
func polardbPlanDiff(t *testing.T, stateStorages, configStorages []interface{}) *sdktf.InstanceDiff {
	t.Helper()
	build := func(storages []interface{}) map[string]interface{} {
		raw := map[string]interface{}{
			"application_type": "polarclaw",
			"architecture":     "x86",
			"region_id":        "cn-beijing",
		}
		if storages != nil {
			raw["storages"] = storages
		}
		return raw
	}
	configRaw := build(configStorages)
	stateData := schema.TestResourceDataRaw(t, ResourceApplication().Schema, build(stateStorages))
	stateData.SetId("app-plan-test")
	state := stateData.State()
	// Read persists these empty collections, so a real state carries their
	// counts; without them the unrelated optional+computed blocks would
	// dominate the diff.
	state.Attributes["endpoints.#"] = "0"
	state.Attributes["mem_application_spec.#"] = "0"
	state.Attributes["security_ip_list.#"] = "0"
	// PlanResourceChange copies the new configuration into the prior state so
	// GetRawConfig sees it; do the same here.
	state.RawConfig = polardbPlanRawConfig(t, configRaw)
	diff, err := ResourceApplication().Diff(context.Background(), state, sdktf.NewResourceConfigRaw(configRaw), nil)
	if err != nil {
		t.Fatalf("plan diff: %v", err)
	}
	return diff
}

// polardbPlanIsEmpty reports whether the plan contains no change at all. The
// SDK returns a nil diff when nothing changed.
func polardbPlanIsEmpty(diff *sdktf.InstanceDiff) bool {
	return diff == nil || diff.Empty()
}

// polardbStoragesMarkedForReplacement reports whether a storages attribute
// carries the replacement marker. The set count itself only changes when the
// number of members changes, so the member attributes are checked.
func polardbStoragesMarkedForReplacement(diff *sdktf.InstanceDiff) bool {
	if diff == nil {
		return false
	}
	for key, attribute := range diff.Attributes {
		if attribute != nil && attribute.RequiresNew && strings.HasPrefix(key, "storages.") {
			return true
		}
	}
	return false
}

// TestUnitPolarDBApplicationStoragesPlanBackfill proves with the SDK's real set
// diff that the attributes Read back-fills do not plan a replacement the
// configuration never asked for.
func TestUnitPolarDBApplicationStoragesPlanBackfill(t *testing.T) {
	backFilled := polardbStorageBackFilledMember()
	for name, configured := range map[string]map[string]interface{}{
		"identical configuration": backFilled,
		"readable attributes omitted": polardbStorageMember(map[string]interface{}{
			"mount_path":           "/data/source",
			"container_mount_path": "/data/container",
			"storage_endpoint":     "oss.example.invalid",
			"endpoint_id":          "pe-1",
			"storage_type":         "essd",
		}),
		"only write-only attributes configured": polardbStorageMember(map[string]interface{}{
			"mount_path":           "/data/source",
			"container_mount_path": "/data/container",
			"storage_endpoint":     "oss.example.invalid",
			"endpoint_id":          "pe-1",
		}),
	} {
		t.Run(name, func(t *testing.T) {
			diff := polardbPlanDiff(t, []interface{}{backFilled}, []interface{}{configured})
			if !polardbPlanIsEmpty(diff) {
				t.Fatalf("plan must be empty when the configuration expresses no change, got %#v", diff.Attributes)
			}
		})
	}

	t.Run("storages block omitted", func(t *testing.T) {
		diff := polardbPlanDiff(t, []interface{}{backFilled}, nil)
		if !polardbPlanIsEmpty(diff) {
			t.Fatalf("plan must be empty when storages is not configured, got %#v", diff.Attributes)
		}
	})
}

// TestUnitPolarDBApplicationStoragesPlanConfiguredChanges proves that changing
// a configured readable attribute is visible in the plan and forces the
// replacement its ForceNew flag promises.
func TestUnitPolarDBApplicationStoragesPlanConfiguredChanges(t *testing.T) {
	backFilled := polardbStorageBackFilledMember()
	writeOnly := polardbStorageMember(map[string]interface{}{
		"mount_path":           "/data/source",
		"container_mount_path": "/data/container",
		"storage_endpoint":     "oss.example.invalid",
		"endpoint_id":          "pe-1",
	})
	for name, change := range map[string]map[string]interface{}{
		"storage_type":              {"storage_type": "standard"},
		"storage_instance_id":       {"storage_instance_id": "pfs-2"},
		"storage_capacity":          {"storage_capacity": "200Gi"},
		"storage_performance_level": {"storage_performance_level": "PL2"},
	} {
		t.Run(name, func(t *testing.T) {
			configured := polardbStorageMember(writeOnly)
			for key, value := range change {
				configured[key] = value
			}
			diff := polardbPlanDiff(t, []interface{}{backFilled}, []interface{}{configured})
			if polardbPlanIsEmpty(diff) {
				t.Fatalf("changing %s planned no change", name)
			}
			if !diff.RequiresNew() {
				t.Fatalf("changing %s must plan a replacement", name)
			}
			if !polardbStoragesMarkedForReplacement(diff) {
				t.Fatalf("storages must carry the replacement marker, got %#v", diff.Attributes)
			}
		})
	}
}

// TestUnitPolarDBApplicationStoragesPlanCapacityForms proves that an equivalent
// capacity spelling plans no change while a real capacity change does.
func TestUnitPolarDBApplicationStoragesPlanCapacityForms(t *testing.T) {
	backFilled := polardbStorageBackFilledMember()
	equivalent := polardbStorageMember(map[string]interface{}{
		"mount_path":           "/data/source",
		"container_mount_path": "/data/container",
		"storage_endpoint":     "oss.example.invalid",
		"endpoint_id":          "pe-1",
		"storage_capacity":     "100",
	})
	if diff := polardbPlanDiff(t, []interface{}{backFilled}, []interface{}{equivalent}); !polardbPlanIsEmpty(diff) {
		t.Fatalf("equivalent capacity spelling must plan no change, got %#v", diff.Attributes)
	}

	changed := polardbStorageMember(equivalent)
	changed["storage_capacity"] = "200"
	diff := polardbPlanDiff(t, []interface{}{backFilled}, []interface{}{changed})
	if polardbPlanIsEmpty(diff) {
		t.Fatal("a real capacity change must be visible in the plan")
	}
	if !diff.RequiresNew() {
		t.Fatal("a capacity change must plan a replacement")
	}
}

// TestUnitPolarDBApplicationStoragesPlanMembersDistinct proves that two members
// that differ only in a readable attribute keep distinct identities and stay
// two members in state and plan.
func TestUnitPolarDBApplicationStoragesPlanMembersDistinct(t *testing.T) {
	first := polardbStorageBackFilledMember()
	second := polardbStorageMember(first)
	second["storage_capacity"] = "200Gi"
	second["storage_instance_id"] = "pfs-2"

	if storageMemberHash(first) == storageMemberHash(second) {
		t.Fatal("members differing only in readable attributes must hash differently")
	}
	stateData := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"application_type": "polarclaw",
		"architecture":     "x86",
		"region_id":        "cn-beijing",
		"storages":         []interface{}{first, second},
	})
	stateData.SetId("app-plan-test")
	if got := stateData.Get("storages").(*schema.Set).Len(); got != 2 {
		t.Fatalf("state holds %d storages members, want 2 (members must not collapse)", got)
	}
	if diff := polardbPlanDiff(t, []interface{}{first, second}, []interface{}{first, second}); !polardbPlanIsEmpty(diff) {
		t.Fatalf("identical two-member configuration must plan no change, got %#v", diff.Attributes)
	}
	if diff := polardbPlanDiff(t, []interface{}{first, second}, []interface{}{first}); polardbPlanIsEmpty(diff) {
		t.Fatal("dropping a configured member must be visible in the plan")
	}
}

// Read never rewrites a configured storage member: its values (including the
// capacity spelling) stay exactly as configured, so the member identity is
// stable and a configuration change stays visible. Only an application with no
// configured storages adopts the response (import / service default), and a
// configured member the response no longer reports is dropped when it can be
// identified.
func TestUnitPolarDBApplicationReadStorageReconciliation(t *testing.T) {
	read := func(t *testing.T, rawStorages []interface{}, storages []interface{}) *schema.ResourceData {
		t.Helper()
		body := map[string]interface{}{"Status": "Activated", "VSwitchId": "vsw-test", "Storages": storages}
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode response: %v", err)
		}
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				_, _ = w.Write(encoded)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		raw := map[string]interface{}{
			"application_type": "polarclaw",
			"architecture":     "x86",
			"region_id":        "cn-beijing",
		}
		if rawStorages != nil {
			raw["storages"] = rawStorages
		}
		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, raw)
		d.SetId("app-storages")
		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}
		return d
	}

	storagesOf := func(t *testing.T, d *schema.ResourceData) []map[string]interface{} {
		t.Helper()
		set, ok := d.Get("storages").(*schema.Set)
		if !ok {
			t.Fatalf("storages is %T, want *schema.Set", d.Get("storages"))
		}
		out := make([]map[string]interface{}, 0, set.Len())
		for _, item := range set.List() {
			out = append(out, item.(map[string]interface{}))
		}
		return out
	}

	t.Run("a configured member keeps its values and fills the omitted readable fields", func(t *testing.T) {
		d := read(t,
			[]interface{}{map[string]interface{}{
				"storage_type": "oss", "storage_instance_id": "pfs-1", "storage_capacity": "100",
				"mount_path": "/data/source", "container_mount_path": "/data/container",
			}},
			[]interface{}{map[string]interface{}{
				"StorageType": "oss", "StorageInstanceId": "pfs-1",
				"StorageCapacity": "100Gi", "StoragePerformanceLevel": "PL1",
			}})
		got := storagesOf(t, d)
		if len(got) != 1 {
			t.Fatalf("storages = %#v, want one member", got)
		}
		for key, want := range map[string]interface{}{
			// Configured values win; the omitted performance level is filled.
			"storage_capacity":          "100",
			"storage_performance_level": "PL1",
			"storage_instance_id":       "pfs-1",
			"mount_path":                "/data/source",
		} {
			if got[0][key] != want {
				t.Errorf("storages[%s] = %#v, want %#v", key, got[0][key], want)
			}
		}
	})

	t.Run("an application without configured storages adopts the response", func(t *testing.T) {
		d := read(t, nil, []interface{}{map[string]interface{}{
			"StorageType": "essd", "StorageInstanceId": "pfs-remote",
			"StorageCapacity": "20Gi", "StoragePerformanceLevel": "PL0",
		}})
		got := storagesOf(t, d)
		if len(got) != 1 {
			t.Fatalf("storages = %#v, want the response member adopted", got)
		}
		for key, want := range map[string]interface{}{
			"storage_type": "essd", "storage_instance_id": "pfs-remote",
			"storage_capacity": "20Gi", "storage_performance_level": "PL0",
		} {
			if got[0][key] != want {
				t.Errorf("storages[%s] = %#v, want %#v", key, got[0][key], want)
			}
		}
	})

	t.Run("same-type same-capacity members keep their own mount fields", func(t *testing.T) {
		d := read(t,
			[]interface{}{
				map[string]interface{}{"storage_type": "essd", "storage_capacity": "100", "mount_path": "/data/one", "container_mount_path": "/one"},
				map[string]interface{}{"storage_type": "essd", "storage_capacity": "100", "mount_path": "/data/two", "container_mount_path": "/two"},
			},
			[]interface{}{
				map[string]interface{}{"StorageType": "essd", "StorageInstanceId": "pfs-b", "StorageCapacity": "100Gi"},
				map[string]interface{}{"StorageType": "essd", "StorageInstanceId": "pfs-a", "StorageCapacity": "100Gi"},
			})
		got := storagesOf(t, d)
		if len(got) != 2 {
			t.Fatalf("storages = %#v, want both members kept", got)
		}
		mounts := map[string]string{}
		for _, item := range got {
			mounts[item["mount_path"].(string)] = item["container_mount_path"].(string)
		}
		if mounts["/data/one"] != "/one" || mounts["/data/two"] != "/two" {
			t.Errorf("members cross-merged their mount fields: %#v", got)
		}
	})

	t.Run("an empty response drops identifiable members", func(t *testing.T) {
		d := read(t,
			[]interface{}{map[string]interface{}{"storage_type": "oss", "storage_instance_id": "pfs-1", "mount_path": "/data/source"}},
			[]interface{}{})
		if got := storagesOf(t, d); len(got) != 0 {
			t.Fatalf("storages = %#v, want the removed member dropped", got)
		}
	})

	t.Run("a member the response no longer reports is dropped", func(t *testing.T) {
		d := read(t,
			[]interface{}{
				map[string]interface{}{"storage_type": "oss", "storage_instance_id": "pfs-1", "mount_path": "/data/one"},
				map[string]interface{}{"storage_type": "essd", "storage_instance_id": "pfs-2", "mount_path": "/data/two"},
			},
			[]interface{}{map[string]interface{}{"StorageType": "oss", "StorageInstanceId": "pfs-1", "StorageCapacity": "100Gi"}})
		got := storagesOf(t, d)
		if len(got) != 1 {
			t.Fatalf("storages = %#v, want only the member the cloud still reports", got)
		}
		if got[0]["storage_instance_id"] != "pfs-1" || got[0]["mount_path"] != "/data/one" {
			t.Errorf("storages[0] = %#v, want pfs-1 with its mount path", got[0])
		}
	})

	t.Run("an ambiguous member is kept instead of paired by order", func(t *testing.T) {
		d := read(t,
			[]interface{}{map[string]interface{}{"storage_type": "oss", "mount_path": "/data/source"}},
			[]interface{}{
				map[string]interface{}{"StorageType": "oss", "StorageInstanceId": "pfs-b"},
				map[string]interface{}{"StorageType": "oss", "StorageInstanceId": "pfs-a"},
			})
		got := storagesOf(t, d)
		if len(got) != 1 {
			t.Fatalf("storages = %#v, want the ambiguous member kept", got)
		}
		if got[0]["storage_instance_id"] != "" || got[0]["mount_path"] != "/data/source" {
			t.Errorf("storages[0] = %#v, want the configured member untouched", got[0])
		}
	})
}

// security_ip_type and the endpoint description are readable through
// DescribeApplicationAttribute; Read restores them so an import does not plan a
// replacement against empty state, and a server-side endpoint description
// change surfaces as drift.
func TestUnitPolarDBApplicationReadWhitelistTypeAndEndpointDescription(t *testing.T) {
	client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeApplicationAttribute":
			fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test",
				"SecurityIPArrays":[
					{"SecurityIPArrayName":"white_list","SecurityIPType":"ipv6","SecurityIPList":"::1"},
					{"SecurityIPArrayName":"default","SecurityIPType":"ipv4","SecurityIPList":"127.0.0.1"}],
				"Endpoints":[{"Description":"server-side-description","NetType":"Private"}]}`)
		case "ListTagResources":
			fmt.Fprint(w, `{"TagResources":{}}`)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"application_type":       "polarclaw",
		"architecture":           "x86",
		"region_id":              "cn-beijing",
		"security_ip_array_name": "white_list",
		"endpoints":              []interface{}{map[string]interface{}{"endpoint_type": "Primary", "description": "configured-description"}},
	})
	d.SetId("app-whitelist")
	if err := ResourceApplication().Read(d, client); err != nil {
		t.Fatalf("resource read returned error: %v", err)
	}

	if got := d.Get("security_ip_type"); got != "ipv6" {
		t.Errorf("security_ip_type = %#v, want ipv6 from the white_list array", got)
	}
	endpoints, ok := d.Get("endpoints").(*schema.Set)
	if !ok || endpoints.Len() != 1 {
		t.Fatalf("endpoints = %#v, want one member", d.Get("endpoints"))
	}
	member := endpoints.List()[0].(map[string]interface{})
	if member["description"] != "server-side-description" {
		t.Errorf("endpoints description = %#v, want the server-side value", member["description"])
	}
	if member["endpoint_type"] != "Primary" {
		t.Errorf("endpoints endpoint_type = %#v, want the configured value kept", member["endpoint_type"])
	}
}

// DescribeApplicationAttribute reports the application whitelist entry and the
// component id, so Read restores them: an import must not lose the whitelist
// name, its IP list, the IP type or the component reference.
func TestUnitPolarDBApplicationReadWhitelistAndComponent(t *testing.T) {
	client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeApplicationAttribute":
			fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test",
				"SecurityIPArrays":[
					{"SecurityIPArrayName":"white_list","SecurityIPType":"ipv4","SecurityIPList":"127.0.0.1, 10.0.0.2"},
					{"SecurityIPArrayName":"default","SecurityIPType":"ipv4","SecurityIPList":"127.0.0.1"}],
				"Components":[{"ComponentType":"polarclaw","ComponentClass":"polar.app.g2.medium","ComponentReplica":1,"ComponentId":"pac-1"}]}`)
		case "ListTagResources":
			fmt.Fprint(w, `{"TagResources":{}}`)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"application_type":       "polarclaw",
		"architecture":           "x86",
		"region_id":              "cn-beijing",
		"security_ip_array_name": "white_list",
	})
	d.SetId("app-whitelist-component")
	if err := ResourceApplication().Read(d, client); err != nil {
		t.Fatalf("resource read returned error: %v", err)
	}

	if got := d.Get("security_ip_array_name"); got != "white_list" {
		t.Errorf("security_ip_array_name = %#v, want white_list", got)
	}
	ipList, _ := d.Get("security_ip_list").([]interface{})
	if len(ipList) != 2 || ipList[0] != "127.0.0.1" || ipList[1] != "10.0.0.2" {
		t.Errorf("security_ip_list = %#v, want the whitelist entry's IPs", d.Get("security_ip_list"))
	}
	if got := d.Get("security_ip_type"); got != "ipv4" {
		t.Errorf("security_ip_type = %#v, want ipv4", got)
	}
	if got := d.Get("component_id"); got != "pac-1" {
		t.Errorf("component_id = %#v, want pac-1", got)
	}
}

// An imported application has no configuration, so Read adopts the reported
// endpoint descriptions; the CreateApplication contract fixes the endpoint type
// to Primary.
func TestUnitPolarDBApplicationReadAdoptsEndpointsOnImport(t *testing.T) {
	client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeApplicationAttribute":
			fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test",
				"Endpoints":[{"Description":"my_endpoint","NetType":"Private"}]}`)
		case "ListTagResources":
			fmt.Fprint(w, `{"TagResources":{}}`)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"application_type": "polarclaw",
		"architecture":     "x86",
		"region_id":        "cn-beijing",
	})
	d.SetId("app-endpoints")
	if err := ResourceApplication().Read(d, client); err != nil {
		t.Fatalf("resource read returned error: %v", err)
	}

	endpoints, ok := d.Get("endpoints").(*schema.Set)
	if !ok || endpoints.Len() != 1 {
		t.Fatalf("endpoints = %#v, want the reported endpoint adopted", d.Get("endpoints"))
	}
	member := endpoints.List()[0].(map[string]interface{})
	if member["endpoint_type"] != "Primary" || member["description"] != "my_endpoint" {
		t.Errorf("endpoints member = %#v, want Primary with the reported description", member)
	}
}

// An import has no configured name; the single non-default (custom) whitelist
// array is the one a configuration would have created, so Read attributes that
// one instead of the implicit default entry.
func TestUnitPolarDBApplicationReadPrefersCustomWhitelistOnImport(t *testing.T) {
	client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DescribeApplicationAttribute":
			fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test",
				"SecurityIPArrays":[
					{"SecurityIPArrayName":"default","SecurityIPType":"ipv4","SecurityIPList":"127.0.0.1"},
					{"SecurityIPArrayName":"test","SecurityIPType":"ipv4","SecurityIPList":"10.0.0.2"}]}`)
		case "ListTagResources":
			fmt.Fprint(w, `{"TagResources":{}}`)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, map[string]interface{}{
		"application_type": "polarclaw",
		"architecture":     "x86",
		"region_id":        "cn-beijing",
	})
	d.SetId("app-import-whitelist")
	if err := ResourceApplication().Read(d, client); err != nil {
		t.Fatalf("resource read returned error: %v", err)
	}

	if got := d.Get("security_ip_array_name"); got != "test" {
		t.Errorf("security_ip_array_name = %#v, want the custom array test", got)
	}
	ipList, _ := d.Get("security_ip_list").([]interface{})
	if len(ipList) != 1 || ipList[0] != "10.0.0.2" {
		t.Errorf("security_ip_list = %#v, want the custom array's IPs", d.Get("security_ip_list"))
	}
}

// When the target whitelist array is located and the response carries an
// explicitly empty IP list, Read clears the stale state; a missing list key
// keeps the configured value.
func TestUnitPolarDBApplicationReadWhitelistEmptyListConverges(t *testing.T) {
	read := func(t *testing.T, raw map[string]interface{}, body string) *schema.ResourceData {
		t.Helper()
		client := polardbAPITestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeApplicationAttribute":
				fmt.Fprint(w, body)
			case "ListTagResources":
				fmt.Fprint(w, `{"TagResources":{}}`)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		d := schema.TestResourceDataRaw(t, ResourceApplication().Schema, raw)
		d.SetId("app-whitelist-empty")
		if err := ResourceApplication().Read(d, client); err != nil {
			t.Fatalf("resource read returned error: %v", err)
		}
		return d
	}

	t.Run("an explicitly empty response clears stale IPs", func(t *testing.T) {
		d := read(t,
			map[string]interface{}{
				"application_type":       "polarclaw",
				"architecture":           "x86",
				"region_id":              "cn-beijing",
				"security_ip_array_name": "white_list",
				"security_ip_list":       []interface{}{"10.0.0.9"},
			},
			`{"Status":"Activated","VSwitchId":"vsw-test",
				"SecurityIPArrays":[{"SecurityIPArrayName":"white_list","SecurityIPType":"ipv4","SecurityIPList":""}]}`)
		ipList, _ := d.Get("security_ip_list").([]interface{})
		if len(ipList) != 0 {
			t.Fatalf("security_ip_list = %#v, want the stale entries cleared", ipList)
		}
	})

	t.Run("a missing list key keeps the configured IPs", func(t *testing.T) {
		d := read(t,
			map[string]interface{}{
				"application_type":       "polarclaw",
				"architecture":           "x86",
				"region_id":              "cn-beijing",
				"security_ip_array_name": "white_list",
				"security_ip_list":       []interface{}{"10.0.0.9"},
			},
			`{"Status":"Activated","VSwitchId":"vsw-test",
				"SecurityIPArrays":[{"SecurityIPArrayName":"white_list","SecurityIPType":"ipv4"}]}`)
		ipList, _ := d.Get("security_ip_list").([]interface{})
		if len(ipList) != 1 || ipList[0] != "10.0.0.9" {
			t.Fatalf("security_ip_list = %#v, want the configured entries kept", ipList)
		}
	})
}
