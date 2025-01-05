package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEventBridgeService() *schema.Resource {
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "eventbridge_post_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("EventBridgeServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_service", action, AlibabaCloudSdkGoERROR)
	}
	action = "GetEventBridgeStatus"
	timeout := time.Now().Add(10 * time.Minute)
	for {
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, nil, true)
			if err != nil {
				wait()
				return resource.RetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_service", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Components", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Components", response)
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
			return WrapError(fmt.Errorf("waiting for EventBridge status to be down timeout. Last response:%s", response))
		}
	}

	d.SetId("EventBridgeServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
