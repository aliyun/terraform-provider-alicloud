package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"table_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_key": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateAllowedStringValue([]string{
								string(IntegerType), string(BinaryType), string(StringType)}),
						},
					},
				},
				MaxItems: 4,
			},
			"time_to_live": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(-1, INT_MAX),
			},
			"max_version": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, INT_MAX),
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

	for _, primaryKey := range d.Get("primary_key").([]interface{}) {
		pk := primaryKey.(map[string]interface{})
		pkValue := otsService.getPrimaryKeyType(pk["type"].(string))
		tableMeta.AddPrimaryKeyColumn(pk["name"].(string), pkValue)
	}
	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = d.Get("time_to_live").(int)
	tableOption.MaxVersion = d.Get("max_version").(int)

	reservedThroughput := new(tablestore.ReservedThroughput)

	createTableRequest := new(tablestore.CreateTableRequest)
	createTableRequest.TableMeta = tableMeta
	createTableRequest.TableOption = tableOption
	createTableRequest.ReservedThroughput = reservedThroughput

	_, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
		return tableStoreClient.CreateTable(createTableRequest)
	})
	if err != nil {
		return fmt.Errorf("failed to create table with error: %s", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, err := parseId(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	describe, err := otsService.DescribeOtsTable(instanceName, tableName)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to describe table with error: %s", err)
	}

	d.Set("instance_name", instanceName)
	d.Set("table_name", describe.TableMeta.TableName)

	var pks []map[string]interface{}
	keys := describe.TableMeta.SchemaEntry
	for _, v := range keys {
		item := make(map[string]interface{})
		item["name"] = *v.Name
		item["type"] = otsService.convertPrimaryKeyType(*v.Type)
		pks = append(pks, item)
	}
	d.Set("primary_key", pks)

	d.Set("time_to_live", describe.TableOption.TimeToAlive)
	d.Set("max_version", describe.TableOption.MaxVersion)

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	// As the issue of ots sdk, time_to_live and max_version need to be updated together at present.
	// For the issue, please refer to https://github.com/aliyun/aliyun-tablestore-go-sdk/issues/18
	if d.HasChange("time_to_live") || d.HasChange("max_version") {
		instanceName, tableName, err := parseId(d, meta)
		if err != nil {
			return err
		}
		client := meta.(*connectivity.AliyunClient)

		updateTableReq := new(tablestore.UpdateTableRequest)
		updateTableReq.TableName = tableName
		tableOption := new(tablestore.TableOption)

		tableOption.TimeToAlive = d.Get("time_to_live").(int)
		tableOption.MaxVersion = d.Get("max_version").(int)

		updateTableReq.TableOption = tableOption
		if _, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.UpdateTable(updateTableReq)
		}); err != nil {
			return fmt.Errorf("failed to update table with error: %s", err)
		}
	}
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, err := parseId(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	req := new(tablestore.DeleteTableRequest)
	req.TableName = tableName
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := otsService.DescribeOtsInstance(instanceName); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting table %s, describing instance %s got an error: %#v.", tableName, instanceName, err))
		}
		_, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.DeleteTable(req)
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), OTSObjectNotExist) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting table %s got an error: %#v.", tableName, err))
		}
		if _, err := otsService.DescribeOtsTable(instanceName, tableName); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting table %s, describing table got an error: %#v.", tableName, err))
		}
		return resource.RetryableError(fmt.Errorf("delete table %s timeout.", tableName))
	})
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
			err = fmt.Errorf("From Provider version 1.10.0, the provider field 'ots_instance_name' has been deprecated and " +
				"you should use resource alicloud_ots_table's new field 'instance_name' and 'table_name' to re-import this resource.")
			return
		}
	} else {
		instanceName = split[0]
		tableName = split[1]
	}

	return
}
