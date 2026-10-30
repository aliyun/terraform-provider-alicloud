package cloudstoragegateway

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAlicloudCloudStorageGatewayExpressSync() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCloudStorageGatewayExpressSyncCreate,
		Read:   ResourceAlicloudCloudStorageGatewayExpressSyncRead,
		Delete: ResourceAlicloudCloudStorageGatewayExpressSyncDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bucket_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"express_sync_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z][\u4e00-\u9fa5a-zA-Z0-9-_.]{0,127}$"), "The length of the name is limited to `1` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, English periods (.), underscores (_), or hyphens (-), and must start with  letters."),
			},
		},
	}
}

func ResourceAlicloudCloudStorageGatewayExpressSyncCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateExpressSync"
	request := make(map[string]interface{})
	var err error
	request["BucketName"] = d.Get("bucket_name")
	if v, ok := d.GetOk("bucket_prefix"); ok {
		request["BucketPrefix"] = v
	}
	request["BucketRegion"] = d.Get("bucket_region")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Name"] = d.Get("express_sync_name")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_express_sync", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["ExpressSyncId"]))
	return ResourceAlicloudCloudStorageGatewayExpressSyncRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayExpressSyncRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	object, err := sgwService.DescribeExpressSyncs(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_express_sync sgwService.DescribeExpressSyncs Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("bucket_name", object["BucketName"])
	d.Set("bucket_prefix", object["BucketPrefix"])
	d.Set("bucket_region", object["BucketRegion"])
	d.Set("description", object["Description"])
	d.Set("express_sync_name", object["Name"])
	return nil
}

func ResourceAlicloudCloudStorageGatewayExpressSyncDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	action := "DeleteExpressSync"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ExpressSyncId": d.Id(),
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
		if errs.IsExpectedErrors(err, []string{"ExpressSyncNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
