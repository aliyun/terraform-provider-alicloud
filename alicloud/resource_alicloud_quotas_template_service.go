// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudQuotasTemplateService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudQuotasTemplateServiceCreate,
		Read:   resourceAliCloudQuotasTemplateServiceRead,
		Update: resourceAliCloudQuotasTemplateServiceUpdate,
		Delete: resourceAliCloudQuotasTemplateServiceDelete,
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
				ValidateFunc: IntInSlice([]int{0, 1, -1}),
			},
		},
	}
}

func resourceAliCloudQuotasTemplateServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyQuotaTemplateServiceStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ServiceStatus"] = d.Get("service_status")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_template_service", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := meta.(*connectivity.AliyunClient).AccountId()
	d.SetId(fmt.Sprintf(accountId))

	return resourceAliCloudQuotasTemplateServiceRead(d, meta)
}

func resourceAliCloudQuotasTemplateServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasTemplateService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_template_service DescribeQuotasTemplateService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["ServiceStatus"] != nil {
		d.Set("service_status", objectRaw["ServiceStatus"])
	}

	return nil
}

func resourceAliCloudQuotasTemplateServiceUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudQuotasTemplateServiceRead(d, meta)
}

func resourceAliCloudQuotasTemplateServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Template Service. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
