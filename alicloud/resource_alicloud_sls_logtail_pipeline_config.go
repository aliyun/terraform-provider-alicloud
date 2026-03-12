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

func resourceAliCloudSlsLogtailPipelineConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsLogtailPipelineConfigCreate,
		Read:   resourceAliCloudSlsLogtailPipelineConfigRead,
		Update: resourceAliCloudSlsLogtailPipelineConfigUpdate,
		Delete: resourceAliCloudSlsLogtailPipelineConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregators": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"config_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flushers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"globals": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"inputs": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"log_sample": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"processors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudSlsLogtailPipelineConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/pipelineconfigs")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project").(string))
	if v, ok := d.GetOk("config_name"); ok {
		request["configName"] = v
	}

	if v, ok := d.GetOk("task"); ok {
		request["task"] = v
	}
	if v, ok := d.GetOk("inputs"); ok {
		inputsMapsArray := convertToInterfaceArray(v)

		request["inputs"] = inputsMapsArray
	}

	if v, ok := d.GetOk("globals"); ok {
		request["global"] = v
	}
	if v, ok := d.GetOk("aggregators"); ok {
		aggregatorsMapsArray := convertToInterfaceArray(v)

		request["aggregators"] = aggregatorsMapsArray
	}

	if v, ok := d.GetOk("log_sample"); ok {
		request["logSample"] = v
	}
	if v, ok := d.GetOk("flushers"); ok {
		flushersMapsArray := convertToInterfaceArray(v)

		request["flushers"] = flushersMapsArray
	}

	if v, ok := d.GetOk("processors"); ok {
		processorsMapsArray := convertToInterfaceArray(v)

		request["processors"] = processorsMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateLogtailPipelineConfig", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_logtail_pipeline_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["configName"]))

	return resourceAliCloudSlsLogtailPipelineConfigRead(d, meta)
}

func resourceAliCloudSlsLogtailPipelineConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsLogtailPipelineConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_logtail_pipeline_config DescribeSlsLogtailPipelineConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("globals", objectRaw["global"])
	d.Set("log_sample", objectRaw["logSample"])
	d.Set("task", objectRaw["task"])

	aggregatorsRaw := objectRaw["aggregators"]
	if err := d.Set("aggregators", aggregatorsRaw); err != nil {
		return err
	}
	flushersRaw := objectRaw["flushers"]
	if err := d.Set("flushers", flushersRaw); err != nil {
		return err
	}
	inputsRaw := objectRaw["inputs"]
	if err := d.Set("inputs", inputsRaw); err != nil {
		return err
	}
	processorsRaw := objectRaw["processors"]
	if err := d.Set("processors", processorsRaw); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("config_name", parts[1])

	return nil
}

func resourceAliCloudSlsLogtailPipelineConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	configName := parts[1]
	action := fmt.Sprintf("/pipelineconfigs/%s", configName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if d.HasChange("task") {
		update = true
	}
	if v, ok := d.GetOk("task"); ok || d.HasChange("task") {
		request["task"] = v
	}
	if d.HasChange("inputs") {
		update = true
	}
	if v, ok := d.GetOk("inputs"); ok || d.HasChange("inputs") {
		inputsMapsArray := convertToInterfaceArray(v)

		request["inputs"] = inputsMapsArray
	}

	if d.HasChange("globals") {
		update = true
	}
	if v, ok := d.GetOk("globals"); ok || d.HasChange("globals") {
		request["global"] = v
	}
	request["configName"] = d.Get("logstore_name")
	if d.HasChange("aggregators") {
		update = true
	}
	if v, ok := d.GetOk("aggregators"); ok || d.HasChange("aggregators") {
		aggregatorsMapsArray := convertToInterfaceArray(v)

		request["aggregators"] = aggregatorsMapsArray
	}

	if d.HasChange("log_sample") {
		update = true
	}
	if v, ok := d.GetOk("log_sample"); ok || d.HasChange("log_sample") {
		request["logSample"] = v
	}
	if d.HasChange("flushers") {
		update = true
	}
	if v, ok := d.GetOk("flushers"); ok || d.HasChange("flushers") {
		flushersMapsArray := convertToInterfaceArray(v)

		request["flushers"] = flushersMapsArray
	}

	if d.HasChange("processors") {
		update = true
	}
	if v, ok := d.GetOk("processors"); ok || d.HasChange("processors") {
		processorsMapsArray := convertToInterfaceArray(v)

		request["processors"] = processorsMapsArray
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateLogtailPipelineConfig", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudSlsLogtailPipelineConfigRead(d, meta)
}

func resourceAliCloudSlsLogtailPipelineConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	configName := parts[1]
	action := fmt.Sprintf("/pipelineconfigs/%s", configName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteLogtailPipelineConfig", action), query, nil, nil, hostMap, false)
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
