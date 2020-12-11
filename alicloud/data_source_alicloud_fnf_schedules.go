package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudFnfSchedules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFnfSchedulesRead,
		Schema: map[string]*schema.Schema{
			"flow_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFnfSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSchedules"
	request := make(map[string]interface{})
	request["FlowName"] = d.Get("flow_name")
	var objects []map[string]interface{}
	var scheduleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		scheduleNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fnf_schedules", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Schedules", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Schedules", response)
	}
	for _, v := range resp.([]interface{}) {
		item := v.(map[string]interface{})
		if scheduleNameRegex != nil {
			if !scheduleNameRegex.MatchString(item["ScheduleName"].(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ScheduleName"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cron_expression":    object["CronExpression"],
			"description":        object["Description"],
			"enable":             object["Enable"],
			"last_modified_time": object["LastModifiedTime"],
			"payload":            object["Payload"],
			"schedule_id":        object["ScheduleId"],
			"id":                 fmt.Sprint(object["ScheduleName"]),
			"schedule_name":      object["ScheduleName"],
		}
		ids = append(ids, fmt.Sprint(object["ScheduleName"]))
		names = append(names, object["ScheduleName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("schedules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
