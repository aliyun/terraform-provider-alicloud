package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsCollationTimeZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsCollationTimeZonesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// response value
			"collation_time_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"standard_time_offset": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsCollationTimeZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeCollationTimeZones"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"SourceIp": client.SourceIp,
	}

	var ids []string
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_collation_time_zones", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.CollationTimeZones.CollationTimeZone", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CollationTimeZones.CollationTimeZone", response)
	}
	result := make([]map[string]interface{}, 0)
	for _, r := range resp.([]interface{}) {
		collationTimeZone := r.(map[string]interface{})
		mapping := map[string]interface{}{}
		if v, ok := collationTimeZone["StandardTimeOffset"]; ok && v != "" {
			mapping["standard_time_offset"] = fmt.Sprint(v)
		}
		if v, ok := collationTimeZone["Description"]; ok && v != "" {
			mapping["description"] = fmt.Sprint(v)
		}
		if v, ok := collationTimeZone["TimeZone"]; ok && v != "" {
			mapping["time_zone"] = fmt.Sprint(v)
		}

		ids = append(ids, fmt.Sprint(collationTimeZone["TimeZone"]))
		result = append(result, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("collation_time_zones", result); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), result)
	}

	return nil
}
