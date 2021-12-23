---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_launch_template"
sidebar_current: "docs-alicloud-resource-ecs-launch-template"
description: |-
  Provides a Alicloud ECS Launch Template resource.
---

# alicloud\_ecs\_launch\_template

Provides a ECS Launch Template resource.

For information about ECS Launch Template and how to use it, see [What is Launch Template](https://www.alibabacloud.com/help/en/doc-detail/74686.htm).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_launch_template" "default" {
  name                          = "tf_test_name"
  description                   = "Test For Terraform"
  image_id                      = "m-bp1i3ucxxxxx"
  host_name                     = "host_name"
  instance_charge_type          = "PrePaid"
  instance_name                 = "instance_name"
  instance_type                 = "instance_type"
  internet_charge_type          = "PayByBandwidth"
  internet_max_bandwidth_in     = "5"
  internet_max_bandwidth_out    = "0"
  io_optimized                  = "optimized"
  key_pair_name                 = "key_pair_name"
  ram_role_name                 = "ram_role_name"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = "5"
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_ids            = ["sg-zkdfjaxxxxxx"]
  system_disk {
    category             = "cloud_ssd"
    description          = "Test For Terraform"
    name                 = "tf_test_name"
    size                 = "40"
    delete_with_instance = "false"
  }

  resource_group_id = "rg-zkdfjaxxxxxx"
  user_data         = "xxxxxxx"
  vswitch_id        = "vw-zwxscaxxxxxx"
  vpc_id            = "vpc-asdfnbgxxxxxxx"
  zone_id           = "cn-hangzhou-i"

  template_tags = {
    Create = "Terraform"
    For    = "Test"
  }

  network_interfaces {
    name              = "eth0"
    description       = "hello1"
    primary_ip        = "10.0.0.2"
    security_group_id = "sg-asdfnbgxxxxxxx"
    vswitch_id        = "vw-zkdfjaxxxxxx"
  }

  data_disks {
    name                 = "disk1"
    description          = "test1"
    delete_with_instance = "true"
    category             = "cloud"
    encrypted            = "false"
    performance_level    = "PL0"
    size                 = "20"
  }
  data_disks {
    name                 = "disk2"
    description          = "test2"
    delete_with_instance = "true"
    category             = "cloud"
    encrypted            = "false"
    performance_level    = "PL0"
    size                 = "20"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_release_time` - (Optional) Instance auto release time. The time is presented using the ISO8601 standard and in UTC time. The format is  YYYY-MM-DDTHH:MM:SSZ.
* `data_disks` - (Optional) The list of data disks created with instance.
* `deployment_set_id` - (Optional) The Deployment Set Id.
* `description` - (Optional) Description of instance launch template version 1. It can be [2, 256] characters in length. It cannot start with "http://" or "https://". The default value is null.
* `enable_vm_os_config` - (Optional) Whether to enable the instance operating system configuration.
* `host_name` - (Optional) Instance host name.It cannot start or end with a period (.) or a hyphen (-) and it cannot have two or more consecutive periods (.) or hyphens (-).For Windows: The host name can be [2, 15] characters in length. It can contain A-Z, a-z, numbers, periods (.), and hyphens (-). It cannot only contain numbers. For other operating systems: The host name can be [2, 64] characters in length. It can be segments separated by periods (.). It can contain A-Z, a-z, numbers, and hyphens (-).
* `image_id` - (Optional) The Image ID.
* `image_owner_alias` - (Optional) Mirror source. Valid values: `system`, `self`, `others`, `marketplace`, `""`. Default to: `""`.
* `instance_charge_type` - (Optional) Billing methods. Valid values: `PostPaid`, `PrePaid`.
* `instance_type` - (Optional) Instance type. For more information, call resource_alicloud_instances to obtain the latest instance type list.
* `internet_charge_type` - (Optional) Internet bandwidth billing method. Valid values: `PayByTraffic`, `PayByBandwidth`.
* `internet_max_bandwidth_in` - (Optional) The maximum inbound bandwidth from the Internet network, measured in Mbit/s. Value range: [1, 200].
* `internet_max_bandwidth_out` - (Optional) Maximum outbound bandwidth from the Internet, its unit of measurement is Mbit/s. Value range: [0, 100].
* `io_optimized` - (Optional) Whether it is an I/O-optimized instance or not. Valid values: `none`, `optimized`.
* `key_pair_name` - (Optional) The name of the key pair.
    - Ignore this parameter for Windows instances. It is null by default. Even if you enter this parameter, only the  Password content is used.
    - The password logon method for Linux instances is set to forbidden upon initialization.
* `launch_template_name` - (Optional, ForceNew) The name of Launch Template.
* `network_interfaces` - (Optional) The list of network interfaces created with instance.
* `network_type` - (Optional) Network type of the instance. Valid values: `classic`, `vpc`.
* `password_inherit` - (Optional) Whether to use the password preset by the mirror.
* `period` - (Optional) The subscription period of the instance. Unit: months. This parameter takes effect and is required only when InstanceChargeType is set to PrePaid. If the DedicatedHostId parameter is specified, the value of the Period parameter must be within the subscription period of the dedicated host.
    - When the PeriodUnit parameter is set to `Week`, the valid values of the Period parameter are `1`, `2`, `3`, and `4`.
    - When the PeriodUnit parameter is set to `Month`, the valid values of the Period parameter are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`, `48`, and `60`.
* `private_ip_address` - (Optional) The private IP address of the instance.
* `ram_role_name` - (Optional) The RAM role name of the instance. You can use the RAM API ListRoles to query instance RAM role names.
* `resource_group_id` - (Optional) The ID of the resource group to which to assign the instance, Elastic Block Storage (EBS) device, and ENI.
* `security_enhancement_strategy` - (Optional) Whether or not to activate the security enhancement feature and install network security software free of charge. Valid values: `Active`, `Deactive`.
* `security_group_id` - (Optional) The security group ID.
* `security_group_ids` - (Optional) The ID of security group N to which to assign the instance.
* `spot_duration` - (Optional, Computed) The protection period of the preemptible instance. Unit: hours. Valid values: `0`, `1`, `2`, `3`, `4`, `5`, and `6`. Default to: `1`.
* `spot_price_limit` -(Optional) Sets the maximum hourly instance price. Supports up to three decimal places.
* `spot_strategy` - (Optional) The spot strategy for a Pay-As-You-Go instance. This parameter is valid and required only when InstanceChargeType is set to PostPaid. Valid values: `NoSpot`, `SpotAsPriceGo`, `SpotWithPriceLimit`.
* `system_disk` - (Optional) The System Disk.
* `template_resource_group_id` - (Optional, ForceNew) The template resource group id.
* `user_data` - (Optional, Computed) The User Data.
* `version_description` - (Optional) The description of the launch template version. The description must be 2 to 256 characters in length and cannot start with http:// or https://.                                    
* `vswitch_id` - (Optional) When creating a VPC-Connected instance, you must specify its VSwitch ID.
* `zone_id` - (Optional) The zone ID of the instance.
* `tags` - (Optional) A mapping of tags to assign to instance, block storage, and elastic network.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `template_tags` - (Optional) A mapping of tags to assign to the launch template.
  

#### Block system_disk

The system_disk supports the following: 

* `category` - (Optional, Computed) The category of the system disk. System disk type. Valid values: `all`, `cloud`, `ephemeral_ssd`, `cloud_essd`, `cloud_efficiency`, `cloud_ssd`, `local_disk`.
* `delete_with_instance` - (Optional) Specifies whether to release the system disk when the instance is released. Default to `true`.
* `description` - (Optional, Computed) System disk description. It cannot begin with http:// or https://.
* `iops` - (Optional) The Iops.
* `name` - (Optional, Computed) System disk name. The name is a string of 2 to 128 characters. It must begin with an English or a Chinese character. It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).
* `performance_level` - (Optional) The performance level of the ESSD used as the system disk. Valid Values: `PL0`, `PL1`, `PL2`, and `PL3`. Default to: `PL0`.
* `size` - (Optional, Computed) Size of the system disk, measured in GB. Value range: [20, 500].

#### Block network_interfaces

The network_interfaces supports the following: 

* `description` - (Optional) The ENI description.
* `name` - (Optional) The ENI name.
* `primary_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID must be one in the same VPC.
* `vswitch_id` - (Optional) The VSwitch ID for ENI. The instance must be in the same zone of the same VPC network as the ENI, but they may belong to different VSwitches.

#### Block data_disks

The data_disks supports the following: 

* `category` - (Optional) The category of the disk.
* `delete_with_instance` - (Optional) Indicates whether the data disk is released with the instance.
* `description` - (Optional) The description of the data disk.
* `encrypted` - (Optional) Encrypted the data in this disk.
* `name` - (Optional) The name of the data disk.
* `performance_level` - (Optional) The performance level of the ESSD used as the data disk.
* `size` - (Optional) The size of the data disk.
* `snapshot_id` - (Optional) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Launch Template.

## Import

ECS Launch Template can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_launch_template.example <id>
```
