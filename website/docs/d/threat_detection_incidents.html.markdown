---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_incidents"
sidebar_current: "docs-alicloud-datasource-threat_detection-incidents"
description: |-
  Provides a list of Threat Detection Incidents owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_incidents

This data source provides a list of Threat Detection Incidents available to the current Alibaba Cloud user.

-> **NOTE:** Available since v1.236.0.

## Example Usage

```terraform
data "alicloud_threat_detection_incidents" "default" {
  incident_status = 0
  threat_level    = ["5", "4"]
  page_size       = 10
}

output "first_incident_uuid" {
  value = data.alicloud_threat_detection_incidents.default.incidents.0.incident_uuid
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Incident UUIDs.
* `incident_uuid` - (Optional) The UUID of the incident.
* `incident_status` - (Optional) The status of the incident. Valid values: `0` (not processed), `1` (in processing), `5` (processing failed), `10` (processed).
* `threat_level` - (Optional) The threat level. Valid values: `5` (serious), `4` (high risk), `3` (medium), `2` (low risk), `1` (information).
* `owner` - (Optional) The list of responsible person aliuids.
* `start_time` - (Optional) The start timestamp, accurate to milliseconds (ms).
* `end_time` - (Optional) The end timestamp, accurate to milliseconds (ms).
* `incident_name` - (Optional) The name of the incident.
* `alert_uuid` - (Optional) The alarm ID.
* `relate_asset_id` - (Optional) The ID of the asset associated with the incident.
* `relate_entity_id` - (Optional) The ID of the entity associated with the incident.
* `order_field_name` - (Optional) The sort field.
* `order_direction` - (Optional) The sort type.
* `role_type` - (Optional) The operating perspective, global or individual.
* `role_for` - (Optional) The user ID of other roles.
* `lang` - (Optional) The language type. Valid values: `zh`, `en`.
* `page_number` - (Optional) The page number, default to `1`.
* `page_size` - (Optional) The number of entries per page, default to `10`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Incident UUIDs.
* `incidents` - A list of Incident Entries. Each element contains the following attributes:
  * `id` - The ID of the incident, which is the same as `incident_uuid`.
  * `incident_uuid` - The UUID of the incident.
  * `incident_name` - The name of the incident.
  * `incident_description` - The description of the incident.
  * `incident_status` - The status of the incident.
  * `owner` - The responsible person aliuid.
  * `threat_level` - The threat level.
  * `threat_score` - The threat score of the incident.
  * `incident_aggregation_type` - The aggregation type of the incident.
  * `incident_tags` - The tags of the incident.
  * `incident_remark` - The remark of the incident.
  * `attck_tactics` - The list of attack stage counts of the alerts associated with the incident.
  * `relate_user_ids` - The list of user IDs associated with the incident.
  * `relate_data_source_ids` - The list of associated data sources.
  * `relate_alert_count` - The number of alerts associated with the incident.
  * `relate_asset_count` - The number of assets associated with the incident.
  * `relate_entity_id` - The ID of the entity associated with the incident.
  * `relate_asset_id` - The ID of the asset associated with the incident.
  * `alert_uuid` - The alarm ID.
  * `role_type` - The operating perspective.
  * `create_time` - The creation time of the incident.
  * `update_time` - The update time of the incident.
  * `start_time` - The start time of the incident.
  * `end_time` - The end time of the incident.
  * `response_time` - The response time of the incident.
  * `detection_rule_id` - The ID of the detection rule.
  * `lang` - The language type.
  * `region_id` - The region of the user.
