package alicloud

import (
	"time"

	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

	return nil
}

func resourceAlicloudPvtzZoneUpdate(d *schema.ResourceData, meta interface{}) error {

	request := pvtz.CreateUpdateZoneRemarkRequest()
	request.ZoneId = d.Id()

	if d.HasChange("remark") {
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
	}
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
