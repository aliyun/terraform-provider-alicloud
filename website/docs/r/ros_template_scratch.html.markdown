---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_template_scratch"
sidebar_current: "docs-alicloud-resource-ros-template-scratch"
description: |-
  Provides a Alicloud ROS Template Scratch resource.
---

# alicloud\_ros\_template\_scratch

Provides a ROS Template Scratch resource.

For information about ROS Template Scratch and how to use it, see [What is Template Scratch](https://www.alibabacloud.com/help/zh/doc-detail/352074.html).

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ros_template_scratch" "example" {
  description           = "tf_testacc"
  template_scratch_type = "ResourceImport"
  preference_parameters {
    parameter_key   = "DeletionPolicy"
    parameter_value = "Retain"
  }
  source_resource_group {
    resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
    resource_type_filter = ["ALIYUN::ECS::VPC"]
  }
}

```

## Argument Reference

The following arguments are supported:

-> **NOTE:** One of `source_resource_group,source_resources,source_tag` must be specified.

* `description` - (Optional) The Description of the Template Scratch.
* `execution_mode` - (Optional) The execution mode. Valid Values: `Async` or `Sync`.
* `logical_id_strategy` - (Optional) Logical ID generation strategy. Valid Values: `LongTypePrefixAndIndexSuffix`, `LongTypePrefixAndHashSuffix` and `ShortTypePrefixAndHashSuffix`.
* `preference_parameters` - (Optional) Priority parameter. See the following `Block preference_parameters`.
* `source_tag` - (Optional) Source tag. See the following `Block source_tag`.
* `source_resource_group` - (Optional) Source resource grouping. See the following `Block source_resource_group`.
* `source_resources` - (Optional) Source resource. See the following `Block source_resources`.
* `template_scratch_type` - (Required, ForceNew) The type of the Template scan. Valid Values: `ResourceImport` or `ArchitectureReplication`.

#### Block source_resources

The source_resources supports the following: 

* `resource_id` - (Required) The ID of the Source Resource.
* `resource_type` - (Required) The type of the Source resource.

#### Block source_resource_group

The source_resource_group supports the following: 

* `resource_group_id` - (Required) The ID of the Source Resource Group.
* `resource_type_filter` - (Optional) Source resource type filter list. If the resource type list is specified, it means to scan the resources of the specified resource type and in the specified resource group; Otherwise, it means to scan all resources in the specified resource group. **NOTE:** A maximum of `20` resource type filter can be configured.

#### Block source_tag

The source_tag supports the following: 

* `resource_tags` - (Required) Source label. **NOTE:** A maximum of 10 source labels can be configured.
* `resource_type_filter` - (Optional) Source resource type filter list. If the resource type list is specified, it means to scan the resources of the specified resource type and in the specified resource group; Otherwise, it means to scan all resources in the specified resource group. **NOTE:** A maximum of `20` resource type filter can be configured.

#### Block preference_parameters

The preference_parameters supports the following: 

* `parameter_key` - (Required) Priority parameter key. For more information about values, see [supplementary instructions for request parameters](https://www.alibabacloud.com/help/zh/doc-detail/358846.html#h2-url-4).
* `parameter_value` - (Required) Priority parameter value. For more information about values, see [supplementary instructions for request parameters](https://www.alibabacloud.com/help/zh/doc-detail/358846.html#h2-url-4).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Template Scratch.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Template Scratch.
* `delete` - (Defaults to 1 mins) Used when delete the Template Scratch.
* `update` - (Defaults to 1 mins) Used when update the Template Scratch.

## Import

ROS Template Scratch can be imported using the id, e.g.

```
$ terraform import alicloud_ros_template_scratch.example <id>
```