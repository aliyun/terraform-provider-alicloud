package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeEventStreaming() *schema.Resource {
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
				ValidateFunc:     validateEventBridgeEventStreamingJson,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"filter_pattern": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateEventBridgeEventStreamingJson,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"sink": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateEventBridgeEventStreamingJson,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"run_options": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validateEventBridgeEventStreamingJson,
				DiffSuppressFunc: eventBridgeEventStreamingJsonDiffSuppress,
			},
			"transforms": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validateEventBridgeEventStreamingJson,
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
				ValidateFunc: StringInSlice([]string{"RUNNING", "PAUSED"}, false),
			},
			"tags": tagsSchema(),
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_event_streaming", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["EventStreamingName"]))

	// The tags passed in the CreateEventStreaming request are ignored by the
	// backend, so tags are always applied through TagResources after creation.
	if _, ok := d.GetOk("tags"); ok {
		if err := eventBridgeService.SetResourceTags(d, "eventstreaming"); err != nil {
			return WrapError(err)
		}
	}

	if err := eventBridgeService.updateEventBridgeEventStreamingStatus(d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudEventBridgeEventStreamingRead(d, meta)
}

func resourceAliCloudEventBridgeEventStreamingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventBridgeServiceV2{client}

	object, err := eventBridgeService.DescribeEventBridgeEventStreaming(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_streaming eventBridgeService.DescribeEventBridgeEventStreaming Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_streaming_name", object["EventStreamingName"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	source, ok, err := eventBridgeEventStreamingJsonAttribute(object["Source"])
	if err != nil {
		return WrapError(err)
	} else if ok {
		d.Set("source", source)
	}

	filterPattern, ok, err := eventBridgeEventStreamingJsonAttribute(object["FilterPattern"])
	if err != nil {
		return WrapError(err)
	} else if ok {
		d.Set("filter_pattern", filterPattern)
	}

	sink, ok, err := eventBridgeEventStreamingJsonAttribute(object["Sink"])
	if err != nil {
		return WrapError(err)
	} else if ok {
		d.Set("sink", sink)
	}

	runOptions, ok, err := eventBridgeEventStreamingJsonAttribute(object["RunOptions"])
	if err != nil {
		return WrapError(err)
	} else if ok {
		d.Set("run_options", runOptions)
	}

	transforms, ok, err := eventBridgeEventStreamingJsonAttribute(object["Transforms"])
	if err != nil {
		return WrapError(err)
	} else if ok {
		d.Set("transforms", transforms)
	}

	objectRaw, err := eventBridgeService.DescribeEventBridgeEventStreamingTags(d.Id())
	if err != nil {
		return WrapError(err)
	}
	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

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
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	// Source, FilterPattern and Sink are required by UpdateEventStreaming, so
	// their current values are always sent together with the optional fields.
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	if d.HasChange("status") {
		if err := eventBridgeService.updateEventBridgeEventStreamingStatus(d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("tags") {
		if err := eventBridgeService.SetResourceTags(d, "eventstreaming"); err != nil {
			return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if NotFoundError(err) || IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventStreamingNotExisted"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eventBridgeService.EventBridgeEventStreamingStateRefreshFunc(d.Id(), "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
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

func validateEventBridgeEventStreamingJson(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value == "" {
		return
	}
	if !json.Valid([]byte(value)) {
		errors = append(errors, fmt.Errorf("%q must be a valid JSON string, got %q", k, value))
	}
	return
}
