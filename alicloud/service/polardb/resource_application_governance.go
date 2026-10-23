// Code in this file belongs to the property_completeness governance pass for
// alicloud_polardb_application: it maps the CreateApplication object/array
// parameters (Storages, MemApplicationSpec, AgenticDBBranchSpec, Endpoints,
// KnowledgeApplicationSpec, DnatEntries) — which the OpenAPI declares with
// style:"json" — onto JSON-string request fields, mirroring the existing
// Components encoding. Each block is a TypeList of maps; marshalBlockItem
// projects one element into an API-shaped map, omitting empty/zero values so
// unspecified optional sub-fields are not sent.
package polardb

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// blockFieldKind enumerates the scalar shapes marshalBlockItem projects.
type blockFieldKind string

const (
	blockString     blockFieldKind = "string"
	blockInt        blockFieldKind = "int"
	blockBool       blockFieldKind = "bool"
	blockStringList blockFieldKind = "stringlist"
)

// blockField ties a schema attribute name to its CreateApplication API name
// and scalar kind.
type blockField struct {
	schemaKey string
	apiName   string
	kind      blockFieldKind
}

// marshalBlockItem projects a single TypeList element (map[string]interface{})
// into an API-shaped map, omitting empty/zero values so optional sub-fields
// the user did not set are not sent to CreateApplication.
func marshalBlockItem(item map[string]interface{}, fields []blockField) map[string]interface{} {
	out := map[string]interface{}{}
	for _, f := range fields {
		switch f.kind {
		case blockString:
			if v, ok := item[f.schemaKey].(string); ok && v != "" {
				out[f.apiName] = v
			}
		case blockInt:
			if v, ok := item[f.schemaKey].(int); ok && v != 0 {
				out[f.apiName] = v
			}
		case blockBool:
			if v, ok := item[f.schemaKey].(bool); ok && v {
				out[f.apiName] = v
			}
		case blockStringList:
			// A nested multi-member field may be TypeList ([]interface{}) or
			// TypeSet (*schema.Set) depending on the block's schema; accept both.
			var raw []interface{}
			switch r := item[f.schemaKey].(type) {
			case []interface{}:
				raw = r
			case *schema.Set:
				raw = r.List()
			}
			if len(raw) > 0 {
				list := make([]string, 0, len(raw))
				for _, e := range raw {
					if s, ok := e.(string); ok && s != "" {
						list = append(list, s)
					}
				}
				if len(list) > 0 {
					out[f.apiName] = list
				}
			}
		}
	}
	return out
}

// marshalBlockList serializes a TypeList of maps into a JSON array string for
// CreateApplication array params (Storages/Endpoints/DnatEntries). The boolean
// is false (and string empty) when the block is unset.
func marshalBlockList(d *schema.ResourceData, key string, fields []blockField) (string, bool, error) {
	v, ok := d.GetOk(key)
	if !ok {
		return "", false, nil
	}
	// TypeList yields []interface{}; TypeSet yields *schema.Set (its List()
	// is the same element slice). Accept both so multi-member order-insensitive
	// blocks can be modeled as TypeSet without a separate marshal path.
	var items []interface{}
	switch raw := v.(type) {
	case []interface{}:
		items = raw
	case *schema.Set:
		items = raw.List()
	default:
		return "", false, nil
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		out = append(out, marshalBlockItem(item.(map[string]interface{}), fields))
	}
	data, err := json.Marshal(out)
	if err != nil {
		return "", false, errs.WrapError(err)
	}
	return string(data), true, nil
}

// marshalBlockObject serializes the first element of a MaxItems=1 TypeList
// into a JSON object string for CreateApplication object params
// (MemApplicationSpec/AgenticDBBranchSpec/KnowledgeApplicationSpec). The
// boolean is false when the block is unset.
func marshalBlockObject(d *schema.ResourceData, key string, fields []blockField) (string, bool, error) {
	v, ok := d.GetOk(key)
	if !ok {
		return "", false, nil
	}
	items := v.([]interface{})
	if len(items) == 0 {
		return "", true, nil
	}
	out := marshalBlockItem(items[0].(map[string]interface{}), fields)
	data, err := json.Marshal(out)
	if err != nil {
		return "", false, errs.WrapError(err)
	}
	return string(data), true, nil
}

// CreateApplication block field tables. API names match the OpenAPI
// CreateApplication parameter sub-properties (style:"json"); schema keys
// follow snake_case provider convention.
var (
	storagesFields = []blockField{
		{schemaKey: "storage_type", apiName: "StorageType", kind: blockString},
		{schemaKey: "storage_instance_id", apiName: "StorageInstanceId", kind: blockString},
		{schemaKey: "endpoint_id", apiName: "EndpointId", kind: blockString},
		{schemaKey: "mount_path", apiName: "MountPath", kind: blockString},
		{schemaKey: "container_mount_path", apiName: "ContainerMountPath", kind: blockString},
		{schemaKey: "storage_capacity", apiName: "StorageCapacity", kind: blockString},
		{schemaKey: "storage_endpoint", apiName: "StorageEndpoint", kind: blockString},
		{schemaKey: "storage_performance_level", apiName: "StoragePerformanceLevel", kind: blockString},
	}
	memApplicationSpecFields = []blockField{
		{schemaKey: "db_name", apiName: "DbName", kind: blockString},
		{schemaKey: "db_password", apiName: "DbPassword", kind: blockString},
		{schemaKey: "db_user", apiName: "DbUser", kind: blockString},
		{schemaKey: "embedder_model", apiName: "EmbedderModel", kind: blockString},
		{schemaKey: "embedder_model_dimension", apiName: "EmbedderModelDimension", kind: blockInt},
		{schemaKey: "graph_llm_model", apiName: "GraphLlmModel", kind: blockString},
		{schemaKey: "llm_model", apiName: "LlmModel", kind: blockString},
		{schemaKey: "memory_type", apiName: "MemoryType", kind: blockString},
		{schemaKey: "project_name", apiName: "ProjectName", kind: blockString},
		{schemaKey: "reranker_model", apiName: "RerankerModel", kind: blockString},
		{schemaKey: "shard", apiName: "Shard", kind: blockInt},
		{schemaKey: "visual_embedder_model", apiName: "VisualEmbedderModel", kind: blockString},
		{schemaKey: "visual_embedder_model_dimension", apiName: "VisualEmbedderModelDimension", kind: blockInt},
		{schemaKey: "visual_llm_modalities", apiName: "VisualLlmModalities", kind: blockStringList},
		{schemaKey: "visual_llm_model", apiName: "VisualLlmModel", kind: blockString},
	}
	agenticDBBranchSpecFields = []blockField{
		{schemaKey: "branch_id", apiName: "BranchId", kind: blockString},
		{schemaKey: "db_cluster_id", apiName: "DBClusterId", kind: blockString},
		{schemaKey: "fork_from_application_id", apiName: "ForkFromApplicationId", kind: blockString},
		{schemaKey: "fork_from_branch", apiName: "ForkFromBranch", kind: blockBool},
		{schemaKey: "project_id", apiName: "ProjectId", kind: blockString},
		{schemaKey: "tenant_id", apiName: "TenantId", kind: blockString},
	}
	endpointsFields = []blockField{
		{schemaKey: "endpoint_type", apiName: "EndpointType", kind: blockString},
		{schemaKey: "description", apiName: "Description", kind: blockString},
	}
	knowledgeApplicationSpecFields = []blockField{
		{schemaKey: "dashboard_password", apiName: "DashboardPassword", kind: blockString},
		{schemaKey: "db_password", apiName: "DbPassword", kind: blockString},
		{schemaKey: "llm_model", apiName: "LlmModel", kind: blockString},
	}
	dnatEntriesFields = []blockField{
		{schemaKey: "port_name", apiName: "PortName", kind: blockString},
		{schemaKey: "front_port", apiName: "FrontPort", kind: blockInt},
	}
)
