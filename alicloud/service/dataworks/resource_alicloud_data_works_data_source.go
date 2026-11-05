// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package dataworks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataSourceCreate,
		Read:   resourceDataSourceRead,
		Update: resourceDataSourceUpdate,
		Delete: resourceDataSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_properties": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"connection_properties_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modify_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modify_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"qualified_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ProjectId"] = d.Get("project_id")
	request["RegionId"] = client.RegionId

	request["Type"] = d.Get("type")
	request["ConnectionPropertiesMode"] = d.Get("connection_properties_mode")
	request["ConnectionProperties"] = d.Get("connection_properties")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Name"] = d.Get("data_source_name")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_data_works_data_source", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ProjectId"], response["Id"]))

	return resourceDataSourceRead(d, meta)
}

func resourceDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := NewService(client)

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDataSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_data_source DescribeDataWorksDataSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if v, ok := objectRaw["ConnectionProperties"]; ok {
		connPropJson, err := json.Marshal(v)
		if err != nil {
			return errs.WrapError(err)
		}
		d.Set("connection_properties", string(connPropJson))
	}
	if v, ok := objectRaw["ConnectionPropertiesMode"]; ok {
		d.Set("connection_properties_mode", v)
	}
	if v, ok := objectRaw["CreateTime"]; ok {
		d.Set("create_time", v)
	}
	if v, ok := objectRaw["CreateUser"]; ok {
		d.Set("create_user", v)
	}
	if v, ok := objectRaw["Name"]; ok {
		d.Set("data_source_name", v)
	}
	if v, ok := objectRaw["Description"]; ok {
		d.Set("description", v)
	}
	if v, ok := objectRaw["ModifyTime"]; ok {
		d.Set("modify_time", v)
	}
	if v, ok := objectRaw["ModifyUser"]; ok {
		d.Set("modify_user", v)
	}
	if v, ok := objectRaw["QualifiedName"]; ok {
		d.Set("qualified_name", v)
	}
	if v, ok := objectRaw["Type"]; ok {
		d.Set("type", v)
	}
	if v, ok := objectRaw["Id"]; ok {
		d.Set("data_source_id", v)
	}
	if v, ok := objectRaw["ProjectId"]; ok {
		d.Set("project_id", v)
	}

	return nil
}

func resourceDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "UpdateDataSource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProjectId"] = parts[0]
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("connection_properties_mode") {
		update = true
	}
	request["ConnectionPropertiesMode"] = d.Get("connection_properties_mode")
	if d.HasChange("connection_properties") {
		update = true
	}
	request["ConnectionProperties"] = d.Get("connection_properties")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceDataSourceRead(d, meta)
}

func resourceDataSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = parts[1]
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, request)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
