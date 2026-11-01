// Package quotas. This file is generated automatically. Please do not modify it manually, thank you!
package quotas

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTemplateService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateServiceCreate,
		Read:   resourceTemplateServiceRead,
		Update: resourceTemplateServiceUpdate,
		Delete: resourceTemplateServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 1, -1}),
			},
		},
	}
}

func resourceTemplateServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyQuotaTemplateServiceStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ServiceStatus"] = d.Get("service_status")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_quotas_template_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	accountId, err := meta.(*connectivity.AliyunClient).AccountId()
	d.SetId(fmt.Sprint(accountId))

	return resourceTemplateServiceRead(d, meta)
}

func resourceTemplateServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasTemplateService(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_template_service DescribeQuotasTemplateService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["ServiceStatus"] != nil {
		d.Set("service_status", objectRaw["ServiceStatus"])
	}

	return nil
}

func resourceTemplateServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyQuotaTemplateServiceStatus"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("service_status") {
		update = true
	}
	request["ServiceStatus"] = d.Get("service_status")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceTemplateServiceRead(d, meta)
}

func resourceTemplateServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Template Service. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
