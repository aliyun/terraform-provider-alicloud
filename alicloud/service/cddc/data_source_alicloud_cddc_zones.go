package cddc

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZonesRead,
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataSourceZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ids := make([]string, 0)
	action := "DescribeRegions"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cddc_zones", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Regions.RDSRegion", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Regions.RDSRegion", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if client.RegionId != item["RegionId"].(string) {
			continue
		}
		if strings.Contains(strings.ToLower(item["ZoneId"].(string)), strings.ToLower("MAZ")) {
			continue
		}
		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"region_id": object["RegionId"],
			"zone_id":   object["ZoneId"],
			"id":        fmt.Sprint(object["ZoneId"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}
	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
