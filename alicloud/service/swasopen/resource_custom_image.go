package swasopen

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCustomImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomImageCreate,
		Read:   resourceCustomImageRead,
		Update: resourceCustomImageUpdate,
		Delete: resourceCustomImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"custom_image_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9a-zA-Z一-龥][一-龥A-Za-z0-9:_-]{2,128}$"), "The name must be `2` to `128` characters in length. It must start with a letter or a number. It can contain letters, digits, colons (:), underscores (_) and hyphens (-)."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Share", "UnShare"}, false),
			},
			"system_snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCustomImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateCustomImage"
	request := make(map[string]interface{})
	var err error
	request["ImageName"] = d.Get("custom_image_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	request["SystemSnapshotId"] = d.Get("system_snapshot_id")
	request["ClientToken"] = helper.BuildClientToken("CreateCustomImage")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_simple_application_server_custom_image", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ImageId"]))

	swasOpenService := NewSwasOpenService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Exist"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, swasOpenService.SimpleApplicationServerCustomImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceCustomImageUpdate(d, meta)
}
func resourceCustomImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := NewSwasOpenService(client)
	object, err := swasOpenService.DescribeSimpleApplicationServerCustomImage(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_custom_image swasOpenService.DescribeSimpleApplicationServerCustomImage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("custom_image_name", object["ImageName"])
	d.Set("description", object["Description"])
	return nil
}
func resourceCustomImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ImageId": d.Id(),
	}
	if d.HasChange("status") {
		update = true
		if v, ok := d.GetOk("status"); ok {
			request["Operation"] = v
		}
	}
	if update {
		action := "ModifyImageShareStatus"
		request["ClientToken"] = helper.BuildClientToken("ModifyImageShareStatus")
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
	return resourceCustomImageRead(d, meta)
}
func resourceCustomImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if v := d.Get("status"); v.(string) == "Share" {
		var response map[string]interface{}
		var err error
		request := map[string]interface{}{
			"ImageId":   d.Id(),
			"Operation": "UnShare",
		}
		action := "ModifyImageShareStatus"
		request["ClientToken"] = helper.BuildClientToken("ModifyImageShareStatus")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
	action := "DeleteCustomImage"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ImageId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("DeleteCustomImage")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, nil, request, true)
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
	swasOpenService := NewSwasOpenService(client)
	stateConf := helper.BuildStateConf([]string{"Exist"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, swasOpenService.SimpleApplicationServerCustomImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
