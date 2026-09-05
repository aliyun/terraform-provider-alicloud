package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionIncidents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionIncidentsRead,
		Schema: map[string]*schema.Schema{
			"incident_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incident_status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"threat_level": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"incident_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relate_asset_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relate_entity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_field_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"role_for": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incidents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threat_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threat_score": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_aggregation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attck_tactics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"relate_user_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"relate_data_source_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"relate_alert_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relate_asset_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relate_entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"relate_asset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"response_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"detection_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lang": {
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

func dataSourceAlicloudThreatDetectionIncidentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{}
	if v, ok := d.GetOk("lang"); ok && v.(string) != "" {
		request["Lang"] = v.(string)
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v.(int)
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v.(int)
	}
	if v, ok := d.GetOk("incident_status"); ok {
		request["IncidentStatus"] = v.(int)
	}
	if v, ok := d.GetOk("incident_name"); ok && v.(string) != "" {
		request["IncidentName"] = v.(string)
	}
	if v, ok := d.GetOk("alert_uuid"); ok && v.(string) != "" {
		request["AlertUuid"] = v.(string)
	}
	if v, ok := d.GetOk("relate_asset_id"); ok && v.(string) != "" {
		request["RelateAssetId"] = v.(string)
	}
	if v, ok := d.GetOk("relate_entity_id"); ok && v.(string) != "" {
		request["RelateEntityId"] = v.(string)
	}
	if v, ok := d.GetOk("order_field_name"); ok && v.(string) != "" {
		request["OrderFieldName"] = v.(string)
	}
	if v, ok := d.GetOk("order_direction"); ok && v.(string) != "" {
		request["OrderDirection"] = v.(string)
	}
	if v, ok := d.GetOk("role_type"); ok {
		request["RoleType"] = v.(int)
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v.(int)
	}
	if v, ok := d.GetOk("threat_level"); ok {
		levels := make([]interface{}, 0)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			levels = append(levels, vv.(string))
		}
		if len(levels) > 0 {
			request["ThreatLevel"] = levels
		}
	}
	if v, ok := d.GetOk("owner"); ok {
		owners := make([]interface{}, 0)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			owners = append(owners, vv.(string))
		}
		if len(owners) > 0 {
			request["Owners"] = owners
		}
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

	var objects []interface{}
	var response map[string]interface{}
	var err error
	action := "ListIncidents"
	wait := incrementalWait(3*time.Second, 3*time.Second)

	pageNumber := d.Get("page_number").(int)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	pageSize := d.Get("page_size").(int)
	if pageSize <= 0 {
		pageSize = 10
	}

	for {
		request["PageNumber"] = pageNumber
		request["PageSize"] = pageSize

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_incidents", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Incidents", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Incidents", response)
		}
		result, _ := resp.([]interface{})

		if isPagingRequest(d) {
			objects = result
			break
		}

		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IncidentUuid"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		incidentUuid := fmt.Sprint(object["IncidentUuid"])
		mapping := map[string]interface{}{
			"id":                        incidentUuid,
			"incident_uuid":             incidentUuid,
			"incident_name":             object["IncidentName"],
			"incident_description":      object["IncidentDescription"],
			"incident_status":           object["IncidentStatus"],
			"owner":                     object["Owner"],
			"threat_level":              object["ThreatLevel"],
			"threat_score":              object["ThreatScore"],
			"incident_aggregation_type": object["IncidentAggregationType"],
			"incident_tags":             object["IncidentTags"],
			"incident_remark":           object["IncidentRemark"],
			"relate_alert_count":        object["RelateAlertCount"],
			"relate_asset_count":        object["RelateAssetCount"],
			"relate_entity_id":          object["RelateEntityId"],
			"relate_asset_id":           object["RelateAssetId"],
			"alert_uuid":                object["AlertUuid"],
			"role_type":                 object["RoleType"],
			"create_time":               object["CreateTime"],
			"update_time":               object["UpdateTime"],
			"start_time":                object["StartTime"],
			"end_time":                  object["EndTime"],
			"response_time":             object["ResponseTime"],
			"detection_rule_id":         object["DetectionRuleId"],
			"lang":                      object["Lang"],
			"region_id":                 object["RegionId"],
		}
		if v, ok := object["AttckTactics"].([]interface{}); ok {
			mapping["attck_tactics"] = v
		}
		if v, ok := object["RelateUserIds"].([]interface{}); ok {
			mapping["relate_user_ids"] = v
		}
		if v, ok := object["RelateDataSourceIds"].([]interface{}); ok {
			mapping["relate_data_source_ids"] = v
		}

		ids = append(ids, incidentUuid)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("incidents", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
