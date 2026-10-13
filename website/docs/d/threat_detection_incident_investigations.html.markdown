---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_incident_investigations"
sidebar_current: "docs-alicloud-datasource-threat_detection-incident-investigations"
description: |-
  Provides a list of Threat Detection Incident Investigations owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_incident_investigations

This data source provides Threat Detection Incident Investigation available to the user.

-> **NOTE:** Available since v1.245.0.

## Example Usage

```terraform
data "alicloud_threat_detection_incident_investigations" "default" {
  incident_uuid = "xxxxx-xxxxx-xxxxx-xxxxx"
}

output "first_investigation_id" {
  value = data.alicloud_threat_detection_incident_investigations.default.investigations.0.incident_investigation_id
}
```

## Argument Reference

The following arguments are supported:
* `incident_uuid` - (Optional, ForceNew) The UUID of the incident to which the investigations belong.
* `lang` - (Optional, ForceNew) The language type. Valid values: `zh`, `en`.
* `role_for` - (Optional, ForceNew) The ID of the other role user.
* `ids` - (Optional, ForceNew, Computed) A list of Incident Investigation IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `page_number` - (Optional, ForceNew) The page number of the list. Default value: `1`.
* `page_size` - (Optional, ForceNew) The number of items per page. Maximum value: `100`. Default value: `100`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Incident Investigation IDs.
* `investigations` - A list of Incident Investigation Entries. Each element contains the following attributes:
  * `incident_investigation_id` - The ID of the incident investigation.
  * `incident_investigation_start_time` - The start time of the incident investigation.
  * `incident_investigation_end_time` - The end time of the incident investigation.
  * `incident_investigation_status` - The status of the incident investigation.
  * `incident_investigation_summary` - The summary of the incident investigation.
  * `incident_investigation_display_id` - The display ID of the incident investigation.
  * `incident_investigation_alert_name` - The name of the alert that triggers the investigation.
  * `incident_investigation_conclusion` - The conclusion of the incident investigation.
  * `incident_uuid` - The UUID of the incident.
