package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudThreatDetectionBaselineStrategies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionBaselineStrategiesRead,
		Schema: map[string]*schema.Schema{
			"custom_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"strategy_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"strategies": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"baseline_strategy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"baseline_strategy_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"custom_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cycle_days": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"cycle_start_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"risk_sub_type_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"target_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionBaselineStrategiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("custom_type"); ok {
		request["CustomType"] = v
	}
	if v, ok := d.GetOk("strategy_ids"); ok {
		request["StrategyIds"] = v
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

	var baselineStrategyNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		baselineStrategyNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "DescribeStrategy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_baseline_strategies", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Strategies", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Strategies", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
				continue
			}
		}

		if baselineStrategyNameRegex != nil && !baselineStrategyNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	sasService := SasService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                     fmt.Sprint(object["Id"]),
			"baseline_strategy_id":   object["Id"],
			"baseline_strategy_name": object["Name"],
			"custom_type":            object["CustomType"],
			"cycle_days":             object["CycleDays"],
			"cycle_start_time":       object["CycleStartTime"],
			"end_time":               object["EndTime"],
			"start_time":             object["StartTime"],
		}

		ids = append(ids, fmt.Sprint(object["Id"]))
		names = append(names, object["Name"])
		id := fmt.Sprint(object["Id"])
		object, err = sasService.DescribeThreatDetectionBaselineStrategy(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["baseline_strategy_id"] = object["Id"]
		mapping["baseline_strategy_name"] = object["Name"]
		mapping["cycle_days"] = object["CycleDays"]
		mapping["cycle_start_time"] = object["CycleStartTime"]
		mapping["end_time"] = object["EndTime"]
		mapping["risk_sub_type_name"] = object["RiskSubTypeName"]
		mapping["start_time"] = object["StartTime"]
		mapping["target_type"] = object["TargetType"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("strategies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
