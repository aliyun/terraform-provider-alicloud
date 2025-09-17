// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudFcv3Trigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3TriggerCreate,
		Read:   resourceAliCloudFcv3TriggerRead,
		Update: resourceAliCloudFcv3TriggerUpdate,
		Delete: resourceAliCloudFcv3TriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"http_trigger": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"invocation_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_config": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"trigger_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"trigger_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudFcv3TriggerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	functionName := d.Get("function_name")
	action := fmt.Sprintf("/2023-03-30/functions/%s/triggers", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("trigger_name"); ok {
		request["triggerName"] = v
	}

	if v, ok := d.GetOk("invocation_role"); ok {
		request["invocationRole"] = v
	}
	request["triggerConfig"] = d.Get("trigger_config")
	request["triggerType"] = d.Get("trigger_type")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("source_arn"); ok {
		request["sourceArn"] = v
	}
	request["qualifier"] = d.Get("qualifier")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_trigger", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", functionName, response["triggerName"]))

	fcv3ServiceV2 := Fcv3ServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 1, fcv3ServiceV2.Fcv3TriggerStateRefreshFunc(d.Id(), "#triggerId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudFcv3TriggerRead(d, meta)
}

func resourceAliCloudFcv3TriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3Trigger(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_trigger DescribeFcv3Trigger Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createdTime"])
	d.Set("description", objectRaw["description"])
	d.Set("invocation_role", objectRaw["invocationRole"])
	d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	d.Set("qualifier", objectRaw["qualifier"])
	d.Set("source_arn", objectRaw["sourceArn"])
	d.Set("status", objectRaw["status"])
	d.Set("target_arn", objectRaw["targetArn"])
	d.Set("trigger_config", objectRaw["triggerConfig"])
	d.Set("trigger_id", objectRaw["triggerId"])
	d.Set("trigger_type", objectRaw["triggerType"])
	d.Set("trigger_name", objectRaw["triggerName"])

	httpTriggerMaps := make([]map[string]interface{}, 0)
	httpTriggerMap := make(map[string]interface{})
	httpTriggerRaw := make(map[string]interface{})
	if objectRaw["httpTrigger"] != nil {
		httpTriggerRaw = objectRaw["httpTrigger"].(map[string]interface{})
	}
	if len(httpTriggerRaw) > 0 {
		httpTriggerMap["url_internet"] = httpTriggerRaw["urlInternet"]
		httpTriggerMap["url_intranet"] = httpTriggerRaw["urlIntranet"]

		httpTriggerMaps = append(httpTriggerMaps, httpTriggerMap)
	}
	if objectRaw["httpTrigger"] != nil {
		if err := d.Set("http_trigger", httpTriggerMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("function_name", parts[0])

	return nil
}

func resourceAliCloudFcv3TriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	functionName := parts[0]
	triggerName := parts[1]
	action := fmt.Sprintf("/2023-03-30/functions/%s/triggers/%s", functionName, triggerName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("invocation_role") {
		update = true
		request["invocationRole"] = d.Get("invocation_role")
	}

	if d.HasChange("trigger_config") {
		update = true
		request["triggerConfig"] = d.Get("trigger_config")
	}

	if d.HasChange("description") {
		update = true
		request["description"] = d.Get("description")
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
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

	return resourceAliCloudFcv3TriggerRead(d, meta)
}

func resourceAliCloudFcv3TriggerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	functionName := parts[0]
	triggerName := parts[1]
	action := fmt.Sprintf("/2023-03-30/functions/%s/triggers/%s", functionName, triggerName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)

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
		if IsExpectedErrors(err, []string{"TriggerNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
