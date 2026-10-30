package cloudstoragegateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAlicloudCloudStorageGatewayGatewayLogging() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCloudStorageGatewayGatewayLoggingCreate,
		Read:   ResourceAlicloudCloudStorageGatewayGatewayLoggingRead,
		Update: ResourceAlicloudCloudStorageGatewayGatewayLoggingUpdate,
		Delete: ResourceAlicloudCloudStorageGatewayGatewayLoggingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
		},
	}
}

func ResourceAlicloudCloudStorageGatewayGatewayLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGatewayLogging"
	request := make(map[string]interface{})
	var err error
	request["GatewayId"] = d.Get("gateway_id")
	request["SlsLogstore"] = d.Get("sls_logstore")
	request["SlsProject"] = d.Get("sls_project")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_logging", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["GatewayId"]))

	return ResourceAlicloudCloudStorageGatewayGatewayLoggingUpdate(d, meta)
}
func ResourceAlicloudCloudStorageGatewayGatewayLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	object, err := sgwService.DescribeCloudStorageGatewayGatewayLogging(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_logging sgwService.DescribeCloudStorageGatewayGatewayLogging Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("gateway_id", d.Id())
	d.Set("sls_logstore", object["SlsLogstore"])
	d.Set("sls_project", object["SlsProject"])
	d.Set("status", object["GatewayLoggingStatus"])
	return nil
}
func ResourceAlicloudCloudStorageGatewayGatewayLoggingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := sgwService.DescribeCloudStorageGatewayGatewayLogging(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["GatewayLoggingStatus"].(string) != target {
			if target == "Disabled" {
				request := map[string]interface{}{
					"GatewayId": d.Id(),
				}
				action := "DisableGatewayLogging"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {

					response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
				if fmt.Sprint(response["Success"]) == "false" {
					return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}
			}
			if target == "Enabled" {
				request := map[string]interface{}{
					"GatewayId": d.Id(),
				}
				action := "EnableGatewayLogging"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {

					response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
				if fmt.Sprint(response["Success"]) == "false" {
					return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}
			}
		}
	}
	d.Partial(false)
	return ResourceAlicloudCloudStorageGatewayGatewayLoggingRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayGatewayLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGatewayLogging"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
