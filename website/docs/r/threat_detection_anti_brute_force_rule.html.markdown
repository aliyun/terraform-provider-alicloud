---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_anti_brute_force_rule"
sidebar_current: "docs-alicloud-resource-threat-detection-anti-brute-force-rule"
description: |-
  Provides a Alicloud Threat Detection Anti Brute Force Rule resource.
---

# alicloud_threat_detection_anti_brute_force_rule

Provides a Threat Detection Anti Brute Force Rule resource.

For information about Threat Detection Anti Brute Force Rule and how to use it, see [What is Anti Brute Force Rule](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-createantibruteforcerule).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_anti_brute_force_rule" "default" {
  anti_brute_force_rule_name = "apispec_test"
  forbidden_time             = 360
  uuid_list                  = ["032b618f-b220-4a0d-bd37-fbdc6ef58b6a"]
  fail_count                 = 80
  span                       = 10
}
```

## Argument Reference

The following arguments are supported:
* `anti_brute_force_rule_name` - (Required) The name of the defense rule.
* `default_rule` - (Optional) Specifies whether to set the defense rule as the default rule.
* `fail_count` - (Required) The threshold for the number of failed user logins when the brute-force defense rule takes effect.
* `forbidden_time` - (Required) The period of time during which logons from an account are not allowed. Unit: minutes.
* `span` - (Required) The period of time during which logon failures from an account are measured. Unit: minutes. If Span is set to 10, the defense rule takes effect when the logon failures measured within 10 minutes reaches the specified threshold. The IP address of attackers cannot be used to log on to the server in the specified period of time.
* `uuid_list` - (Required) An array consisting of the UUIDs of servers to which the defense rule is applied.**The binding status must be Enterprise Edition.**

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the defense rule.
* `anti_brute_force_rule_id` - The ID of the defense rule.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anti Brute Force Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Anti Brute Force Rule.
* `update` - (Defaults to 5 mins) Used when update the Anti Brute Force Rule.

## Import

Threat Detection Anti Brute Force Rule can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_anti_brute_force_rule.example <id>
```