// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package fcv3

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAsyncInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAsyncInvokeConfigCreate,
		Read:   resourceAsyncInvokeConfigRead,
		Update: resourceAsyncInvokeConfigUpdate,
		Delete: resourceAsyncInvokeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"async_task": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_success": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"on_failure": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"function_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_async_event_age_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 86400),
			},
			"max_async_retry_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAsyncInvokeConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	functionName := d.Get("function_name")
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	if v, ok := d.GetOkExists("async_task"); ok {
		request["asyncTask"] = v
	}
	destinationConfig := make(map[string]interface{})

	if v := d.Get("destination_config"); !helper.IsNil(v) {
		onFailure := make(map[string]interface{})
		destination1, _ := jsonpath.Get("$[0].on_failure[0].destination", d.Get("destination_config"))
		if destination1 != nil && destination1 != "" {
			onFailure["destination"] = destination1
		}

		if len(onFailure) > 0 {
			destinationConfig["onFailure"] = onFailure
		}
		onSuccess := make(map[string]interface{})
		destination3, _ := jsonpath.Get("$[0].on_success[0].destination", d.Get("destination_config"))
		if destination3 != nil && destination3 != "" {
			onSuccess["destination"] = destination3
		}

		if len(onSuccess) > 0 {
			destinationConfig["onSuccess"] = onSuccess
		}

		request["destinationConfig"] = destinationConfig
	}

	if v, ok := d.GetOkExists("max_async_event_age_in_seconds"); ok && v.(int) > 0 {
		request["maxAsyncEventAgeInSeconds"] = v
	}
	if v, ok := d.GetOkExists("max_async_retry_attempts"); ok {
		request["maxAsyncRetryAttempts"] = v
	}
	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fcv3_async_invoke_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(functionName))

	return resourceAsyncInvokeConfigRead(d, meta)
}

func resourceAsyncInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeAsyncInvokeConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_async_invoke_config DescribeAsyncInvokeConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("async_task", objectRaw["asyncTask"])
	d.Set("create_time", objectRaw["createdTime"])
	d.Set("function_arn", objectRaw["functionArn"])
	d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	d.Set("max_async_event_age_in_seconds", objectRaw["maxAsyncEventAgeInSeconds"])
	d.Set("max_async_retry_attempts", objectRaw["maxAsyncRetryAttempts"])

	destinationConfigMaps := make([]map[string]interface{}, 0)
	destinationConfigMap := make(map[string]interface{})
	destinationConfigRaw := make(map[string]interface{})
	if objectRaw["destinationConfig"] != nil {
		destinationConfigRaw = objectRaw["destinationConfig"].(map[string]interface{})
	}
	if len(destinationConfigRaw) > 0 {

		onFailureMaps := make([]map[string]interface{}, 0)
		onFailureMap := make(map[string]interface{})
		onFailureRaw := make(map[string]interface{})
		if destinationConfigRaw["onFailure"] != nil {
			onFailureRaw = destinationConfigRaw["onFailure"].(map[string]interface{})
		}
		if len(onFailureRaw) > 0 {
			onFailureMap["destination"] = onFailureRaw["destination"]

			onFailureMaps = append(onFailureMaps, onFailureMap)
		}
		destinationConfigMap["on_failure"] = onFailureMaps
		onSuccessMaps := make([]map[string]interface{}, 0)
		onSuccessMap := make(map[string]interface{})
		onSuccessRaw := make(map[string]interface{})
		if destinationConfigRaw["onSuccess"] != nil {
			onSuccessRaw = destinationConfigRaw["onSuccess"].(map[string]interface{})
		}
		if len(onSuccessRaw) > 0 {
			onSuccessMap["destination"] = onSuccessRaw["destination"]

			onSuccessMaps = append(onSuccessMaps, onSuccessMap)
		}
		destinationConfigMap["on_success"] = onSuccessMaps
		destinationConfigMaps = append(destinationConfigMaps, destinationConfigMap)
	}
	if err := d.Set("destination_config", destinationConfigMaps); err != nil {
		return err
	}

	d.Set("function_name", d.Id())

	return nil
}

func resourceAsyncInvokeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if d.HasChange("async_task") {
		update = true

		if v, ok := d.GetOkExists("async_task"); ok {
			request["asyncTask"] = v
		}
	}

	if d.HasChange("destination_config") {
		update = true
	}
	destinationConfig := make(map[string]interface{})

	if v := d.Get("destination_config"); !helper.IsNil(v) || d.HasChange("destination_config") {
		onFailure := make(map[string]interface{})
		destination1, _ := jsonpath.Get("$[0].on_failure[0].destination", d.Get("destination_config"))
		if destination1 != nil && destination1 != "" {
			onFailure["destination"] = destination1
		}

		if len(onFailure) > 0 {
			destinationConfig["onFailure"] = onFailure
		}

		onSuccess := make(map[string]interface{})
		destination3, _ := jsonpath.Get("$[0].on_success[0].destination", d.Get("destination_config"))
		if destination3 != nil && destination3 != "" {
			onSuccess["destination"] = destination3
		}

		if len(onSuccess) > 0 {
			destinationConfig["onSuccess"] = onSuccess
		}

		request["destinationConfig"] = destinationConfig
	}

	if d.HasChange("max_async_event_age_in_seconds") {
		update = true

		if v, ok := d.GetOkExists("max_async_event_age_in_seconds"); ok && v.(int) > 0 {
			request["maxAsyncEventAgeInSeconds"] = v
		}
	}

	if d.HasChange("max_async_retry_attempts") {
		update = true

		if v, ok := d.GetOkExists("max_async_retry_attempts"); ok {
			request["maxAsyncRetryAttempts"] = v
		}
	}

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
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

	return resourceAsyncInvokeConfigRead(d, meta)
}

func resourceAsyncInvokeConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
