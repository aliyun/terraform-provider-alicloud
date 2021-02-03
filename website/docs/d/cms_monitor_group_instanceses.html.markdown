---
subcategory: "Cloud Monitor"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group_instanceses"
sidebar_current: "docs-alicloud-datasource-cms-monitor-group-instanceses"
description: |-
  Provides a list of Cms Monitor Group Instanceses to the user.
---

# alicloud\_cms\_monitor\_group\_instanceses

This data source provides the Cms Monitor Group Instanceses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.115.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_monitor_group_instanceses" "example" {
  ids      = ["example_value"]
}

output "first_cms_monitor_group_instances_id" {
  value = data.alicloud_cms_monitor_group_instanceses.example.instanceses.0.instances
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Required, ForceNew) A list of Monitor Group Instances IDs.
* `keyword` - (Optional, ForceNew) The keyword.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instanceses` - A list of Cms Monitor Group Instanceses. Each element contains the following attributes:
    * `instances` - Instance information added to the Cms Group.
        * `category` - The category of instance.
        * `instance_id` - The id of instance.
        * `instance_name` - The name of instance.
        * `region_id` - The region id of instance.