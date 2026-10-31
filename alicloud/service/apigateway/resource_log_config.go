package apigateway

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: ResourceLogConfigCreate,
		Read:   ResourceLogConfigRead,
		Update: ResourceLogConfigUpdate,
		Delete: ResourceLogConfigDelete,
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

func ResourceLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLogConfig"
	request := make(map[string]interface{})
	var err error

	request["SlsProject"] = d.Get("sls_project")
	request["SlsLogStore"] = d.Get("sls_log_store")
	request["LogType"] = d.Get("log_type")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_log_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LogType"]))

	return ResourceLogConfigRead(d, meta)
}

func ResourceLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	object, err := cloudApiService.DescribeApiGatewayLogConfig(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("sls_project", object["SlsProject"])
	d.Set("sls_log_store", object["SlsLogStore"])
	d.Set("log_type", object["LogType"])

	return nil
}

func ResourceLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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

	return ResourceLogConfigRead(d, meta)
}

func ResourceLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLogConfig"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"LogType": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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

	return nil
}
