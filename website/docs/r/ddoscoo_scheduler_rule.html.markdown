---
subcategory: "DdosCoo"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_scheduler_rule"
sidebar_current: "docs-alicloud-resource-ddoscoo-scheduler-rule"
description: |-
  Provides a Alicloud DdosCoo Scheduler Rule resource.
---

# alicloud\_ddoscoo\_scheduler\_rule

Provides a DdosCoo Scheduler Rule resource. For information about DdosCoo Scheduler Rule and how to use it, see[What is DdosCoo Scheduler Rule](https://www.alibabacloud.com/help/en/doc-detail/157481.htm).

-> **NOTE:** Available in 1.86.0+

## Example Usage

Basic Usage

```
resource "alicloud_ddoscoo_scheduler_rule" "example" {
  rule_name = "tf-testacc7929727******"
  rule_type =3
  rules{
    priority = "100"
    region_id = "cn-hangzhou"
    type = "A"
    value="170.33.2.125"
    value_type = "3"
  }
  rules{
    priority = "50"
    region_id = "cn-hangzhou"
    type = "A"
    value= "170.33.14.193"
    value_type = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, ForceNew) The name of the rule.
* `rule_type` - (Required) The rule type. Valid values:
    `2`: tiered protection.
    `3`: globalization acceleration.
    `6`: Cloud product interaction.
* `rules` - (Required) The details of the common filter interaction rule, expressed as a JSON string. The structure is as follows:
    `Type`: String type, required, the address format of the linkage resource. Valid values:
        `A`: IP address.
        `CNAME`: Domain name.
    `Value`: String type, required, link address of resource.
    `Priority`: the priority of the rule. This parameter is required and of Integer type. Valid values: 0~100 the larger the value, the higher the priority.
    `ValueType`: Required. The type of the linked resource. It is an Integer. Valid values:
        `1`: Anti-DDoS Pro.
        `2`: (Tiered protection) cloud resource IP.
        `3`: (sea acceleration) MCA IP address.
        `6`: (Cloud product linkage) cloud resource IP.
    `RegionId`: String type, optional (Required when ValueType is 2) the ID of the region.
* `resource_group_id` - (Optional) The ID of the resource group to which the anti-DDoS pro instance belongs in resource management. By default, no value is specified, indicating that the domains in the default resource group are listed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of scheduler rule. The value is `<rule_name>`.
* `cname` - The cname is the traffic scheduler corresponding to rules.


## Import

DdosCoo Scheduler Rule can be imported using the id or the rule name, e.g.

```
$ terraform import alicloud_ddoscoo_scheduler_rule.example fbb20dc77e8fc******
```
