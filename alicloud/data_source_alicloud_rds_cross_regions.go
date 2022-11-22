package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsCrossRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsCrossRegionsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regions": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudRdsCrossRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeAvailableCrossRegion"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var response map[string]interface{}
	var s []map[string]interface{}
	var ids []string
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_cross_region", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Regions.Region", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Regions.Region", response)
	}
	for _, r := range resp.([]interface{}) {
		mapping := map[string]interface{}{
			"id": fmt.Sprint(r),
		}
		ids = append(ids, fmt.Sprint(r))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("regions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
