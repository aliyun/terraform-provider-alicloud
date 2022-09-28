package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOtsSecondaryIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsSecondaryIndexCreate,
		Read:   resourceAliyunOtsSecondaryIndexRead,
		Delete: resourceAliyunOtsSecondaryIndexDelete,
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
			"index_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(Local), string(Global)}, false),
			},
			"include_base_data": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"primary_keys": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 4,
			},
			"defined_columns": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 32,
			},
		},
	}
}

func parseArgs(d *schema.ResourceData) *SecIndexResourceArgs {
	args := &SecIndexResourceArgs{
		instanceName:    d.Get("instance_name").(string),
		tableName:       d.Get("table_name").(string),
		includeBaseData: d.Get("include_base_data").(bool),
		indexName:       d.Get("index_name").(string),
		indexType:       SecondaryIndexTypeString(d.Get("index_type").(string)),
		primaryKeys:     Unique(Interface2StrSlice(d.Get("primary_keys").([]interface{}))),
	}

	if v, ok := d.GetOk("defined_columns"); ok && len(v.([]interface{})) > 0 {
		args.definedColumns = Unique(Interface2StrSlice(v.([]interface{})))
	}
	return args
}

type SecIndexResourceArgs struct {
	instanceName    string
	tableName       string
	includeBaseData bool
	indexName       string
	indexType       SecondaryIndexTypeString
	primaryKeys     []string
	definedColumns  []string
}

func resourceAliyunOtsSecondaryIndexCreate(d *schema.ResourceData, meta interface{}) error {
	args := parseArgs(d)

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	// check table exists
	tableResp, err := otsService.LoopWaitTable(args.instanceName, args.tableName)
	if err != nil {
		return WrapError(err)
	}
	// serverside arguments check
	if err := checkArgs(tableResp, args); err != nil {
		return err
	}
	// build request
	idxType, err := ConvertSecIndexTypeString(args.indexType)
	if err != nil {
		return WrapError(err)
	}
	req := &tablestore.CreateIndexRequest{
		MainTableName:   args.tableName,
		IncludeBaseData: args.includeBaseData,

		IndexMeta: &tablestore.IndexMeta{
			IndexName:      args.indexName,
			Primarykey:     args.primaryKeys,
			IndexType:      idxType,
			DefinedColumns: args.definedColumns,
		},
	}

	var reqClient *tablestore.TableStoreClient
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(args.instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			reqClient = tableStoreClient
			return tableStoreClient.CreateIndex(req)
		})
		defer func() {
			addDebug("CreateTableSecondaryIndex", raw, reqClient, req)
		}()

		if err != nil {
			if IsExpectedErrors(err, OtsSecondaryIndexIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_secondary_index", "CreateIndex", AliyunTablestoreGoSdk)
	}

	d.SetId(ID(args.instanceName, args.tableName, args.indexName, string(args.indexType)))
	return resourceAliyunOtsSecondaryIndexRead(d, meta)
}

func checkArgs(tableResp *tablestore.DescribeTableResponse, args *SecIndexResourceArgs) error {
	dataPks := simplifyPK(tableResp.TableMeta.SchemaEntry)
	dataCols := simplifyCol(tableResp.TableMeta.DefinedColumns)

	allCards := append(dataPks, dataCols...)
	if !IsSubCollection(args.primaryKeys, allCards) {
		return WrapError(fmt.Errorf("some primary keys not exist in table: %s/%s", args.primaryKeys, args.tableName))
	}
	if !IsSubCollection(args.definedColumns, allCards) {
		return WrapError(fmt.Errorf("some defined columns not exist in table: %s/%s", args.definedColumns, args.tableName))
	}
	if args.indexName == args.tableName {
		return WrapError(fmt.Errorf("index name cannot be the same as table: %s/%s", args.indexName, args.tableName))
	}
	if args.indexType == Local && args.primaryKeys[0] != dataPks[0] {
		return WrapError(fmt.Errorf("when using a local secondary index, the first primary key of the index must be "+
			"the same as the first primary key of the data table: %s/%s", args.primaryKeys, dataPks))
	}
	if tableResp.TableOption.TimeToAlive != -1 {
		return WrapError(fmt.Errorf("when creating a secondary index, the TimeToAlive of the table must be -1: %v", tableResp.TableOption.TimeToAlive))
	}
	if tableResp.TableOption.MaxVersion != 1 {
		return WrapError(fmt.Errorf("when creating a secondary index, the table's MaxVersion must be 1: %v", tableResp.TableOption.MaxVersion))
	}
	return nil
}

func simplifyPK(keys []*tablestore.PrimaryKeySchema) []string {
	var pks = make([]string, 0, len(keys))
	for _, key := range keys {
		pks = append(pks, *key.Name)
	}
	return pks
}

func simplifyCol(columns []*tablestore.DefinedColumnSchema) []string {
	var cols = make([]string, 0, len(columns))
	for _, col := range columns {
		cols = append(cols, col.Name)
	}
	return cols
}

func simplifySecIndex(indexes []*tablestore.IndexMeta) []string {
	var ii = make([]string, 0, len(indexes))
	for _, idx := range indexes {
		ii = append(ii, idx.IndexName)
	}
	return ii
}

func resourceAliyunOtsSecondaryIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	idx, err := otsService.DescribeOtsSecondaryIndex(d.Id())
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

	if err := d.Set("instance_name", idx.InstanceName); err != nil {
		return WrapError(err)
	}
	if err := d.Set("table_name", idx.TableName); err != nil {
		return WrapError(err)
	}
	if err := d.Set("index_name", idx.Index.IndexName); err != nil {
		return WrapError(err)
	}
	indexType, err := ConvertSecIndexType(idx.Index.IndexType)
	if err != nil {
		return WrapError(err)
	}
	if err := d.Set("index_type", string(indexType)); err != nil {
		return WrapError(err)
	}
	if err := d.Set("primary_keys", idx.Index.Primarykey); err != nil {
		return WrapError(err)
	}
	if err := d.Set("defined_columns", idx.Index.DefinedColumns); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunOtsSecondaryIndexDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, indexName, _, err := ParseIndexId(d.Id())
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)

	req := &tablestore.DeleteIndexRequest{
		MainTableName: tableName,
		IndexName:     indexName,
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		var requestInfo *tablestore.TableStoreClient
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.DeleteIndex(req)
		})
		defer func() {
			addDebug("DeleteTableSecondaryIndex", raw, requestInfo, req)
		}()

		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTableSecondaryIndex", AliyunTablestoreGoSdk)
	}

	otsService := OtsService{client}
	return WrapError(otsService.WaitForSecondaryIndex(instanceName, tableName, indexName, Deleted, DefaultTimeout))
}
