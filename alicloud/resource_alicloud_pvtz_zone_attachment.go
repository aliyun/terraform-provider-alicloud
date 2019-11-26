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
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("vpcs") {
				d.SetNewComputed("vpc_ids")
			} else if d.HasChange("vpc_ids") {
				d.SetNewComputed("vpcs")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_ids": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "The attribute vpc_ids has been deprecated on pvtz_zone resource. Replace it with vpcs which supports vpc bindings for different regions",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
			"vpcs": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"vpc_ids"},
				Computed:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
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
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	request := pvtz.CreateBindZoneVpcRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Id()

	vpcIdMap := make(map[string]string)

	vpcIds := d.Get("vpc_ids").(*schema.Set).List()

	if d.HasChange("vpc_ids") && len(vpcIds) != 0 {
		bindZoneVpcVpcs := make([]pvtz.BindZoneVpcVpcs, len(vpcIds))
		for i, v := range vpcIds {
			bindZoneVpcVpcs[i].VpcId = v.(string)
			bindZoneVpcVpcs[i].RegionId = client.RegionId
			vpcIdMap[bindZoneVpcVpcs[i].VpcId] = bindZoneVpcVpcs[i].VpcId
		}
		request.Vpcs = &bindZoneVpcVpcs
	} else {
		vpcs := d.Get("vpcs").(*schema.Set).List()
		bindZoneVpcVpcs := make([]pvtz.BindZoneVpcVpcs, len(vpcs))
		for i, v := range vpcs {
			vpc := v.(map[string]interface{})

			bindZoneVpcVpcs[i].VpcId = vpc["vpc_id"].(string)
			regionId := vpc["region_id"].(string)
			if regionId == "" {
				regionId = client.RegionId
			}
			bindZoneVpcVpcs[i].RegionId = regionId
			vpcIdMap[bindZoneVpcVpcs[i].VpcId] = bindZoneVpcVpcs[i].VpcId
		}
		request.Vpcs = &bindZoneVpcVpcs
	}

	request.UserClientIp = d.Get("user_client_ip").(string)
	request.Lang = d.Get("lang").(string)

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
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
		return WrapError(err)
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
	vpcMaps := make([]map[string]interface{}, 0)
	for _, vpc := range vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
		vpcMap := map[string]interface{}{
			"vpc_id":    vpc.VpcId,
			"region_id": vpc.RegionId,
		}
		vpcMaps = append(vpcMaps, vpcMap)
	}

	d.Set("zone_id", d.Id())
	if err := d.Set("vpc_ids", vpcIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("vpcs", vpcMaps); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudPvtzZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	request := pvtz.CreateBindZoneVpcRequest()
	request.RegionId = client.RegionId
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
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
