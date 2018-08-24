package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudPvtzZoneAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneAttachmentCreate,
		Update: resourceAlicloudPvtzZoneAttachmentUpdate,
		Read:   resourceAlicloudPvtzZoneAttachmentRead,
		Delete: resourceAlicloudPvtzZoneAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func resourceAlicloudPvtzZoneAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	zone, err := client.DescribePvtzZoneInfo(d.Get("zone_id").(string))
	if err != nil {
		return err
	}

	d.SetId(zone.ZoneId)

	return resourceAlicloudPvtzZoneAttachmentUpdate(d, meta)
}

func resourceAlicloudPvtzZoneAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("vpc_ids") {
		client := meta.(*AliyunClient)
		conn := client.pvtzconn

		args := pvtz.CreateBindZoneVpcRequest()
		args.ZoneId = d.Id()

		o, n := d.GetChange("vpc_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		bindZoneVpcs := ns.Difference(os).List()

		vpcs := make([]pvtz.BindZoneVpcVpcs, len(bindZoneVpcs))
		for i, e := range bindZoneVpcs {
			vpcId := e.(string)
			v, _ := client.DescribeVpc(vpcId)

			regionId := v.RegionId

			vpcs[i].RegionId = regionId
			vpcs[i].VpcId = vpcId
		}

		args.Vpcs = &vpcs

		_, err := conn.BindZoneVpc(args)
		if nil != err {
			return fmt.Errorf("bindZoneVpc error:%#v", err)
		}
	}

	return resourceAlicloudPvtzZoneAttachmentRead(d, meta)
}

func resourceAlicloudPvtzZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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

	var vpcIds []string
	vpcs := response.BindVpcs.Vpc
	for _, vpc := range vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	d.Set("zone_id", d.Id())
	d.Set("vpc_ids", vpcIds)

	return nil
}

func resourceAlicloudPvtzZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

	request := pvtz.CreateBindZoneVpcRequest()
	request.ZoneId = d.Id()
	vpcs := make([]pvtz.BindZoneVpcVpcs, 0)
	request.Vpcs = &vpcs

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.BindZoneVpc(request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error unbind zone vpc failed: %#v", err))
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
