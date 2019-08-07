package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func resourceAlicloudPvtzZoneAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	zone, err := pvtzService.DescribePvtzZone(d.Get("zone_id").(string))
	if err != nil {
		return WrapError(err)
	}

	d.SetId(zone.ZoneId)

	return resourceAlicloudPvtzZoneAttachmentUpdate(d, meta)
}

func resourceAlicloudPvtzZoneAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("vpc_ids") {
		client := meta.(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		pvtzService := PvtzService{client}

		request := pvtz.CreateBindZoneVpcRequest()
		request.ZoneId = d.Id()

		o, n := d.GetChange("vpc_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		bindZoneVpcs := ns.Difference(os).List()

		vpcIdMap := make(map[string]string)
		vpcs := make([]pvtz.BindZoneVpcVpcs, len(bindZoneVpcs))
		for i, e := range bindZoneVpcs {
			vpcId := e.(string)
			object, err := vpcService.DescribeVpc(vpcId)
			if err != nil {
				return WrapError(err)
			}

			regionId := object.RegionId

			vpcs[i].RegionId = regionId
			vpcs[i].VpcId = vpcId
			vpcIdMap[vpcId] = vpcId
		}

		request.Vpcs = &vpcs
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.BindZoneVpc(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy, ZoneNotExists}) {
					time.Sleep(5 * time.Second)
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

		if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudPvtzZoneAttachmentRead(d, meta)
}

func resourceAlicloudPvtzZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	object, err := pvtzService.DescribePvtzZoneAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	vpcs := object.BindVpcs.Vpc
	vpcIds := make([]string, 0)
	for _, vpc := range vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	d.Set("zone_id", d.Id())
	if err := d.Set("vpc_ids", vpcIds); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudPvtzZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	request := pvtz.CreateBindZoneVpcRequest()
	request.ZoneId = d.Id()
	vpcs := make([]pvtz.BindZoneVpcVpcs, 0)
	request.Vpcs = &vpcs

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.BindZoneVpc(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{PvtzThrottlingUser, PvtzSystemBusy}) {
				time.Sleep(time.Duration(2) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(pvtzService.WaitForPvtzZoneAttachment(d.Id(), Deleted, DefaultTimeout))
}
