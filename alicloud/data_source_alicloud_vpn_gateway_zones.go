// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudVPNGatewayZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVPNGatewayZoneRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudVPNGatewayZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeVpnGatewayAvailableZones"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	if v, ok := d.GetOk("spec"); ok {
		query["Spec"] = StringPointer(v.(string))
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcGet("Vpc", "2016-04-28", action, query, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	resp, _ := jsonpath.Get("$.AvailableZoneIdList", response)

	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	availableZoneIdList1Raw := resp.([]interface{})
	zonesMaps := make([]map[string]interface{}, 0)
	if availableZoneIdList1Raw != nil {
		for _, availableZoneIdListChild1Raw := range availableZoneIdList1Raw {
			zonesMap := make(map[string]interface{})
			availableZoneIdListChild1Raw := availableZoneIdListChild1Raw.(map[string]interface{})
			zonesMap["zone_id"] = availableZoneIdListChild1Raw["ZoneId"]
			zonesMap["zone_name"] = availableZoneIdListChild1Raw["ZoneName"]

			ids = append(ids, fmt.Sprint(availableZoneIdListChild1Raw["ZoneId"]))
			zonesMaps = append(zonesMaps, zonesMap)
		}
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("zones", zonesMaps); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
