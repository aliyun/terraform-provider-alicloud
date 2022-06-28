---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_hosts"
sidebar_current: "docs-alicloud-datasource-cddc-dedicated-hosts"
description: |-
  Provides a list of Cddc Dedicated Hosts to the user.
---

# alicloud\_cddc\_dedicated\_hosts

This data source provides the Cddc Dedicated Hosts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cddc_dedicated_hosts" "ids" {
  dedicated_host_group_id = "example_value"
  ids                     = ["example_value-1", "example_value-2"]
}
output "cddc_dedicated_host_id_1" {
  value = data.alicloud_cddc_dedicated_hosts.ids.hosts.0.id
}

data "alicloud_cddc_dedicated_hosts" "status" {
  dedicated_host_group_id = "example_value"
  ids                     = ["example_value-1", "example_value-2"]
  status                  = "1"
}
output "cddc_dedicated_host_id_2" {
  value = data.alicloud_cddc_dedicated_hosts.status.hosts.0.id
}

data "alicloud_cddc_dedicated_hosts" "zoneId" {
  dedicated_host_group_id = "example_value"
  ids                     = ["example_value-1", "example_value-2"]
  zone_id                 = "example_value"
}
output "cddc_dedicated_host_id_3" {
  value = data.alicloud_cddc_dedicated_hosts.zoneId.hosts.0.id
}

data "alicloud_cddc_dedicated_hosts" "allocationStatus" {
  dedicated_host_group_id = "example_value"
  ids                     = ["example_value-1", "example_value-2"]
  allocation_status       = "Allocatable"
}
output "cddc_dedicated_host_id_4" {
  value = data.alicloud_cddc_dedicated_hosts.allocationStatus.hosts.0.id
}

data "alicloud_cddc_dedicated_hosts" "hostType" {
  dedicated_host_group_id = "example_value"
  ids                     = ["example_value-1", "example_value-2"]
  host_type               = "dhg_cloud_ssd"
}
output "cddc_dedicated_host_id_5" {
  value = data.alicloud_cddc_dedicated_hosts.hostType.hosts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `allocation_status` - (Optional, ForceNew) Specifies whether instances can be created on the host. Valid values: `Allocatable` or `Suspended`. `Allocatable`: Instances can be created on the host. `Suspended`: Instances cannot be created on the host.
* `dedicated_host_group_id` - (Required, ForceNew) The ID of the dedicated cluster.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `host_type` - (Optional, ForceNew) The storage type of the host. Valid values: `dhg_local_ssd` or `dhg_cloud_ssd`. `dhg_local_ssd`: specifies that the host uses local SSDs. `dhg_cloud_ssd`: specifies that the host uses enhanced SSDs (ESSDs).
* `ids` - (Optional, ForceNew, Computed)  A list of Dedicated Host IDs.
* `order_id` - (Optional, ForceNew) The ID of the order.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The state of the host. Valid values: 
  * `0:` The host is being created. 
  * `1`: The host is running. 
  * `2`: The host is faulty. 
  * `3`: The host is ready for deactivation. 
  * `4`: The host is being maintained. 
  * `5`: The host is deactivated. 
  * `6`: The host is restarting. 
  * `7`: The host is locked.
* `zone_id` - (Optional, ForceNew) The ID of the zone.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `hosts` - A list of Cddc Dedicated Hosts. Each element contains the following attributes:
  * `bastion_instance_id` - The ID of the bastion host with which the host is associated.
  * `cpu_allocation_ratio` - The numeric value of the CPU over commit ratio of the dedicated cluster.
  * `cpu_used` - The number of CPU cores used by the host.
  * `create_time` - The time when the host was created. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
  * `dedicated_host_group_id` - The ID of the dedicated cluster in which the host is created.
  * `dedicated_host_id` - The ID of the host.
  * `disk_allocation_ratio` - The disk usage in percentage.
  * `ecs_class_code` - The Elastic Compute Service (ECS) instance type.
  * `end_time` - The time when the host expires. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
  * `engine` - The type of the database engine that is used by the host.
  * `expired_time` - The time when the host expires. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
  * `host_class` - The instance type of the host.
  * `host_cpu` - The number of CPU cores specified for the host. Unit: `core`.
  * `host_mem` - The memory of the host. Unit: `GB`.
  * `host_name` - The name of the host.
  * `host_storage` - The total storage capacity of the host. Unit: `GB`.
  * `host_type` - The storage type of the host.
  * `id` - The ID of the Dedicated Host. The value formats as `<dedicated_host_group_id>:<dedicated_host_id>`.
  * `image_category` - The image type of the host.
  * `ip_address` - The IP address of the host.
  * `mem_allocation_ratio` - The memory usage in percentage.
  * `memory_used` - The amount of memory used by the host. Unit: `GB`.
  * `open_permission` - Indicates whether you have the OS permissions on the host. Valid values: `0`: You do not have the OS permissions on the host. `1`: You have the OS permissions on the host.
  * `allocation_status` - Specifies whether instances can be created on the host. Valid values: `1` or `0`. `1`: Instances can be created on the host. `0`: Instances cannot be created on the host.
  * `status` - The state of the host.
  * `storage_used` - The storage usage of the host. Unit: `GB`.
  * `tags` - The tag of the resource.
    * `tag_key` - The key of the tags.
    * `tag_value` - The value of the tags.
  * `vpc_id` - The ID of the virtual private cloud (VPC) to which the host is connected.
  * `vswitch_id` - The ID of the vSwitch.
  * `zone_id` - The zone ID of the host.