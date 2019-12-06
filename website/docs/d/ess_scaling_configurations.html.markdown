---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_configurations"
sidebar_current: "docs-alicloud_ess_scaling_configurations"
description: |-
    Provides a list of scaling configurations available to the user.
---

# alicloud_ess_scaling_configurations

This data source provides available scaling configuration resources. 

## Example Usage

```
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

* `scaling_group_id` - (Optional) Scaling group id the scaling configurations belong to.
* `name_regex` - (Optional) A regex string to filter resulting scaling configurations by name.
* `ids` - (Optional) A list of scaling configuration IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling configuration ids.
* `names` - A list of scaling configuration names.
* `scaling_configurations` - A list of scaling rules. Each element contains the following attributes:
  * `id` - ID of the scaling rule.
  * `scaling_group_id` - ID of the scaling group.
  * `name` - Name of the scaling configuration.
  * `image_id` - Image ID of the scaling configuration.
  * `instance_type` - Instance type of the scaling configuration.
  * `security_group_id` - Security group ID of the scaling configuration.
  * `internet_charge_type` - Internet charge type of the scaling configuration.
  * `internet_max_bandwidth_in` - Internet max bandwidth in of the scaling configuration.
  * `internet_max_bandwidth_out` - Internet max bandwidth of the scaling configuration.
  * `system_disk_category` - System disk category of the scaling configuration.
  * `system_disk_size` - System disk size of the scaling configuration.
  * `data_disks` - Data disks of the scaling configuration.
    * `size` - Size of data disk.
    * `category` - Category of data disk.
    * `snapshot_id` - Size of data disk.
    * `device` - Device attribute of data disk.
    * `delete_with_instance` - Delete_with_instance attribute of data disk.
  * `lifecycle_state` - Lifecycle state of the scaling configuration.
  * `creation_time` - Creation time of the scaling configuration.
  
