---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_regions"
sidebar_current: "docs-alicloud-datasource-ros-regions"
description: |-
  Provides a list of Ros Regions to the user.
---

# alicloud\_ros\_regions

This data source provides the Ros Regions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_regions" "all" {}

output "ros_region_region_id_1" {
  value = data.alicloud_ros_regions.all.regions.0.region_id
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of Ros Regions. Each element contains the following attributes:
	* `local_name` - The name of the region.
	* `region_endpoint` - The endpoint of the region.
	* `region_id` - The ID of the region.
