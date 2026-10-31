---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group_instanceses"
sidebar_current: "docs-alicloud-datasource-cms-monitor-group-instanceses"
description: |-
  Provides a list of Cms Monitor Group Instanceses to the user.
---

# alicloud\_cms\_monitor\_group\_instanceses

This data source provides the Cms Monitor Group Instanceses of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.115.0.

-> **NOTE:** This data source is the legacy alias of `alicloud_cms_monitor_group_instances` and shares the same implementation.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_monitor_group" "example" {
  monitor_group_name = "tf-example-accmonitorgroup"
}

data "alicloud_cms_monitor_group_instanceses" "example" {
  ids = [alicloud_cms_monitor_group.example.id]
}

output "first_cms_monitor_group_instances" {
  value = data.alicloud_cms_monitor_group_instanceses.example.instances.0.instances
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Required) A list of Monitor Group Instances IDs.
* `keyword` - (Optional) The keyword.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Cms Monitor Group Instances. Each element contains the following attributes:
  * `instances` - Instance information added to the Cms Group.
    * `category` - The category of instance.
    * `instance_id` - The id of instance.
    * `instance_name` - The name of instance.
    * `region_id` - The region id of instance.
