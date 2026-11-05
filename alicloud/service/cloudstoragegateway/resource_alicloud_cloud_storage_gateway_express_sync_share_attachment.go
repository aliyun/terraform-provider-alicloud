package cloudstoragegateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentCreate,
		Read:   ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentRead,
		Delete: ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"express_sync_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"share_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	action := "AddSharesToExpressSync"
	request := make(map[string]interface{})
	var err error
	request["ExpressSyncId"] = d.Get("express_sync_id")
	sharesMaps := make(map[string]interface{})
	gatewayId := fmt.Sprint(d.Get("gateway_id"))
	shareName := fmt.Sprint(d.Get("share_name"))
	sharesMaps[gatewayId] = shareName

	if v, err := convert.ConvertMaptoJsonString(sharesMaps); err != nil {
		return errs.WrapError(err)
	} else {
		request["GatewayShares"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_express_sync_share_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["ExpressSyncId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	d.SetId(fmt.Sprint(request["ExpressSyncId"], ":", gatewayId, ":", shareName))
	return ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	object, err := sgwService.DescribeExpressSyncShares(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_express_sync sgwService.DescribeExpressSyncShares Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("express_sync_id", object["ExpressSyncId"])
	d.Set("gateway_id", object["GatewayId"])
	d.Set("share_name", object["ShareName"])
	return nil
}
func ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	action := "RemoveSharesFromExpressSync"
	request := make(map[string]interface{})
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request["ExpressSyncId"] = parts[0]
	sharesMaps := make(map[string]interface{})
	gatewayId := fmt.Sprint(parts[1])
	shareName := fmt.Sprint(parts[2])
	sharesMaps[gatewayId] = shareName

	if v, err := convert.ConvertMaptoJsonString(sharesMaps); err != nil {
		return errs.WrapError(err)
	} else {
		request["GatewayShares"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
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
		if errs.IsExpectedErrors(err, []string{"ExpressSyncNotExist", "GatewayNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(parts[0], fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
