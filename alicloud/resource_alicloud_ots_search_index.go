package alicloud

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOtsSearchIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsSearchIndexCreate,
		Read:   resourceAliyunOtsSearchIndexRead,
		Delete: resourceAliyunOtsSearchIndexDelete,
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
							Elem: &schema.Resource{
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
											string(OtsSearchTypeNested)},
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
								},
							},
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

func parseSearchIndexResourceArgs(d *schema.ResourceData) (*SearchIndexResourceArgs, error) {
	indexSchema, err := parseSearchIndexSchema(d)
	if err != nil {
		return nil, WrapError(err)
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
		return nil, WrapError(err)
	}
	// optional
	indexSetting := parseIndexSetting(schemaArg)
	// optional
	indexSort, err := parseIndexSort(schemaArg)
	if err != nil {
		return nil, WrapError(err)
	}

	return &tablestore.IndexSchema{
		FieldSchemas: fieldSchemas,
		IndexSetting: indexSetting,
		IndexSort:    indexSort,
	}, nil
}

func parseFieldSchemas(schemaArg map[string]interface{}) ([]*tablestore.FieldSchema, error) {
	var fieldSchemas []*tablestore.FieldSchema

	fieldSchemasArg := schemaArg["field_schema"].([]interface{})
	for _, fs := range fieldSchemasArg {
		fsMap := fs.(map[string]interface{})
		fieldSchema, err := parseFieldSchema(fsMap)
		if err != nil {
			return nil, WrapError(err)
		}

		fieldSchemas = append(fieldSchemas, fieldSchema)
	}
	return fieldSchemas, nil
}

func parseIndexSort(schemaArg map[string]interface{}) (*search.Sort, error) {
	var indexSort *search.Sort
	if v, ok := schemaArg["index_sort"]; ok {
		var sorts []search.Sorter

		indexSortArg := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := indexSortArg["sorter"]; ok {
			sortersArg := v.([]interface{})
			for _, s := range sortersArg {
				sorterArg := s.(map[string]interface{})
				sort, err := parseIndexFieldSort(sorterArg)
				if err != nil {
					return nil, WrapError(err)
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
		indexSettingArg := v.(*schema.Set).List()[0].(map[string]interface{})
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
			return nil, WrapError(err)
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
			return nil, WrapError(err)
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
	return fieldSchema, nil
}

func parseIndexFieldSort(sorterArg map[string]interface{}) (search.Sorter, error) {
	sortFieldType, err := ConvertSearchIndexSortFieldTypeString(SearchIndexSortFieldTypeString(sorterArg["sorter_type"].(string)))
	if err != nil {
		return nil, WrapError(err)
	}

	orderType, err := ConvertSearchIndexOrderTypeString(SearchIndexOrderTypeString(sorterArg["order"].(string)))
	if err != nil {
		return nil, WrapError(err)
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
			return nil, WrapError(err)
		}
		sort.Mode = &mode

		return sort, nil
	default:
		return nil, WrapError(errors.New(fmt.Sprintf("not find search index sort field type [PrimaryKeySort|FieldSort]: %v", sortFieldType)))
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
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	// check table exists
	tableResp, err := otsService.LoopWaitTable(args.instanceName, args.tableName)
	if err != nil {
		return WrapError(err)
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
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(args.instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			reqClient = tableStoreClient
			return tableStoreClient.CreateSearchIndex(req)
		})
		defer func() {
			addDebug("CreateTableSearchIndex", raw, reqClient, req)
		}()

		if err != nil {
			if IsExpectedErrors(err, OtsSearchIndexIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_search_index", "CreateSearchIndex", AliyunTablestoreGoSdk)
	}

	d.SetId(ID(args.instanceName, args.tableName, args.indexName, SearchIndexTypeHolder))
	return resourceAliyunOtsSearchIndexRead(d, meta)
}

func (args *SearchIndexResourceArgs) checkArgs(tableResp *tablestore.DescribeTableResponse) error {

	if args.indexName == args.tableName {
		return WrapError(fmt.Errorf("index name cannot be the same as table: %s/%s", args.indexName, args.tableName))
	}

	if tableResp.TableOption.TimeToAlive != -1 {
		return WrapError(fmt.Errorf("when creating a search index, the TimeToAlive of the table must be -1: %v", tableResp.TableOption.TimeToAlive))
	}
	if tableResp.TableOption.MaxVersion != 1 {
		return WrapError(fmt.Errorf("when creating a search index, the table's MaxVersion must be 1: %v", tableResp.TableOption.MaxVersion))
	}

	for _, fs := range args.schema.FieldSchemas {
		if fs.FieldType == tablestore.FieldType_NESTED && *fs.EnableSortAndAgg {
			return WrapError(fmt.Errorf("search index nested type field do not support enable_sort_and_agg: %s/%s", args.indexName, *fs.FieldName))
		}
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
				return WrapError(fmt.Errorf("search index router field must be in primary key: %s/%s", args.indexName, router))
			}
		}
	}

	if args.schema.IndexSort != nil {
		if len(args.schema.IndexSort.Sorters) > 0 {
			for _, fs := range args.schema.FieldSchemas {
				if fs.FieldType == tablestore.FieldType_NESTED {
					return WrapError(fmt.Errorf("search index with nested field type does not support index_sort: %s/%s", args.indexName, *fs.FieldName))
				}
			}
		}
	}

	return nil
}

func resourceAliyunOtsSearchIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	idx, err := otsService.DescribeOtsSearchIndex(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if idx == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("index_id", d.Id()); err != nil {
		return WrapError(err)
	}
	if err := d.Set("create_time", idx.CreateTime); err != nil {
		return WrapError(err)
	}

	phase, err := ConvertSearchIndexSyncPhase(idx.SyncStat.SyncPhase)
	if err != nil {
		return WrapError(err)
	}
	if err := d.Set("sync_phase", string(phase)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("current_sync_timestamp", *idx.SyncStat.CurrentSyncTimestamp); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunOtsSearchIndexDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, indexName, _, err := ParseIndexId(d.Id())
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	if err := otsService.DeleteSearchIndex(instanceName, tableName, indexName); err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteSearchIndex", AliyunTablestoreGoSdk)
	}

	return WrapError(otsService.WaitForSearchIndex(instanceName, tableName, indexName, Deleted, DefaultTimeout))
}
