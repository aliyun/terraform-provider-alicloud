package alicloud

import (
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEmrInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"destination_resource": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Zone",
					"Network",
					"InstanceType",
					"SystemDisk",
					"DataDisk",
				}),
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"support_local_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"PostPaid",
					"PrePaid",
				}),
			},
			"support_node_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_storage_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEmrInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
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

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ListEmrAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_instance_types",
			request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	supportedResources := make(map[string][]emr.SupportedResource)
	resourceResponse, _ := raw.(*emr.ListEmrAvailableResourceResponse)
	for _, zoneInfo := range resourceResponse.EmrZoneInfoList.EmrZoneInfo {
		resourceInfo := zoneInfo.EmrResourceInfoList.EmrResourceInfo
		if len(resourceInfo) == 1 {
			supportedResources[zoneInfo.ZoneId] = resourceInfo[0].SupportedResourceList.SupportedResource
		}
	}

	return emrClusterInstanceTypesAttributes(d, supportedResources)
}

func emrClusterInstanceTypesAttributes(d *schema.ResourceData,
	supportedResources map[string][]emr.SupportedResource) error {
	var ids []string
	var zoneIDs []string
	var s []map[string]interface{}

	for k, v := range supportedResources {
		// ignore empty zoneId or empty emr instance type of the specific zoneId
		if k == "" || len(v) == 0 {
			continue
		}

		zoneIDs = append(zoneIDs, k)
	}
	sort.Strings(zoneIDs)
	localStorage := d.Get("support_local_storage").(bool)
	supportNodeType := d.Get("support_node_type").([]interface{})
	nodeTypeFilter := func(filter []interface{}, source []string) bool {
		if len(source) == 0 {
			return false
		}
		sourceMapping := make(map[string]struct{})
		for _, s := range source {
			sourceMapping[s] = struct{}{}
		}
		for _, f := range filter {
			if _, ok := sourceMapping[f.(string)]; !ok {
				return false
			}
		}
		return true
	}

	for _, zoneID := range zoneIDs {
		mapping := map[string]interface{}{
			"zone_id": zoneID,
		}
		if v, ok := supportedResources[zoneID]; ok {
			for _, tpe := range v {
				if nodeTypeFilter(supportNodeType, tpe.SupportNodeTypeList.SupportNodeType) == false {
					continue
				}

				if localStorage == true && tpe.EmrInstanceType.LocalStorageAmount > 0 {
					mapping["id"] = tpe.EmrInstanceType.InstanceType
					mapping["local_storage_capacity"] = tpe.EmrInstanceType.LocalStorageCapacity
					ids = append(ids, tpe.EmrInstanceType.InstanceType)
					s = append(s, mapping)
					break
				} else if localStorage == false && tpe.EmrInstanceType.LocalStorageAmount == 0 {
					mapping["id"] = tpe.EmrInstanceType.InstanceType
					ids = append(ids, tpe.EmrInstanceType.InstanceType)
					s = append(s, mapping)
					break
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("types", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
