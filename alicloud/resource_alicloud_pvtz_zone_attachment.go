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

	zone, err := pvtzService.DescribePvtzZoneInfo(d.Get("zone_id").(string))
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

		args := pvtz.CreateBindZoneVpcRequest()
		args.ZoneId = d.Id()

		o, n := d.GetChange("vpc_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		bindZoneVpcs := ns.Difference(os).List()

		vpcIdMap := make(map[string]string)
		vpcs := make([]pvtz.BindZoneVpcVpcs, len(bindZoneVpcs))
		for i, e := range bindZoneVpcs {
			vpcId := e.(string)
			v, err := vpcService.DescribeVpc(vpcId)
			if err != nil {
				return WrapError(err)
			}

			regionId := v.RegionId

			vpcs[i].RegionId = regionId
			vpcs[i].VpcId = vpcId
			vpcIdMap[vpcId] = vpcId
		}

		args.Vpcs = &vpcs
		invoker := PvtzInvoker()
		invoker.AddCatcher(Catcher{ZoneNotExists, 30, 3})
		if err := invoker.Run(func() error {
			_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.BindZoneVpc(args)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if err := pvtzService.WaitZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudPvtzZoneAttachmentRead(d, meta)
}

func resourceAlicloudPvtzZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	response, err := pvtzService.DescribePvtzZoneAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	vpcs := response.BindVpcs.Vpc
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.BindZoneVpc(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{PvtzThrottlingUser, PvtzSystemBusy}) {
				time.Sleep(time.Duration(2) * time.Second)
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			if !IsExceptedErrors(err, []string{PvtzInternalError}) {
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
		}

		if _, err := pvtzService.DescribePvtzZoneAttachment(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))

	})
}
