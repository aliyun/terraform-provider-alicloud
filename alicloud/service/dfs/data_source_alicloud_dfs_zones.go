package dfs

import (
	"fmt"
	"strconv"
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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol_type": {
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

func dataSourceZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "GetRegion"
	request := make(map[string]interface{})
	request["InputRegionId"] = client.RegionId

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dfs_zones", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.AvailableZones", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AvailableZones", response)
	}

	s := make([]map[string]interface{}, 0)
	for _, object := range resp.([]interface{}) {
		if zoneMap, ok := object.(map[string]interface{}); ok {

			mapping := map[string]interface{}{
				"zone_id": fmt.Sprint(zoneMap["ZoneId"]),
			}

			options := make([]interface{}, 0)
			for _, optionArg := range zoneMap["Options"].([]interface{}) {
				optionMap := optionArg.(map[string]interface{})
				singleObj := make(map[string]interface{}, 0)
				singleObj["storage_type"] = fmt.Sprint(optionMap["StorageType"])
				singleObj["protocol_type"] = fmt.Sprint(optionMap["ProtocolType"])
				options = append(options, singleObj)
			}
			mapping["options"] = options

			s = append(s, mapping)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("zones", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
