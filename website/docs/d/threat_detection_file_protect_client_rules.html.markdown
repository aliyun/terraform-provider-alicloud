---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_file_protect_client_rules"
sidebar_current: "docs-alicloud-datasource-threat-detection-file-protect-client-rules"
description: |-
  Provides a list of Threat Detection File Protect Client Rule owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_file_protect_client_rules

This data source provides Threat Detection File Protect Client Rule available to the user.[What is File Protect Client Rule](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateFileProtectClientRule)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_threat_detection_file_protect_client_rule" "default" {
  rule_name   = var.name
  rule_action = "pass"
  status      = 0
  alert_level = 1
  platform    = "linux"
  file_paths  = ["/opt/a", "/tmp/d"]
  file_ops    = ["READ", "WRITE"]
  proc_paths  = ["/usr/bin/java"]
}

data "alicloud_threat_detection_file_protect_client_rules" "default" {
  ids = [alicloud_threat_detection_file_protect_client_rule.default.id]
}

output "file_protect_client_rule_id" {
  value = data.alicloud_threat_detection_file_protect_client_rules.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:
* `alert_level` - (Optional, Int) Filter rules by alert level. Valid values: `0` (no alert), `1` (info), `2` (suspicious), `3` (critical).
* `platform` - (Optional) Filter rules by operating system type. Valid values: `windows`, `linux`.
* `rule_action` - (Optional) Filter rules by action. Valid values: `pass`, `monitor`, `block`.
* `rule_name` - (Optional) Filter rules by rule name.
* `ids` - (Optional, Computed) A list of File Protect Client Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of File Protect Client Rule IDs.
* `rules` - A list of File Protect Client Rule Entries. Each element contains the following attributes:
  * `id` - The ID of the rule.
  * `alert_level` - The alert level of the rule. `0`: no alert, `1`: info, `2`: suspicious, `3`: critical.
  * `exclude_users` - The list of users to exclude.
  * `file_ops` - The operations that are performed on the files.
  * `file_paths` - The paths to the files that are monitored.
  * `file_types` - The list of file types to monitor.
  * `platform` - The type of the operating system.
  * `proc_paths` - The paths to the processes that are monitored.
  * `rule_action` - The action of the rule.
  * `rule_name` - The name of the rule.
  * `status` - The status of the rule. `0`: disabled, `1`: enabled.
  * `switch_id` - The switch ID of the rule.
