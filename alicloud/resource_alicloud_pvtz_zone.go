package alicloud

import (
	"fmt"
	"time"

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

	raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
		return pvtzClient.AddZone(args)
	})
	if err != nil {
		return fmt.Errorf("AddZone got an error:%#v", err)
	}
	response, _ := raw.(*pvtz.AddZoneResponse)
	if response == nil {
		return fmt.Errorf("AddZone got a nil response: %#v", response)
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

	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(request)
		})

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting zone failed: %#v", err))
		}

		if _, err := pvtzService.DescribePvtzZoneInfo(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil

	})

}
