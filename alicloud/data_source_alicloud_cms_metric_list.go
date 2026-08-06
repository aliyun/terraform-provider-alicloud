package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// data source alicloud_cms_metric_list wraps the DescribeMetricList RPC
// (Cms/2019-01-01). It returns the metric datapoints of the given metric
// name and namespace within a time range. The result is paginated by
// NextToken and Datapoints is a JSON string that must be unmarshalled.

const cmsDescribeMetricListDefaultLength = 1000

func dataSourceAlicloudCmsMetricList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMetricListRead,
		Schema: map[string]*schema.Schema{
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dimensions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"express": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"next_token": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datapoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"average": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"sum": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"maximum": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"minimum": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actual_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCmsMetricListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeMetricList"
	request := make(map[string]interface{})
	request["MetricName"] = d.Get("metric_name")
	request["Namespace"] = d.Get("metric_namespace")
	if v, ok := d.GetOk("dimensions"); ok {
		request["Dimensions"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("express"); ok {
		request["Express"] = v
	}
	if v, ok := d.GetOk("length"); ok {
		request["Length"] = v.(int)
	}
	if v, ok := d.GetOk("next_token"); ok {
		request["NextToken"] = v
	}
	if _, ok := request["Length"]; !ok {
		request["Length"] = cmsDescribeMetricListDefaultLength
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	currentToken := fmt.Sprint(request["NextToken"])
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError", "BadRequest"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_metric_list", action, AlibabaCloudSdkGoERROR)
		}
		if code, ok := response["Code"].(string); ok && code != "" && code != "200" {
			msg := fmt.Sprint(response["Message"])
			return WrapError(fmt.Errorf("%s failed, code: %s, message: %s", action, code, msg))
		}

		datapointsStr, ok := response["Datapoints"].(string)
		var datapoints []map[string]interface{}
		if ok && datapointsStr != "" {
			if jsonErr := json.Unmarshal([]byte(datapointsStr), &datapoints); jsonErr != nil {
				return WrapErrorf(jsonErr, "unmarshal Datapoints for %s", action)
			}
		}
		objects = append(objects, datapoints...)

		nextToken, ok := response["NextToken"].(string)
		if !ok || nextToken == "" || nextToken == currentToken {
			break
		}
		currentToken = nextToken
		request["NextToken"] = nextToken
	}

	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, dp := range objects {
		item := map[string]interface{}{
			"uuid":        fmt.Sprint(dp["uuid"]),
			"metric_name": fmt.Sprint(dp["MetricName"]),
			"timestamp":   fmt.Sprint(dp["Timestamp"]),
			"average":     dp["Average"],
			"sum":         dp["Sum"],
			"maximum":     dp["Maximum"],
			"minimum":     dp["Minimum"],
			"count":       dp["Count"],
			"user_id":     fmt.Sprint(dp["userId"]),
			"instance_id": fmt.Sprint(dp["instanceId"]),
		}
		s = append(s, item)
		if id := fmt.Sprint(dp["uuid"]); id != "" && id != "<nil>" {
			ids = append(ids, id)
		}
	}

	d.SetId(dataResourceIdHash([]string{d.Get("metric_name").(string), d.Get("metric_namespace").(string)}))
	if err := d.Set("datapoints", s); err != nil {
		return WrapError(err)
	}
	if v, ok := response["NextToken"]; ok {
		if err := d.Set("next_token", fmt.Sprint(v)); err != nil {
			return WrapError(err)
		}
	}
	if v, ok := response["Period"]; ok {
		if err := d.Set("actual_period", fmt.Sprint(v)); err != nil {
			return WrapError(err)
		}
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
