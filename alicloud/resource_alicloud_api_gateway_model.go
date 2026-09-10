package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudApiGatewayModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayModelCreate,
		Read:   resourceAlicloudApiGatewayModelRead,
		Update: resourceAlicloudApiGatewayModelUpdate,
		Delete: resourceAlicloudApiGatewayModelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudApiGatewayModelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateModel"
	request := make(map[string]interface{})
	var err error

	request["GroupId"] = d.Get("group_id")
	request["ModelName"] = d.Get("model_name")
	request["Schema"] = d.Get("schema")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_model", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["GroupId"], ":", response["ModelName"]))

	return resourceAlicloudApiGatewayModelRead(d, meta)
}

func resourceAlicloudApiGatewayModelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	object, err := cloudApiService.DescribeApiGatewayModel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_id", object["GroupId"])
	d.Set("model_name", object["ModelName"])
	d.Set("schema", object["Schema"])
	d.Set("description", object["Description"])

	return nil
}

func resourceAlicloudApiGatewayModelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId":   parts[0],
		"ModelName": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("schema") {
		update = true
	}
	if v, ok := d.GetOk("schema"); ok {
		request["Schema"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyModel"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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

	return resourceAlicloudApiGatewayModelRead(d, meta)
}

func resourceAlicloudApiGatewayModelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "DeleteModel"

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId":   parts[0],
		"ModelName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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

	return nil
}
