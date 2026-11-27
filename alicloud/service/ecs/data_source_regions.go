package ecs

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	legacyclient "github.com/aliyun/terraform-provider-alicloud/alicloud/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"current": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			//Computed value
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	currentRegion := client.RegionId

	// Neither branch sends AcceptLanguage: the legacy path matches the old
	// ecsClient call (the old SDK request carries no AcceptLanguage field), and
	// the V4 path follows the shared gateway configuration, which injects
	// AcceptLanguage=en-US for every migrated product. As a result local_name
	// stays in the local language under v1 signing and switches to English
	// under v4 signing — an inherent property of the EO migration.
	response, err := client.LegacySdkRpc(request, "/", ecs.CreateDescribeRegionsResponse())
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_regions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), response, request)
	resp := &ecs.DescribeRegionsResponse{}
	if err := legacyclient.LegacyRpcResponse(response, resp); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_regions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	// The rehydrated envelope is trusted only after the top-level Regions key
	// is confirmed present: without this check a payload with a different
	// shape decodes into an empty list and would be reported to the user as
	// "no matching regions found" instead of a structure mismatch.
	if _, ok := response["Regions"]; !ok {
		return WrapErrorf(fmt.Errorf("missing Regions in response"), DataDefaultErrorMsg, "alicloud_regions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if len(resp.Regions.Region) == 0 {
		return fmt.Errorf("no matching regions found")
	}
	name, nameOk := d.GetOk("name")
	current := false
	if v, exist := d.GetOk("current"); exist {
		current = v.(bool)
	}
	var filterRegions []ecs.Region
	for _, region := range resp.Regions.Region {
		if current {
			if nameOk && name.(string) != currentRegion {
				return fmt.Errorf("name doesn't match current region: %#v, please input again.", currentRegion)
			}
			if region.RegionId == currentRegion {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		if nameOk {
			if name.(string) == region.RegionId {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		filterRegions = append(filterRegions, region)
	}
	if len(filterRegions) < 1 {
		return fmt.Errorf("Your query region returned no results. Please change your search criteria and try again.")
	}

	return regionsDescriptionAttributes(d, filterRegions)
}

func regionsDescriptionAttributes(d *schema.ResourceData, regions []ecs.Region) error {
	var ids []string
	var s []map[string]interface{}
	for _, region := range regions {
		mapping := map[string]interface{}{
			"id":         region.RegionId,
			"region_id":  region.RegionId,
			"local_name": region.LocalName,
		}

		log.Printf("[DEBUG] alicloud_regions - adding region mapping: %v", mapping)
		ids = append(ids, string(region.RegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("regions", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
