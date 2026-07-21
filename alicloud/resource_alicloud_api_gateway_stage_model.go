package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudApiGatewayStageModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayStageModelCreate,
		Read:   resourceAlicloudApiGatewayStageModelRead,
		Update: resourceAlicloudApiGatewayStageModelUpdate,
		Delete: resourceAlicloudApiGatewayStageModelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"stage_model_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[A-Z0-9]{2,10}$"), "must be 2-10 uppercase letters or digits"),
					validation.StringNotInSlice([]string{"RELEASE", "PRE", "TEST"}, false),
				),
			},
			"stage_model_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
			},
			"stage_model_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudApiGatewayStageModelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})

	request["StageName"] = d.Get("stage_model_name").(string)
	request["StageAlias"] = d.Get("stage_model_alias").(string)
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	var response map[string]interface{}
	action := "CreateStageModel"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_stage_model", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.StageModelId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_api_gateway_stage_model")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudApiGatewayStageModelRead(d, meta)
}

func resourceAlicloudApiGatewayStageModelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	object, err := cloudApiService.DescribeApiGatewayStageModel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_stage_model cloudApiService.DescribeApiGatewayStageModel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("stage_model_name", object["StageName"])
	d.Set("stage_model_alias", object["StageAlias"])
	d.Set("description", object["Description"])
	d.Set("stage_model_id", object["StageModelId"])
	d.Set("type", object["Type"])
	d.Set("created_time", object["CreatedTime"])
	d.Set("modified_time", object["ModifiedTime"])
	return nil
}

func resourceAlicloudApiGatewayStageModelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	update := false

	request := map[string]interface{}{
		"StageModelId": d.Id(),
	}

	if d.HasChange("stage_model_alias") {
		update = true
		request["StageAlias"] = d.Get("stage_model_alias").(string)
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		action := "ModifyStageModel"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudApiGatewayStageModelRead(d, meta)
}

func resourceAlicloudApiGatewayStageModelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"StageModelId": d.Id(),
	}

	action := "DeleteStageModel"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundStageModel"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
