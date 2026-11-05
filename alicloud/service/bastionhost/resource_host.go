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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostCreate,
		Read:   resourceHostRead,
		Update: resourceHostUpdate,
		Delete: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"active_address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_private_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_public_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Linux", "Windows"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ecs", "Local", "Rds"}, false),
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHost"
	request := make(map[string]interface{})
	var err error
	request["ActiveAddressType"] = d.Get("active_address_type")
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["HostName"] = d.Get("host_name")
	if v, ok := d.GetOk("host_private_address"); ok {
		request["HostPrivateAddress"] = v
	}
	if v, ok := d.GetOk("host_public_address"); ok {
		request["HostPublicAddress"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("instance_region_id"); ok {
		request["InstanceRegionId"] = v
	}
	request["OSType"] = d.Get("os_type")
	request["RegionId"] = client.RegionId
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_instance_id"); ok {
		request["SourceInstanceId"] = v
	}
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_bastionhost_host", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostId"]))

	return resourceHostRead(d, meta)
}
func resourceHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := NewService(client)
	object, err := yundunBastionhostService.DescribeBastionhostHost(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host yundunBastionhostService.DescribeBastionhostHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("host_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("active_address_type", object["ActiveAddressType"])
	d.Set("comment", object["Comment"])
	d.Set("host_name", object["HostName"])
	d.Set("host_private_address", object["HostPrivateAddress"])
	d.Set("host_public_address", object["HostPublicAddress"])
	d.Set("os_type", object["OSType"])
	d.Set("source", object["Source"])
	d.Set("source_instance_id", object["SourceInstanceId"])
	return nil
}
func resourceHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"HostIds":    convert.ConvertListToJsonString([]interface{}{parts[1]}),
		"InstanceId": parts[0],
	}
	if d.HasChange("active_address_type") {
		update = true
	}
	request["ActiveAddressType"] = d.Get("active_address_type")
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyHostsActiveAddressType"
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
	update = false
	modifyHostReq := map[string]interface{}{
		"HostId":     parts[1],
		"InstanceId": parts[0],
	}
	if d.HasChange("comment") {
		update = true
		if v, ok := d.GetOk("comment"); ok {
			modifyHostReq["Comment"] = v
		}
	}
	if d.HasChange("host_name") {
		update = true
		modifyHostReq["HostName"] = d.Get("host_name")
	}
	if d.HasChange("host_private_address") {
		update = true
		if v, ok := d.GetOk("host_private_address"); ok {
			modifyHostReq["HostPrivateAddress"] = v
		}
	}
	if d.HasChange("host_public_address") {
		update = true
		if v, ok := d.GetOk("host_public_address"); ok {
			modifyHostReq["HostPublicAddress"] = v
		}
	}
	if d.HasChange("os_type") {
		update = true
		modifyHostReq["OSType"] = d.Get("os_type")
	}
	modifyHostReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyHost"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, modifyHostReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyHostReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceHostRead(d, meta)
}
func resourceHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteHost"
	var response map[string]interface{}
	request := map[string]interface{}{
		"HostId":     parts[1],
		"InstanceId": parts[0],
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
		if errs.IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
