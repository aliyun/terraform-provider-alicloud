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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_anti_brute_force_rule&exampleId=88eb579a-c215-af50-3c55-a95206a550031c805008&activeTab=example&spm=docs.r.threat_detection_anti_brute_force_rule.0.88eb579ac2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `protocol_type` - (Optional, Set, Available since v1.255.0) The types of protocols supported for interception by the brute force attack rule creation. See [`protocol_type`](#protocol_type) below.
* `span` - (Required, Int) The maximum period of time during which failed logon attempts from an account can occur. Unit: minutes. Valid values:

  -`1`
  -`2`
  -`5`
  -`10`
  -`15`

-> **NOTE:**   To configure a defense rule, you must specify the Span, FailCount, and ForbiddenTime parameters. If the number of failed logon attempts from an account within the minutes specified by Span exceeds the value specified by FailCount, the account cannot be used for logons within the minutes specified by ForbiddenTime.

* `uuid_list` - (Required, Set) The UUIDs of the servers to which you want to apply the defense rule.

### `protocol_type`

The protocol_type supports the following:
* `rdp` - (Optional) Whether to enable RDP interception. Default value: `on`. Valid values: `on`, `off`.
* `sql_server` - (Optional) Whether to enable the SqlServer interception method. Default value: `off`. Valid values: `on`, `off`.
* `ssh` - (Optional) Whether to enable SSH interception. Default value: `on`. Valid values: `on`, `off`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anti Brute Force Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Anti Brute Force Rule.
* `update` - (Defaults to 5 mins) Used when update the Anti Brute Force Rule.

## Import

Threat Detection Anti Brute Force Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_anti_brute_force_rule.example <id>
```
