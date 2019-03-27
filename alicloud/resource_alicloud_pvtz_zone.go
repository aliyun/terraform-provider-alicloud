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

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.ZoneName = v.(string)
	}
	// API AddZone has a throttling limitation 5qps which one use only can send 5 requests in one second.
	invoker := PvtzInvoker()
	var raw interface{}
	var err error
	if err := invoker.Run(func() error {
		raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZone(request)
		})
		addDebug(request.GetActionName(), raw)
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*pvtz.AddZoneResponse)

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
		invoker := PvtzInvoker()
		err := invoker.Run(func() error {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.UpdateZoneRemark(request)
			})
			addDebug(request.GetActionName(), raw)
			return err
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
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw)
		return nil

	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(pvtzService.WaitForPvtzZone(d.Id(), Deleted, DefaultTimeout))
}
