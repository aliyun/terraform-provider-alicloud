package alicloud

import (
	"fmt"
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

	args := pvtz.CreateAddZoneRequest()

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.ZoneName = v.(string)
	}
	// API AddZone has a throttling limitation 5qps which one use only can send 5 requests in one second.
	invoker := NewPvtzInvoker()
	var raw interface{}
	if err := invoker.Run(func() error {
		rsp, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZone(args)
		})
		raw = rsp
		return err
	}); err != nil {
		return WrapError(args.GetActionName(), args.ZoneName, APIERROR, err)
	}
	response, _ := raw.(*pvtz.AddZoneResponse)
	if response == nil {
		return WrapError(args.GetActionName(), args.ZoneName, SDKERROR, fmt.Errorf("AddZone got a nil response: %#v", response))
	}

	d.SetId(response.ZoneId)

	return resourceAlicloudPvtzZoneUpdate(d, meta)

}

func resourceAlicloudPvtzZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = d.Id()

	raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
		return pvtzClient.DescribeZoneInfo(request)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return err
	}
	response, _ := raw.(*pvtz.DescribeZoneInfoResponse)

	d.Set("name", response.ZoneName)
	d.Set("remark", response.Remark)
	d.Set("creation_time", response.CreateTime)
	d.Set("update_time", response.UpdateTime)
	d.Set("is_ptr", response.IsPtr)
	d.Set("record_count", response.RecordCount)

	return nil
}

func resourceAlicloudPvtzZoneUpdate(d *schema.ResourceData, meta interface{}) error {

	request := pvtz.CreateUpdateZoneRemarkRequest()
	request.ZoneId = d.Id()

	if d.HasChange("remark") {
		request.Remark = d.Get("remark").(string)

		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.UpdateZoneRemark(request)
		})
		if err != nil {
			return err
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{PvtzThrottlingUser}) {
				return resource.RetryableError(WrapError(request.GetActionName(), d.Id(), APIERROR, err))
			}
			return resource.NonRetryableError(WrapError(request.GetActionName(), d.Id(), APIERROR, err))
		}

		if _, err := pvtzService.DescribePvtzZoneInfo(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError("DescribePvtzZoneInfo", d.Id(), ProviderERROR, err))
		}

		return nil

	})

}
