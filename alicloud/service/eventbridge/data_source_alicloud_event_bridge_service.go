package eventbridge

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceEventBridgeService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEventBridgeServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudEventBridgeServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("EventBridgeServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateInstance"
	request := map[string]interface{}{
		"ProductCode":      "eventbridge",
		"ProductType":      "eventbridge_post_public_cn",
		"SubscriptionType": "PayAsYouGo",
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "eventbridge_post_public_intl"
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "eventbridge_post_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("EventBridgeServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_event_bridge_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	action = "GetEventBridgeStatus"
	timeout := time.Now().Add(10 * time.Minute)
	for {
		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, nil, true)
			if err != nil {
				wait()
				return retry.RetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_event_bridge_service", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Components", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Components", response)
		}
		result, _ := resp.([]interface{})
		up := true
		for _, v := range result {
			if v, ok := v.(map[string]interface{})["Status"]; !ok || fmt.Sprint(v) != "UP" {
				up = false
				break
			}
		}
		if up {
			break
		}
		if time.Now().After(timeout) {
			return errs.WrapError(fmt.Errorf("waiting for EventBridge status to be down timeout. Last response:%s", response))
		}
	}

	d.SetId("EventBridgeServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
