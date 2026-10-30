package nas

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudNasRecycleBin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasRecycleBinCreate,
		Read:   resourceAlicloudNasRecycleBinRead,
		Update: resourceAlicloudNasRecycleBinUpdate,
		Delete: resourceAlicloudNasRecycleBinDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserved_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 180),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudNasRecycleBinCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "EnableRecycleBin"
	request := make(map[string]interface{})
	var err error
	request["FileSystemId"] = d.Get("file_system_id")
	if v, ok := d.GetOk("reserved_days"); ok {
		request["ReservedDays"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_nas_recycle_bin", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"]))

	return resourceAlicloudNasRecycleBinRead(d, meta)
}
func resourceAlicloudNasRecycleBinRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasRecycleBin(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_recycle_bin nasService.DescribeNasRecycleBin Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("file_system_id", d.Id())
	reserved_daysVal, err := convert.FormatInt(object["ReservedDays"])
	if err != nil {
		return err
	}
	d.Set("reserved_days", reserved_daysVal)
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudNasRecycleBinUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
	}
	var response map[string]interface{}
	if d.HasChange("reserved_days") {
		if v, ok := d.GetOk("reserved_days"); ok {
			request["ReservedDays"] = v
		}
	}
	action := "UpdateRecycleBinAttribute"
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcGet("NAS", "2017-06-26", action, request, nil)
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
	return resourceAlicloudNasRecycleBinRead(d, meta)
}
func resourceAlicloudNasRecycleBinDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DisableAndCleanRecycleBin"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcGet("NAS", "2017-06-26", action, request, nil)
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
