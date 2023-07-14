---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group_instances"
sidebar_current: "docs-alicloud-resource-cms-monitor-group-instances"
description: |-
  Provides a Alicloud Cloud Monitor Service Monitor Group Instances resource.
---

# alicloud_cms_monitor_group_instances

Provides a Cloud Monitor Service Monitor Group Instances resource.

For information about Cloud Monitor Service Monitor Group Instances and how to use it, see [What is Monitor Group Instances](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createmonitorgroupinstances).

-> **NOTE:** Available since v1.115.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = var.name
}
data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_cms_monitor_group_instances" "example" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id   = alicloud_vpc.default.id
    instance_name = var.name
    region_id     = data.alicloud_regions.default.regions.0.id
    category      = "vpc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The id of Cms Group.
* `instances` - (Required) Instance information added to the Cms Group. See [`instances`](#instances) below. 

### `instances`

The instances supports the following: 

* `category` - (Required) The category of instance.
* `instance_id` - (Required) The id of instance.
* `instance_name` - (Required) The name of instance.
* `region_id` - (Required) The region id of instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitor Group Instances. Value as `group_id`.

## Import

Cloud Monitor Service Monitor Group Instances can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_monitor_group_instances.example <group_id>
```
