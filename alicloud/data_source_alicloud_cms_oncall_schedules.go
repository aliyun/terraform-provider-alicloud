// Package alicloud. This file is hand-written from the Cms 2024-03-30 OpenAPI definition.
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsOncallSchedules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsOncallSchedulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"oncall_schedule_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
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
						"oncall_schedule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oncall_schedule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shift_robot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rotations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"active_days": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"contacts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"rotation_end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rotation_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rotation_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"shift_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"shift_recurrence_frequency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_zone": {
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

func dataSourceAliCloudCmsOncallSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("oncall_schedule_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	action := "/oncallSchedules"
	var response map[string]interface{}
	query := make(map[string]*string)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	if v, ok := d.GetOk("source"); ok {
		query["source"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			var err error
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.oncallSchedules[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			itemId := fmt.Sprint(item["oncallScheduleId"])
			if nameRegex != nil {
				if !nameRegex.MatchString(fmt.Sprint(item["oncallScheduleName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[itemId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0, len(objects))
	s := make([]map[string]interface{}, 0, len(objects))
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["oncall_schedule_id"] = fmt.Sprint(objectRaw["oncallScheduleId"])
		mapping["oncall_schedule_name"] = objectRaw["oncallScheduleName"]
		mapping["shift_robot_id"] = objectRaw["shiftRobotId"]
		mapping["rotations"] = flattenOncallScheduleRotations(objectRaw["rotations"])

		ids = append(ids, mapping["oncall_schedule_id"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("schedules", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}

	return nil
}
