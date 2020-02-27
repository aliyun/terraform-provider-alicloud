package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEssAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssAlarmsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "custom"}, false),
			},
			"states": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OK", "ALARM", "INSUFFICIENT_DATA"}, false),
			},
			"alarms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"alarm_actions": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
							MaxItems: 5,
							MinItems: 1,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  System,
							ValidateFunc: validation.StringInSlice([]string{
								string(System),
								string(Custom),
							}, false),
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  FiveMinite,
							ValidateFunc: validation.IntInSlice([]int{
								int(OneMinite),
								int(TwoMinite),
								int(FiveMinite),
								int(FifteenMinite),
							}),
						},
						"statistics": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Average,
							ValidateFunc: validation.StringInSlice([]string{
								string(Average),
								string(Minimum),
								string(Maximum),
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Gte,
							ValidateFunc: validation.StringInSlice([]string{
								string(Gt),
								string(Gte),
								string(Lt),
								string(Lte),
							}, false),
						},
						"evaluation_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"cloud_monitor_group_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dimensions": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDescribeAlarmsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}
	if isenable, ok := d.GetOk("IsEnable"); ok && isenable.(requests.Boolean) != "" {
		request.IsEnable = isenable.(requests.Boolean)
	}
	if metric_type, ok := d.GetOk("Metric_type"); ok && metric_type.(string) != "" {
		request.MetricType = metric_type.(string)
	}
	var allAlarms []ess.Alarm
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeAlarms(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_alarms", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeAlarmsResponse)
		if len(response.AlarmList.Alarm) < 1 {
			break
		}
		allAlarms = append(allAlarms, response.AlarmList.Alarm...)
		if len(response.AlarmList.Alarm) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	var filteredAlarms = make([]ess.Alarm, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, alarm := range allAlarms {
			if okNameRegex && nameRegex != "" {
				var r = regexp.MustCompile(nameRegex.(string))
				if r != nil && !r.MatchString(alarm.AlarmTaskName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[alarm.AlarmTaskId]; !ok {
					continue
				}
			}
			filteredAlarms = append(filteredAlarms, alarm)
		}
	} else {
		filteredAlarms = allAlarms
	}

	return alarmsDescriptionAttribute(d, filteredAlarms, meta)
}

func alarmsDescriptionAttribute(d *schema.ResourceData, alarms []ess.Alarm, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, alarm := range alarms {
		mapping := map[string]interface{}{
			"state":               alarm.State,
			"task_id":             alarm.AlarmTaskId,
			"name":                alarm.Name,
			"scaling_group_id":    alarm.ScalingGroupId,
			"metric_name":         alarm.MetricName,
			"description":         alarm.Description,
			"enable":              alarm.Enable,
			"alarm_actions":       alarm.AlarmActions,
			"metric_type":         alarm.MetricType,
			"period":              alarm.Period,
			"statistics":          alarm.Statistics,
			"threshold":           alarm.Threshold,
			"comparison_operator": alarm.ComparisonOperator,
			"evaluation_count":    alarm.EvaluationCount,
			"dimensions":          alarm.Dimensions,
		}
		ids = append(ids, alarm.AlarmTaskId)
		names = append(names, alarm.Name)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("alarms", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
