// Code in this file belongs to the property_completeness governance pass for
// alicloud_polardb_application: it maps the CreateApplication object/array
// parameters (Storages, MemApplicationSpec, AgenticDBBranchSpec, Endpoints,
// KnowledgeApplicationSpec, DnatEntries) — which the OpenAPI declares with
// style:"json" — onto JSON-string request fields, mirroring the existing
// Components encoding. Blocks whose members are order-insensitive collections
// (Storages, Endpoints, DnatEntries) are modeled as TypeSet, the remaining
// MaxItems=1 object blocks as TypeList; the marshal helpers below accept both
// element shapes. marshalBlockItem projects one element into an API-shaped
// map, omitting empty/zero values so unspecified optional sub-fields are not
// sent.
package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/go-cty/cty"
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

// marshalBlockItem projects a single block element (map[string]interface{},
// from either a TypeList or a TypeSet) into an API-shaped map, omitting
// empty/zero values so optional sub-fields the user did not set are not sent
// to CreateApplication.
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

// marshalBlockList serializes an order-insensitive block collection (TypeSet)
// or a list block (TypeList) of maps into a JSON array string for
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
		// A block whose sub-attributes are all optional and non-computed yields
		// a nil element when the user sets the block without any field; skip it
		// instead of asserting on nil.
		element, ok := item.(map[string]interface{})
		if !ok || element == nil {
			continue
		}
		out = append(out, marshalBlockItem(element, fields))
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
	// All sub-attributes are optional and non-computed, so an unset block can
	// surface a single nil element instead of a populated map.
	element, ok := items[0].(map[string]interface{})
	if !ok || element == nil {
		return "", true, nil
	}
	out := marshalBlockItem(element, fields)
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

// Readable DescribeApplicationAttribute projections. The query API returns only
// part of each create-only block; Read overlays exactly these fields and keeps
// every write-only member from the configuration.
var (
	storagesReadableFields = []blockField{
		{schemaKey: "storage_type", apiName: "StorageType", kind: blockString},
		{schemaKey: "storage_instance_id", apiName: "StorageInstanceId", kind: blockString},
		{schemaKey: "storage_capacity", apiName: "StorageCapacity", kind: blockString},
		{schemaKey: "storage_performance_level", apiName: "StoragePerformanceLevel", kind: blockString},
	}
	memApplicationSpecReadableFields = []blockField{
		{schemaKey: "db_name", apiName: "DbName", kind: blockString},
		{schemaKey: "db_user", apiName: "UserName", kind: blockString},
		{schemaKey: "embedder_model", apiName: "EmbedderModelName", kind: blockString},
		{schemaKey: "graph_llm_model", apiName: "GraphLlmModelName", kind: blockString},
		{schemaKey: "llm_model", apiName: "LlmModelName", kind: blockString},
		{schemaKey: "memory_type", apiName: "MemoryType", kind: blockString},
		{schemaKey: "project_name", apiName: "ProjectName", kind: blockString},
		{schemaKey: "reranker_model", apiName: "RerankerModelName", kind: blockString},
		{schemaKey: "visual_embedder_model", apiName: "VisualEmbedderModelName", kind: blockString},
		{schemaKey: "visual_embedder_model_dimension", apiName: "VisualEmbedderDimensions", kind: blockInt},
		{schemaKey: "visual_llm_modalities", apiName: "VisualLlmModalities", kind: blockStringList},
		{schemaKey: "visual_llm_model", apiName: "VisualLlmModelName", kind: blockString},
	}
)

// expandApplicationParameters projects the parameters block into the
// Parameters JSON list that CreateApplication and ModifyApplicationParameter
// both accept (style:"json"). The boolean is false when the block is unset.
func expandApplicationParameters(d *schema.ResourceData) ([]map[string]interface{}, bool) {
	v, ok := d.GetOk("parameters")
	if !ok {
		return nil, false
	}
	params, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	parameters := make([]map[string]interface{}, 0, len(params))
	for _, param := range params {
		// An all-optional, non-computed block can yield a nil element when the
		// user declares it without any sub-field.
		item, ok := param.(map[string]interface{})
		if !ok || item == nil {
			continue
		}
		paramMap := make(map[string]interface{})
		if v, ok := item["parameter_name"].(string); ok {
			paramMap["ParameterName"] = v
		}
		if v, ok := item["parameter_value"].(string); ok {
			paramMap["ParameterValue"] = v
		}
		parameters = append(parameters, paramMap)
	}
	return parameters, true
}

// marshalApplicationParameters renders the parameters block as the JSON string
// both write calls send. The boolean is false when the block is unset.
func marshalApplicationParameters(d *schema.ResourceData) (string, bool, error) {
	parameters, ok := expandApplicationParameters(d)
	if !ok {
		return "", false, nil
	}
	data, err := json.Marshal(parameters)
	if err != nil {
		return "", false, errs.WrapError(err)
	}
	return string(data), true, nil
}

// validateApplicationTypeSpecs enforces the CreateApplication requirements that
// depend on application_type: mem0 must carry MemApplicationSpec with its three
// model names, knowledge must carry KnowledgeApplicationSpec.LlmModel. The
// blocks stay optional for the other application types.
func validateApplicationTypeSpecs(d *schema.ResourceData) error {
	applicationType := strings.ToLower(strings.TrimSpace(d.Get("application_type").(string)))
	switch applicationType {
	case "mem0":
		missing := make([]string, 0, 4)
		spec := firstBlockItem(d.Get("mem_application_spec"))
		if spec == nil {
			missing = append(missing, "mem_application_spec (MemApplicationSpec)")
		} else {
			for _, required := range []struct{ key, api string }{
				{"embedder_model", "EmbedderModel"},
				{"llm_model", "LlmModel"},
				{"reranker_model", "RerankerModel"},
			} {
				if v, _ := spec[required.key].(string); strings.TrimSpace(v) == "" {
					missing = append(missing, "mem_application_spec."+required.key+" ("+required.api+")")
				}
			}
		}
		if len(missing) > 0 {
			return fmt.Errorf("application_type \"mem0\" requires %s to be set before calling CreateApplication", strings.Join(missing, ", "))
		}
	case "knowledge":
		spec := firstBlockItem(d.Get("knowledge_application_spec"))
		if spec == nil {
			return fmt.Errorf("application_type \"knowledge\" requires knowledge_application_spec (KnowledgeApplicationSpec) to be set before calling CreateApplication")
		}
		if v, _ := spec["llm_model"].(string); strings.TrimSpace(v) == "" {
			return fmt.Errorf("application_type \"knowledge\" requires knowledge_application_spec.llm_model (LlmModel) to be set before calling CreateApplication")
		}
	}
	return nil
}

// firstBlockItem returns the first populated member of a MaxItems=1 block.
func firstBlockItem(configured interface{}) map[string]interface{} {
	items, ok := configured.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	item, ok := items[0].(map[string]interface{})
	if !ok || item == nil {
		return nil
	}
	return item
}

// blockItems normalizes a TypeSet/TypeList block into its member maps.
func blockItems(configured interface{}) []map[string]interface{} {
	var raw []interface{}
	switch v := configured.(type) {
	case *schema.Set:
		raw = v.List()
	case []interface{}:
		raw = v
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok && m != nil {
			out = append(out, m)
		}
	}
	return out
}

// storageMemberAttributes are every attribute of a storages member. All of them
// take part in the Set hash: SDKv2's diffSet returns before any per-field diff
// when the hashes of the old and new sets match (helper/schema/schema.go), so a
// member attribute left out of the hash makes a configuration change to it
// invisible and lets two members that differ only in that attribute collapse
// into one. Including the attributes Read back-fills (type, instance id,
// capacity, performance level) is what keeps such changes visible; the
// back-fill-only difference that produces is dropped at plan time by
// resourceApplicationCustomizeDiff.
var storageMemberAttributes = []string{
	"storage_type",
	"storage_instance_id",
	"endpoint_id",
	"mount_path",
	"container_mount_path",
	"storage_capacity",
	"storage_endpoint",
	"storage_performance_level",
}

// storageMemberIdentityAttributes identify a storages member across the
// configuration and the state, in matching order: the write-only attributes
// (never rewritten by Read) plus the reported type and instance id.
var storageMemberIdentityAttributes = []string{
	"storage_instance_id",
	"storage_type",
	"endpoint_id",
	"mount_path",
	"container_mount_path",
	"storage_endpoint",
}

// storageMemberHash is the Set hash of the storages block.
func storageMemberHash(v interface{}) int {
	item, _ := v.(map[string]interface{})
	var identity strings.Builder
	for _, key := range storageMemberAttributes {
		value, _ := item[key].(string)
		identity.WriteString(key)
		identity.WriteByte('=')
		identity.WriteString(value)
		identity.WriteByte(';')
	}
	return schema.HashString(identity.String())
}

// resourceApplicationCustomizeDiff keeps the storages plan honest.
//
// storageMemberHash covers every member attribute, so a configuration change to
// any of them changes the member identity and the ForceNew flags plan the
// replacement. The same attributes are back-filled by Read: an application
// whose configuration leaves storage_type, storage_instance_id,
// storage_capacity or storage_performance_level unset receives the
// service-reported value in its state. Hashing that value would plan a
// replacement on every run, so the difference is dropped when it is
// explainable without a configuration change (same member count, every
// configured value equal to the state value after capacity normalization).
// Anything the configuration expresses differently — a changed value, an added
// or removed member, an unknown configuration value — keeps the SDK diff, so a
// real replacement stays visible.
func resourceApplicationCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Id() == "" {
		// Create: there is no back-filled state to reconcile against.
		return nil
	}
	oldValue, _ := d.GetChange("storages")
	oldItems := blockItems(oldValue)
	if len(oldItems) == 0 {
		// Import/first read: nothing has been back-filled yet.
		return nil
	}
	configuredItems, known := configuredStorageItems(d)
	if !known {
		// The configuration value cannot be resolved yet; keep the SDK diff
		// rather than risk hiding a change.
		return nil
	}
	if storagePlanIsBackfillOnly(configuredItems, oldItems) {
		return d.Clear("storages")
	}
	return nil
}

// configuredStorageItems returns the storages members the configuration
// declares, keeping only the attributes it explicitly sets. GetRawConfig is the
// source because an Optional+Computed attribute the block omits is null there
// while the planned state may carry the back-filled value. When the raw
// configuration is unavailable the planned set is used and every non-empty
// value counts as configured. The boolean is false when a value is unknown, in
// which case no conclusion can be drawn.
func configuredStorageItems(d *schema.ResourceDiff) ([]map[string]interface{}, bool) {
	rawConfig := d.GetRawConfig()
	if rawConfig == cty.NilVal || rawConfig.IsNull() {
		return storageConfiguredFromPlanned(d.Get("storages")), true
	}
	if !rawConfig.IsKnown() {
		return nil, false
	}
	rawStorages := rawConfig.GetAttr("storages")
	if rawStorages.IsNull() {
		// The block is absent. storages is Optional+Computed, so an absent
		// block keeps the state members: nothing the configuration expresses
		// changed.
		return nil, true
	}
	if !rawStorages.IsKnown() {
		return nil, false
	}
	items := make([]map[string]interface{}, 0, rawStorages.LengthInt())
	for iterator := rawStorages.ElementIterator(); iterator.Next(); {
		_, element := iterator.Element()
		if element.IsNull() {
			continue
		}
		if !element.IsKnown() {
			return nil, false
		}
		item := make(map[string]interface{}, len(storageMemberAttributes))
		for _, key := range storageMemberAttributes {
			value := element.GetAttr(key)
			if value.IsNull() {
				continue
			}
			if !value.IsKnown() {
				return nil, false
			}
			item[key] = value.AsString()
		}
		items = append(items, item)
	}
	return items, true
}

// storageConfiguredFromPlanned is the fallback used when the raw configuration
// is unavailable: the planned set carries the configured values, and an
// attribute the block omits is empty.
func storageConfiguredFromPlanned(planned interface{}) []map[string]interface{} {
	items := blockItems(planned)
	configured := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		entry := make(map[string]interface{}, len(storageMemberAttributes))
		for _, key := range storageMemberAttributes {
			if value, _ := item[key].(string); strings.TrimSpace(value) != "" {
				entry[key] = value
			}
		}
		configured = append(configured, entry)
	}
	return configured
}

// storagePlanIsBackfillOnly reports whether the storages difference between the
// configuration and the state is explainable without a configuration change:
// the same number of members, every member pairable, and every configured
// attribute equal to the paired member's state value.
func storagePlanIsBackfillOnly(configured, previous []map[string]interface{}) bool {
	if len(configured) != len(previous) {
		return false
	}
	paired := make([]bool, len(previous))
	for _, configuredItem := range configured {
		index := pairConfiguredStorageMember(configuredItem, previous, paired)
		if index < 0 {
			return false
		}
		paired[index] = true
		for key, raw := range configuredItem {
			configuredValue, _ := raw.(string)
			previousValue, _ := previous[index][key].(string)
			if normalizeStorageValue(key, configuredValue) != normalizeStorageValue(key, previousValue) {
				return false
			}
		}
	}
	return true
}

// pairConfiguredStorageMember returns the index of the state member a
// configured member corresponds to, or -1 when it cannot be paired. Identity
// attributes are preferred, then an equal capacity, then the only remaining
// member; an ambiguous candidate is not guessed.
func pairConfiguredStorageMember(configured map[string]interface{}, previous []map[string]interface{}, paired []bool) int {
	pick := func(match func(map[string]interface{}) bool) int {
		candidate := -1
		for i, item := range previous {
			if paired[i] || !match(item) {
				continue
			}
			if candidate >= 0 {
				return -1
			}
			candidate = i
		}
		return candidate
	}
	hasIdentity := false
	for _, key := range storageMemberIdentityAttributes {
		if _, ok := configured[key]; ok {
			hasIdentity = true
			break
		}
	}
	if hasIdentity {
		if index := pick(func(item map[string]interface{}) bool {
			for _, key := range storageMemberIdentityAttributes {
				raw, ok := configured[key]
				if !ok {
					continue
				}
				configuredValue, _ := raw.(string)
				previousValue, _ := item[key].(string)
				if normalizeStorageValue(key, configuredValue) != normalizeStorageValue(key, previousValue) {
					return false
				}
			}
			return true
		}); index >= 0 {
			return index
		}
	}
	if rawCapacity, ok := configured["storage_capacity"]; ok {
		configuredCapacity, _ := rawCapacity.(string)
		if configuredCapacity != "" {
			if index := pick(func(item map[string]interface{}) bool {
				previousCapacity, _ := item["storage_capacity"].(string)
				return previousCapacity != "" && equalStorageCapacity(configuredCapacity, previousCapacity)
			}); index >= 0 {
				return index
			}
		}
	}
	return pick(func(map[string]interface{}) bool { return true })
}

// normalizeStorageValue gives a storages attribute the form used for
// comparison: capacities are unit-insensitive ("100" equals "100Gi"), every
// other attribute is compared after trimming spaces.
func normalizeStorageValue(key, value string) string {
	if key == "storage_capacity" {
		return normalizedStorageCapacity(value)
	}
	return strings.TrimSpace(value)
}

// mergeApplicationStorages reconciles the configured storages with the list
// DescribeApplicationAttribute returns.
//
//   - With no configured member (import, or an application that relies on the
//     service default) the response is adopted with its readable fields.
//   - A matched member keeps its configured values and fills the readable
//     fields the configuration left empty, so a normal apply and an import
//     build the same state. A configured value is never rewritten: the service
//     reporting a different value for it would plan a replacement on every
//     apply.
//   - A configured member the response no longer reports is dropped only when
//     it can be identified (instance id, or a storage type); an ambiguous
//     member is kept instead of being paired by list order.
func mergeApplicationStorages(configured interface{}, remote []interface{}) []map[string]interface{} {
	configuredItems := blockItems(configured)
	remoteItems := make([]map[string]interface{}, 0, len(remote))
	for _, entry := range remote {
		if item, ok := entry.(map[string]interface{}); ok && item != nil {
			remoteItems = append(remoteItems, item)
		}
	}
	if len(configuredItems) == 0 {
		adopted := make([]map[string]interface{}, 0, len(remoteItems))
		for _, remoteItem := range remoteItems {
			item := map[string]interface{}{}
			setBlockItemFromResponse(item, remoteItem, storagesReadableFields)
			if len(item) > 0 {
				adopted = append(adopted, item)
			}
		}
		return adopted
	}
	matched := matchStorageMembers(configuredItems, remoteItems)
	merged := make([]map[string]interface{}, 0, len(configuredItems))
	for i, configuredItem := range configuredItems {
		if idx := matched[i]; idx >= 0 {
			merged = append(merged, mergeStorageReadableFields(configuredItem, remoteItems[idx]))
			continue
		}
		if storageMemberRemoved(configuredItem, remoteItems) {
			// The response is the complete list and no longer reports this
			// member: it was removed outside Terraform.
			continue
		}
		merged = append(merged, configuredItem)
	}
	return merged
}

// mergeStorageReadableFields fills the readable fields a configured member left
// empty with the values the response reports, so the state of a matched member
// carries the same information whether it came from an apply or from an import.
func mergeStorageReadableFields(configured, remote map[string]interface{}) map[string]interface{} {
	item := make(map[string]interface{}, len(configured)+len(storagesReadableFields))
	for key, value := range configured {
		item[key] = value
	}
	for _, field := range storagesReadableFields {
		if current, _ := item[field.schemaKey].(string); current != "" {
			continue
		}
		value, ok := remote[field.apiName].(string)
		if !ok {
			continue
		}
		if value = strings.TrimSpace(value); value != "" {
			item[field.schemaKey] = value
		}
	}
	return item
}

// storageMemberRemoved reports whether a configured member that could not be
// paired is definitely absent from the response. A configured StorageInstanceId
// identifies the member outright, so its absence means removal; without an id
// the member is only treated as removed when the response reports no member of
// its storage type at all. Anything less definite keeps the member instead of
// guessing by list order.
func storageMemberRemoved(item map[string]interface{}, remote []map[string]interface{}) bool {
	if id, _ := item["storage_instance_id"].(string); id != "" {
		return true
	}
	storageType, _ := item["storage_type"].(string)
	if storageType == "" {
		return false
	}
	for _, remoteItem := range remote {
		if remoteType, _ := remoteItem["StorageType"].(string); remoteType == storageType {
			return false
		}
	}
	return true
}

// matchStorageMembers pairs configured members with response members and
// returns the response index per configured member (-1 when unpaired). A shared
// StorageInstanceId identifies a member. Otherwise the pairing stays inside one
// StorageType and uses the capacity only when it selects exactly one candidate;
// an ambiguous group is left unpaired rather than guessed by list order.
func matchStorageMembers(configured, remote []map[string]interface{}) []int {
	matched := make([]int, len(configured))
	for i := range matched {
		matched[i] = -1
	}
	usedRemote := make([]bool, len(remote))
	for i, configuredItem := range configured {
		configuredID, _ := configuredItem["storage_instance_id"].(string)
		if configuredID == "" {
			continue
		}
		for j, remoteItem := range remote {
			if usedRemote[j] {
				continue
			}
			if remoteID, _ := remoteItem["StorageInstanceId"].(string); remoteID == configuredID {
				matched[i] = j
				usedRemote[j] = true
				break
			}
		}
	}
	for i, configuredItem := range configured {
		if matched[i] >= 0 {
			continue
		}
		configuredType, _ := configuredItem["storage_type"].(string)
		if configuredType == "" {
			continue
		}
		candidates := make([]int, 0, len(remote))
		for j, remoteItem := range remote {
			if usedRemote[j] {
				continue
			}
			if remoteType, _ := remoteItem["StorageType"].(string); remoteType == configuredType {
				candidates = append(candidates, j)
			}
		}
		if len(candidates) == 1 {
			matched[i] = candidates[0]
			usedRemote[candidates[0]] = true
			continue
		}
		if len(candidates) == 0 {
			continue
		}
		configuredCapacity, _ := configuredItem["storage_capacity"].(string)
		capacityMatches := make([]int, 0, 1)
		if configuredCapacity != "" {
			for _, j := range candidates {
				remoteCapacity, _ := remote[j]["StorageCapacity"].(string)
				if remoteCapacity != "" && equalStorageCapacity(configuredCapacity, remoteCapacity) {
					capacityMatches = append(capacityMatches, j)
				}
			}
		}
		if len(capacityMatches) == 1 {
			matched[i] = capacityMatches[0]
			usedRemote[capacityMatches[0]] = true
		}
	}
	return matched
}

// equalStorageCapacity compares two capacity values ignoring unit spelling:
// "20", "20Gi" and "20gib" describe the same size while "100" and "20Gi" do
// not. It is only used to attribute a response member, never to rewrite a
// configured value.
func equalStorageCapacity(left, right string) bool {
	return normalizedStorageCapacity(left) == normalizedStorageCapacity(right)
}

func normalizedStorageCapacity(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.TrimSuffix(normalized, "ib")
	normalized = strings.TrimSuffix(normalized, "b")
	normalized = strings.TrimSuffix(normalized, "i")
	normalized = strings.TrimSuffix(normalized, "g")
	return strings.TrimSpace(normalized)
}

// applicationWhitelistFromResponse attributes the application-level whitelist
// entry to restore: the array named like the configuration, else the "default"
// array, else the only array. It returns the entry name, the comma separated
// IP list split into values, and the IP type, so Read can restore the
// whitelist state after an import.
func applicationWhitelistFromResponse(arrays []interface{}, configuredName string) (name string, ipList []string, ipListPresent bool, ipType string, found bool) {
	entries := make([]map[string]interface{}, 0, len(arrays))
	for _, raw := range arrays {
		if item, ok := raw.(map[string]interface{}); ok && item != nil {
			entries = append(entries, item)
		}
	}
	if len(entries) == 0 {
		return "", nil, false, "", false
	}
	pick := -1
	if configuredName != "" {
		for i, item := range entries {
			if entryName, _ := item["SecurityIPArrayName"].(string); entryName == configuredName {
				pick = i
				break
			}
		}
		if pick < 0 {
			return "", nil, false, "", false
		}
	}
	if pick < 0 {
		// Without a configured name the application's own (non-default) arrays
		// are the ones a configuration would have created, so they win over the
		// implicit default entry when there is exactly one.
		custom := make([]int, 0, len(entries))
		for i, item := range entries {
			if entryName, _ := item["SecurityIPArrayName"].(string); entryName != "" && entryName != "default" {
				custom = append(custom, i)
			}
		}
		if len(custom) == 1 {
			pick = custom[0]
		}
	}
	if pick < 0 {
		for i, item := range entries {
			if entryName, _ := item["SecurityIPArrayName"].(string); entryName == "default" {
				pick = i
				break
			}
		}
	}
	if pick < 0 && len(entries) == 1 {
		pick = 0
	}
	if pick < 0 {
		return "", nil, false, "", false
	}
	entry := entries[pick]
	name, _ = entry["SecurityIPArrayName"].(string)
	ipType, _ = entry["SecurityIPType"].(string)
	if rawList, ok := entry["SecurityIPList"].(string); ok {
		ipListPresent = true
		for _, ip := range strings.Split(rawList, ",") {
			if trimmed := strings.TrimSpace(ip); trimmed != "" {
				ipList = append(ipList, trimmed)
			}
		}
	}
	return name, ipList, ipListPresent, strings.TrimSpace(ipType), true
}

// mergeApplicationEndpoints reconciles the configured endpoints with the
// response. The query API returns an endpoint description but no endpoint
// type, so a response member can only be attributed by an equal description; a
// single configured member paired with a single response member adopts the
// response description so a server-side change surfaces as drift instead of
// being hidden. Response-only members are not invented because their type is
// unknown. It returns nil when there is nothing identity-safe to write.
func mergeApplicationEndpoints(configured interface{}, remote []interface{}) []map[string]interface{} {
	configuredItems := blockItems(configured)
	remoteItems := make([]map[string]interface{}, 0, len(remote))
	descriptions := make(map[string]bool, len(remote))
	for _, raw := range remote {
		item, ok := raw.(map[string]interface{})
		if !ok || item == nil {
			continue
		}
		remoteItems = append(remoteItems, item)
		if description, _ := item["Description"].(string); strings.TrimSpace(description) != "" {
			descriptions[strings.TrimSpace(description)] = true
		}
	}
	if len(remoteItems) == 0 {
		return nil
	}
	if len(configuredItems) == 0 {
		// Import: the CreateApplication contract fixes the endpoint type to
		// Primary, so the reported descriptions can be rebuilt as members and
		// the description state is restored instead of ignored.
		adopted := make([]map[string]interface{}, 0, len(remoteItems))
		for _, remoteItem := range remoteItems {
			description, _ := remoteItem["Description"].(string)
			if description = strings.TrimSpace(description); description == "" {
				continue
			}
			adopted = append(adopted, map[string]interface{}{"endpoint_type": "Primary", "description": description})
		}
		if len(adopted) == 0 {
			return nil
		}
		return adopted
	}
	merged := make([]map[string]interface{}, 0, len(configuredItems))
	for _, item := range configuredItems {
		member := make(map[string]interface{}, len(item))
		for key, value := range item {
			member[key] = value
		}
		description, _ := member["description"].(string)
		description = strings.TrimSpace(description)
		if description != "" && descriptions[description] {
			member["description"] = description
		} else if len(configuredItems) == 1 && len(remoteItems) == 1 {
			if remoteDescription, _ := remoteItems[0]["Description"].(string); strings.TrimSpace(remoteDescription) != "" {
				member["description"] = strings.TrimSpace(remoteDescription)
			}
		}
		merged = append(merged, member)
	}
	return merged
}

// mergeApplicationMemSpec keeps the configured member's write-only fields
// (db_password, shard, embedder_model_dimension) and overlays the readable
// MemApplicationAttribute values. It reports false when there is nothing to
// write, so an absent attribute neither creates an empty block nor clears a
// configured one.
func mergeApplicationMemSpec(configured interface{}, response map[string]interface{}) (map[string]interface{}, bool) {
	item := map[string]interface{}{}
	if configuredItem := firstBlockItem(configured); configuredItem != nil {
		for k, v := range configuredItem {
			item[k] = v
		}
	}
	setBlockItemFromResponse(item, response, memApplicationSpecReadableFields)
	if len(item) == 0 {
		return nil, false
	}
	return item, true
}

// setBlockItemFromResponse copies the readable fields present in an API
// response object onto a schema-shaped block member. String values are
// trimmed — the API pads model/project names with whitespace — and an empty
// response value never overwrites a populated member.
func setBlockItemFromResponse(item map[string]interface{}, response map[string]interface{}, fields []blockField) {
	for _, f := range fields {
		switch f.kind {
		case blockString:
			value, ok := response[f.apiName].(string)
			if !ok {
				continue
			}
			if value = strings.TrimSpace(value); value != "" {
				item[f.schemaKey] = value
			}
		case blockInt:
			switch value := response[f.apiName].(type) {
			case json.Number:
				if number, err := value.Int64(); err == nil && number != 0 {
					item[f.schemaKey] = int(number)
				}
			case int:
				if value != 0 {
					item[f.schemaKey] = value
				}
			case float64:
				if value != 0 {
					item[f.schemaKey] = int(value)
				}
			}
		case blockStringList:
			values, ok := response[f.apiName].([]interface{})
			if !ok {
				continue
			}
			list := make([]string, 0, len(values))
			for _, element := range values {
				if text, ok := element.(string); ok && strings.TrimSpace(text) != "" {
					list = append(list, strings.TrimSpace(text))
				}
			}
			if len(list) > 0 {
				item[f.schemaKey] = list
			}
		}
	}
}
