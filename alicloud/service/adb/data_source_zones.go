package adb

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAdbZones() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAdbZonesRead,

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

func DataSourceAdbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := adb.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeRegions(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_adb_zones", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	regions, _ := raw.(*adb.DescribeRegionsResponse)
	if len(regions.Regions.Region) <= 0 {
		return errs.WrapError(fmt.Errorf("[ERROR] There is no available region for adb."))
	}
	for _, r := range regions.Regions.Region {
		for _, zone := range r.Zones.Zone {
			if multi && strings.Contains(zone.ZoneId, helper.MultiIzSymbol) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zone.ZoneId)
				continue
			}
			if !multi && !strings.Contains(zone.ZoneId, helper.MultiIzSymbol) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zone.ZoneId)
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
				"multi_zone_ids": helper.SplitMultiZoneId(zoneId),
			}
			s = append(s, mapping)
		}
	}
	d.SetId(helper.DataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", zoneIds); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
