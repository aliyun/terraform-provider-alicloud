---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group_instances"
sidebar_current: "docs-alicloud-resource-cms-monitor-group-instances"
description: |-
  Provides a Alicloud Cloud Monitor Service Monitor Group Instances resource.
---

# alicloud\_cms\_monitor\_group\_instances

Provides a Cloud Monitor Service Monitor Group Instances resource.

For information about Cloud Monitor Service Monitor Group Instances and how to use it, see [What is Monitor Group Instances](https://www.alibabacloud.com/help/en/doc-detail/115031.htm).

-> **NOTE:** Available in v1.115.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testacc-vpcname"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = "tf-testaccmonitorgroup"
}

resource "alicloud_cms_monitor_group_instances" "example" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id   = alicloud_vpc.default.id
    instance_name = "tf-testacc-vpcname"
    region_id     = "cn-hangzhou"
    category      = "vpc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The id of Cms Group.
* `instances` - (Required) Instance information added to the Cms Group.

#### Block instances

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

```
$ terraform import alicloud_cms_monitor_group_instances.example <group_id>
```
