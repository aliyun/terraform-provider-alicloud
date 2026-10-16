package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError", "BadRequest"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			// DescribeMetricList returns error 10002 ("the metric ... is not
			// exist") when the requested metric name does not exist for the
			// given namespace. For this query-only data source that is a valid
			// empty result, so set the id and return instead of surfacing the
			// error.
			if IsExpectedErrors(err, []string{"10002"}) || IsExpectedErrorsMessage(err, "is not exist") {
				d.SetId(dataResourceIdHash([]string{d.Get("metric_name").(string), d.Get("metric_namespace").(string)}))
				return nil
			}
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
	for _, dp := range objects {
		// The API does not echo MetricName in datapoint objects, so backfill
		// the requested metric name; datapoints only exist for the queried
		// metric because both MetricName and Namespace are required request
		// parameters. Never render missing keys as "<nil>".
		metricName := d.Get("metric_name").(string)
		if v := dp["MetricName"]; v != nil {
			metricName = cmsMetricListString(v)
		}
		item := map[string]interface{}{
			"uuid":        cmsMetricListString(dp["uuid"]),
			"metric_name": metricName,
			"timestamp":   cmsMetricListTimestampString(dp["timestamp"]),
			"average":     dp["Average"],
			"sum":         dp["Sum"],
			"maximum":     dp["Maximum"],
			"minimum":     dp["Minimum"],
			"count":       dp["Count"],
			"user_id":     cmsMetricListString(dp["userId"]),
			"instance_id": cmsMetricListString(dp["instanceId"]),
		}
		s = append(s, item)
	}

	d.SetId(dataResourceIdHash([]string{d.Get("metric_name").(string), d.Get("metric_namespace").(string)}))
	if err := d.Set("datapoints", s); err != nil {
		return WrapError(err)
	}
	// NextToken is null on the last page; always store a string so the state
	// never holds "<nil>" and a re-read does not resend a bogus NextToken.
	nextToken := ""
	if v, ok := response["NextToken"]; ok {
		if str, ok := v.(string); ok {
			nextToken = str
		}
	}
	if err := d.Set("next_token", nextToken); err != nil {
		return WrapError(err)
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

func cmsMetricListString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// cmsMetricListTimestampString renders the millisecond epoch as a decimal
// string: json.Unmarshal decodes the JSON number into float64 and fmt.Sprint
// would emit scientific notation (e.g. "1.54877766e+12").
func cmsMetricListTimestampString(v interface{}) string {
	switch t := v.(type) {
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case string:
		return t
	default:
		return ""
	}
}
