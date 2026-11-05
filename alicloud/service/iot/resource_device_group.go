package iot

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeviceGroupCreate,
		Read:   resourceDeviceGroupRead,
		Update: resourceDeviceGroupUpdate,
		Delete: resourceDeviceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 30),
			},
			"iot_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"super_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDeviceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDeviceGroup"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("group_desc"); ok {
		request["GroupDesc"] = v
	}
	request["GroupName"] = d.Get("group_name")
	if v, ok := d.GetOk("iot_instance_id"); ok {
		request["IotInstanceId"] = v
	}
	if v, ok := d.GetOk("super_group_id"); ok {
		request["SuperGroupId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Iot", "2018-01-20", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_iot_device_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["GroupId"]))

	return resourceDeviceGroupRead(d, meta)
}
func resourceDeviceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	iotService := NewService(client)
	object, err := iotService.DescribeDeviceGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_iot_device_group iotService.DescribeDeviceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	descmap := object["Data"].(map[string]interface{})
	if val, exist := descmap["GroupDesc"]; exist {
		d.Set("group_desc", val)
	}

	d.Set("group_name", object["Data"].(map[string]interface{})["GroupName"])
	return nil
}
func resourceDeviceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}
	if d.HasChange("group_desc") {
		update = true
		if v, ok := d.GetOk("group_desc"); ok {
			request["GroupDesc"] = v
		}
	}
	if d.HasChange("iot_instance_id") {
		update = true
		if v, ok := d.GetOk("iot_instance_id"); ok {
			request["IotInstanceId"] = v
		}
	}
	if update {
		action := "UpdateDeviceGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Iot", "2018-01-20", action, nil, request, false)
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
	return resourceDeviceGroupRead(d, meta)
}
func resourceDeviceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeviceGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}

	if v, ok := d.GetOk("iot_instance_id"); ok {
		request["IotInstanceId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Iot", "2018-01-20", action, nil, request, false)
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
	return nil
}
