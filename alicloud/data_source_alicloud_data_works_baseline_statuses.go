package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksBaselineStatuses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksBaselineStatusesRead,
		Schema: map[string]*schema.Schema{
			"bizdate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"baseline_types": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"finish_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buffer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finish_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finish_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bizdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exp_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sla_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_cast": {
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

func dataSourceAlicloudDataWorksBaselineStatusesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListBaselineStatuses"
	request := make(map[string]interface{})
	request["Bizdate"] = d.Get("bizdate").(string)
	if v, ok := d.GetOk("baseline_types"); ok {
		request["BaselineTypes"] = v
	}
	if v, ok := d.GetOk("owner"); ok {
		request["Owner"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("finish_status"); ok {
		request["FinishStatus"] = v
	}
	if v, ok := d.GetOk("search_text"); ok {
		request["SearchText"] = v
	}
	if v, ok := d.GetOk("topic_id"); ok {
		request["TopicId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_baseline_statuses", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.BaselineStatuses", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.BaselineStatuses", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["BaselineId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":            fmt.Sprint(object["BaselineId"]),
			"baseline_id":   fmt.Sprint(object["BaselineId"]),
			"baseline_name": fmt.Sprint(object["BaselineName"]),
			"baseline_type": fmt.Sprint(object["BaselineType"]),
			"buffer":        fmt.Sprint(object["Buffer"]),
			"status":        fmt.Sprint(object["Status"]),
			"owner":         fmt.Sprint(object["Owner"]),
			"priority":      fmt.Sprint(object["Priority"]),
			"finish_status": fmt.Sprint(object["FinishStatus"]),
			"finish_time":   fmt.Sprint(object["FinishTime"]),
			"project_id":    fmt.Sprint(object["ProjectId"]),
			"bizdate":       fmt.Sprint(object["Bizdate"]),
			"exp_time":      fmt.Sprint(object["ExpTime"]),
			"in_group_id":   fmt.Sprint(object["InGroupId"]),
			"sla_time":      fmt.Sprint(object["SlaTime"]),
			"end_cast":      fmt.Sprint(object["EndCast"]),
			"region_id":     client.RegionId,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("statuses", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
