// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec Mission definition.
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsMissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMissionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digital_employee_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"missions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"digital_employee_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"variables": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notification_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"workspace": {
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

func dataSourceAlicloudCmsMissionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
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

	action := "/missions"
	var response map[string]interface{}
	query := make(map[string]*string)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	if v, ok := d.GetOk("digital_employee_name"); ok {
		query["digitalEmployeeName"] = StringPointer(v.(string))
	}
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
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

		if raw, ok := response["missions"]; ok && raw != nil {
			if list, ok := raw.([]interface{}); ok {
				for _, item := range list {
					m, ok := item.(map[string]interface{})
					if !ok {
						continue
					}
					missionName := fmt.Sprint(m["name"])
					if nameRegex != nil && !nameRegex.MatchString(missionName) {
						continue
					}
					if len(idsMap) > 0 {
						if _, ok := idsMap[missionName]; !ok {
							continue
						}
					}
					objects = append(objects, m)
				}
			}
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
		mapping["name"] = objectRaw["name"]
		mapping["display_name"] = objectRaw["displayName"]
		mapping["digital_employee_name"] = objectRaw["digitalEmployeeName"]
		mapping["enabled"] = objectRaw["enabled"]
		mapping["create_time"] = objectRaw["createTime"]
		mapping["update_time"] = objectRaw["updateTime"]
		if v, ok := objectRaw["variables"]; ok && v != nil {
			mapping["variables"] = flattenCmsMissionMap(v)
		}
		if v, ok := objectRaw["notificationPolicy"]; ok && v != nil {
			mapping["notification_policy"] = flattenCmsMissionNotificationPolicy(v)
		}
		ids = append(ids, fmt.Sprint(mapping["name"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("missions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
