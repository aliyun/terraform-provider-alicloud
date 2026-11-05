package hbase

import (
	"fmt"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceHBaseZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHBaseZonesRead,

		Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHBaseZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var zoneIds []string

	request := hbase.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
		return hbaseClient.DescribeRegions(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_hbase_zones", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	regions, _ := raw.(*hbase.DescribeRegionsResponse)
	if len(regions.Regions.Region) <= 0 {
		return errs.WrapError(fmt.Errorf("[ERROR] There is no available region for HBase."))
	}
	for _, r := range regions.Regions.Region {
		for _, zonid := range r.Zones.Zone {
			if r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zonid.Id)
			}
		}
	}
	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
	}

	var s []map[string]interface{}
	for _, zoneId := range zoneIds {
		mapping := map[string]interface{}{"id": zoneId}
		s = append(s, mapping)
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
