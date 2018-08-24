package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"remark": &schema.Schema{
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
			"is_ptr": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"record_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPvtzZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

	args := pvtz.CreateAddZoneRequest()

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.ZoneName = v.(string)
	}

	response, err := conn.AddZone(args)
	if err != nil {
		return fmt.Errorf("AddZone got an error:%#v", err)
	}

	if response == nil {
		return fmt.Errorf("AddZone got a nil response: %#v", response)
	}

	d.SetId(response.ZoneId)

	return resourceAlicloudPvtzZoneUpdate(d, meta)

}

func resourceAlicloudPvtzZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = d.Id()

	response, err := conn.DescribeZoneInfo(request)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return err
	}

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

		if _, err := meta.(*AliyunClient).pvtzconn.UpdateZoneRemark(request); err != nil {
			return err
		}
	}

	return resourceAlicloudPvtzZoneRead(d, meta)
}

func resourceAlicloudPvtzZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteZone(request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting zone failed: %#v", err))
		}

		if _, err := client.DescribePvtzZoneInfo(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil

	})

}
