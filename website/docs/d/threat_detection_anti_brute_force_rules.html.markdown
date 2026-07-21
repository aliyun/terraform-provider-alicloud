---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_anti_brute_force_rules"
sidebar_current: "docs-alicloud-datasource-threat-detection-anti-brute-force-rules"
description: |-
  Provides a list of Threat Detection Anti-Brute Force Rule owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_anti_brute_force_rules

This data source provides Threat Detection Anti Brute Force Rule available to the user.[What is Anti Brute Force Rule](https://www.alibabacloud.com/help/en/security-center/latest/api-sas-2018-12-03-createantibruteforcerule)

-> **NOTE:** Available since v1.195.0.

## Example Usage

```terraform
variable "name" {
  default = "example_value"
}

resource "alicloud_threat_detection_anti_brute_force_rule" "default" {
  anti_brute_force_rule_name = var.name
  forbidden_time             = 360
  uuid_list = [
  "7567806c-4ec5-4597-9543-7c9543381a13"]
  fail_count = 80
  span       = 10
}

data "alicloud_threat_detection_anti_brute_force_rules" "default" {
  ids        = ["${alicloud_threat_detection_anti_brute_force_rule.default.id}"]
  name_regex = alicloud_threat_detection_anti_brute_force_rule.default.name
}

output "alicloud_threat_detection_anti_brute_force_rule_example_id" {
  value = data.alicloud_threat_detection_anti_brute_force_rules.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Anti-Brute Force Rule IDs.
* `names` - (Optional, ForceNew) The name of the Anti-Brute Force Rule. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by the name of the defense rule.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Anti Brute Force Rule IDs.
* `names` - A list of name of Anti Brute Force Rules.
* `rules` - A list of Anti Brute Force Rule Entries. Each element contains the following attributes:
    * `id` - The ID of the defense rule.
    * `anti_brute_force_rule_id` - The ID of the defense rule.
    * `anti_brute_force_rule_name` - The name of the defense rule.
    * `default_rule` - Specifies whether to set the defense rule as the default rule.
    * `fail_count` - The threshold for the number of failed user logins when the brute-force defense rule takes effect.
    * `forbidden_time` - The period of time during which logons from an account are not allowed. Unit: minutes.
    * `span` - The period of time during which logon failures from an account are measured. Unit: minutes. If Span is set to 10, the defense rule takes effect when the logon failures measured within 10 minutes reaches the specified threshold. The IP address of attackers cannot be used to log on to the server in the specified period of time.
    * `uuid_list` - An array consisting of the UUIDs of servers to which the defense rule is applied.
