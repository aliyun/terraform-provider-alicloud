package ots

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// searchIndexResourceSchema composes the search index resource schema.
func searchIndexResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateOTSInstanceName,
		},
		"table_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateOTSTableName,
		},
		"index_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateOTSIndexName,
		},
		"time_to_live": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  -1,
			ForceNew: true,
			// 86400s = 1d
			ValidateFunc: validation.Any(validation.IntInSlice([]int{-1}), validation.IntAtLeast(86400)),
		},
		"schema": searchIndexSchemaAttributeSchema(),
		"index_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"create_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"sync_phase": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"current_sync_timestamp": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

// searchIndexSchemaAttributeSchema composes the `schema` container from the
// independently declared schema modules.
func searchIndexSchemaAttributeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		Set:      searchIndexSchemaHash,
		Elem:     searchIndexSchemaResource(),
	}
}

func searchIndexSchemaResource() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		// Some parameters such as IndexOptions and AnalyzerParameter in field_schema are not supported for the time being,
		// because there is no description of these parameters in the official documentation.
		"field_schema":  searchIndexFieldSchemaList(),
		"index_setting": searchIndexSettingSchema(),
		"index_sort":    searchIndexSortSchema(),
	}}
}

// searchIndexSchemaHash makes only the top-level field order insignificant.
// The field list itself retains its configuration order for RawConfig correlation.
func searchIndexSchemaHash(value interface{}) int {
	block := value.(map[string]interface{})
	normalized := make(map[string]interface{}, len(block))
	for name, attribute := range block {
		normalized[name] = attribute
	}
	if fields, ok := block["field_schema"].([]interface{}); ok {
		ordered := append([]interface{}(nil), fields...)
		sort.SliceStable(ordered, func(i, j int) bool {
			return ordered[i].(map[string]interface{})["field_name"].(string) < ordered[j].(map[string]interface{})["field_name"].(string)
		})
		normalized["field_schema"] = ordered
	}
	return schema.HashResource(searchIndexSchemaResource())(normalized)
}

// searchIndexFieldSchemaList returns the `field_schema` list of a top-level
// field.
func searchIndexFieldSchemaList() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem:     searchIndexFieldSchema(),
	}
}

// searchIndexFieldSchema composes a top-level field from the attributes every
// field accepts plus the modules only a top-level field accepts.
func searchIndexFieldSchema() *schema.Resource {
	return &schema.Resource{Schema: composeSchemaMaps(
		searchIndexCommonFieldSchema(),
		map[string]*schema.Schema{
			"field_type":        searchIndexFieldTypeSchema(),
			"vector_options":    searchIndexVectorOptionsSchema(),
			"sub_field_schemas": searchIndexSubFieldSchemaSet(),
		},
	)}
}

// searchIndexSubFieldSchema composes a Nested sub-field from the attributes
// every field accepts. A sub-field cannot nest another level, and the service
// rejects Vector as a sub-field.
func searchIndexSubFieldSchema() *schema.Resource {
	return &schema.Resource{Schema: composeSchemaMaps(
		searchIndexCommonFieldSchema(),
		map[string]*schema.Schema{
			"field_type": searchIndexSubFieldTypeSchema(),
		},
	)}
}

// composeSchemaMaps returns a new map holding every attribute of the given
// maps so schema blocks can be assembled from independent modules.
func composeSchemaMaps(maps ...map[string]*schema.Schema) map[string]*schema.Schema {
	composed := make(map[string]*schema.Schema)
	for _, attributes := range maps {
		for name, attribute := range attributes {
			composed[name] = attribute
		}
	}
	return composed
}

// searchIndexCommonFieldSchema returns the attributes shared by a top-level
// field and a Nested sub-field.
func searchIndexCommonFieldSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"is_array": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"index": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},
		"analyzer": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(OtsSearchSingleWord),
				string(OtsSearchSplit),
				string(OtsSearchMinWord),
				string(OtsSearchMaxWord),
				string(OtsSearchFuzzy)},
				false),
		},
		"enable_sort_and_agg": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"store": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"date_format_list": {
			// A set: the service treats the configured formats as a group and the
			// order carries no meaning, so reordering it must not create a diff.
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"is_virtual_field": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"source_field_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
	}
}

// searchIndexFieldTypeSchema returns the field_type attribute of a top-level
// field, which may also be Nested or Vector.
func searchIndexFieldTypeSchema() *schema.Schema {
	fieldTypes := searchIndexLeafFieldTypes()
	fieldTypes = append(fieldTypes, string(OtsSearchTypeNested), string(SearchIndexFieldTypeVector))
	return searchIndexFieldTypeAttributeSchema(fieldTypes)
}

// searchIndexSubFieldTypeSchema returns the field_type attribute of a Nested
// sub-field. Only a top-level field may be Nested, and the service rejects
// Vector as a sub-field, so both stay top-level only.
func searchIndexSubFieldTypeSchema() *schema.Schema {
	return searchIndexFieldTypeAttributeSchema(searchIndexLeafFieldTypes())
}

func searchIndexFieldTypeAttributeSchema(fieldTypes []string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(fieldTypes, false),
	}
}

// searchIndexLeafFieldTypes returns the field types a field can declare at
// any level.
func searchIndexLeafFieldTypes() []string {
	return []string{
		string(OtsSearchTypeLong),
		string(OtsSearchTypeDouble),
		string(OtsSearchTypeBoolean),
		string(OtsSearchTypeKeyword),
		string(OtsSearchTypeText),
		string(OtsSearchTypeDate),
		string(OtsSearchTypeGeoPoint),
		string(SearchIndexFieldTypeFuzzyKeyword),
	}
}

// searchIndexVectorOptionsSchema returns the `vector_options` block of a
// Vector field. Its members stay Optional because the repository's
// BreakingChange gate rejects newly required attributes; the parser rejects an
// incomplete block before the create request instead.
func searchIndexVectorOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"data_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(SearchIndexVectorDataTypeFloat32)}, false),
				},
				"dimension": {
					Type:         schema.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 4096),
				},
				"metric_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(SearchIndexVectorMetricEuclidean),
						string(SearchIndexVectorMetricCosine),
						string(SearchIndexVectorMetricDotProduct)}, false),
				},
			},
		},
	}
}

// searchIndexSubFieldSchemaSet returns the `sub_field_schemas` set of a
// Nested field. A set: the sub-field order carries no meaning.
func searchIndexSubFieldSchemaSet() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem:     searchIndexSubFieldSchema(),
	}
}

func searchIndexSettingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"routing_fields": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func searchIndexSortSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sorter": {
					Type:     schema.TypeList,
					Required: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sorter_type": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(OtsSearchPrimaryKeySort),
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(OtsSearchPrimaryKeySort), string(OtsSearchFieldSort)}, false),
							},
							"order": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(OtsSearchSortOrderAsc),
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(OtsSearchSortOrderAsc), string(OtsSearchSortOrderDesc)}, false),
							},
							"field_name": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
							},
							"mode": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(OtsSearchModeMin), string(OtsSearchModeMax), string(OtsSearchModeAvg)}, false),
							},
						},
					},
				},
			},
		},
	}
}
func ResourceOtsSearchIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsSearchIndexCreate,
		Read:   resourceAliyunOtsSearchIndexRead,
		Delete: resourceAliyunOtsSearchIndexDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: searchIndexResourceSchema(),
	}
}

func parseSearchIndexResourceArgs(d *schema.ResourceData) (*SearchIndexResourceArgs, error) {
	indexSchema, err := parseSearchIndexSchema(d)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	args := &SearchIndexResourceArgs{
		instanceName: d.Get("instance_name").(string),
		tableName:    d.Get("table_name").(string),
		indexName:    d.Get("index_name").(string),
		schema:       indexSchema,
	}

	if v, ok := d.GetOk("time_to_live"); ok {
		args.ttl = int32(v.(int))
	}
	return args, nil
}

func parseSearchIndexSchema(d *schema.ResourceData) (*tablestore.IndexSchema, error) {
	// A search index has exactly one schema. Rejecting extra blocks keeps the
	// materialized block aligned with the raw configuration, which is required to
	// tell an omitted boolean from an explicit false.
	schemaBlocks := d.Get("schema").(*schema.Set)
	if schemaBlocks.Len() != 1 {
		return nil, fmt.Errorf("only one schema block is supported, got %d", schemaBlocks.Len())
	}
	schemaArg := schemaBlocks.List()[0].(map[string]interface{})
	rawSchemaBlock, rawSchemaBlockPresent := searchIndexRawSchemaBlock(d)

	// required
	fieldSchemas, err := parseFieldSchemas(schemaArg, rawSchemaBlock, rawSchemaBlockPresent)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	// optional
	indexSetting := parseIndexSetting(schemaArg)
	// optional
	indexSort, err := parseIndexSort(schemaArg)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	return &tablestore.IndexSchema{
		FieldSchemas: fieldSchemas,
		IndexSetting: indexSetting,
		IndexSort:    indexSort,
	}, nil
}

func parseFieldSchemas(schemaArg map[string]interface{}, rawSchemaBlock cty.Value, rawSchemaBlockPresent bool) ([]*tablestore.FieldSchema, error) {
	var fieldSchemas []*tablestore.FieldSchema

	fieldSchemasArg := schemaArg["field_schema"].([]interface{})
	for index, fs := range fieldSchemasArg {
		fsMap := fs.(map[string]interface{})
		rawField, rawFieldPresent := searchIndexRawListElement(rawSchemaBlock, rawSchemaBlockPresent, "field_schema", index)
		fieldSchema, err := parseFieldSchema(fsMap, rawField, rawFieldPresent, false)
		if err != nil {
			return nil, errs.WrapError(err)
		}

		fieldSchemas = append(fieldSchemas, fieldSchema)
	}
	return fieldSchemas, nil
}

// searchIndexRawSchemaBlock returns the raw configuration of the single `schema`
// block so optional boolean arguments can be told apart from their zero values.
func searchIndexRawSchemaBlock(d *schema.ResourceData) (cty.Value, bool) {
	return searchIndexRawFirstElement(d.GetRawConfig(), "schema")
}

// searchIndexRawFirstElement returns the first element of a collection attribute
// of a raw object.
func searchIndexRawFirstElement(parent cty.Value, attribute string) (cty.Value, bool) {
	if parent.IsNull() || !parent.IsKnown() || !parent.Type().IsObjectType() || !parent.Type().HasAttribute(attribute) {
		return cty.NilVal, false
	}
	collection := parent.GetAttr(attribute)
	if collection.IsNull() || !collection.IsKnown() || !collection.CanIterateElements() {
		return cty.NilVal, false
	}
	iterator := collection.ElementIterator()
	if !iterator.Next() {
		return cty.NilVal, false
	}
	_, element := iterator.Element()
	if iterator.Next() {
		// The parser consumes one schema block; with more than one block the raw
		// configuration cannot be matched to it reliably.
		return cty.NilVal, false
	}
	if element.IsNull() || !element.IsKnown() || !element.Type().IsObjectType() {
		return cty.NilVal, false
	}
	return element, true
}

// searchIndexRawListElement returns the raw configuration of the index-th element
// of a list attribute inside a raw object.
func searchIndexRawListElement(parent cty.Value, present bool, attribute string, index int) (cty.Value, bool) {
	if !present || parent.IsNull() || !parent.IsKnown() || !parent.Type().IsObjectType() || !parent.Type().HasAttribute(attribute) {
		return cty.NilVal, false
	}
	collection := parent.GetAttr(attribute)
	if collection.IsNull() || !collection.IsKnown() || !collection.CanIterateElements() {
		return cty.NilVal, false
	}
	elements := collection.AsValueSlice()
	if index < 0 || index >= len(elements) {
		return cty.NilVal, false
	}
	element := elements[index]
	if element.IsNull() || !element.IsKnown() || !element.Type().IsObjectType() {
		return cty.NilVal, false
	}
	return element, true
}

type searchIndexRawSubFieldSet struct {
	byName      map[string]cty.Value
	unknownName []cty.Value
}

// searchIndexRawSubFieldElements separates raw sub-fields by whether their names can correlate set elements.
func searchIndexRawSubFieldElements(rawField cty.Value, present bool) searchIndexRawSubFieldSet {
	elements := searchIndexRawSubFieldSet{byName: make(map[string]cty.Value)}
	if !present || rawField.IsNull() || !rawField.IsKnown() || !rawField.Type().IsObjectType() || !rawField.Type().HasAttribute("sub_field_schemas") {
		return elements
	}
	collection := rawField.GetAttr("sub_field_schemas")
	if collection.IsNull() || !collection.IsKnown() || !collection.CanIterateElements() {
		return elements
	}
	for iterator := collection.ElementIterator(); iterator.Next(); {
		_, element := iterator.Element()
		if element.IsNull() || !element.IsKnown() || !element.Type().IsObjectType() || !element.Type().HasAttribute("field_name") {
			continue
		}
		name := element.GetAttr("field_name")
		if name.Type() != cty.String {
			continue
		}
		if name.IsNull() {
			continue
		}
		if !name.IsKnown() {
			elements.unknownName = append(elements.unknownName, element)
			continue
		}
		elements.byName[name.AsString()] = element
	}
	return elements
}

// searchIndexUncorrelatedStore returns a value only when all unmatched elements prove identical store semantics.
func searchIndexUncorrelatedStore(rawFields []cty.Value, materialized []interface{}, byName map[string]cty.Value) (bool, bool, error) {
	allOmitted := true
	allExplicitFalse := true
	for _, rawField := range rawFields {
		if !rawField.Type().IsObjectType() || !rawField.Type().HasAttribute("store") {
			return false, false, fmt.Errorf("cannot correlate raw sub-field store configuration")
		}
		store := rawField.GetAttr("store")
		if store.Type() != cty.Bool {
			return false, false, fmt.Errorf("cannot correlate raw sub-field store configuration")
		}
		if store.IsNull() {
			allExplicitFalse = false
			continue
		}
		allOmitted = false
		if store.IsKnown() && store.True() {
			allExplicitFalse = false
		}
	}
	if allOmitted {
		return false, false, nil
	}
	allMaterializedFalse := true
	for _, subField := range materialized {
		subFieldMap := subField.(map[string]interface{})
		if _, ok := byName[subFieldMap["field_name"].(string)]; ok {
			continue
		}
		if subFieldMap["store"].(bool) {
			allMaterializedFalse = false
			break
		}
	}
	if allExplicitFalse && allMaterializedFalse {
		return true, false, nil
	}
	return false, false, fmt.Errorf("cannot reliably correlate sub_field_schemas with unknown field_name values; use known field names to preserve explicit store values")
}

// searchIndexStoreConfigured reports whether the `store` argument is present in
// the raw configuration of a field and, when it is, its value. An unknown raw
// value has already been resolved in the materialized configuration at apply time.
func searchIndexStoreConfigured(fsMap map[string]interface{}, rawField cty.Value, present bool) (bool, bool) {
	if !present || rawField.IsNull() || !rawField.IsKnown() || !rawField.Type().IsObjectType() || !rawField.Type().HasAttribute("store") {
		return false, false
	}
	value := rawField.GetAttr("store")
	if value.IsNull() || value.Type() != cty.Bool {
		return false, false
	}
	if !value.IsKnown() {
		store, ok := fsMap["store"].(bool)
		return ok, store
	}
	return true, value.True()
}

// searchIndexKeepsLegacyStore reports whether the legacy behaviour of always
// sending an explicit `store = false` for an omitted argument still applies.
// Only the top-level field types that existed before the Nested sub-field support
// keep it; sub-fields, virtual fields and the newly supported Vector and
// FuzzyKeyword types follow the SDK default when the argument is omitted.
func searchIndexKeepsLegacyStore(isSubField bool, fieldType tablestore.FieldType, isVirtualField bool) bool {
	if isSubField || isVirtualField {
		return false
	}
	switch fieldType {
	case tablestore.FieldType_VECTOR, tablestore.FieldType_FUZZY_KEYWORD:
		return false
	default:
		return true
	}
}

func parseIndexSort(schemaArg map[string]interface{}) (*search.Sort, error) {
	var indexSort *search.Sort
	if v, ok := schemaArg["index_sort"]; ok {
		var sorts []search.Sorter
		valSet := v.(*schema.Set)
		if valSet == nil || valSet.Len() == 0 {
			return indexSort, nil
		}
		indexSortArg := valSet.List()[0].(map[string]interface{})
		if v, ok := indexSortArg["sorter"]; ok {
			sortersArg := v.([]interface{})
			for _, s := range sortersArg {
				sorterArg := s.(map[string]interface{})
				sort, err := parseIndexFieldSort(sorterArg)
				if err != nil {
					return nil, errs.WrapError(err)
				}

				sorts = append(sorts, sort)
			}

		}

		//default sort
		if len(sorts) < 1 {
			asc := search.SortOrder_ASC
			sorts = append(sorts, &search.PrimaryKeySort{
				Order: &asc,
			})
		}
		indexSort = &search.Sort{
			Sorters: sorts,
		}
	}
	return indexSort, nil
}

func parseIndexSetting(schemaArg map[string]interface{}) *tablestore.IndexSetting {
	if v, ok := schemaArg["index_setting"]; ok {
		valSet := v.(*schema.Set)
		if valSet == nil || valSet.Len() == 0 {
			return nil
		}
		indexSettingArg := valSet.List()[0].(map[string]interface{})
		if v, ok := indexSettingArg["routing_fields"]; ok {
			routersArg := v.([]interface{})
			var routerPKs []string
			for _, router := range routersArg {
				routerPK := router.(string)
				routerPKs = append(routerPKs, routerPK)
			}

			if len(routerPKs) > 0 {
				return &tablestore.IndexSetting{
					RoutingFields: routerPKs,
				}
			}
		}

	}
	return nil
}

func parseFieldSchema(fsMap map[string]interface{}, rawField cty.Value, rawFieldPresent bool, isSubField bool) (*tablestore.FieldSchema, error) {
	fieldName := fsMap["field_name"].(string)
	// fieldTypeName keeps the spelling the configuration uses (Text,
	// FuzzyKeyword) so the validation errors match the documented values.
	fieldTypeName := fsMap["field_type"].(string)
	fieldType, err := ConvertSearchIndexFieldTypeString(SearchIndexFieldTypeString(fieldTypeName))
	if err != nil {
		return nil, errs.WrapError(err)
	}
	isNested := fieldType == tablestore.FieldType_NESTED
	isVirtualField := fsMap["is_virtual_field"].(bool)
	sourceFieldName := fsMap["source_field_name"].(string)

	fieldSchema := &tablestore.FieldSchema{
		FieldName: &fieldName,
		FieldType: fieldType,
	}

	if v, ok := fsMap["sub_field_schemas"]; ok {
		subFieldSchemas := v.(*schema.Set)
		if !isNested && subFieldSchemas.Len() > 0 {
			return nil, fmt.Errorf("field %q of type %s does not support sub_field_schemas", fieldName, fieldTypeName)
		}
		rawSubFields := searchIndexRawSubFieldElements(rawField, rawFieldPresent)
		materializedSubFields := subFieldSchemas.List()
		uncorrelatedCount := 0
		for _, subField := range materializedSubFields {
			subFieldMap := subField.(map[string]interface{})
			if _, ok := rawSubFields.byName[subFieldMap["field_name"].(string)]; !ok {
				uncorrelatedCount++
			}
		}
		for _, subField := range materializedSubFields {
			subFieldMap := subField.(map[string]interface{})
			rawSubField, rawSubFieldPresent := rawSubFields.byName[subFieldMap["field_name"].(string)]
			if !rawSubFieldPresent && uncorrelatedCount == 1 && len(rawSubFields.unknownName) == 1 {
				rawSubField, rawSubFieldPresent = rawSubFields.unknownName[0], true
			}
			if !rawSubFieldPresent && len(rawSubFields.unknownName) > 0 && !subFieldMap["store"].(bool) {
				configured, value, err := searchIndexUncorrelatedStore(rawSubFields.unknownName, materializedSubFields, rawSubFields.byName)
				if err != nil {
					return nil, fmt.Errorf("field %q: %w", fieldName, err)
				}
				if configured {
					rawSubField = cty.ObjectVal(map[string]cty.Value{"store": cty.BoolVal(value)})
					rawSubFieldPresent = true
				}
			}
			if !rawSubFieldPresent {
				rawSubField = cty.NilVal
			}
			parsedSubField, err := parseFieldSchema(subFieldMap, rawSubField, rawSubFieldPresent, true)
			if err != nil {
				return nil, err
			}
			fieldSchema.FieldSchemas = append(fieldSchema.FieldSchemas, parsedSubField)
		}
	}
	if isNested && len(fieldSchema.FieldSchemas) == 0 {
		return nil, fmt.Errorf("field %q of type Nested requires sub_field_schemas", fieldName)
	}

	isArray := fsMap["is_array"].(bool)
	if isNested {
		if isArray {
			return nil, fmt.Errorf("field %q of type Nested does not support is_array", fieldName)
		}
	} else {
		if isArray && (fieldType == tablestore.FieldType_VECTOR || fieldType == tablestore.FieldType_FUZZY_KEYWORD) {
			return nil, fmt.Errorf("field %q of type %s does not support is_array", fieldName, fieldTypeName)
		}
		// Non-Nested fields keep the wire shape of the released provider, which
		// always sent the value including the implicit false.
		fieldSchema.IsArray = &isArray
	}

	index := fsMap["index"].(bool)
	if isVirtualField && !index {
		return nil, fmt.Errorf("virtual field %q requires index = true", fieldName)
	}
	fieldSchema.Index = &index

	if analyzer := fsMap["analyzer"].(string); analyzer != "" {
		if isNested {
			return nil, fmt.Errorf("field %q of type Nested does not support analyzer", fieldName)
		}
		if fieldType != tablestore.FieldType_TEXT && (isSubField || fieldType == tablestore.FieldType_VECTOR || fieldType == tablestore.FieldType_FUZZY_KEYWORD) {
			return nil, fmt.Errorf("field %q of type %s does not support analyzer", fieldName, fieldTypeName)
		}
		analyzerType, err := ConvertSearchIndexAnalyzerTypeString(SearchIndexAnalyzerTypeString(analyzer))
		if err != nil {
			return nil, errs.WrapError(err)
		}
		fieldSchema.Analyzer = &analyzerType
	}

	enableSortAndAgg := fsMap["enable_sort_and_agg"].(bool)
	if enableSortAndAgg {
		switch fieldType {
		case tablestore.FieldType_NESTED:
			return nil, fmt.Errorf("field %q of type Nested does not support enable_sort_and_agg", fieldName)
		case tablestore.FieldType_TEXT, tablestore.FieldType_FUZZY_KEYWORD:
			return nil, fmt.Errorf("field %q of type %s does not support enable_sort_and_agg", fieldName, fieldTypeName)
		}
	}
	if !isNested {
		// Non-Nested fields keep the wire shape of the released provider, which
		// always sent the value including the implicit false.
		fieldSchema.EnableSortAndAgg = &enableSortAndAgg
	}

	storeConfigured, storeValue := searchIndexStoreConfigured(fsMap, rawField, rawFieldPresent)
	if !storeConfigured && fsMap["store"].(bool) {
		// `store` has no schema default and is never computed, so a materialized
		// true can only come from an explicit configuration, even when the raw
		// configuration cannot be correlated with the field.
		storeConfigured, storeValue = true, true
	}
	if isNested {
		if storeConfigured && storeValue {
			return nil, fmt.Errorf("field %q of type Nested does not support store", fieldName)
		}
	} else if storeConfigured {
		fieldSchema.Store = &storeValue
	} else if searchIndexKeepsLegacyStore(isSubField, fieldType, isVirtualField) {
		// The top-level field types that existed before the Nested sub-field
		// support always sent an explicit `store = false` when the argument was
		// omitted. Keep that wire shape so the fix does not change the behaviour
		// of existing configurations.
		store := false
		fieldSchema.Store = &store
	}

	if v, ok := fsMap["date_format_list"]; ok {
		dateFormats := v.(*schema.Set)
		if fieldType != tablestore.FieldType_DATE && dateFormats.Len() > 0 {
			return nil, fmt.Errorf("field %q of type %s does not support date_format_list", fieldName, fieldTypeName)
		}
		for _, dateFormat := range dateFormats.List() {
			fieldSchema.DateFormats = append(fieldSchema.DateFormats, dateFormat.(string))
		}
	}
	if isSubField && fieldType == tablestore.FieldType_DATE && len(fieldSchema.DateFormats) == 0 {
		return nil, fmt.Errorf("field %q of type Date requires date_format_list", fieldName)
	}

	if isVirtualField {
		if isNested {
			return nil, fmt.Errorf("field %q of type Nested cannot be a virtual field", fieldName)
		}
		if sourceFieldName == "" {
			return nil, fmt.Errorf("virtual field %q requires source_field_name", fieldName)
		}
		fieldSchema.IsVirtualField = &isVirtualField
		fieldSchema.SourceFieldNames = []string{sourceFieldName}
	} else if sourceFieldName != "" {
		return nil, fmt.Errorf("field %q sets source_field_name but is_virtual_field is not true", fieldName)
	}

	if v, ok := fsMap["vector_options"]; ok {
		vectorOptionsArg := v.([]interface{})
		if len(vectorOptionsArg) > 0 && vectorOptionsArg[0] != nil {
			if fieldType != tablestore.FieldType_VECTOR {
				return nil, fmt.Errorf("field %q of type %s does not support vector_options", fieldName, fieldTypeName)
			}
			options := vectorOptionsArg[0].(map[string]interface{})
			vectorOptions := &tablestore.VectorOptions{}
			if dataTypeStr := options["data_type"].(string); dataTypeStr != "" {
				dataType, err := ConvertSearchIndexVectorDataTypeString(SearchIndexVectorDataTypeString(dataTypeStr))
				if err != nil {
					return nil, errs.WrapError(err)
				}
				vectorOptions.VectorDataType = &dataType
			} else {
				return nil, fmt.Errorf("vector_options requires data_type")
			}
			if dimension := options["dimension"].(int); dimension > 0 {
				dimensionValue := int32(dimension)
				vectorOptions.Dimension = &dimensionValue
			} else {
				return nil, fmt.Errorf("vector_options requires dimension")
			}
			if metricTypeStr := options["metric_type"].(string); metricTypeStr != "" {
				metricType, err := ConvertSearchIndexVectorMetricTypeString(SearchIndexVectorMetricTypeString(metricTypeStr))
				if err != nil {
					return nil, errs.WrapError(err)
				}
				vectorOptions.VectorMetricType = &metricType
			} else {
				return nil, fmt.Errorf("vector_options requires metric_type")
			}
			fieldSchema.VectorOptions = vectorOptions
		} else if fieldType == tablestore.FieldType_VECTOR {
			return nil, fmt.Errorf("field %q of type Vector requires vector_options", fieldName)
		}
	}

	return fieldSchema, nil
}

func parseIndexFieldSort(sorterArg map[string]interface{}) (search.Sorter, error) {
	sortFieldType, err := ConvertSearchIndexSortFieldTypeString(SearchIndexSortFieldTypeString(sorterArg["sorter_type"].(string)))
	if err != nil {
		return nil, errs.WrapError(err)
	}

	orderType, err := ConvertSearchIndexOrderTypeString(SearchIndexOrderTypeString(sorterArg["order"].(string)))
	if err != nil {
		return nil, errs.WrapError(err)
	}

	switch sort := sortFieldType.(type) {
	case *search.PrimaryKeySort:
		sort.Order = &orderType
		return sort, nil
	case *search.FieldSort:
		sort.Order = &orderType
		// field_name and mode are required when sortFieldType is FieldSort
		sort.FieldName = sorterArg["field_name"].(string)
		mode, err := ConvertSearchIndexSortModeString(SearchIndexSortModeString(sorterArg["mode"].(string)))
		if err != nil {
			return nil, errs.WrapError(err)
		}
		sort.Mode = &mode

		return sort, nil
	default:
		return nil, errs.WrapError(fmt.Errorf("not find search index sort field type [PrimaryKeySort|FieldSort]: %v", sortFieldType))
	}
}

type SearchIndexResourceArgs struct {
	instanceName string
	tableName    string
	indexName    string
	ttl          int32
	schema       *tablestore.IndexSchema
}

func resourceAliyunOtsSearchIndexCreate(d *schema.ResourceData, meta interface{}) error {
	args, err := parseSearchIndexResourceArgs(d)
	if err != nil {
		return errs.WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	// check table exists
	tableResp, err := otsService.LoopWaitTable(args.instanceName, args.tableName)
	if err != nil {
		return errs.WrapError(err)
	}
	// serverside arguments check
	if err := args.checkArgs(tableResp); err != nil {
		return err
	}
	// build request
	req := &tablestore.CreateSearchIndexRequest{
		TableName:   args.tableName,
		IndexName:   args.indexName,
		IndexSchema: args.schema,
		TimeToLive:  &args.ttl,
	}

	var reqClient *tablestore.TableStoreClient
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := client.WithTableStoreClient(args.instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			reqClient = tableStoreClient
			return tableStoreClient.CreateSearchIndex(req)
		})
		defer func() {
			logs.AddDebug("CreateTableSearchIndex", raw, reqClient, req)
		}()

		if err != nil {
			if errs.IsExpectedErrors(err, errs.OtsSearchIndexIsTemporarilyUnavailable) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ots_search_index", "CreateSearchIndex", errs.AliyunTablestoreGoSdk)
	}

	d.SetId(ID(args.instanceName, args.tableName, args.indexName, SearchIndexTypeHolder))
	return resourceAliyunOtsSearchIndexRead(d, meta)
}

func (args *SearchIndexResourceArgs) checkArgs(tableResp *tablestore.DescribeTableResponse) error {

	if args.indexName == args.tableName {
		return errs.WrapError(fmt.Errorf("index name cannot be the same as table: %s/%s", args.indexName, args.tableName))
	}

	if tableResp.TableOption.TimeToAlive != -1 {
		return errs.WrapError(fmt.Errorf("when creating a search index, the TimeToAlive of the table must be -1: %v", tableResp.TableOption.TimeToAlive))
	}
	if tableResp.TableOption.MaxVersion != 1 {
		return errs.WrapError(fmt.Errorf("when creating a search index, the table's MaxVersion must be 1: %v", tableResp.TableOption.MaxVersion))
	}

	if args.schema.IndexSetting != nil {
		for _, router := range args.schema.IndexSetting.RoutingFields {
			var routerInPk bool
			for _, pk := range tableResp.TableMeta.SchemaEntry {
				if router == *pk.Name {
					routerInPk = true
					break
				}
			}
			if !routerInPk {
				return errs.WrapError(fmt.Errorf("search index router field must be in primary key: %s/%s", args.indexName, router))
			}
		}
	}

	if args.schema.IndexSort != nil {
		if len(args.schema.IndexSort.Sorters) > 0 {
			for _, fs := range args.schema.FieldSchemas {
				if fs.FieldType == tablestore.FieldType_NESTED {
					return errs.WrapError(fmt.Errorf("search index with nested field type does not support index_sort: %s/%s", args.indexName, *fs.FieldName))
				}
			}
		}
	}

	return nil
}

func flattenSearchIndexSchema(indexSchema *tablestore.IndexSchema) ([]interface{}, error) {
	if indexSchema == nil {
		return nil, fmt.Errorf("search index schema is empty")
	}
	fields, err := flattenSearchIndexFields(indexSchema.FieldSchemas, false)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{"field_schema": fields}
	if indexSchema.IndexSetting != nil && len(indexSchema.IndexSetting.RoutingFields) > 0 {
		result["index_setting"] = []interface{}{map[string]interface{}{
			"routing_fields": indexSchema.IndexSetting.RoutingFields,
		}}
	}
	if indexSchema.IndexSort != nil && len(indexSchema.IndexSort.Sorters) > 0 {
		sorters, err := flattenSearchIndexSorters(indexSchema.IndexSort)
		if err != nil {
			return nil, err
		}
		result["index_sort"] = []interface{}{map[string]interface{}{"sorter": sorters}}
	}
	return []interface{}{result}, nil
}

func flattenSearchIndexFields(fields []*tablestore.FieldSchema, subField bool) ([]interface{}, error) {
	result := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		if field == nil || field.FieldName == nil {
			return nil, fmt.Errorf("search index contains a field without a name")
		}
		fieldType, err := flattenSearchIndexFieldType(field.FieldType)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", *field.FieldName, err)
		}
		if subField && (field.FieldType == tablestore.FieldType_NESTED || field.FieldType == tablestore.FieldType_VECTOR) {
			return nil, fmt.Errorf("field %q has an unsupported nested field type %s", *field.FieldName, fieldType)
		}
		if field.IndexOptions != nil || field.AnalyzerParameter != nil {
			return nil, fmt.Errorf("field %q uses index options or analyzer parameters not supported by the provider", *field.FieldName)
		}
		if len(field.SourceFieldNames) > 1 {
			return nil, fmt.Errorf("field %q has more than one source field", *field.FieldName)
		}
		item := map[string]interface{}{
			"field_name": *field.FieldName,
			"field_type": fieldType,
		}
		if field.IsArray != nil {
			item["is_array"] = *field.IsArray
		}
		if field.Index != nil {
			item["index"] = *field.Index
		} else if field.FieldType == tablestore.FieldType_NESTED {
			// Nested fields are always indexed, and both the request and the
			// response omit the value, so keep the schema default explicit.
			item["index"] = true
		}
		if field.EnableSortAndAgg != nil {
			item["enable_sort_and_agg"] = *field.EnableSortAndAgg
		}
		if field.Store != nil {
			item["store"] = *field.Store
		}
		if field.IsVirtualField != nil {
			item["is_virtual_field"] = *field.IsVirtualField
		}
		if len(field.SourceFieldNames) == 1 {
			item["source_field_name"] = field.SourceFieldNames[0]
		}
		if len(field.DateFormats) > 0 {
			item["date_format_list"] = field.DateFormats
		}
		if field.Analyzer != nil {
			analyzer, err := flattenSearchIndexAnalyzer(*field.Analyzer)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", *field.FieldName, err)
			}
			item["analyzer"] = analyzer
		}
		if len(field.FieldSchemas) > 0 {
			subFields, err := flattenSearchIndexFields(field.FieldSchemas, true)
			if err != nil {
				return nil, err
			}
			item["sub_field_schemas"] = subFields
		}
		if field.VectorOptions != nil {
			if subField {
				return nil, fmt.Errorf("field %q has unsupported vector options", *field.FieldName)
			}
			options, err := flattenSearchIndexVectorOptions(field.VectorOptions)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", *field.FieldName, err)
			}
			item["vector_options"] = []interface{}{options}
		}
		result = append(result, item)
	}
	return result, nil
}

func flattenSearchIndexFieldType(fieldType tablestore.FieldType) (string, error) {
	types := map[tablestore.FieldType]string{
		tablestore.FieldType_LONG:          string(OtsSearchTypeLong),
		tablestore.FieldType_DOUBLE:        string(OtsSearchTypeDouble),
		tablestore.FieldType_BOOLEAN:       string(OtsSearchTypeBoolean),
		tablestore.FieldType_KEYWORD:       string(OtsSearchTypeKeyword),
		tablestore.FieldType_TEXT:          string(OtsSearchTypeText),
		tablestore.FieldType_DATE:          string(OtsSearchTypeDate),
		tablestore.FieldType_GEO_POINT:     string(OtsSearchTypeGeoPoint),
		tablestore.FieldType_NESTED:        string(OtsSearchTypeNested),
		tablestore.FieldType_VECTOR:        string(SearchIndexFieldTypeVector),
		tablestore.FieldType_FUZZY_KEYWORD: string(SearchIndexFieldTypeFuzzyKeyword),
	}
	if result, ok := types[fieldType]; ok {
		return result, nil
	}
	return "", fmt.Errorf("unsupported field type %v", fieldType)
}

func flattenSearchIndexAnalyzer(analyzer tablestore.Analyzer) (string, error) {
	analyzers := map[tablestore.Analyzer]string{
		tablestore.Analyzer_SingleWord: string(OtsSearchSingleWord),
		tablestore.Analyzer_Split:      string(OtsSearchSplit),
		tablestore.Analyzer_MinWord:    string(OtsSearchMinWord),
		tablestore.Analyzer_MaxWord:    string(OtsSearchMaxWord),
		tablestore.Analyzer_Fuzzy:      string(OtsSearchFuzzy),
	}
	if result, ok := analyzers[analyzer]; ok {
		return result, nil
	}
	return "", fmt.Errorf("unsupported analyzer %q", analyzer)
}

func flattenSearchIndexVectorOptions(options *tablestore.VectorOptions) (map[string]interface{}, error) {
	if options.VectorDataType == nil || options.Dimension == nil || options.VectorMetricType == nil {
		return nil, fmt.Errorf("vector options are incomplete")
	}
	if *options.VectorDataType != tablestore.VectorDataType_FLOAT_32 {
		return nil, fmt.Errorf("unsupported vector data type %q", *options.VectorDataType)
	}
	metrics := map[tablestore.VectorMetricType]string{
		tablestore.VectorMetricType_EUCLIDEAN:   string(SearchIndexVectorMetricEuclidean),
		tablestore.VectorMetricType_COSINE:      string(SearchIndexVectorMetricCosine),
		tablestore.VectorMetricType_DOT_PRODUCT: string(SearchIndexVectorMetricDotProduct),
	}
	metric, ok := metrics[*options.VectorMetricType]
	if !ok {
		return nil, fmt.Errorf("unsupported vector metric type %q", *options.VectorMetricType)
	}
	return map[string]interface{}{
		"data_type":   string(SearchIndexVectorDataTypeFloat32),
		"dimension":   int(*options.Dimension),
		"metric_type": metric,
	}, nil
}

func flattenSearchIndexSorters(sort *search.Sort) ([]interface{}, error) {
	if sort.DisableDefaultPkSorter != nil && *sort.DisableDefaultPkSorter {
		return nil, fmt.Errorf("search index disables the default primary-key sorter, which the provider does not support")
	}
	result := make([]interface{}, 0, len(sort.Sorters))
	for _, sorter := range sort.Sorters {
		item := map[string]interface{}{"order": string(OtsSearchSortOrderAsc)}
		switch value := sorter.(type) {
		case *search.PrimaryKeySort:
			item["sorter_type"] = string(OtsSearchPrimaryKeySort)
			if value.Order != nil {
				order, err := flattenSearchIndexSortOrder(*value.Order)
				if err != nil {
					return nil, err
				}
				item["order"] = order
			}
		case *search.FieldSort:
			if value.NestedFilter != nil || value.MissingValue != nil || value.MissingField != nil || len(value.MissingFields) > 0 {
				return nil, fmt.Errorf("field sorter %q uses options not supported by the provider", value.FieldName)
			}
			item["sorter_type"] = string(OtsSearchFieldSort)
			item["field_name"] = value.FieldName
			if value.Order != nil {
				order, err := flattenSearchIndexSortOrder(*value.Order)
				if err != nil {
					return nil, err
				}
				item["order"] = order
			}
			if value.Mode != nil {
				modes := map[search.SortMode]string{
					search.SortMode_Min: string(OtsSearchModeMin),
					search.SortMode_Max: string(OtsSearchModeMax),
					search.SortMode_Avg: string(OtsSearchModeAvg),
				}
				mode, ok := modes[*value.Mode]
				if !ok {
					return nil, fmt.Errorf("field sorter %q has unsupported mode %v", value.FieldName, *value.Mode)
				}
				item["mode"] = mode
			}
		default:
			return nil, fmt.Errorf("unsupported search index sorter %T", sorter)
		}
		result = append(result, item)
	}
	return result, nil
}

func flattenSearchIndexSortOrder(order search.SortOrder) (string, error) {
	switch order {
	case search.SortOrder_ASC:
		return string(OtsSearchSortOrderAsc), nil
	case search.SortOrder_DESC:
		return string(OtsSearchSortOrderDesc), nil
	default:
		return "", fmt.Errorf("unsupported sort order %v", order)
	}
}

func resourceAliyunOtsSearchIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	idx, err := otsService.DescribeOtsSearchIndex(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	if idx == nil {
		d.SetId("")
		return nil
	}
	instanceName, tableName, indexName, _, err := ParseIndexId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	// The response is not a faithful round trip of the configuration: it omits
	// booleans such as `store`, materializes the implicit primary-key sorter and
	// drops the FieldSort `mode`. Writing it back on every refresh would plan a
	// replacement for an unchanged configuration, so the schema is only filled
	// from the response when the state has none yet, which is the import case.
	if current := d.Get("schema").(*schema.Set); current == nil || current.Len() == 0 {
		flattenedSchema, err := flattenSearchIndexSchema(idx.Schema)
		if err != nil {
			return errs.WrapError(err)
		}
		if err := d.Set("schema", flattenedSchema); err != nil {
			return errs.WrapError(err)
		}
	}

	if err := d.Set("instance_name", instanceName); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("table_name", tableName); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("index_name", indexName); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("time_to_live", idx.TimeToLive); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("index_id", d.Id()); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("create_time", idx.CreateTime); err != nil {
		return errs.WrapError(err)
	}

	if idx.SyncStat == nil {
		return errs.WrapError(fmt.Errorf("search index sync status is empty"))
	}
	phase, err := ConvertSearchIndexSyncPhase(idx.SyncStat.SyncPhase)
	if err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("sync_phase", string(phase)); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("current_sync_timestamp", idx.SyncStat.CurrentSyncTimestamp); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func resourceAliyunOtsSearchIndexDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, indexName, _, err := ParseIndexId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)

	if err := otsService.DeleteSearchIndex(instanceName, tableName, indexName); err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteSearchIndex", errs.AliyunTablestoreGoSdk)
	}

	return errs.WrapError(otsService.WaitForSearchIndex(instanceName, tableName, indexName, Deleted, DefaultTimeout))
}
