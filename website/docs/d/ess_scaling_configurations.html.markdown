---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_configurations"
sidebar_current: "docs-alicloud_ess_scaling_configurations"
description: |-
    Provides a list of scaling configurations available to the user.
---

# alicloud_ess_scaling_configurations

This data source provides available scaling configuration resources. 

## Example Usage

```terraform
data "alicloud_ess_scaling_configurations" "scalingconfigurations_ds" {
  scaling_group_id = "scaling_group_id"
  ids              = ["scaling_configuration_id1", "scaling_configuration_id2"]
  name_regex       = "scaling_configuration_name"
}

output "first_scaling_rule" {
  value = "${data.alicloud_ess_scaling_configurations.scalingconfigurations_ds.configurations.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional, ForceNew) Scaling group id the scaling configurations belong to.
* `name_regex` - (Optional, ForceNew) A regex string to filter resulting scaling configurations by name.
* `ids` - (Optional, ForceNew) A list of scaling configuration IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling configuration ids.
* `names` - A list of scaling configuration names.
* `configurations` - A list of scaling rules. Each element contains the following attributes:
  * `id` - ID of the scaling rule.
  * `scaling_group_id` - ID of the scaling group.
  * `name` - Name of the scaling configuration.
  * `image_id` - Image ID of the scaling configuration.
  * `instance_type` - Instance type of the scaling configuration.
  * `security_group_id` - Security group ID of the scaling configuration.
  * `internet_charge_type` - Internet charge type of the scaling configuration.
  * `internet_max_bandwidth_in` - Internet max bandwidth in of the scaling configuration.
  * `internet_max_bandwidth_out` - Internet max bandwidth of the scaling configuration.
  * `credit_specification` - Performance mode of the t5 burstable instance.
  * `system_disk_category` - System disk category of the scaling configuration.
  * `system_disk_size` - System disk size of the scaling configuration.
  * `system_disk_performance_level` - The performance level of the ESSD used as the system disk.
  * `data_disks` - Data disks of the scaling configuration.
    * `size` - Size of data disk.
    * `category` - Category of data disk.
    * `snapshot_id` - Size of data disk.
    * `device` - Device attribute of data disk.
    * `delete_with_instance` - Delete_with_instance attribute of data disk.
    * `performance_level` - The performance level of the ESSD used as data disk.
  * `lifecycle_state` - Lifecycle state of the scaling configuration.
  * `creation_time` - Creation time of the scaling configuration.
  * `instance_name` - (Optional,Available in 1.143.0+) InstanceName of an ECS instance.
  * `host_name` - (Optional,Available in 1.143.0+) Hostname of an ECS instance.
  * `spot_strategy` - (Optional, Available in 1.151.0+) The spot strategy for a Pay-As-You-Go instance.
  * `spot_price_limit` - (Optional, Available in 1.151.0+) The maximum price hourly for instance types.
    * `instance_type` - Resource type of an ECS instance.
    * `price_limit` - Price limit hourly of instance type.
  * `instance_pattern_info` - (Optional, Available since v1.239.0+) intelligent configuration mode. In this mode, you only need to specify the number of vCPUs, memory size, instance family, and maximum price. The system selects an instance type that is provided at the lowest price based on your configurations to create ECS instances. This mode is available only for scaling groups that reside in virtual private clouds (VPCs). This mode helps reduce the failures of scale-out activities caused by insufficient inventory of instance types. 
    * `cores` - The number of vCPUs that are specified for an instance type in instancePatternInfo.
    * `instance_family_level` - The instance family level in instancePatternInfo.
    * `max_price` - The maximum hourly price for a pay-as-you-go instance or a preemptible instance in instancePatternInfo.
    * `memory` - The memory size that is specified for an instance type in instancePatternInfo.
    * `burstable_performance` - Specifies whether to include burstable instance types.  Valid values: Exclude, Include, Required.
    * `excluded_instance_types` - Instance type N that you want to exclude. You can use wildcard characters, such as an asterisk (*), to exclude an instance type or an instance family.
    * `architectures` -  Architecture N of instance type N. Valid values: X86, Heterogeneous, BareMetal, Arm, SuperComputeCluster.
