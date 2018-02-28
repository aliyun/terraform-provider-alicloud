---
layout: "alicloud"
page_title: "Alicloud: alicloud_zones"
sidebar_current: "docs-alicloud-datasource-zones"
description: |-
    Provides a list of Availability Zones which can be used by an Alicloud account.
---

# alicloud\_zones

The Zones data source allows access to the list of Alicloud Zones which can be accessed by an Alicloud account within the region configured in the provider.

## Example Usage

```
# Declare the data source
data "alicloud_zones" "default" {
	"available_instance_type"= "ecs.n4.large"
	"available_disk_category"= "cloud_ssd"
}

# Create ecs instance with the first matched zone

resource "alicloud_instance" "instance" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"

  # Other properties...
}

```

## Argument Reference

The following arguments are supported:

* `available_instance_type` - (Optional) Limit search to specific instance type.
* `available_resource_creation` - (Optional) Limit search to specific resource type. The following values are allowed `Instance`, `Disk`, `VSwitch` and `Rds`.
* `available_disk_category` - (Optional) Limit search to specific disk category. Can be either `cloud`, `cloud_efficiency`, `cloud_ssd`.
* `output_file` - (Optional) The name of file that can save zones data source after running `terraform plan`.

~> **NOTE:** Available disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. So many available zones haven't support it. Recommend `cloud_efficiency` and `cloud_ssd`.

## Attributes Reference

A list of zones will be exported and its every element contains the following attributes:

* `id` - ID of the zone.
* `local_name` - Name of the zone in the local language.
* `available_instance_types` - Instance types allowed.
* `available_resource_creation` - Type of resource that can be created.
* `available_disk_categories` - Set of supported disk categories.
