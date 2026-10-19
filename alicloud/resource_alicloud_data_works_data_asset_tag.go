// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksDataAssetTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDataAssetTagCreate,
		Read:   resourceAliCloudDataWorksDataAssetTagRead,
		Update: resourceAliCloudDataWorksDataAssetTagUpdate,
		Delete: resourceAliCloudDataWorksDataAssetTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[\w.-]+$`), "Data Asset Tag Key"),
			},
			"managers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Boolean", "Int", "String", "Double"}, false),
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAliCloudDataWorksDataAssetTagCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDataAssetTag"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["Key"] = d.Get("key")
	request["ValueType"] = d.Get("value_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("values"); ok && len(v.([]interface{})) > 0 {
		valuesJson, err := json.Marshal(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Values"] = string(valuesJson)
	}
	if v, ok := d.GetOk("managers"); ok && len(v.([]interface{})) > 0 {
		managersJson, err := json.Marshal(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Managers"] = string(managersJson)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"9990020002", "9990040003"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_data_asset_tag", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("key")))

	return resourceAliCloudDataWorksDataAssetTagRead(d, meta)
}

func resourceAliCloudDataWorksDataAssetTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDataAssetTag(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_data_asset_tag DescribeDataWorksDataAssetTag Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key", objectRaw["Key"])
	d.Set("value_type", objectRaw["ValueType"])
	if v, ok := objectRaw["Description"]; ok && v != nil {
		d.Set("description", v)
	}
	if v, ok := objectRaw["Category"]; ok && v != nil {
		d.Set("category", v)
	}
	if v, ok := objectRaw["CreateTime"]; ok && v != nil {
		d.Set("create_time", fmt.Sprint(v))
	}
	if v, ok := objectRaw["ModifyTime"]; ok && v != nil {
		d.Set("modify_time", fmt.Sprint(v))
	}
	if v, ok := objectRaw["Values"]; ok && v != nil {
		d.Set("values", v)
	}
	if v, ok := objectRaw["Managers"]; ok && v != nil {
		d.Set("managers", v)
	}

	return nil
}

func resourceAliCloudDataWorksDataAssetTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateDataAssetTag"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["Key"] = d.Get("key")
	update := false

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("values") {
		update = true
		valuesJson, err := json.Marshal(d.Get("values").([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Values"] = string(valuesJson)
	}

	if !d.IsNewResource() && d.HasChange("managers") {
		update = true
		managersJson, err := json.Marshal(d.Get("managers").([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Managers"] = string(managersJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudDataWorksDataAssetTagRead(d, meta)
}

func resourceAliCloudDataWorksDataAssetTagDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDataAssetTag"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["Key"] = d.Get("key")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
