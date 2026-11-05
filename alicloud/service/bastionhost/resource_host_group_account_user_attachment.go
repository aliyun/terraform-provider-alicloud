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

func ResourceHostGroupAccountUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostGroupAccountUserAttachmentCreate,
		Read:   resourceHostGroupAccountUserAttachmentRead,
		Update: resourceHostGroupAccountUserAttachmentUpdate,
		Delete: resourceHostGroupAccountUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"host_account_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceHostGroupAccountUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprint(d.Get("instance_id"), ":", d.Get("user_id"), ":", d.Get("host_group_id")))

	return resourceHostGroupAccountUserAttachmentUpdate(d, meta)
}
func resourceHostGroupAccountUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	object, err := yundunBastionhostService.DescribeBastionhostHostGroupAccountUserAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_group_account_user_attachment yundunBastionhostService.DescribeBastionhostHostGroupAccountUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("host_group_id", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("user_id", parts[1])
	d.Set("host_account_names", object)
	return nil
}
func resourceHostGroupAccountUserAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("host_account_names") {
		parts, err := helper.ParseResourceId(d.Id(), 3)
		if err != nil {
			return errs.WrapError(err)
		}
		action := "AttachHostGroupAccountsToUser"
		request := make(map[string]interface{})

		oraw, nraw := d.GetChange("host_account_names")
		request["InstanceId"] = parts[0]
		request["RegionId"] = client.RegionId
		request["UserId"] = parts[1]

		if oraw != nil && len(oraw.(*schema.Set).List()) > 0 {
			action = "DetachHostGroupAccountsFromUser"
			hostRequestMaps := make([]map[string]interface{}, 0)
			hostRequestMap := make(map[string]interface{}, 0)
			hostRequestMap["HostGroupId"] = parts[2]
			hostRequestMap["HostAccountNames"] = oraw.(*schema.Set).List()
			hostRequestMaps = append(hostRequestMaps, hostRequestMap)
			if v, err := convert.ConvertListMapToJsonString(hostRequestMaps); err != nil {
				return errs.WrapError(err)
			} else {
				request["HostGroups"] = v
			}
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				_, err := client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}

				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if nraw != nil && len(nraw.(*schema.Set).List()) > 0 {
			action = "AttachHostGroupAccountsToUser"
			hostRequestMaps := make([]map[string]interface{}, 0)
			hostRequestMap := make(map[string]interface{}, 0)
			hostRequestMap["HostGroupId"] = parts[2]
			hostRequestMap["HostAccountNames"] = nraw.(*schema.Set).List()
			hostRequestMaps = append(hostRequestMaps, hostRequestMap)
			if v, err := convert.ConvertListMapToJsonString(hostRequestMaps); err != nil {
				return errs.WrapError(err)
			} else {
				request["HostGroups"] = v
			}
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				_, err := client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}

				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return resourceHostGroupAccountUserAttachmentRead(d, meta)
}
func resourceHostGroupAccountUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DetachHostGroupAccountsFromUser"
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	request["RegionId"] = client.RegionId
	hostRequestMaps := make([]map[string]interface{}, 0)
	hostRequestMap := make(map[string]interface{}, 0)
	hostRequestMap["HostGroupId"] = parts[2]
	hostRequestMap["HostAccountNames"] = d.Get("host_account_names").(*schema.Set).List()
	hostRequestMaps = append(hostRequestMaps, hostRequestMap)
	if v, err := convert.ConvertListMapToJsonString(hostRequestMaps); err != nil {
		return errs.WrapError(err)
	} else {
		request["HostGroups"] = v
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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
