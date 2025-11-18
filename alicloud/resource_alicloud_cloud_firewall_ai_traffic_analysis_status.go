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

func resourceAliCloudCloudFirewallAiTrafficAnalysisStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallAiTrafficAnalysisStatusCreate,
		Read:   resourceAliCloudCloudFirewallAiTrafficAnalysisStatusRead,
		Update: resourceAliCloudCloudFirewallAiTrafficAnalysisStatusUpdate,
		Delete: resourceAliCloudCloudFirewallAiTrafficAnalysisStatusDelete,
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
				ValidateFunc: StringInSlice([]string{"Open", "Close"}, false),
			},
		},
	}
}

func resourceAliCloudCloudFirewallAiTrafficAnalysisStatusCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_ai_traffic_analysis_status", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudCloudFirewallAiTrafficAnalysisStatusRead(d, meta)
}

func resourceAliCloudCloudFirewallAiTrafficAnalysisStatusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallAiTrafficAnalysisStatus(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_ai_traffic_analysis_status DescribeCloudFirewallAiTrafficAnalysisStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", convertCloudFirewallAiTrafficAnalysisStatusStatusResponse(objectRaw["Status"]))

	return nil
}

func resourceAliCloudCloudFirewallAiTrafficAnalysisStatusUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceAliCloudCloudFirewallAiTrafficAnalysisStatusRead(d, meta)
}

func resourceAliCloudCloudFirewallAiTrafficAnalysisStatusDelete(d *schema.ResourceData, meta interface{}) error {
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
