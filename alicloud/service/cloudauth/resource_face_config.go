package cloudauth

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
)

func ResourceFaceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceFaceConfigCreate,
		Read:   resourceFaceConfigRead,
		Update: resourceFaceConfigUpdate,
		Delete: resourceFaceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"biz_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"biz_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFaceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFaceConfig"
	request := make(map[string]interface{})
	var err error
	request["BizName"] = d.Get("biz_name")
	request["BizType"] = d.Get("biz_type")
	request["SourceIp"] = helper.BuildClientToken(action)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cloudauth", "2019-03-07", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloudauth_face_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["BizType"]))

	return resourceFaceConfigRead(d, meta)
}

func resourceFaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudauthService := NewService(client)
	object, err := cloudauthService.DescribeFaceConfig(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloudauth_face_config cloudauthService.DescribeFaceConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("biz_name", object["BizName"])
	d.Set("biz_type", object["BizType"])
	d.Set("gmt_modified", object["GmtUpdated"])
	return nil
}

func resourceFaceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"BizType": d.Id(),
	}
	if d.HasChange("biz_name") {
		update = true
	}
	request["BizName"] = d.Get("biz_name")
	if update {
		action := "UpdateFaceConfig"
		request["SourceIp"] = helper.BuildClientToken(action)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cloudauth", "2019-03-07", action, nil, request, false)
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
	return resourceFaceConfigRead(d, meta)
}

func resourceFaceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceFaceConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
