---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_incident"
sidebar_current: "docs-alicloud-datasource-threat_detection-incident"
description: |-
  Provides a Threat Detection Incident resource.
---

# alicloud_threat_detection_incident

This data source provides the Threat Detection Incident of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.236.0.

## Example Usage

```terraform
data "alicloud_threat_detection_incident" "default" {
  incident_uuid = "xxxx-xxxx-xxxx-xxxx"
}

output "incident_name" {
  value = data.alicloud_threat_detection_incident.default.incident_name
}
```

## Argument Reference

The following arguments are supported:

* `incident_uuid` - (Required) The UUID of the incident.
* `lang` - (Optional) The language type. Valid values: `zh`, `en`.
* `role_for` - (Optional) The user ID of other roles.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the incident, which is the same as `incident_uuid`.
* `incident_description` - The description of the incident.
* `attck_tactics` - The list of attack stage counts of the alerts associated with the incident.
* `owner` - The responsible person aliuid.
* `incident_tags` - The tags of the incident.
* `incident_status` - The status of the incident. Valid values: `0` (not processed), `1` (in processing), `5` (processing failed), `10` (processed).
* `role_type` - The operating perspective, global or individual.
* `incident_aggregation_type` - The aggregation type of the incident.
* `threat_level` - The threat level. Valid values: `5` (serious), `4` (high risk), `3` (medium), `2` (low risk), `1` (information).
* `threat_score` - The threat score of the incident, ranging from 0 to 100.
* `relate_user_ids` - The list of user IDs associated with the incident.
* `relate_alert_count` - The number of alerts associated with the incident.
* `relate_data_source_ids` - The list of associated data sources.
* `relate_asset_count` - The number of assets associated with the incident.
* `relate_entity_id` - The ID of the entity associated with the incident.
* `relate_asset_id` - The ID of the asset associated with the incident.
* `alert_uuid` - The alarm ID.
* `incident_name` - The name of the incident.
* `incident_remark` - The remark of the incident.
* `create_time` - The creation time of the incident.
* `update_time` - The update time of the incident.
* `start_time` - The start time of the incident.
* `end_time` - The end time of the incident.
* `response_time` - The response time of the incident.
* `detection_rule_id` - The ID of the detection rule.
* `region_id` - The region of the user.
