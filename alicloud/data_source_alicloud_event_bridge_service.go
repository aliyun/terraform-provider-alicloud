package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
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
	var response map[string]interface{}
	action := "CreateInstance"
	request := map[string]interface{}{
		"ProductCode":      "eventbridge",
		"SubscriptionType": "PayAsYouGo",
	}
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
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
	conn, err = client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	timeout := time.Now().Add(10 * time.Minute)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-04-01"), StringPointer("AK"), nil, nil, &runtime)
			if err != nil {
				wait()
				return resource.RetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_rules", action, AlibabaCloudSdkGoERROR)
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
