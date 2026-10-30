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
)

func ResourceAlicloudCloudStorageGatewayStorageBundle() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCloudStorageGatewayStorageBundleCreate,
		Read:   ResourceAlicloudCloudStorageGatewayStorageBundleRead,
		Update: ResourceAlicloudCloudStorageGatewayStorageBundleUpdate,
		Delete: ResourceAlicloudCloudStorageGatewayStorageBundleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_bundle_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceAlicloudCloudStorageGatewayStorageBundleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateStorageBundle"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["BackendBucketRegionId"] = client.RegionId
	request["Name"] = d.Get("storage_bundle_name")
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
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_storage_bundle", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["StorageBundleId"]))

	return ResourceAlicloudCloudStorageGatewayStorageBundleRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayStorageBundleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	object, err := sgwService.DescribeCloudStorageGatewayStorageBundle(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_storage_bundle sgwService.DescribeCloudStorageGatewayStorageBundle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("storage_bundle_name", object["Name"])
	return nil
}
func ResourceAlicloudCloudStorageGatewayStorageBundleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"StorageBundleId": d.Id(),
	}
	if d.HasChange("storage_bundle_name") {
		update = true
	}
	request["Name"] = d.Get("storage_bundle_name")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if update {
		action := "ModifyStorageBundle"
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
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return ResourceAlicloudCloudStorageGatewayStorageBundleRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayStorageBundleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStorageBundle"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"StorageBundleId": d.Id(),
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
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"StorageBundleNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
