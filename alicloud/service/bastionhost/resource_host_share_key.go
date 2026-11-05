package bastionhost

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

func ResourceHostShareKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostShareKeyCreate,
		Read:   resourceHostShareKeyRead,
		Update: resourceHostShareKeyUpdate,
		Delete: resourceHostShareKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"host_share_key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pass_phrase": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"host_share_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHostShareKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHostShareKey"
	request := make(map[string]interface{})
	var err error
	request["HostShareKeyName"] = d.Get("host_share_key_name")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("pass_phrase"); ok {
		request["PassPhrase"] = v
	}
	request["PrivateKey"] = d.Get("private_key")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_bastionhost_host_share_key", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostShareKeyId"]))

	return resourceHostShareKeyUpdate(d, meta)
}
func resourceHostShareKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	object, err := yundunBastionhostService.DescribeBastionhostHostShareKey(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_share_key yundunBastionhostService.DescribeBastionhostHostShareKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("host_share_key_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("host_share_key_name", object["HostShareKeyName"])
	d.Set("private_key_finger_print", object["PrivateKeyFingerPrint"])
	return nil
}
func resourceHostShareKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}
	if !d.IsNewResource() && d.HasChange("host_share_key_name") {
		update = true
		request["HostShareKeyName"] = d.Get("host_share_key_name")
	}
	if !d.IsNewResource() && d.HasChange("private_key") {
		update = true
		request["PrivateKey"] = d.Get("private_key")
	}
	request["RegionId"] = client.RegionId
	if update {
		if v, ok := d.GetOk("pass_phrase"); ok {
			request["PassPhrase"] = v
		}
		action := "ModifyHostShareKey"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
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
	}
	return resourceHostShareKeyRead(d, meta)
}
func resourceHostShareKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteHostShareKey"
	var response map[string]interface{}
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}

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
