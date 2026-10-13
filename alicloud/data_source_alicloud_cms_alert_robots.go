package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAlertRobots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertRobotsRead,
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"DING", "WEIXIN", "FEISHU", "SLACK", "TEAMS", "DING_COOL_APP"}, false),
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"robots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_robot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_robot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"digital_employee_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"robot_sign_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCmsAlertRobotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/robots"

	var allRobots []map[string]interface{}
	pageNumber := 1
	pageSize := 100

	for {
		query := make(map[string]*string)
		query["pageNumber"] = StringPointer(fmt.Sprintf("%d", pageNumber))
		query["pageSize"] = StringPointer(fmt.Sprintf("%d", pageSize))

		if v, ok := d.GetOk("type"); ok {
			types, _ := json.Marshal([]string{v.(string)})
			query["types"] = StringPointer(string(types))
		}
		if v, ok := d.GetOk("workspace"); ok {
			query["workspace"] = StringPointer(v.(string))
		}
		if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
			idList := make([]string, len(v.([]interface{})))
			for i, id := range v.([]interface{}) {
				idList[i] = id.(string)
			}
			robotIds, _ := json.Marshal(idList)
			query["robotIds"] = StringPointer(string(robotIds))
		}

		var response map[string]interface{}
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_robots", action, AlibabaCloudSdkGoERROR)
		}

		robots, ok := response["robots"].([]interface{})
		if !ok {
			break
		}
		for _, r := range robots {
			if robot, ok := r.(map[string]interface{}); ok {
				allRobots = append(allRobots, robot)
			}
		}

		total := 0
		if v, ok := response["total"].(float64); ok {
			total = int(v)
		}
		if len(allRobots) >= total || len(robots) == 0 {
			break
		}
		pageNumber++
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(fmt.Errorf("name_regex format error: %s", err))
		}
		nameRegex = r
	}

	var filterType string
	if v, ok := d.GetOk("type"); ok {
		filterType = v.(string)
	}

	var filterIds []string
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			filterIds = append(filterIds, id.(string))
		}
	}

	var matchedRobots []map[string]interface{}
	var ids []string
	for _, robot := range allRobots {
		name, _ := robot["name"].(string)
		robotType, _ := robot["type"].(string)
		robotId, _ := robot["robotId"].(string)

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		if filterType != "" && robotType != filterType {
			continue
		}
		if len(filterIds) > 0 {
			found := false
			for _, id := range filterIds {
				if id == robotId {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		matchedRobots = append(matchedRobots, robot)
		ids = append(ids, robotId)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)

	var robotMaps []map[string]interface{}
	for _, robot := range matchedRobots {
		// Defensive nil->"" conversion: the List API may omit fields (e.g.
		// robotSignKey is nil when the robot was created without a sign key,
		// webhook is write-only and never returned). A nil interface{} value
		// inside the nested map, combined with a Sensitive schema flag, can
		// prevent d.Set from persisting the data source to state under the
		// legacy SDK; coerce every field to a concrete string here.
		alertRobotId, _ := robot["robotId"].(string)
		alertRobotName, _ := robot["name"].(string)
		robotType, _ := robot["type"].(string)
		lang, _ := robot["lang"].(string)
		webhook, _ := robot["webhook"].(string)
		workspace, _ := robot["workspace"].(string)
		digitalEmployeeName, _ := robot["digitalEmployeeName"].(string)
		robotSignKey, _ := robot["robotSignKey"].(string)

		robotMaps = append(robotMaps, map[string]interface{}{
			"alert_robot_id":        alertRobotId,
			"alert_robot_name":      alertRobotName,
			"type":                  robotType,
			"lang":                  lang,
			"url":                   webhook,
			"workspace":             workspace,
			"digital_employee_name": digitalEmployeeName,
			"robot_sign_key":        robotSignKey,
		})
	}

	if err := d.Set("robots", robotMaps); err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok {
		if err := writeToFile(v.(string), robotMaps); err != nil {
			log.Printf("[WARN] write output file failed: %s", err)
		}
	}

	return nil
}
