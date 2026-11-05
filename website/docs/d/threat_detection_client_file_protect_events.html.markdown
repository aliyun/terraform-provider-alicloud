---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_client_file_protect_events"
sidebar_current: "docs-alicloud-datasource-threat-detection-client-file-protect-events"
description: |-
  Provides a list of Threat Detection Client File Protect Events owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_client_file_protect_events

This data source provides the file protect client events of Cloud Security Center (Threat Detection) available to the user. File protect client events are system-generated security detections and are read-only. [ListFileProtectClientEvent](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-listfileprotectclientevent)

-> **NOTE:** Available since v2.3.0.

## Example Usage

```terraform
data "alicloud_threat_detection_client_file_protect_events" "default" {
  rule_name    = "demo-file-protect-rule"
  alert_levels = [2, 3]
}

output "first_event_id" {
  value = data.alicloud_threat_detection_client_file_protect_events.default.events.0.id
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Optional) The name of the file protect rule.
* `status` - (Optional) The status of the event.
* `alert_levels` - (Optional) A list of alert levels to filter. Valid values: `0` (no alert), `1` (reminder), `2` (suspicious) and `3` (high risk).
* `instance_id` - (Optional) The ID of the instance.
* `uuid` - (Optional) The UUID of the instance.
* `instance_name` - (Optional) The name of the instance.
* `internet_ip` - (Optional) The internet IP of the instance.
* `intranet_ip` - (Optional) The intranet IP of the instance.
* `start_time` - (Optional) The start time of the query range, as a Unix timestamp in milliseconds.
* `end_time` - (Optional) The end time of the query range, as a Unix timestamp in milliseconds.
* `operation` - (Optional) The file operation.
* `proc_path` - (Optional) The path of the process.
* `file_path` - (Optional) The path of the file.
* `ids` - (Optional) A list of event IDs. The filter is applied on the provider side after all pages are read.
* `page_size` - (Optional) The number of entries per page when querying the server. If unset, the provider uses the default page size of the product. The provider reads all pages regardless of this value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, derived from the configured filters so it stays stable across refreshes.
* `ids` - A list of event IDs that match the configured filters.
* `events` - A list of file protect client events. Each element contains the following attributes:
  * `id` - The ID of the event.
  * `uuid` - The UUID of the instance.
  * `rule_name` - The name of the file protect rule.
  * `status` - The status of the event.
  * `alert_level` - The alert level of the event.
  * `instance_name` - The name of the instance.
  * `internet_ip` - The internet IP of the instance.
  * `intranet_ip` - The intranet IP of the instance.
  * `file_path` - The path of the file.
  * `operation` - The file operation.
  * `remark` - The remark of the event.
  * `proc_permission` - The permission of the process.
  * `cmd_line` - The command line of the process.
  * `proc_path` - The path of the process.
  * `handle_time` - The time the event was handled, as a Unix timestamp in milliseconds.
  * `first_time` - The first time the event occurred, as a Unix timestamp in milliseconds.
  * `latest_time` - The latest time the event occurred, as a Unix timestamp in milliseconds.
  * `process_id` - The ID of the process.
  * `platform` - The platform of the instance.
  * `rule_action` - The action of the file protect rule.
  * `count` - The number of times the event occurred.
