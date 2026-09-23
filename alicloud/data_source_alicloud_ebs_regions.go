package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEbsRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEbsRegionsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEbsRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeRegions"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
	}

	regionId, regionIdOk := d.GetOk("region_id")

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	response, err = client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ebs_zones", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Regions", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Regions", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {

		item := v.(map[string]interface{})
		if regionIdOk && regionId.(string) != "" && regionId.(string) != item["RegionId"].(string) {
			continue
		}

		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"region_id": object["RegionId"],
		}

		zones := make([]map[string]interface{}, 0)
		if zoneArgs, ok := object["Zones"]; ok {
			for _, zoneArg := range zoneArgs.([]interface{}) {
				zoneMap := zoneArg.(map[string]interface{})
				zones = append(zones, map[string]interface{}{
					"zone_id": zoneMap["ZoneId"],
				})
			}
		}
		mapping["zones"] = zones

		s = append(s, mapping)
	}

	d.SetId(time.Now().String())

	if err := d.Set("regions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
