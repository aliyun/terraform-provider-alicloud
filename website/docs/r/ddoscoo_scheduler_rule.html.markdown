---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_scheduler_rule"
sidebar_current: "docs-alicloud-resource-ddoscoo-scheduler-rule"
description: |-
  Provides a Alicloud DdosCoo Scheduler Rule resource.
---

# alicloud_ddoscoo_scheduler_rule

Provides a DdosCoo Scheduler Rule resource. For information about DdosCoo Scheduler Rule and how to use it, see[What is DdosCoo Scheduler Rule](https://www.alibabacloud.com/help/en/ddos-protection/latest/api-ddoscoo-2020-01-01-createschedulerrule).

-> **NOTE:** Available since v1.86.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddoscoo_scheduler_rule&exampleId=67208afe-8b60-eef9-2dd4-a9595fdbb85b98c6cdb5&activeTab=example&spm=docs.r.ddoscoo_scheduler_rule.0.67208afe8b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
resource "alicloud_ddoscoo_scheduler_rule" "example" {
  rule_name = var.name
  rule_type = 3
  rules {
    priority   = 100
    region_id  = "cn-hangzhou"
    type       = "A"
    value      = "127.0.0.1"
    value_type = 3
  }
  rules {
    priority   = 50
    region_id  = "cn-hangzhou"
    type       = "A"
    value      = "127.0.0.0"
    value_type = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, ForceNew) The name of the rule.
* `param` - (Optional) The scheduling rule for the Global Accelerator instance that interacts with Anti-DDoS Pro or Anti-DDoS Premium.
* `resource_group_id` - (Optional) The ID of the resource group to which the anti-DDoS pro instance belongs in resource management. By default, no value is specified, indicating that the domains in the default resource group are listed.
* `rule_type` - (Required) The rule type. Valid values:
    `2`: tiered protection.
    `3`: globalization acceleration.
    `6`: Cloud product interaction.
* `rules` - (Required) The information about the scheduling rules. See [`rules`](#rules) below.

### `rules`

The rules supports the following:

* `type` - (Optional) The address type of the interaction resource. Valid values:
    `A`: IPv4 address.
    `CNAME`: CNAME record.
* `value` - (Optional) The address of the interaction resource.
* `priority` - (Optional) The priority of the rule.
* `value_type` - (Optional) Required. The type of the linked resource. It is an Integer. Valid values:
    `1`: The IP address of Anti-DDoS Pro or Anti-DDoS Premium
    `2`: the IP address of the interaction resource (in the tiered protection scenario)
    `3`: the IP address used to accelerate access (in the network acceleration scenario)
    `6` the IP address of the interaction resource (in the cloud service interaction scenario)
* `region_id` - (Optional) The region where the interaction resource that is used in the scheduling rule is deployed. **NOTE:** This parameter is returned only if the RuleType parameter is set to 2.
* `status` - (Optional) The status of the scheduling rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of scheduler rule. The value is `<rule_name>`.
* `cname` - The cname is the traffic scheduler corresponding to rules.

## Timeouts

-> **NOTE:** Available in 1.163.0+

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the scheduler rule.
* `update` - (Defaults to 1 mins) Used when update the scheduler rule.
* `delete` - (Defaults to 1 mins) Used when delete the scheduler rule.


## Import

DdosCoo Scheduler Rule can be imported using the id or the rule name, e.g.

```shell
$ terraform import alicloud_ddoscoo_scheduler_rule.example fbb20dc77e8fc******
```
