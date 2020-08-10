package alicloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudHBaseZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHBaseZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHBaseZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := hbase.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
		return hbaseClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbase_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	regions, _ := raw.(*hbase.DescribeRegionsResponse)
	if len(regions.Regions.Region) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available region for HBase."))
	}
	for _, r := range regions.Regions.Region {
		for _, zonid := range r.Zones.Zone {
			if multi && strings.Contains(zonid.Id, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zonid.Id)
				continue
			}
			if !multi && !strings.Contains(zonid.Id, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zonid.Id)
				continue
			}
		}
	}
	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
	}

	var s []map[string]interface{}
	if !multi {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{"id": zoneId}
			s = append(s, mapping)
		}
	} else {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{
				"id":             zoneId,
				"multi_zone_ids": splitMultiZoneId(zoneId),
			}
			s = append(s, mapping)
		}
	}
	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
