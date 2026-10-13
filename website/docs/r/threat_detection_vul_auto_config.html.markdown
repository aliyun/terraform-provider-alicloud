---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vul_auto_config"
sidebar_current: "docs-alicloud-resource-threat-detection-vul-auto-config"
description: |-
  Provides a Alicloud Threat Detection Vul Auto Config resource.
---

# alicloud_threat_detection_vul_auto_config

Provides a Threat Detection Vul Auto Config resource.

Vulnerability automatic repair configuration. You can use this resource to manage the automatic fix configuration for vulnerabilities detected by Security Center (Sas).

For information about Threat Detection Vul Auto Config and how to use it, see [AddOrUpdateAutoFixConfig](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-addorupdateautofixconfig).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_vul_auto_config" "default" {
  type              = "vul"
  start_time        = 1700000000
  all_uuid          = 1
  need_snapshot     = 0
  enable            = 1
  period_unit       = "day"
  necessity         = "asap"
  target_start_time = 0
  target_end_time   = 23
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The task type of the automatic vulnerability repair configuration.
* `start_time` - (Required) The task start time. Unix timestamp in seconds.
* `all_uuid` - (Required) Whether to select all servers. Valid values: `0`, `1`.
* `need_snapshot` - (Required) Whether to take a snapshot before the fix. Valid values: `0`, `1`.
* `enable` - (Required) The status of the rule. Valid values: `0` (disabled), `1` (enabled).
* `period_unit` - (Optional) The interval unit of the task.
* `necessity` - (Optional) The vulnerability severity level to be automatically repaired.
* `target_start_time` - (Optional) The start hour of the daily time range.
* `target_end_time` - (Optional) The end hour of the daily time range.
* `snapshot_name` - (Optional) The name of the snapshot.
* `snapshot_time` - (Optional) The retention period of the snapshot.
* `rules` - (Optional) The list of automatic repair rules.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vul Auto Config. The value is the `config_id`.
* `config_id` - The config ID of the automatic vulnerability repair configuration.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 3 mins) Used when create the Vul Auto Config.
* `update` - (Defaults to 3 mins) Used when update the Vul Auto Config.
* `delete` - (Defaults to 3 mins) Used when delete the Vul Auto Config.

## Import

Threat Detection Vul Auto Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_vul_auto_config.example <config_id>
```

-> **NOTE:** The ThreatDetection VulAutoConfig API does not provide a delete operation. `AddOrUpdateAutoFixConfig` is an upsert and there is no `DeleteAutoFixConfig` API. Destroying the resource via Terraform removes it from state only; the cloud-side configuration is preserved. To disable the rule, set `enable = 0`.
