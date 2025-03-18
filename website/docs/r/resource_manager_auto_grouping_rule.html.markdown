---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_auto_grouping_rule"
description: |-
  Provides a Alicloud Resource Manager Auto Grouping Rule resource.
---

# alicloud_resource_manager_auto_grouping_rule

Provides a Resource Manager Auto Grouping Rule resource.

Auto grouping rules of resource group.

For information about Resource Manager Auto Grouping Rule and how to use it, see [What is Auto Grouping Rule](https://www.alibabacloud.com/help/en/resource-management/resource-group/developer-reference/api-resourcemanager-2020-03-31-createautogroupingrule-rg).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_resource_manager_auto_grouping_rule" "default" {
  rule_contents {
    target_resource_group_condition = <<EOF
    {
        "children": [
      {
        "desired": "rg-aek*****sbvy",
        "featurePath": "$.resourceGroupId",
        "featureSource": "RESOURCE",
        "operator": "StringEquals"
      }
        ],
        "operator": "and"
    }
    EOF
    auto_grouping_scope_condition   = <<EOF
    {
        "children": [
      {
        "desired": "name_a",
        "featurePath": "$.resourceName",
        "featureSource": "RESOURCE",
        "operator": "StringEqualsAny"
      }
        ],
        "operator": "and"
    }
    EOF
  }
  rule_desc                = var.name
  rule_type                = "custom_condition"
  region_ids_scope         = "cn-hangzhou,cn-shanghai"
  resource_ids_scope       = "imock-xxxxxx"
  resource_group_ids_scope = "rg-aek22*****3sbvz"
  resource_types_scope     = "ecs.instance,vpc.vpc"
  rule_name                = var.name
}
```

## Argument Reference

The following arguments are supported:
* `exclude_region_ids_scope` - (Optional) The IDs of regions to be excluded. Separate multiple IDs with commas (,).
* `exclude_resource_group_ids_scope` - (Optional) The IDs of resource groups to be excluded. Separate multiple IDs with commas (,).
* `exclude_resource_ids_scope` - (Optional) The IDs of resources to be excluded. Separate multiple IDs with commas (,).
* `exclude_resource_types_scope` - (Optional) The resource types to be excluded. Separate multiple resource types with commas (,).
* `region_ids_scope` - (Optional) The IDs of regions. Separate multiple IDs with commas (,).
* `resource_group_ids_scope` - (Optional) The IDs of resource groups. Separate multiple IDs with commas (,).
* `resource_ids_scope` - (Optional) The IDs of resources. Separate multiple IDs with commas (,).
* `resource_types_scope` - (Optional) The resource types. Separate multiple resource types with commas (,).
* `rule_contents` - (Required, Set) The content records of the rule. See [`rule_contents`](#rule_contents) below.
* `rule_desc` - (Optional) The description of the rule.
* `rule_name` - (Required) The name of the rule.
* `rule_type` - (Required, ForceNew) The type of the rule. Valid values:
  - `custom_condition`: Custom transfer rule.
  - `associated_transfer`: Transfer rule for associated resources.

### `rule_contents`

The rule_contents supports the following:
* `auto_grouping_scope_condition` - (Optional) The condition for the range of resources to be automatically transferred.
* `target_resource_group_condition` - (Required) The condition for the destination resource group.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Auto Grouping Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Grouping Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Grouping Rule.
* `update` - (Defaults to 5 mins) Used when update the Auto Grouping Rule.

## Import

Resource Manager Auto Grouping Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_auto_grouping_rule.example <id>
```
