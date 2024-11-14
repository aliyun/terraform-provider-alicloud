// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudFcv3AsyncInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3AsyncInvokeConfigCreate,
		Read:   resourceAliCloudFcv3AsyncInvokeConfigRead,
		Update: resourceAliCloudFcv3AsyncInvokeConfigUpdate,
		Delete: resourceAliCloudFcv3AsyncInvokeConfigDelete,
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
				ValidateFunc: IntBetween(0, 86400),
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

func resourceAliCloudFcv3AsyncInvokeConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	functionName := d.Get("function_name")
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	if v, ok := d.GetOkExists("async_task"); ok {
		request["asyncTask"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("destination_config"); !IsNil(v) {
		onFailure := make(map[string]interface{})
		destination1, _ := jsonpath.Get("$[0].on_failure[0].destination", d.Get("destination_config"))
		if destination1 != nil && destination1 != "" {
			onFailure["destination"] = destination1
		}

		objectDataLocalMap["onFailure"] = onFailure
		onSuccess := make(map[string]interface{})
		destination3, _ := jsonpath.Get("$[0].on_success[0].destination", d.Get("destination_config"))
		if destination3 != nil && destination3 != "" {
			onSuccess["destination"] = destination3
		}

		objectDataLocalMap["onSuccess"] = onSuccess

		request["destinationConfig"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("max_async_event_age_in_seconds"); ok && v.(int) > 0 {
		request["maxAsyncEventAgeInSeconds"] = v
	}
	if v, ok := d.GetOkExists("max_async_retry_attempts"); ok {
		request["maxAsyncRetryAttempts"] = v
	}
	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_async_invoke_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(functionName))

	return resourceAliCloudFcv3AsyncInvokeConfigRead(d, meta)
}

func resourceAliCloudFcv3AsyncInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3AsyncInvokeConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_async_invoke_config DescribeFcv3AsyncInvokeConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["asyncTask"] != nil {
		d.Set("async_task", objectRaw["asyncTask"])
	}
	if objectRaw["createdTime"] != nil {
		d.Set("create_time", objectRaw["createdTime"])
	}
	if objectRaw["functionArn"] != nil {
		d.Set("function_arn", objectRaw["functionArn"])
	}
	if objectRaw["lastModifiedTime"] != nil {
		d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	}
	if objectRaw["maxAsyncEventAgeInSeconds"] != nil {
		d.Set("max_async_event_age_in_seconds", objectRaw["maxAsyncEventAgeInSeconds"])
	}
	if objectRaw["maxAsyncRetryAttempts"] != nil {
		d.Set("max_async_retry_attempts", objectRaw["maxAsyncRetryAttempts"])
	}

	destinationConfigMaps := make([]map[string]interface{}, 0)
	destinationConfigMap := make(map[string]interface{})
	destinationConfig1Raw := make(map[string]interface{})
	if objectRaw["destinationConfig"] != nil {
		destinationConfig1Raw = objectRaw["destinationConfig"].(map[string]interface{})
	}
	if len(destinationConfig1Raw) > 0 {

		onFailureMaps := make([]map[string]interface{}, 0)
		onFailureMap := make(map[string]interface{})
		onFailure1Raw := make(map[string]interface{})
		if destinationConfig1Raw["onFailure"] != nil {
			onFailure1Raw = destinationConfig1Raw["onFailure"].(map[string]interface{})
		}
		if len(onFailure1Raw) > 0 {
			onFailureMap["destination"] = onFailure1Raw["destination"]

			onFailureMaps = append(onFailureMaps, onFailureMap)
		}
		destinationConfigMap["on_failure"] = onFailureMaps
		onSuccessMaps := make([]map[string]interface{}, 0)
		onSuccessMap := make(map[string]interface{})
		onSuccess1Raw := make(map[string]interface{})
		if destinationConfig1Raw["onSuccess"] != nil {
			onSuccess1Raw = destinationConfig1Raw["onSuccess"].(map[string]interface{})
		}
		if len(onSuccess1Raw) > 0 {
			onSuccessMap["destination"] = onSuccess1Raw["destination"]

			onSuccessMaps = append(onSuccessMaps, onSuccessMap)
		}
		destinationConfigMap["on_success"] = onSuccessMaps
		destinationConfigMaps = append(destinationConfigMaps, destinationConfigMap)
	}
	if objectRaw["destinationConfig"] != nil {
		if err := d.Set("destination_config", destinationConfigMaps); err != nil {
			return err
		}
	}

	d.Set("function_name", d.Id())

	return nil
}

func resourceAliCloudFcv3AsyncInvokeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if d.HasChange("async_task") {
		update = true
		request["asyncTask"] = d.Get("async_task")
	}

	if d.HasChange("destination_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("destination_config"); !IsNil(v) {
			onFailure := make(map[string]interface{})
			destination1, _ := jsonpath.Get("$[0].on_failure[0].destination", v)
			if destination1 != nil && (d.HasChange("destination_config.0.on_failure.0.destination") || destination1 != "") {
				onFailure["destination"] = destination1
			}

			objectDataLocalMap["onFailure"] = onFailure
			onSuccess := make(map[string]interface{})
			destination3, _ := jsonpath.Get("$[0].on_success[0].destination", v)
			if destination3 != nil && (d.HasChange("destination_config.0.on_success.0.destination") || destination3 != "") {
				onSuccess["destination"] = destination3
			}

			objectDataLocalMap["onSuccess"] = onSuccess

			request["destinationConfig"] = objectDataLocalMap
		}
	}

	if d.HasChange("max_async_event_age_in_seconds") {
		update = true
		request["maxAsyncEventAgeInSeconds"] = d.Get("max_async_event_age_in_seconds")
	}

	if d.HasChange("max_async_retry_attempts") {
		update = true
		request["maxAsyncRetryAttempts"] = d.Get("max_async_retry_attempts")
	}

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	return resourceAliCloudFcv3AsyncInvokeConfigRead(d, meta)
}

func resourceAliCloudFcv3AsyncInvokeConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
