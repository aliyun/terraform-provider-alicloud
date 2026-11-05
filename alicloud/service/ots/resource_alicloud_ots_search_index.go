package ots

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceOtsSearchIndex() *schema.Resource {
	return &schema.Resource{
		Create:        resourceAliyunOtsSearchIndexCreate,
		Read:          resourceAliyunOtsSearchIndexRead,
		Delete:        resourceAliyunOtsSearchIndexDelete,
		CustomizeDiff: validateSearchIndexSchemaCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
			"schema": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Some parameters such as IndexOptions and AnalyzerParameter in field_schema are not supported for the time being,
						// because there is no description of these parameters in the official documentation.
						"field_schema": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     otsSearchIndexFieldSchema(maxSearchIndexFieldNestingDepth),
						},

						"index_setting": {
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
						},

						"index_sort": {
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
						},
					},
				},
			},

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
		},
	}
}

// maxSearchIndexFieldNestingDepth limits how many sub_field_schemas levels the
// static Terraform schema exposes. Terraform requires a finite schema tree, so
// arbitrarily deep Nested fields cannot be represented recursively.
const maxSearchIndexFieldNestingDepth = 2

func otsSearchIndexFieldSchema(remainingNestingDepth int) *schema.Resource {
	fieldSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"field_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(OtsSearchTypeLong),
					string(OtsSearchTypeDouble),
					string(OtsSearchTypeBoolean),
					string(OtsSearchTypeKeyword),
					string(OtsSearchTypeText),
					string(OtsSearchTypeDate),
					string(OtsSearchTypeGeoPoint),
					string(OtsSearchTypeNested),
					string(OtsSearchTypeVector),
					string(OtsSearchTypeFuzzyKeyword)},
					false),
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
			"vector_options": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(OtsSearchVectorDataTypeFloat32)}, false),
						},
						"dimension": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 4096),
						},
						"metric_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(OtsSearchVectorMetricEuclidean),
								string(OtsSearchVectorMetricCosine),
								string(OtsSearchVectorMetricDotProduct)}, false),
						},
					},
				},
			},
		},
	}

	if remainingNestingDepth > 0 {
		fieldSchema.Schema["sub_field_schemas"] = &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem:     otsSearchIndexFieldSchema(remainingNestingDepth - 1),
		}
	}

	return fieldSchema
}

func validateSearchIndexSchemaCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// TypeSet diff readers do not reliably mark unknown nested values as
	// computed, so gate on the raw configuration instead of NewValueKnown.
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsWhollyKnown() {
		return nil
	}
	raw, ok := d.GetOk("schema")
	if !ok {
		return nil
	}
	schemaSet, ok := raw.(*schema.Set)
	if !ok || schemaSet.Len() == 0 {
		return nil
	}
	schemaArg, ok := schemaSet.List()[0].(map[string]interface{})
	if !ok {
		return nil
	}
	_, err := parseFieldSchemas(schemaArg)
	return err
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
	// only one schema is valid for a search index
	schemaArg := d.Get("schema").(*schema.Set).List()[0].(map[string]interface{})

	// required
	fieldSchemas, err := parseFieldSchemas(schemaArg)
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

func parseFieldSchemas(schemaArg map[string]interface{}) ([]*tablestore.FieldSchema, error) {
	fieldSchemasArg, ok := schemaArg["field_schema"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("field_schema is required")
	}
	return parseFieldSchemaList(fieldSchemasArg)
}

func parseFieldSchemaList(fieldSchemasArg []interface{}) ([]*tablestore.FieldSchema, error) {
	var fieldSchemas []*tablestore.FieldSchema

	for _, fs := range fieldSchemasArg {
		fsMap, ok := fs.(map[string]interface{})
		if !ok {
			continue
		}
		fieldSchema, err := parseFieldSchema(fsMap)
		if err != nil {
			return nil, errs.WrapError(err)
		}

		fieldSchemas = append(fieldSchemas, fieldSchema)
	}
	return fieldSchemas, nil
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

func parseFieldSchema(fsMap map[string]interface{}) (*tablestore.FieldSchema, error) {
	fieldName := fsMap["field_name"].(string)

	// required
	fieldSchema := &tablestore.FieldSchema{
		FieldName: &fieldName,
	}

	// optionals
	if v, ok := fsMap["field_type"]; ok {
		vv, err := ConvertSearchIndexFieldTypeString(SearchIndexFieldTypeString(v.(string)))
		if err != nil {
			return nil, errs.WrapError(err)
		}
		fieldSchema.FieldType = vv
	}
	if v, ok := fsMap["is_array"]; ok {
		isArray := v.(bool)
		fieldSchema.IsArray = &isArray
	}
	if v, ok := fsMap["index"]; ok {
		index := v.(bool)
		fieldSchema.Index = &index
	}
	if v, ok := fsMap["analyzer"]; ok && fsMap["analyzer"] != "" {
		analyzer, err := ConvertSearchIndexAnalyzerTypeString(SearchIndexAnalyzerTypeString(v.(string)))
		if err != nil {
			return nil, errs.WrapError(err)
		}
		fieldSchema.Analyzer = &analyzer
	}
	if v, ok := fsMap["enable_sort_and_agg"]; ok {
		enableSortAndAgg := v.(bool)
		fieldSchema.EnableSortAndAgg = &enableSortAndAgg
	}
	if v, ok := fsMap["store"]; ok {
		store := v.(bool)
		fieldSchema.Store = &store
	}
	if v, ok := fsMap["date_format_list"]; ok {
		if dateFormatSet, ok := v.(*schema.Set); ok && dateFormatSet.Len() > 0 {
			dateFormats := make([]string, 0, dateFormatSet.Len())
			for _, dateFormat := range dateFormatSet.List() {
				dateFormats = append(dateFormats, dateFormat.(string))
			}
			fieldSchema.DateFormats = dateFormats
		}
	}
	if v, ok := fsMap["is_virtual_field"]; ok {
		isVirtualField := v.(bool)
		fieldSchema.IsVirtualField = &isVirtualField
	}
	if v, ok := fsMap["source_field_name"]; ok {
		if sourceFieldName, ok := v.(string); ok && sourceFieldName != "" {
			fieldSchema.SourceFieldNames = []string{sourceFieldName}
		}
	}
	if v, ok := fsMap["vector_options"]; ok {
		if vectorOptionsArg, ok := v.([]interface{}); ok && len(vectorOptionsArg) > 0 && vectorOptionsArg[0] != nil {
			vectorOptions, err := parseVectorOptions(vectorOptionsArg[0].(map[string]interface{}))
			if err != nil {
				return nil, errs.WrapError(err)
			}
			fieldSchema.VectorOptions = vectorOptions
		}
	}
	if v, ok := fsMap["sub_field_schemas"]; ok {
		if subFieldSchemaSet, ok := v.(*schema.Set); ok && subFieldSchemaSet.Len() > 0 {
			subFieldSchemas, err := parseFieldSchemaList(subFieldSchemaSet.List())
			if err != nil {
				return nil, errs.WrapError(err)
			}
			fieldSchema.FieldSchemas = subFieldSchemas
		}
	}

	if err := validateSearchIndexFieldSchema(fieldSchema); err != nil {
		return nil, errs.WrapError(err)
	}
	return fieldSchema, nil
}

func parseVectorOptions(vectorOptionsArg map[string]interface{}) (*tablestore.VectorOptions, error) {
	dataType, err := ConvertSearchIndexVectorDataTypeString(SearchIndexVectorDataTypeString(vectorOptionsArg["data_type"].(string)))
	if err != nil {
		return nil, errs.WrapError(err)
	}
	metricType, err := ConvertSearchIndexVectorMetricTypeString(SearchIndexVectorMetricTypeString(vectorOptionsArg["metric_type"].(string)))
	if err != nil {
		return nil, errs.WrapError(err)
	}
	dimension := int32(vectorOptionsArg["dimension"].(int))

	return &tablestore.VectorOptions{
		VectorDataType:   &dataType,
		VectorMetricType: &metricType,
		Dimension:        &dimension,
	}, nil
}

func validateSearchIndexFieldSchema(fieldSchema *tablestore.FieldSchema) error {
	fieldName := ""
	if fieldSchema.FieldName != nil {
		fieldName = *fieldSchema.FieldName
	}

	switch fieldSchema.FieldType {
	case tablestore.FieldType_NESTED:
		if len(fieldSchema.FieldSchemas) == 0 {
			return fmt.Errorf("field %q of type Nested requires sub_field_schemas", fieldName)
		}
		if fieldSchema.EnableSortAndAgg != nil && *fieldSchema.EnableSortAndAgg {
			return fmt.Errorf("field %q of type Nested does not support enable_sort_and_agg", fieldName)
		}
	case tablestore.FieldType_TEXT, tablestore.FieldType_FUZZY_KEYWORD:
		if fieldSchema.EnableSortAndAgg != nil && *fieldSchema.EnableSortAndAgg {
			return fmt.Errorf("field %q of type %s does not support enable_sort_and_agg", fieldName, fieldSchema.FieldType.String())
		}
	case tablestore.FieldType_VECTOR:
		if fieldSchema.VectorOptions == nil {
			return fmt.Errorf("field %q of type Vector requires vector_options", fieldName)
		}
	case tablestore.FieldType_DATE:
		if len(fieldSchema.DateFormats) == 0 {
			return fmt.Errorf("field %q of type Date requires date_format_list", fieldName)
		}
	}

	isVirtualField := fieldSchema.IsVirtualField != nil && *fieldSchema.IsVirtualField
	if isVirtualField && len(fieldSchema.SourceFieldNames) == 0 {
		return fmt.Errorf("virtual field %q requires source_field_name", fieldName)
	}
	if !isVirtualField && len(fieldSchema.SourceFieldNames) > 0 {
		return fmt.Errorf("field %q sets source_field_name but is_virtual_field is not true", fieldName)
	}
	return nil
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
			if fieldName, ok := findNestedSearchIndexField(args.schema.FieldSchemas); ok {
				return errs.WrapError(fmt.Errorf("search index with nested field type does not support index_sort: %s/%s", args.indexName, fieldName))
			}
		}
	}

	return nil
}

func findNestedSearchIndexField(fieldSchemas []*tablestore.FieldSchema) (string, bool) {
	for _, fs := range fieldSchemas {
		if fs.FieldType == tablestore.FieldType_NESTED {
			if fs.FieldName != nil {
				return *fs.FieldName, true
			}
			return "", true
		}
		if fieldName, ok := findNestedSearchIndexField(fs.FieldSchemas); ok {
			return fieldName, ok
		}
	}
	return "", false
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

	if err := d.Set("index_id", d.Id()); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("create_time", idx.CreateTime); err != nil {
		return errs.WrapError(err)
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
