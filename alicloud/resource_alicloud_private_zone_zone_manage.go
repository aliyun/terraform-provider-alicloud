package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudPrivateZoneZoneManage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPrivateZoneZoneManageCreate,
		Read:   resourceAlicloudPrivateZoneZoneManageRead,
		Update: resourceAlicloudPrivateZoneZoneManageUpdate,
		Delete: resourceAlicloudPrivateZoneZoneManageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en", "ja"}, false),
			},
			"proxy_pattern": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ZONE", "RECORD"}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPrivateZoneZoneManageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := pvtz.CreateAddZoneRequest()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("proxy_pattern"); ok {
		request.ProxyPattern = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = v.(string)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZone(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*pvtz.AddZoneResponse)
		d.SetId(fmt.Sprintf("%v", response.ZoneId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_private_zone_zone_manage", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudPrivateZoneZoneManageUpdate(d, meta)
}
func resourceAlicloudPrivateZoneZoneManageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePrivateZoneZoneManage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("proxy_pattern", object.ProxyPattern)
	d.Set("remark", object.Remark)
	d.Set("zone_name", object.ZoneName)
	return nil
}
func resourceAlicloudPrivateZoneZoneManageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	update := false
	request := pvtz.CreateSetProxyPatternRequest()
	request.ZoneId = d.Id()
	if !d.IsNewResource() && d.HasChange("proxy_pattern") {
		update = true
	}
	request.ProxyPattern = d.Get("proxy_pattern").(string)
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if !d.IsNewResource() && d.HasChange("user_client_ip") {
		update = true
		request.UserClientIp = d.Get("user_client_ip").(string)
	}
	if update {
		err := resource.Retry(15*time.Second, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.SetProxyPattern(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("proxy_pattern")
		d.SetPartial("lang")
		d.SetPartial("user_client_ip")
	}
	update = false
	updateZoneRemarkReq := pvtz.CreateUpdateZoneRemarkRequest()
	updateZoneRemarkReq.ZoneId = d.Id()
	updateZoneRemarkReq.Lang = d.Get("lang").(string)
	if d.HasChange("remark") {
		update = true
		updateZoneRemarkReq.Remark = d.Get("remark").(string)
	}
	updateZoneRemarkReq.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		err := resource.Retry(15*time.Second, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.UpdateZoneRemark(updateZoneRemarkReq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(updateZoneRemarkReq.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateZoneRemarkReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudPrivateZoneZoneManageRead(d, meta)
}
func resourceAlicloudPrivateZoneZoneManageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	err := resource.Retry(15*time.Second, func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
