// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudfirewall

import (
	"fmt"
	"log"
	"time"

	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	validators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCloudFirewallAiTrafficAnalysisStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAiTrafficAnalysisStatusCreate,
		Read:   resourceAiTrafficAnalysisStatusRead,
		Update: resourceAiTrafficAnalysisStatusUpdate,
		Delete: resourceAiTrafficAnalysisStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Open", "Close"}, false),
			},
		},
	}
}

func resourceAiTrafficAnalysisStatusCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UpdateAITrafficAnalysisStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("status"); ok {
		request["Status"] = convertCloudFirewallAiTrafficAnalysisStatusStatusRequest(v.(string))
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_ai_traffic_analysis_status", action, errs.AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAiTrafficAnalysisStatusRead(d, meta)
}

func resourceAiTrafficAnalysisStatusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallAiTrafficAnalysisStatus(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_ai_traffic_analysis_status DescribeCloudFirewallAiTrafficAnalysisStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("status", convertCloudFirewallAiTrafficAnalysisStatusStatusResponse(objectRaw["Status"]))

	return nil
}

func resourceAiTrafficAnalysisStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAITrafficAnalysisStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("status") {
		update = true
		request["Status"] = convertCloudFirewallAiTrafficAnalysisStatusStatusRequest(d.Get("status").(string))
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceAiTrafficAnalysisStatusRead(d, meta)
}

func resourceAiTrafficAnalysisStatusDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Ai Traffic Analysis Status. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertCloudFirewallAiTrafficAnalysisStatusStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "open":
		return "Open"
	case "close":
		return "Close"
	}
	return source
}
func convertCloudFirewallAiTrafficAnalysisStatusStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Open":
		return "open"
	case "Close":
		return "close"
	}
	return source
}
