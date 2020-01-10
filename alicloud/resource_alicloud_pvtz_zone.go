package alicloud

import (
	"time"

	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudPvtzZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneCreate,
		Read:   resourceAlicloudPvtzZoneRead,
		Update: resourceAlicloudPvtzZoneUpdate,
		Delete: resourceAlicloudPvtzZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_pattern": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ZONE", "RECORD"}, false),
				Default:      "ZONE",
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en", "jp"}, false),
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_ptr": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPvtzZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := pvtz.CreateAddZoneRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.ZoneName = v.(string)
	}
	if v, ok := d.GetOk("proxy_pattern"); ok && v.(string) != "" {
		request.ProxyPattern = v.(string)
	}
	if v, ok := d.GetOk("user_client_ip"); ok && v.(string) != "" {
		request.UserClientIp = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok && v.(string) != "" {
		request.Lang = v.(string)
	}
	// API AddZone has a throttling limitation 5qps which one use only can send 5 requests in one second.
	var response *pvtz.AddZoneResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZone(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*pvtz.AddZoneResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.ZoneId)

	return resourceAlicloudPvtzZoneUpdate(d, meta)

}

func resourceAlicloudPvtzZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZone(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	d.Set("name", object.ZoneName)
	d.Set("remark", object.Remark)
	d.Set("creation_time", object.CreateTime)
	d.Set("update_time", object.UpdateTime)
	d.Set("is_ptr", object.IsPtr)
	d.Set("record_count", object.RecordCount)
	d.Set("proxy_pattern", object.ProxyPattern)

	return nil
}

func resourceAlicloudPvtzZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("remark") {
		request := pvtz.CreateUpdateZoneRemarkRequest()
		request.ZoneId = d.Id()
		request.Remark = d.Get("remark").(string)

		client := meta.(*connectivity.AliyunClient)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.UpdateZoneRemark(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudPvtzZoneRead(d, meta)
	}
	if d.HasChange("proxy_pattern") || d.HasChange("user_client_ip") || d.HasChange("lang") {
		request := pvtz.CreateSetProxyPatternRequest()
		request.ZoneId = d.Id()
		request.ProxyPattern = d.Get("proxy_pattern").(string)
		request.UserClientIp = d.Get("user_client_ip").(string)
		request.Lang = d.Get("lang").(string)

		client := meta.(*connectivity.AliyunClient)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.SetProxyPattern(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("proxy_pattern")
		d.SetPartial("user_client_ip")
		d.SetPartial("lang")
	}
	d.Partial(false)
	return resourceAlicloudPvtzZoneRead(d, meta)
}

func resourceAlicloudPvtzZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	runtime.Caller(0)
	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{PvtzThrottlingUser, PvtzSystemBusy, ZoneVpcExists}) {
				time.Sleep(time.Duration(2) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil

	})

	if err != nil {
		if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(pvtzService.WaitForPvtzZone(d.Id(), Deleted, DefaultTimeout))
}
