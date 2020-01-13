package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"regexp"
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
			"RegionId": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"AlarmTaskId": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"IsEnable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"Metric_type": {
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "custom"}, false),
			},
			"alarms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"RegionId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"AlarmTaskId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"IsEnable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"Metric_type": {
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
			"state":            alarm.State,
			"task_id":          alarm.AlarmTaskId,
			"name":             alarm.Name,
			"scaling_group_id": alarm.ScalingGroupId,
			"metric_name":      alarm.MetricName,
			"description":      alarm.Description,
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
