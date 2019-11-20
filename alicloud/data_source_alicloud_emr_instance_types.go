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
						"local_instance_type": {
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_instance_types", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	types := make(map[string][]emr.EmrInstanceType)
	resourceResponse, _ := raw.(*emr.ListEmrAvailableResourceResponse)
	for _, zoneInfo := range resourceResponse.EmrZoneInfoList.EmrZoneInfo {
		instanceTypes := make([]emr.EmrInstanceType, 0)
		resourceInfo := zoneInfo.EmrResourceInfoList.EmrResourceInfo
		if len(resourceInfo) == 1 {
			for _, res := range resourceInfo[0].SupportedResourceList.SupportedResource {
				instanceTypes = append(instanceTypes, res.EmrInstanceType)
			}
		}
		types[zoneInfo.ZoneId] = instanceTypes
	}

	return emrClusterInstanceTypesAttributes(d, types)
}

func emrClusterInstanceTypesAttributes(d *schema.ResourceData, types map[string][]emr.EmrInstanceType) error {
	var ids []string
	var zoneIDs []string
	var s []map[string]interface{}

	for k, v := range types {
		// ignore empty zoneId or empty emr instance type of the specific zoneId
		if k == "" || len(v) == 0 {
			continue
		}

		zoneIDs = append(zoneIDs, k)
	}
	sort.Strings(zoneIDs)
	localStorage := d.Get("support_local_storage").(bool)
	for _, zoneID := range zoneIDs {
		mapping := map[string]interface{}{
			"zone_id": zoneID,
		}
		localFlag := make(map[string]string)
		if v, ok := types[zoneID]; ok {
			for _, tpe := range v {
				if localStorage == false {
					if tpe.LocalStorageAmount != 0 {
						continue
					}
					mapping["id"] = tpe.InstanceType
					ids = append(ids, tpe.InstanceType)
					s = append(s, mapping)
					break
				} else {
					if _, ok := localFlag["cloud"]; !ok && tpe.LocalStorageAmount == 0 {
						mapping["id"] = tpe.InstanceType
						localFlag["cloud"] = tpe.InstanceType
					} else if _, ok := localFlag["local"]; !ok && tpe.LocalStorageAmount > 0 {
						mapping["local_instance_type"] = tpe.InstanceType
						mapping["local_storage_capacity"] = tpe.LocalStorageCapacity
						localFlag["local"] = tpe.InstanceType
					}
					if len(localFlag) == 2 {
						s = append(s, mapping)
						ids = append(ids, localFlag["cloud"], localFlag["local"])
						break
					}
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
