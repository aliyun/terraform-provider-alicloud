package eventbridge

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEventBridgeEventSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeEventSourceCreate,
		Read:   resourceAliCloudEventBridgeEventSourceRead,
		Update: resourceAliCloudEventBridgeEventSourceUpdate,
		Delete: resourceAliCloudEventBridgeEventSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.269.0 and will be removed in the future. Please use 'alicloud_event_bridge_event_source_v2' instead.",
		Schema: map[string]*schema.Schema{
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// lintignore: S006
			"external_source_config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"linked_external_source": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEventBridgeEventSourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEventSource"
	request := make(map[string]interface{})
	var err error

	request["EventBusName"] = d.Get("event_bus_name")
	request["EventSourceName"] = d.Get("event_source_name")

	if v, ok := d.GetOk("external_source_type"); ok {
		request["ExternalSourceType"] = v
	}

	if v, ok := d.GetOk("external_source_config"); ok {
		externalSourceConfigJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}

		request["ExternalSourceConfig"] = externalSourceConfigJson
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("linked_external_source"); ok {
		request["LinkedExternalSource"] = v
	}

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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_event_bridge_event_source", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["EventSourceName"]))

	return resourceAliCloudEventBridgeEventSourceRead(d, meta)
}

func resourceAliCloudEventBridgeEventSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventbridgeService{client}

	object, err := eventBridgeService.DescribeEventBridgeEventSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_source eventbridgeService.DescribeEventBridgeEventSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("event_bus_name", object["EventBusName"])
	d.Set("event_source_name", object["Name"])
	d.Set("external_source_type", object["ExternalSourceType"])
	d.Set("external_source_config", object["ExternalSourceConfig"])
	d.Set("description", object["Description"])

	return nil
}

func resourceAliCloudEventBridgeEventSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}

	update := false

	request := map[string]interface{}{
		"EventBusName":    d.Get("event_bus_name"),
		"EventSourceName": d.Id(),
	}

	if d.HasChange("external_source_type") {
		update = true
	}
	if v, ok := d.GetOk("external_source_type"); ok {
		request["ExternalSourceType"] = v
	}

	if d.HasChange("external_source_config") {
		update = true
	}
	if v, ok := d.GetOk("external_source_config"); ok {
		externalSourceConfigJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}

		request["ExternalSourceConfig"] = externalSourceConfigJson
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("linked_external_source") {
		update = true

		if v, ok := d.GetOkExists("linked_external_source"); ok {
			request["LinkedExternalSource"] = v
		}
	}

	if update {
		action := "UpdateEventSource"

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Code"]) != "Success" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	return resourceAliCloudEventBridgeEventSourceRead(d, meta)
}

func resourceAliCloudEventBridgeEventSourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventSource"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"EventSourceName": d.Id(),
	}

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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
