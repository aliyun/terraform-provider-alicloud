package bastionhost

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

func ResourceHostAccountShareKeyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostAccountShareKeyAttachmentCreate,
		Read:   resourceHostAccountShareKeyAttachmentRead,
		Delete: resourceHostAccountShareKeyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"host_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_share_key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceHostAccountShareKeyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachHostAccountsToHostShareKey"
	request := make(map[string]interface{})
	var err error

	request["HostAccountIds"] = convert.ConvertListToJsonString([]interface{}{d.Get("host_account_id")})
	request["HostShareKeyId"] = d.Get("host_share_key_id")
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_bastionhost_host_account_share_key_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["HostShareKeyId"], ":", d.Get("host_account_id")))

	return resourceHostAccountShareKeyAttachmentRead(d, meta)
}
func resourceHostAccountShareKeyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	_, err := yundunBastionhostService.DescribeBastionhostHostAccountShareKeyAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_account_share_key_attachment yundunBastionhostService.DescribeBastionhostHostAccountShareKeyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("host_account_id", parts[2])
	d.Set("host_share_key_id", parts[1])
	d.Set("instance_id", parts[0])
	return nil
}
func resourceHostAccountShareKeyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DetachHostAccountsFromHostShareKey"
	var response map[string]interface{}
	request := map[string]interface{}{
		"HostAccountIds": convert.ConvertListToJsonString([]interface{}{parts[2]}),
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
