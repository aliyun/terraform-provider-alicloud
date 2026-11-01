package emr

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceEmrDiskTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmrDiskTypesRead,

		Schema: map[string]*schema.Schema{
			"destination_resource": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Zone",
					"Network",
					"InstanceType",
					"SystemDisk",
					"DataDisk",
				}, false),
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PostPaid",
					"PrePaid",
				}, false),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEmrDiskTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := emr.CreateListEmrAvailableResourceRequest()
	if dstRes, ok := d.GetOk("destination_resource"); ok {
		request.DestinationResource = strings.TrimSpace(dstRes.(string))
	}
	if typ, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = strings.TrimSpace(typ.(string))
	}
	if chargeType, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = strings.TrimSpace(chargeType.(string))
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = strings.TrimSpace(instanceType.(string))
	}
	if zoneID, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = strings.TrimSpace(zoneID.(string))
	}

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ListEmrAvailableResource(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_emr_disk_types", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	var resources []emr.SupportedResource
	resourceResponse, _ := raw.(*emr.ListEmrAvailableResourceResponse)
	for _, zoneInfo := range resourceResponse.EmrZoneInfoList.EmrZoneInfo {
		resourceInfo := zoneInfo.EmrResourceInfoList.EmrResourceInfo
		if len(resourceInfo) > 0 {
			resources = resourceInfo[0].SupportedResourceList.SupportedResource
		}
	}

	return emrClusterDiskTypesAttributes(d, resources)
}

func emrClusterDiskTypesAttributes(d *schema.ResourceData, resources []emr.SupportedResource) error {
	var ids []string
	var s []map[string]interface{}

	for _, resource := range resources {
		mapping := map[string]interface{}{
			"min":   resource.Min,
			"max":   resource.Max,
			"value": resource.Value,
		}
		ids = append(ids, resource.Value)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("types", s); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
