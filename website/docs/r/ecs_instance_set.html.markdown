---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_instance_set"
sidebar_current: "docs-alicloud-ecs-instance-set"
description: |-
  Provides an ECS Instance Set resource.
---

# alicloud\_ecs\_instance\_set

Provides a ECS Instance Set resource.

For information about ECS Instance Set and how to use it, see [What is Instance Set](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/runinstances).

-> **NOTE:** Available in v1.173.0+.

-> **NOTE:** This resource is used to batch create a group of instance resources with the same configuration. However, this resource is not recommended. `alicloud_instance` is preferred.

-> **NOTE:** In the instances managed by this resource, names are automatically generated based on `instance_name` and `unique_suffix`.

-> **NOTE:** Only `tags` support batch modification.

## Example Usage

```
variable "name" {
  default = "tf-testaccecsset"
}


data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}


resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecs_instance_set" "beijing_k" {
  amount                        = 100
  image_id                      = data.alicloud_images.default.images[0].id
  instance_type                 = data.alicloud_instance_types.default.instance_types[0].id
  instance_name                 = var.name
  instance_charge_type          = "PostPaid"
  system_disk_performance_level = "PL0"
  system_disk_category          = "cloud_essd"
  system_disk_size              = 200
  vswitch_id                    = data.alicloud_vswitches.default.ids[0]
  security_group_ids            = alicloud_security_group.default.*.id
  zone_id                       = data.alicloud_zones.default.zones[0].id
}
```


## Argument Reference

The following arguments are supported:

* `amount` - (Optional, ForceNew) The number of instances that you want to create. Valid values: `1` to `100`.
* `resource_group_id` - (Optional, ForceNew) The ID of resource group which the instance belongs.
* `hpc_cluster_id` - (Optional, ForceNew) The ID of the Elastic High Performance Computing (E-HPC) cluster to which to assign the instance.
* `description` - (Optional, ForceNew) The description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with `http://` or `https://`.
* `security_group_ids` - (Required, ForceNew)  A list of security group ids to associate with.
* `image_id` - (Required, ForceNew) The Image to use for the instance.
* `instance_type` - (Required, ForceNew) The type of instance to start. 
* `security_enhancement_strategy` - (Optional, ForceNew) The security enhancement strategy.
    - `Active`: Enable security enhancement strategy, it only works on system images.
    - `Deactive`: Disable security enhancement strategy, it works on all images.
* `instance_name` - (Optional, ForceNew) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen, and must not begin with `http://` or `https://`.
* `password` - (Optional, Sensitive, ForceNew) The password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols.
* `password_inherit` - (Optional, ForceNew) Whether to use the password preset in the image.
* `zone_id` - (Optional, ForceNew, Computed)The ID of the zone in which to create the instance.
* `host_name` - (Optional, ForceNew) The hostname of instance.
* `auto_release_time` - (Optional,ForceNew ) The automatic release time of the `PostPaid` instance.
* `data_disks` - (Optional, ForceNew) The list of data disks created with instance. See the following `Block data_disks`.
* `internet_charge_type` - (Optional, ForceNew, Computed) The Internet charge type of the instance. Valid values are `PayByBandwidth`, `PayByTraffic`.
* `internet_max_bandwidth_out` - (Optional, ForceNew, Computed) The Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value values: `1` to `100`.
* `system_disk_category` - (Optional, ForceNew, Computed) The category of the system disk. Valid values are `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`.
* `system_disk_description` - (Optional, ForceNew) The description of the system disk. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `system_disk_name` - (Optional, ForceNew) The name of the system disk.
* `system_disk_size` - (Optional, ForceNew, Computed) The size of the system disk, measured in GiB. Value range:  values: `20` to `500`.
* `system_disk_performance_level` (Optional, ForceNew, Computed) The performance level of the ESSD used as the system disk. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `system_disk_auto_snapshot_policy_id` - (Optional) The ID of the automatic snapshot policy applied to the system disk.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch in VPC. 
* `ram_role_name` - (Optional, ForceNew) The Instance RAM role name. 
* `key_pair_name` - (Optional, ForceNew) The name of key pair that can login ECS instance successfully without password.
* `spot_strategy` - (Optional, ForceNew, Computed) The spot strategy of a Pay-As-You-Go instance, and it takes effect only when parameter `instance_charge_type` is 'PostPaid'.
    - `NoSpot`: A regular Pay-As-You-Go instance.
    - `SpotWithPriceLimit`: A price threshold for a spot instance.
    - `SpotAsPriceGo`: A price that is based on the highest Pay-As-You-Go instance
* `spot_price_limit` - (Optional, ForceNew, Computed) The hourly price threshold of a instance, and it takes effect only when parameter 'spot_strategy' is 'SpotWithPriceLimit'. Three decimals is allowed at most.
* `dedicated_host_id` - (Optional, ForceNew) The ID of the dedicated host on which to create the instance. If the `dedicated_host_id` is specified, the `spot_strategy` and `spot_price_limit`  are ignored. This is because preemptible instances cannot be created on dedicated hosts.
* `launch_template_name` - (Optional, ForceNew) The name of the launch template. To use a launch template to create an instance, you must use the `launch_template_id` or `launch_template_name` parameter to specify the launch template.
* `launch_template_id` - (Optional, ForceNew) The ID of the launch template.
* `launch_template_version` - (Optional, ForceNew) The version of the launch template.
* `period_unit` - (Optional, ForceNew) The duration unit that you will buy the resource. It is valid when `instance_charge_type` is 'PrePaid'. Valid value: `Week`, `Month`.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal for the instance. This parameter is valid only when the `instance_charge_type` is set to `PrePaid`.
* `auto_renew_period` - (Optional, ForceNew) Auto renewal period of an instance, in the unit of month. It is valid when `instance_charge_type` is `PrePaid`.
    - When `period_unit` is `Month`, Valid values: `1`, `2`, `3`, `6`, `12`.
    - When `period_unit` is `Week`, Valid values: `1`, `2`, `3`.
* `instance_charge_type` - (Optional, ForceNew) The billing method of the instance. Valid values: `PrePaid`, `PostPaid`.
* `period` - (Optional, ForceNew) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`.
    - When `period_unit` is `Month`, Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`, `48`, `60`.
    - When `period_unit` is `Week`, Valid values: `1`, `2`, `3`.
* `deletion_protection` - (Optional, ForceNew, Computed) Whether to enable release protection for the instance.
* `deployment_set_id` - (Optional, ForceNew) The ID of the deployment set to which to deploy the instance.
* `network_interfaces` - (Optional, ForceNew) A list of NetworkInterface. See the following `Block network_interfaces`.
* `unique_suffix` - (Optional, ForceNew) Whether to automatically append incremental suffixes to the hostname specified by the HostName parameter and to the instance name specified by the InstanceName parameter when you batch create instances. The incremental suffixes can range from `001` to `999`.
* `exclude_instance_filter` - (Optional, Available in v1.176.0+) The instances that need to be excluded from the Instance Set. See the following `Block exclude_instance_filter`.
* `boot_check_os_with_assistant`  - (Optional, Available in v1.177.0+) Indicate how to check instance ready to use.
  - `false`: Default value. Means that the instances are ready when their DescribeInstances status is Running, at which time guestOS(Ecs os) may not be ready yet.
  - `true`: Checking instance ready with Ecs assistant, which means guestOs boots successfully. Premise is that the specified image `image_id` has built-in Ecs assistant. Most of the public images have assistant installed already.


#### Block data_disks

The `data_disks` supports the following:

* `disk_name` - (Optional, ForceNew) The name of the data disk.
* `disk_size` - (Required, ForceNew) The size of the data disk. Unit: GiB.
  - When `disk_category` is `cloud_efficiency`, Valid values: `20` to `32768`.
  - When `disk_category` is `cloud_ssd`, Valid values: `20` to `32768`.
  - When `disk_category` is `cloud_essd`, Valid values: `20` to `32768`.
  - When `disk_category` is `cloud`, Valid values: `5` to `200`.
* `disk_category` - (Optional, ForceNew, Computed) The category of the disk. Valid values: `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`.
* `performance_level` - (Optional, ForceNew, Computed) The performance level of the ESSD used as data disk. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `kms_key_id` - (Optional, ForceNew) The KMS key ID corresponding to the data disk.
* `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
* `auto_snapshot_policy_id` - (Optional, ForceNew) The ID of the automatic snapshot policy applied to the system disk.
* `disk_description` - (Optional, ForceNew) The description of the data disk.

#### Block network_interfaces

The `network_interfaces` supports the following:

* `security_group_id` -(Required, ForceNew) The ID of the security group to which to assign secondary ENI.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch to which to connect ENI.
* `description` - (Optional, ForceNew) The description of ENI.
* `network_interface_name` - (Optional, ForceNew) The name of ENI.
* `primary_ip_address` - (Optional, ForceNew) The primary private IP address of ENI.

#### Block exclude_instance_filter

The `exclude_instance_filter` supports the following:

* `key` - (Required) The type of the excluded. Valid values: `InstanceId`, `InstanceName`.
* `value` - (Required) The value of the excluded. The identification of the excluded instances. It is a list of instance Ids or names.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value of the id is the Base64 encoding of `instance_ids`.
* `instance_ids` -  A list of ECS Instance ID.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the ECS Instance Set.
* `update` - (Defaults to 30 mins) Used when update the ECS Instance Set.
* `delete` - (Defaults to 30 mins) Used when delete the ECS Instance Set.

