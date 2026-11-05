package oos

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosApplicationCreate,
		Read:   resourceAlicloudOosApplicationRead,
		Update: resourceAlicloudOosApplicationUpdate,
		Delete: resourceAlicloudOosApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceAlicloudOosApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateApplication"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Name"] = d.Get("application_name")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["ClientToken"] = helper.BuildClientToken("CreateApplication")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oos_application", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseApplication := response["Application"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseApplication["Name"]))

	return resourceAlicloudOosApplicationRead(d, meta)
}
func resourceAlicloudOosApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosApplication(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_application oosService.DescribeOosApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("application_name", object["Name"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("tags", tags.TagsToMap(object["Tags"]))
	return nil
}
func resourceAlicloudOosApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			respJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return errs.WrapError(err)
			}
			request["Tags"] = respJson
		}
	}
	if update {
		action := "UpdateApplication"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
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
	return resourceAlicloudOosApplicationRead(d, meta)
}
func resourceAlicloudOosApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApplication"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"DeleteConflict.Application"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Application"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
