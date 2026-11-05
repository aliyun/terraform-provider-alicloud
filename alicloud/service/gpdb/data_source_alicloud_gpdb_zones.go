package gpdb

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZonesRead,

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

func dataSourceZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := gpdb.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
		return gpdbClient.DescribeRegions(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_gpdb_zones", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*gpdb.DescribeRegionsResponse)
	if len(response.Regions.Region) <= 0 {
		return errs.WrapError(fmt.Errorf("[ERROR] There is no available region for gpdb."))
	}
	for _, r := range response.Regions.Region {
		for _, zoneId := range r.Zones.Zone {
			if multi && strings.Contains(zoneId.ZoneId, helper.MultiIzSymbol) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zoneId.ZoneId)
				continue
			}
			if !multi && !strings.Contains(zoneId.ZoneId, helper.MultiIzSymbol) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zoneId.ZoneId)
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
