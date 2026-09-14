package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudApiGatewayTrafficControls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApiGatewayTrafficControlsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"traffic_control_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"traffic_control_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_control_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_control_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_control_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudApiGatewayTrafficControlsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeTrafficControls"
	request := make(map[string]interface{})
	query := make(map[string]interface{})

	if v, ok := d.GetOk("traffic_control_id"); ok {
		query["TrafficControlId"] = v
	}
	if v, ok := d.GetOk("traffic_control_name"); ok {
		query["TrafficControlName"] = v
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	pageNumber := 1
	pageSize := 50

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	for {
		query["PageNumber"] = fmt.Sprint(pageNumber)
		query["PageSize"] = fmt.Sprint(pageSize)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_traffic_controls", action, AlibabaCloudSdkGoERROR)
		}

		totalCount := formatInt(response["TotalCount"])
		// When no traffic controls exist (e.g. during destroy-phase refresh),
		// the response omits the TrafficControls key. Skip the jsonpath lookup
		// to avoid an "unknown key" error and return an empty list.
		if totalCount == 0 {
			break
		}

		resp, err := jsonpath.Get("$.TrafficControls.TrafficControl", response)
		if err != nil {
			// An out-of-range page omits the TrafficControls key; treat it as the end of pagination.
			break
		}
		result, _ := resp.([]interface{})

		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TrafficControlId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                   fmt.Sprint(object["TrafficControlId"]),
			"traffic_control_id":   fmt.Sprint(object["TrafficControlId"]),
			"traffic_control_name": object["TrafficControlName"],
			"traffic_control_unit": object["TrafficControlUnit"],
			"api_default":          object["ApiDefault"],
			"user_default":         object["UserDefault"],
			"app_default":          object["AppDefault"],
			"description":          object["Description"],
			"create_time":          object["CreatedTime"],
			"modified_time":        object["ModifiedTime"],
			"region_id":            client.RegionId,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("controls", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
