package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiLlmTraceEval() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiLlmTraceEvalCreate,
		Read:   resourceAliCloudPaiLlmTraceEvalRead,
		Update: resourceAliCloudPaiLlmTraceEvalUpdate,
		Delete: resourceAliCloudPaiLlmTraceEvalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"eval_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gmt_create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"record_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPaiLlmTraceEvalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/api/v1/PAILLMTrace/eval")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("app_name"); ok {
		request["AppName"] = v
	}
	if v, ok := d.GetOk("data_source"); ok {
		request["DataSource"] = v
	}
	request["Description"] = d.Get("description")
	request["EvaluationName"] = d.Get("eval_name")
	if v, ok := d.GetOk("evaluation_data"); ok {
		var evalData interface{}
		if err := json.Unmarshal([]byte(v.(string)), &evalData); err == nil {
			request["EvaluationData"] = evalData
		} else {
			request["EvaluationData"] = v
		}
	}
	if v, ok := d.GetOk("metadata"); ok {
		var metadata interface{}
		if err := json.Unmarshal([]byte(v.(string)), &metadata); err == nil {
			request["Metadata"] = metadata
		} else {
			request["Metadata"] = v
		}
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("PaiLLMTrace", "2024-03-11", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_llm_trace_eval", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["EvaluationId"]))

	return resourceAliCloudPaiLlmTraceEvalRead(d, meta)
}

func resourceAliCloudPaiLlmTraceEvalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paillmtraceServiceV2 := PaillmtraceServiceV2{client}

	objectRaw, err := paillmtraceServiceV2.DescribePaillmtraceEval(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_llm_trace_eval DescribePaillmtraceEval Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("data_source", objectRaw["DataSource"])
	d.Set("description", objectRaw["Description"])
	d.Set("eval_id", d.Id())
	d.Set("eval_name", objectRaw["EvaluationName"])
	d.Set("gmt_create_time", objectRaw["GmtCreateTime"])
	d.Set("record_count", objectRaw["RecordCount"])
	d.Set("region_id", client.RegionId)

	return nil
}

func resourceAliCloudPaiLlmTraceEvalUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	var err error

	action := fmt.Sprintf("/api/v1/PAILLMTrace/eval")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["EvaluationId"] = d.Id()

	if d.HasChange("app_name") {
		update = true
	}
	if v, ok := d.GetOk("app_name"); ok || d.HasChange("app_name") {
		request["AppName"] = v
	}
	if d.HasChange("data_source") {
		update = true
	}
	if v, ok := d.GetOk("data_source"); ok || d.HasChange("data_source") {
		request["DataSource"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("eval_name") {
		update = true
	}
	request["EvaluationName"] = d.Get("eval_name")
	if d.HasChange("metadata") {
		update = true
	}
	if v, ok := d.GetOk("metadata"); ok || d.HasChange("metadata") {
		var metadata interface{}
		if err := json.Unmarshal([]byte(v.(string)), &metadata); err == nil {
			request["Metadata"] = metadata
		} else {
			request["Metadata"] = v
		}
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("PaiLLMTrace", "2024-03-11", action, query, nil, body, true)
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

	return resourceAliCloudPaiLlmTraceEvalRead(d, meta)
}

func resourceAliCloudPaiLlmTraceEvalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/api/v1/PAILLMTrace/eval")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["EvaluationId"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("PaiLLMTrace", "2024-03-11", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"NotFound", "EntityNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
