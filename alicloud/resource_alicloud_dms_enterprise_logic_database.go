package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDmsEnterpriseLogicDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseLogicDatabaseCreate,
		Read:   resourceAlicloudDmsEnterpriseLogicDatabaseRead,
		Update: resourceAlicloudDmsEnterpriseLogicDatabaseUpdate,
		Delete: resourceAlicloudDmsEnterpriseLogicDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Required: true,
				Type:     schema.TypeString,
			},
			"database_ids": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"env_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"logic": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"logic_database_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"owner_id_list": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_name_list": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"schema_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"search_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudDmsEnterpriseLogicDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("alias"); ok {
		request["Alias"] = v
	}

	if v, ok := d.GetOk("database_ids"); ok {
		request["DatabaseIds"] = convertListToJsonString(v.([]interface{}))
	}

	var response map[string]interface{}
	action := "CreateLogicDatabase"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_logic_database", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.LogicDbId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_dms_enterprise_logic_database")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudDmsEnterpriseLogicDatabaseRead(d, meta)
}

func resourceAlicloudDmsEnterpriseLogicDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}

	object, err := dmsEnterpriseService.DescribeDmsEnterpriseLogicDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_logic_database dmsEnterpriseService.DescribeDmsEnterpriseLogicDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("logic_database_id", object["DatabaseId"])
	d.Set("alias", object["Alias"])
	databaseIds, _ := jsonpath.Get("$.DatabaseIds.DatabaseIds", object)
	d.Set("database_ids", databaseIds)
	d.Set("db_type", object["DbType"])
	d.Set("env_type", object["EnvType"])
	d.Set("logic", object["Logic"])
	ownerIdList, _ := jsonpath.Get("$.OwnerIdList.OwnerIds", object)
	d.Set("owner_id_list", ownerIdList)
	ownerNameList, _ := jsonpath.Get("$.OwnerNameList.OwnerNames", object)
	d.Set("owner_name_list", ownerNameList)
	d.Set("schema_name", object["SchemaName"])
	d.Set("search_name", object["SearchName"])

	return nil
}

func resourceAlicloudDmsEnterpriseLogicDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"LogicDbId": d.Id(),
		"RegionId":  client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("alias") {
		update = true
	}
	request["Alias"] = d.Get("alias")
	if d.HasChange("database_ids") {
		update = true
		if v, ok := d.GetOk("database_ids"); ok {
			request["DatabaseIds"] = convertListToJsonString(v.([]interface{}))
		}
	}

	if update {
		action := "EditLogicDatabase"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudDmsEnterpriseLogicDatabaseRead(d, meta)
}

func resourceAlicloudDmsEnterpriseLogicDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"LogicDbId": d.Id(),
		"RegionId":  client.RegionId,
	}

	action := "DeleteLogicDatabase"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
