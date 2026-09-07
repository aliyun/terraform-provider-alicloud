package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudApiGatewayBackendModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayBackendModelCreate,
		Read:   resourceAlicloudApiGatewayBackendModelRead,
		Update: resourceAlicloudApiGatewayBackendModelUpdate,
		Delete: resourceAlicloudApiGatewayBackendModelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backend_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backend_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "VPC", "FC_EVENT", "FC_EVENT_V3", "FC_HTTP", "FC_HTTP_V3", "OSS", "MOCK", "EVENTBRIDGE"}, false),
			},
			"stage_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backend_model_data": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == new {
						return true
					}
					var oldMap, newMap map[string]interface{}
					if err := json.Unmarshal([]byte(old), &oldMap); err != nil {
						return false
					}
					if err := json.Unmarshal([]byte(new), &newMap); err != nil {
						return false
					}
					return reflect.DeepEqual(oldMap, newMap)
				},
			},
			"backend_model_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudApiGatewayBackendModelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	request["BackendId"] = d.Get("backend_id").(string)
	request["BackendType"] = d.Get("backend_type").(string)
	request["StageName"] = d.Get("stage_name").(string)
	request["BackendModelData"] = d.Get("backend_model_data").(string)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	var response map[string]interface{}
	action := "CreateBackendModel"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_backend_model", action, AlibabaCloudSdkGoERROR)
	}

	backendId := d.Get("backend_id").(string)
	stageName := d.Get("stage_name").(string)
	d.SetId(fmt.Sprintf("%s%s%s", backendId, COLON_SEPARATED, stageName))

	return resourceAlicloudApiGatewayBackendModelRead(d, meta)
}

func resourceAlicloudApiGatewayBackendModelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	object, err := cloudApiService.DescribeApiGatewayBackendModel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_backend_model cloudApiService.DescribeApiGatewayBackendModel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backend_id", object["BackendId"])
	d.Set("backend_type", object["BackendType"])
	d.Set("stage_name", object["StageName"])
	d.Set("description", object["Description"])
	d.Set("backend_model_data", object["BackendModelData"])
	d.Set("backend_model_id", object["BackendModelId"])

	return nil
}

func resourceAlicloudApiGatewayBackendModelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	backendId := parts[0]
	stageName := parts[1]

	request := map[string]interface{}{
		"BackendId":      backendId,
		"StageName":      stageName,
		"BackendType":    d.Get("backend_type").(string),
		"BackendModelId": d.Get("backend_model_id").(string),
	}

	if d.HasChange("backend_model_data") {
		update = true
	}
	request["BackendModelData"] = d.Get("backend_model_data").(string)

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		action := "ModifyBackendModel"
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

	return resourceAlicloudApiGatewayBackendModelRead(d, meta)
}

func resourceAlicloudApiGatewayBackendModelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	backendId := parts[0]
	stageName := parts[1]

	request := map[string]interface{}{
		"BackendId":      backendId,
		"StageName":      stageName,
		"BackendModelId": d.Get("backend_model_id").(string),
	}

	action := "DeleteBackendModel"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"NotFoundBackend", "NotFoundBackendModel"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
