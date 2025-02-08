package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudApiGatewayLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayLogConfigCreate,
		Read:   resourceAlicloudApiGatewayLogConfigRead,
		Update: resourceAlicloudApiGatewayLogConfigUpdate,
		Delete: resourceAlicloudApiGatewayLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sls_project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sls_log_store": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROVIDER"}, false),
			},
		},
	}
}

func resourceAlicloudApiGatewayLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLogConfig"
	request := make(map[string]interface{})
	var err error

	request["SlsProject"] = d.Get("sls_project")
	request["SlsLogStore"] = d.Get("sls_log_store")
	request["LogType"] = d.Get("log_type")

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_log_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LogType"]))

	return resourceAlicloudApiGatewayLogConfigRead(d, meta)
}

func resourceAlicloudApiGatewayLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	object, err := cloudApiService.DescribeApiGatewayLogConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("sls_project", object["SlsProject"])
	d.Set("sls_log_store", object["SlsLogStore"])
	d.Set("log_type", object["LogType"])

	return nil
}

func resourceAlicloudApiGatewayLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"LogType": d.Id(),
	}

	if d.HasChange("sls_project") || d.HasChange("sls_log_store") {
		update = true
	}
	request["SlsProject"] = d.Get("sls_project")
	request["SlsLogStore"] = d.Get("sls_log_store")

	if update {
		action := "ModifyLogConfig"
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

	return resourceAlicloudApiGatewayLogConfigRead(d, meta)
}

func resourceAlicloudApiGatewayLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLogConfig"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"LogType": d.Id(),
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
