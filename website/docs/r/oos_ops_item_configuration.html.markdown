---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_ops_item_configuration"
sidebar_current: "docs-alicloud-resource-oos-ops-item-configuration"
description: |-
  Provides a Alicloud OOS OpsItem Configuration resource.
---

# alicloud_oos_ops_item_configuration

Provides an OOS OpsItem Configuration resource.

OpsItem Configuration is a named configuration entry that stores operational data for an OpsItem. Once created, the configuration can be referenced and updated by its system-assigned `configuration_id`.

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_oos_ops_item_configuration" "example" {
  ops_item_configuration_name = "tf-example"
  data_source                 = jsonencode({ Type = "ALIYUN/CMS" })
}
```

With a resource group

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_ops_item_configuration" "example" {
  ops_item_configuration_name = "tf-example"
  data_source                 = jsonencode({ Type = "ALIYUN/CMS" })
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ops_item_configuration_name` - (Required) The name of the OpsItem configuration.
* `data_source` - (Required) The data source of the OpsItem configuration, JSON-encoded. `Type` is required and must be one of `ALIYUN/ROS`, `ALIYUN/ITSM`, `ALIYUN/ARMS`, `ALIYUN/CMS`, `ALIYUN/CIAM`, `ALIYUN/OOS`. Optional fields: `Arn`, `Path`, `ExcludePath`, `NodeId`, `Interval`, `FilterCondition`, `Arguments`.
* `resource_group_id` - (Optional) The ID of the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform. Its value is the same as `configuration_id`.
* `configuration_id` - The first-level ID of the OpsItem configuration, assigned by the system on creation.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.

## Import

OOS OpsItem Configuration can be imported using the `id`, e.g.

```shell
$ terraform import alicloud_oos_ops_item_configuration.example <configuration_id>
```
