package eventbridge

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceEventBridgeEventStreaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeEventStreamingCreate,
		Read:   resourceAliCloudEventBridgeEventStreamingRead,
		Update: resourceAliCloudEventBridgeEventStreamingUpdate,
		Delete: resourceAliCloudEventBridgeEventStreamingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"event_streaming_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"filter_pattern": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"sink": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"run_options": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"transforms": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"RUNNING", "PAUSED"}, false),
			},
			"tags": tags.TagsSchemaWithElements(),
		},
	}
}

func resourceAliCloudEventBridgeEventStreamingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventBridgeServiceV2{client}
	var response map[string]interface{}
	action := "CreateEventStreaming"
	request := make(map[string]interface{})
	var err error

	request["EventStreamingName"] = d.Get("event_streaming_name")
	request["Source"] = d.Get("source")
	request["FilterPattern"] = d.Get("filter_pattern")
	request["Sink"] = d.Get("sink")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("run_options"); ok {
		request["RunOptions"] = v
	}
	if v, ok := d.GetOk("transforms"); ok {
		request["Transforms"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_event_bridge_event_streaming", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["EventStreamingName"]))

	// The tags passed in the CreateEventStreaming request are ignored by the
	// backend, so tags are always applied through TagResources after creation.
	if _, ok := d.GetOk("tags"); ok {
		if err := eventBridgeService.SetResourceTags(d, "eventstreaming"); err != nil {
			return errs.WrapError(err)
		}
	}

	if err := eventBridgeService.updateEventBridgeEventStreamingStatus(d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return errs.WrapError(err)
	}

	return resourceAliCloudEventBridgeEventStreamingRead(d, meta)
}

func resourceAliCloudEventBridgeEventStreamingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventBridgeServiceV2{client}

	object, err := eventBridgeService.DescribeEventBridgeEventStreaming(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_streaming eventBridgeService.DescribeEventBridgeEventStreaming Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("event_streaming_name", object["EventStreamingName"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	source, ok, err := eventBridgeEventStreamingJsonAttribute(object["Source"])
	if err != nil {
		return errs.WrapError(err)
	} else if ok {
		d.Set("source", source)
	}

	filterPattern, ok, err := eventBridgeEventStreamingJsonAttribute(object["FilterPattern"])
	if err != nil {
		return errs.WrapError(err)
	} else if ok {
		d.Set("filter_pattern", filterPattern)
	}

	sink, ok, err := eventBridgeEventStreamingJsonAttribute(object["Sink"])
	if err != nil {
		return errs.WrapError(err)
	} else if ok {
		d.Set("sink", sink)
	}

	runOptions, ok, err := eventBridgeEventStreamingJsonAttribute(object["RunOptions"])
	if err != nil {
		return errs.WrapError(err)
	} else if ok {
		d.Set("run_options", runOptions)
	}

	transforms, ok, err := eventBridgeEventStreamingJsonAttribute(object["Transforms"])
	if err != nil {
		return errs.WrapError(err)
	} else if ok {
		d.Set("transforms", transforms)
	}

	objectRaw, err := eventBridgeService.DescribeEventBridgeEventStreamingTags(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEventBridgeEventStreamingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventBridgeServiceV2{client}
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"EventStreamingName": d.Id(),
	}

	if d.HasChange("description") {
		update = true
	}
	// Send Description even when it is empty so removing description from the
	// configuration clears it on the remote EventStreaming resource.
	request["Description"] = d.Get("description")

	// Send the current Source, FilterPattern and Sink values with every update
	// to keep the resource configuration stable across partial updates.
	if d.HasChange("source") || d.HasChange("filter_pattern") || d.HasChange("sink") || d.HasChange("run_options") || d.HasChange("transforms") {
		update = true
	}
	request["Source"] = d.Get("source")
	request["FilterPattern"] = d.Get("filter_pattern")
	request["Sink"] = d.Get("sink")
	if v, ok := d.GetOk("run_options"); ok {
		request["RunOptions"] = v
	}
	if v, ok := d.GetOk("transforms"); ok {
		request["Transforms"] = v
	}

	if update {
		action := "UpdateEventStreaming"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
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

	if d.HasChange("status") {
		if err := eventBridgeService.updateEventBridgeEventStreamingStatus(d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("tags") {
		if err := eventBridgeService.SetResourceTags(d, "eventstreaming"); err != nil {
			return errs.WrapError(err)
		}
	}

	return resourceAliCloudEventBridgeEventStreamingRead(d, meta)
}

func resourceAliCloudEventBridgeEventStreamingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventBridgeServiceV2{client}
	action := "DeleteEventStreaming"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"EventStreamingName": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
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
		if errs.NotFoundError(err) || errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventStreamingNotExisted"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eventBridgeService.EventBridgeEventStreamingStateRefreshFunc(d.Id(), "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

// eventBridgeEventStreamingJsonAttribute renders an EventStreaming API attribute
// as a JSON string. The API returns FilterPattern as a string and the other
// structured attributes as objects or arrays.
func eventBridgeEventStreamingJsonAttribute(value interface{}) (string, bool, error) {
	if value == nil {
		return "", false, nil
	}

	if valueString, ok := value.(string); ok {
		return valueString, true, nil
	}

	jsonString, err := json.Marshal(value)
	if err != nil {
		return "", false, fmt.Errorf("marshal attribute of alicloud_event_bridge_event_streaming failed: %#v", err)
	}

	return string(jsonString), true, nil
}

func eventBridgeEventStreamingJsonDiffSuppress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return oldValue == newValue
	}

	var oldJsonObject interface{}
	var newJsonObject interface{}
	if err := json.Unmarshal([]byte(oldValue), &oldJsonObject); err != nil {
		return oldValue == newValue
	}
	if err := json.Unmarshal([]byte(newValue), &newJsonObject); err != nil {
		return oldValue == newValue
	}

	return reflect.DeepEqual(oldJsonObject, newJsonObject)
}
