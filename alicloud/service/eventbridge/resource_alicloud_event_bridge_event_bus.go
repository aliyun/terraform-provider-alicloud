package eventbridge

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

func ResourceEventBridgeEventBus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeEventBusCreate,
		Read:   resourceAlicloudEventBridgeEventBusRead,
		Update: resourceAlicloudEventBridgeEventBusUpdate,
		Delete: resourceAlicloudEventBridgeEventBusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bus_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringLenBetween(0, 127),
			},
		},
	}
}

func resourceAlicloudEventBridgeEventBusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEventBus"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["EventBusName"] = d.Get("event_bus_name")
	request["ClientToken"] = helper.BuildClientToken("CreateEventBus")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_event_bridge_event_bus", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return errs.WrapError(fmt.Errorf("CreateEventBus failed, response: %v", response))
	}

	d.SetId(fmt.Sprint(request["EventBusName"]))

	return resourceAlicloudEventBridgeEventBusRead(d, meta)
}
func resourceAlicloudEventBridgeEventBusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	object, err := eventbridgeService.DescribeEventBridgeEventBus(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_bus eventbridgeService.DescribeEventBridgeEventBus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("event_bus_name", d.Id())
	d.Set("description", object["Description"])
	return nil
}
func resourceAlicloudEventBridgeEventBusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	if d.HasChange("description") {
		request := map[string]interface{}{
			"EventBusName": d.Id(),
		}
		request["Description"] = d.Get("description")
		action := "UpdateEventBus"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, false)
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
		if fmt.Sprint(response["Code"]) != "Success" {
			return errs.WrapError(fmt.Errorf("UpdateEventBus failed, response: %v", response))
		}
	}
	return resourceAlicloudEventBridgeEventBusRead(d, meta)
}
func resourceAlicloudEventBridgeEventBusDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventBus"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"EventBusName": d.Id(),
	}

	request["ClientToken"] = helper.BuildClientToken("DeleteEventBus")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EventBusNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
