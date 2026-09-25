package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionIncidentInvestigations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionIncidentInvestigationsRead,
		Schema: map[string]*schema.Schema{
			"incident_uuid": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"role_for": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"ids": {
				Optional: true,
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
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  100,
			},
			"investigations": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incident_investigation_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_investigation_start_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"incident_investigation_end_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"incident_investigation_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_investigation_summary": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_investigation_display_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_investigation_alert_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_investigation_conclusion": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"incident_uuid": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionIncidentInvestigationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOkExists("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("incident_uuid"); ok {
		request["IncidentUuid"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListIncidentInvestigations"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_incident_investigations", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.IncidentInvestigations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.IncidentInvestigations", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IncidentInvestigationId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"incident_investigation_id":         fmt.Sprint(object["IncidentInvestigationId"]),
			"incident_investigation_start_time": object["IncidentInvestigationStartTime"],
			"incident_investigation_end_time":   object["IncidentInvestigationEndTime"],
			"incident_investigation_status":     object["IncidentInvestigationStatus"],
			"incident_investigation_summary":    object["IncidentInvestigationSummary"],
			"incident_investigation_display_id": object["IncidentInvestigationDisplayId"],
			"incident_investigation_alert_name": object["IncidentInvestigationAlertName"],
			"incident_investigation_conclusion": object["IncidentInvestigationConclusion"],
			"incident_uuid":                     fmt.Sprint(object["IncidentUuid"]),
		}

		ids = append(ids, fmt.Sprint(object["IncidentInvestigationId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("investigations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
