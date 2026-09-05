package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionIncident() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionIncidentRead,
		Schema: map[string]*schema.Schema{
			"incident_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_for": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"incident_description": {
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
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incident_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incident_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"incident_aggregation_type": {
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
			"relate_user_ids": {
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
			"relate_data_source_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"incident_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incident_remark": {
				Type:     schema.TypeString,
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
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionIncidentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetIncident"
	request := map[string]interface{}{
		"IncidentUuid": d.Get("incident_uuid").(string),
	}
	if v, ok := d.GetOk("lang"); ok && v.(string) != "" {
		request["Lang"] = v.(string)
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v.(int)
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_incident", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Incident", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Incident", response)
	}
	object, ok := resp.(map[string]interface{})
	if !ok || object == nil {
		d.SetId("")
		return nil
	}

	incidentUuid := fmt.Sprint(object["IncidentUuid"])
	if incidentUuid == "" {
		d.SetId("")
		return nil
	}
	d.SetId(incidentUuid)

	d.Set("incident_description", object["IncidentDescription"])
	d.Set("owner", object["Owner"])
	d.Set("incident_tags", object["IncidentTags"])
	d.Set("incident_status", object["IncidentStatus"])
	d.Set("role_type", object["RoleType"])
	d.Set("incident_aggregation_type", object["IncidentAggregationType"])
	d.Set("threat_level", object["ThreatLevel"])
	d.Set("threat_score", object["ThreatScore"])
	d.Set("relate_alert_count", object["RelateAlertCount"])
	d.Set("relate_asset_count", object["RelateAssetCount"])
	d.Set("relate_entity_id", object["RelateEntityId"])
	d.Set("relate_asset_id", object["RelateAssetId"])
	d.Set("alert_uuid", object["AlertUuid"])
	d.Set("incident_name", object["IncidentName"])
	d.Set("incident_remark", object["IncidentRemark"])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])
	d.Set("start_time", object["StartTime"])
	d.Set("end_time", object["EndTime"])
	d.Set("response_time", object["ResponseTime"])
	d.Set("detection_rule_id", object["DetectionRuleId"])
	d.Set("region_id", object["RegionId"])
	d.Set("lang", object["Lang"])

	if v, ok := object["AttckTactics"].([]interface{}); ok {
		d.Set("attck_tactics", v)
	}
	if v, ok := object["RelateUserIds"].([]interface{}); ok {
		d.Set("relate_user_ids", v)
	}
	if v, ok := object["RelateDataSourceIds"].([]interface{}); ok {
		d.Set("relate_data_source_ids", v)
	}

	return nil
}
