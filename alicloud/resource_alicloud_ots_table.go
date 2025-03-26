package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOtsTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsTableCreate,
		Read:   resourceAliyunOtsTableRead,
		Update: resourceAliyunOtsTableUpdate,
		Delete: resourceAliyunOtsTableDelete,
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
			"primary_key": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(IntegerType), string(BinaryType), string(StringType)}, false),
						},
					},
				},
				MaxItems: 4,
				ForceNew: true,
			},
			"defined_column": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(DefinedColumnInteger), string(DefinedColumnString), string(DefinedColumnBinary), string(DefinedColumnDouble), string(DefinedColumnBoolean)}, false),
						},
					},
				},
				MaxItems: 32,
			},
			"time_to_live": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(-1, INT_MAX),
			},
			"max_version": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, INT_MAX),
			},
			"allow_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deviation_cell_version_in_sec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringConvertInt64(),
				Default:      "86400",
			},
			"enable_sse": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sse_key_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SseKMSService), string(SseByOk)}, false),
			},
			"sse_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return string(SseByOk) != d.Get("sse_key_type").(string)
				},
			},
			"sse_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return string(SseByOk) != d.Get("sse_key_type").(string)
				},
			},
		},
	}
}

func resourceAliyunOtsTableCreate(d *schema.ResourceData, meta interface{}) error {
	tableMeta := new(tablestore.TableMeta)
	instanceName := d.Get("instance_name").(string)
	tableName := d.Get("table_name").(string)
	tableMeta.TableName = tableName
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	if err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, e := otsService.DescribeOtsInstance(instanceName)
		if e != nil {
			if NotFoundError(e) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		return nil
	}); err != nil {
		return WrapError(err)
	}
	for _, primaryKey := range d.Get("primary_key").([]interface{}) {
		pk := primaryKey.(map[string]interface{})
		pkValue := otsService.getPrimaryKeyType(pk["type"].(string))
		tableMeta.AddPrimaryKeyColumn(pk["name"].(string), pkValue)
	}

	if v, ok := d.GetOk("defined_column"); ok {
		definedColumns := v.([]interface{})
		for _, definedColumn := range definedColumns {
			columnArgs := definedColumn.(map[string]interface{})
			columnType, err := ParseDefinedColumnType(columnArgs["type"].(string))
			if err != nil {
				return WrapError(err)
			}
			tableMeta.AddDefinedColumn(columnArgs["name"].(string), columnType)
		}
	}

	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = d.Get("time_to_live").(int)
	tableOption.MaxVersion = d.Get("max_version").(int)
	allowUpdate := d.Get("allow_update").(bool)
	tableOption.AllowUpdate = &allowUpdate
	if deviation, ok := d.GetOk("deviation_cell_version_in_sec"); ok {
		tableOption.DeviationCellVersionInSec, _ = strconv.ParseInt(deviation.(string), 10, 64)
	}
	reservedThroughput := new(tablestore.ReservedThroughput)

	request := new(tablestore.CreateTableRequest)
	request.TableMeta = tableMeta
	request.TableOption = tableOption
	request.ReservedThroughput = reservedThroughput

	if enableSSE, ok := d.GetOkExists("enable_sse"); ok {
		sseSpec := new(tablestore.SSESpecification)
		sseSpec.Enable = enableSSE.(bool)

		if sseKeyType, ok2 := d.GetOk("sse_key_type"); ok2 {
			var typ tablestore.SSEKeyType
			switch sseKeyType.(string) {
			case string(SseKMSService):
				typ = tablestore.SSE_KMS_SERVICE
			case string(SseByOk):
				typ = tablestore.SSE_BYOK
			default:
				return WrapError(Error("unknown sse key type: %s", sseKeyType.(string)))
			}
			sseSpec.KeyType = &typ
		}
		if sseKeyId, ok3 := d.GetOk("sse_key_id"); ok3 {
			keyId := sseKeyId.(string)
			sseSpec.KeyId = &keyId
		}

		if sseRoleArn, ok4 := d.GetOk("sse_role_arn"); ok4 {
			roleArn := sseRoleArn.(string)
			sseSpec.RoleArn = &roleArn
		}
		request.SSESpecification = sseSpec
	}

	var requestinfo *tablestore.TableStoreClient
	if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestinfo = tableStoreClient
			return tableStoreClient.CreateTable(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateTable", raw, requestinfo, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_table", "CreateTable", AliyunTablestoreGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	instanceName, _, err := parseId(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	tableResp, err := otsService.DescribeOtsTable(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if tableResp == nil {
		d.SetId("")
		return nil
	}

	d.Set("instance_name", instanceName)
	d.Set("table_name", tableResp.TableMeta.TableName)

	var pks []map[string]interface{}
	keys := tableResp.TableMeta.SchemaEntry
	for _, v := range keys {
		item := make(map[string]interface{})
		item["name"] = *v.Name
		item["type"] = otsService.convertPrimaryKeyType(*v.Type)
		pks = append(pks, item)
	}
	d.Set("primary_key", pks)

	var columns []map[string]interface{}
	for _, column := range tableResp.TableMeta.DefinedColumns {
		columnType, err := ConvertDefinedColumnType(column.ColumnType)
		if err != nil {
			return WrapError(err)
		}
		item := map[string]interface{}{
			"name": column.Name,
			"type": columnType,
		}
		columns = append(columns, item)
	}
	d.Set("defined_column", columns)

	d.Set("time_to_live", tableResp.TableOption.TimeToAlive)
	d.Set("max_version", tableResp.TableOption.MaxVersion)
	d.Set("allow_update", *tableResp.TableOption.AllowUpdate)
	d.Set("deviation_cell_version_in_sec", strconv.FormatInt(tableResp.TableOption.DeviationCellVersionInSec, 10))

	if tableResp.SSEDetails != nil && tableResp.SSEDetails.Enable {
		d.Set("enable_sse", tableResp.SSEDetails.Enable)
		d.Set("sse_key_type", tableResp.SSEDetails.KeyType.String())
		d.Set("sse_key_id", tableResp.SSEDetails.KeyId)
		d.Set("sse_role_arn", tableResp.SSEDetails.RoleArn)

	}

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	// As the issue of ots sdk, time_to_live and max_version need to be updated together at present.
	// For the issue, please refer to https://github.com/aliyun/aliyun-tablestore-go-sdk/issues/18
	instanceName, tableName, err := parseId(d, meta)
	if err != nil {
		return err
	}
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("time_to_live") || d.HasChange("max_version") || d.HasChange("deviation_cell_version_in_sec") || d.HasChange("allow_update") {
		request := new(tablestore.UpdateTableRequest)
		request.TableName = tableName
		tableOption := new(tablestore.TableOption)

		tableOption.TimeToAlive = d.Get("time_to_live").(int)
		tableOption.MaxVersion = d.Get("max_version").(int)
		if deviation, ok := d.GetOk("deviation_cell_version_in_sec"); ok {
			tableOption.DeviationCellVersionInSec, _ = strconv.ParseInt(deviation.(string), 10, 64)
		}
		allowUpdate := d.Get("allow_update").(bool)
		tableOption.AllowUpdate = &allowUpdate

		request.TableOption = tableOption
		var requestinfo *tablestore.TableStoreClient
		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				requestinfo = tableStoreClient
				return tableStoreClient.UpdateTable(request)
			})
			if err != nil {
				if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("UpdateTable", raw, requestinfo, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTable", AliyunTablestoreGoSdk)
		}
	}
	if d.HasChange("defined_column") {
		o, n := d.GetChange("defined_column")
		statedColumns, err := parseColsFromConfig(o.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		declareColumns, err := parseColsFromConfig(n.([]interface{}))
		if err != nil {
			return WrapError(err)
		}

		if needAddColumns, fetchErr := fetchNeedAddColumns(declareColumns, statedColumns); fetchErr != nil {
			return fetchErr
		} else if err := updateDefinedColumns(client, instanceName, tablestore.AddDefinedColumnRequest{TableName: tableName}, needAddColumns); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "AddTableDefineColumn", AliyunTablestoreGoSdk)
		}

		if needDeleteColumns, fetchErr := fetchNeedDeleteColumns(statedColumns, declareColumns); fetchErr != nil {
			return fetchErr
		} else if err := updateDefinedColumns(client, instanceName, tablestore.DeleteDefinedColumnRequest{TableName: tableName}, needDeleteColumns); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTableDefineColumn", AliyunTablestoreGoSdk)
		}

	}
	return resourceAliyunOtsTableRead(d, meta)
}

func parseColsFromConfig(cols []interface{}) ([]*tablestore.DefinedColumnSchema, error) {
	otsCols := make([]*tablestore.DefinedColumnSchema, 0, len(cols))
	for _, col := range cols {
		colMap := col.(map[string]interface{})
		columnType, err := ParseDefinedColumnType(colMap["type"].(string))
		if err != nil {
			return nil, WrapError(err)
		}
		otsCols = append(otsCols, &tablestore.DefinedColumnSchema{
			Name:       colMap["name"].(string),
			ColumnType: columnType,
		})
	}
	return otsCols, nil
}

func fetchNeedDeleteColumns(statedColumns []*tablestore.DefinedColumnSchema, declareColumns []*tablestore.DefinedColumnSchema) ([]*tablestore.DefinedColumnSchema, error) {
	var needDeleteColumns []*tablestore.DefinedColumnSchema
	for _, statedCol := range statedColumns {
		switch FindDefinedColumn(declareColumns, statedCol) {
		case ExistEqual:
			continue
		case ExistNotEqual:
			return nil, WrapError(fmt.Errorf("modifying defined column type is not supported: %v", statedCol))
		case NotExist:
			needDeleteColumns = append(needDeleteColumns, statedCol)
		}

	}
	return needDeleteColumns, nil
}

func fetchNeedAddColumns(declareColumns []*tablestore.DefinedColumnSchema, statedColumns []*tablestore.DefinedColumnSchema) ([]*tablestore.DefinedColumnSchema, error) {
	var needAddColumns []*tablestore.DefinedColumnSchema
	for _, declareCol := range declareColumns {
		switch FindDefinedColumn(statedColumns, declareCol) {
		case ExistEqual:
			continue
		case ExistNotEqual:
			return nil, WrapError(fmt.Errorf("modifying defined column type is not supported: %v", declareCol))
		case NotExist:
			needAddColumns = append(needAddColumns, declareCol)
		}
	}
	return needAddColumns, nil
}

func updateDefinedColumns(client *connectivity.AliyunClient, instance string, req interface{}, columns []*tablestore.DefinedColumnSchema) error {
	var clientInfo *tablestore.TableStoreClient

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(instance, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			clientInfo = tableStoreClient

			switch request := req.(type) {
			case tablestore.AddDefinedColumnRequest:
				for _, column := range columns {
					request.AddDefinedColumn(column.Name, column.ColumnType)
				}
				return tableStoreClient.AddDefinedColumn(&request)
			case tablestore.DeleteDefinedColumnRequest:
				for _, column := range columns {
					request.DefinedColumns = append(request.DefinedColumns, column.Name)
				}
				return tableStoreClient.DeleteDefinedColumn(&request)
			default:
				return nil, WrapError(fmt.Errorf("unexpected defined column request type %T: %v", req, req))
			}
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("UpdateTableDefineColumn", raw, clientInfo, req)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, err := parseId(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	req := new(tablestore.DeleteTableRequest)
	req.TableName = tableName
	var requestCli *tablestore.TableStoreClient
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestCli = tableStoreClient
			return tableStoreClient.DeleteTable(req)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteTable", raw, requestCli, req)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTable", AliyunTablestoreGoSdk)
	}
	return WrapError(otsService.WaitForOtsTable(instanceName, tableName, Deleted, DefaultTimeout))
}

func parseId(d *schema.ResourceData, meta interface{}) (instanceName, tableName string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) == 1 {
		// For compatibility
		if meta.(*connectivity.AliyunClient).OtsInstanceName != "" {
			tableName = split[0]
			instanceName = meta.(*connectivity.AliyunClient).OtsInstanceName
			d.SetId(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
		} else {
			err = WrapError(Error("From Provider version 1.10.0, the provider field 'ots_instance_name' has been deprecated and " +
				"you should use resource alicloud_ots_table's new field 'instance_name' and 'table_name' to re-import this resource."))
			return
		}
	} else {
		instanceName = split[0]
		tableName = split[1]
	}

	return
}
