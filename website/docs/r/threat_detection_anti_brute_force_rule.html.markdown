---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_anti_brute_force_rule"
description: |-
  Provides a Alicloud Threat Detection Anti Brute Force Rule resource.
---

# alicloud_threat_detection_anti_brute_force_rule

Provides a Threat Detection Anti Brute Force Rule resource.

Anti-brute force cracking rules.

For information about Threat Detection Anti Brute Force Rule and how to use it, see [What is Anti Brute Force Rule](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createantibruteforcerule).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_anti_brute_force_rule" "default" {
  anti_brute_force_rule_name = "apispec_example"
  forbidden_time             = 360
  uuid_list                  = ["032b618f-b220-4a0d-bd37-fbdc6ef58b6a"]
  fail_count                 = 80
  span                       = 10
}
```

## Argument Reference

The following arguments are supported:
* `anti_brute_force_rule_name` - (Required) The name of the defense rule.
* `default_rule` - (Optional) Specifies whether to set the defense rule as the default rule. Valid values:

  -`true`: yes
  -`false`: no

-> **NOTE:**   If no defense rule is created for a server, the default rule is applied to the server.

* `fail_count` - (Required, Int) FailCount
* `forbidden_time` - (Required, Int) The period of time during which logons from an account are not allowed. Unit: minutes. Valid values:

  - `5`: 5 minutes
  - `15`: 15 minutes
  - `30`: 30 minutes
  - `60`: 1 hour
  - `120`: 2 hours
  - `360`: 6 hours
  - `720`: 12 hours
  - `1440`: 24 hours
  - `10080`: 7 days
  - `52560000`: permanent
* `span` - (Required, Int) The maximum period of time during which failed logon attempts from an account can occur. Unit: minutes. Valid values:

  -`1`
  -`2`
  -`5`
  -`10`
  -`15`

-> **NOTE:**   To configure a defense rule, you must specify the Span, FailCount, and ForbiddenTime parameters. If the number of failed logon attempts from an account within the minutes specified by Span exceeds the value specified by FailCount, the account cannot be used for logons within the minutes specified by ForbiddenTime.

* `uuid_list` - (Required, Set) The UUIDs of the servers to which you want to apply the defense rule.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anti Brute Force Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Anti Brute Force Rule.
* `update` - (Defaults to 5 mins) Used when update the Anti Brute Force Rule.

## Import

Threat Detection Anti Brute Force Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_anti_brute_force_rule.example <id>
```