package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEcdZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdZonesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
	}
}

func dataSourceAlicloudEcdZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZones"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_zones", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Zones", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Zones", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"zone_id": object["ZoneId"],
		}
		s = append(s, mapping)
		ids = append(ids, fmt.Sprint(mapping["zone_id"]))
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
