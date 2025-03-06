---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_launch_template"
sidebar_current: "docs-alicloud-resource-ecs-launch-template"
description: |-
  Provides a Alicloud ECS Launch Template resource.
---

# alicloud_ecs_launch_template

Provides a ECS Launch Template resource.

For information about ECS Launch Template and how to use it, see [What is Launch Template](https://www.alibabacloud.com/help/en/doc-detail/74686.htm).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_launch_template&exampleId=5466eb32-9bb1-6f2c-d8d0-d78c6ec2ff1525b4b0ef&activeTab=example&spm=docs.r.ecs_launch_template.0.5466eb329b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  security_group_name = "terraform-example"
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_ecs_launch_template" "default" {
  launch_template_name          = "terraform-example"
  description                   = "terraform-example"
  image_id                      = data.alicloud_images.default.images.0.id
  host_name                     = "terraform-example"
  instance_charge_type          = "PrePaid"
  instance_name                 = "terraform-example"
  instance_type                 = data.alicloud_instance_types.default.instance_types.0.id
  internet_charge_type          = "PayByBandwidth"
  internet_max_bandwidth_in     = "5"
  internet_max_bandwidth_out    = "5"
  io_optimized                  = "optimized"
  key_pair_name                 = "key_pair_name"
  ram_role_name                 = "ram_role_name"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = "5"
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_ids            = [alicloud_security_group.default.id]
  system_disk {
    category             = "cloud_ssd"
    description          = "Test For Terraform"
    name                 = "terraform-example"
    size                 = "40"
    delete_with_instance = "false"
  }

  user_data  = "xxxxxxx"
  vswitch_id = alicloud_vswitch.default.id
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_zones.default.zones.0.id

  template_tags = {
    Create = "Terraform"
    For    = "example"
  }

  network_interfaces {
    name              = "eth0"
    description       = "hello1"
    primary_ip        = "10.0.0.2"
    security_group_id = alicloud_security_group.default.id
    vswitch_id        = alicloud_vswitch.default.id
  }

  data_disks {
    name                 = "disk1"
    description          = "description"
    delete_with_instance = "true"
    category             = "cloud"
    encrypted            = "false"
    performance_level    = "PL0"
    size                 = "20"
  }
  data_disks {
    name                 = "disk2"
    description          = "description2"
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

* `auto_renew` - (Optional, ForceNew, Computed, Available since v1.226.0) Specifies whether to enable auto-renewal for the instance. This parameter is valid only if `internet_charge_type` is set to `PrePaid`.
* `auto_renew_period` - (Optional, ForceNew, Computed, Available since v1.226.0) The auto-renewal period of the instance. Valid values when `period_unit` is set to `Month`: 1, 2, 3, 6, 12, 24, 36, 48, and 60. Default value: 1.
* `auto_release_time` - (Optional) Instance auto release time. The time is presented using the ISO8601 standard and in UTC time. The format is  YYYY-MM-DDTHH:MM:SSZ.
* `data_disks` - (Optional) The list of data disks created with instance. See [`data_disks`](#data_disks) below.
* `deployment_set_id` - (Optional) The Deployment Set Id.
* `description` - (Optional) Description of instance launch template version 1. It can be [2, 256] characters in length. It cannot start with "http://" or "https://". The default value is null.
* `enable_vm_os_config` - (Optional) Whether to enable the instance operating system configuration.
* `host_name` - (Optional) Instance host name.It cannot start or end with a period (.) or a hyphen (-) and it cannot have two or more consecutive periods (.) or hyphens (-).For Windows: The host name can be [2, 15] characters in length. It can contain A-Z, a-z, numbers, periods (.), and hyphens (-). It cannot only contain numbers. For other operating systems: The host name can be [2, 64] characters in length. It can be segments separated by periods (.). It can contain A-Z, a-z, numbers, and hyphens (-).
* `image_id` - (Optional) The Image ID.
* `image_owner_alias` - (Optional) Mirror source. Valid values: `system`, `self`, `others`, `marketplace`, `""`. Default to: `""`.
* `instance_name` - (Optional) The name of the instance. The name must be `2` to `128` characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), periods (.), commas (,), brackets ([]), and hyphens (-).
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
* `network_interfaces` - (Optional) The list of network interfaces created with instance. See [`network_interfaces`](#network_interfaces) below.
* `network_type` - (Optional) Network type of the instance. Valid values: `classic`, `vpc`.
* `password_inherit` - (Optional) Whether to use the password preset by the mirror.
* `period_unit` - (Optional, ForceNew, Computed, Available since v1.226.0) The unit of the subscription period. Default value: `Month`. Valid values: `Week`, `Month`.
* `period` - (Optional, ForceNew, Computed) The subscription period of the instance. Unit: months. This parameter takes effect and is required only when InstanceChargeType is set to PrePaid. If the DedicatedHostId parameter is specified, the value of the Period parameter must be within the subscription period of the dedicated host.
  - When the `period_unit` is set to `Week`, the valid values of the Period parameter are `1`, `2`, `3`.
  - When the `period_unit` is set to `Month`, the valid values of the Period parameter are `1`, `2`, `3`, `6`, `12`, `24`, `36`, `48`, and `60`.
* `private_ip_address` - (Optional) The private IP address of the instance.
* `ram_role_name` - (Optional) The RAM role name of the instance. You can use the RAM API ListRoles to query instance RAM role names.
* `resource_group_id` - (Optional) The ID of the resource group to which to assign the instance, Elastic Block Storage (EBS) device, and ENI.
* `security_enhancement_strategy` - (Optional) Whether or not to activate the security enhancement feature and install network security software free of charge. Valid values: `Active`, `Deactive`.
* `security_group_id` - (Optional) The security group ID.
* `security_group_ids` - (Optional) The ID of security group N to which to assign the instance.
* `spot_duration` - (Optional, Computed) The protection period of the preemptible instance. Unit: hours. Valid values: `0`, `1`, `2`, `3`, `4`, `5`, and `6`. Default to: `1`.
* `spot_price_limit` -(Optional) Sets the maximum hourly instance price. Supports up to three decimal places.
* `spot_strategy` - (Optional) The spot strategy for a Pay-As-You-Go instance. This parameter is valid and required only when InstanceChargeType is set to PostPaid. Valid values: `NoSpot`, `SpotAsPriceGo`, `SpotWithPriceLimit`.
* `system_disk` - (Optional) The System Disk. See [`system_disk`](#system_disk) below.
* `template_resource_group_id` - (Optional, ForceNew) The template resource group id.
* `user_data` - (Optional, Computed) The User Data.
* `version_description` - (Optional) The description of the launch template version. The description must be 2 to 256 characters in length and cannot start with http:// or https://.                                    
* `vpc_id` - (Optional) The ID of the VPC.
* `vswitch_id` - (Optional) When creating a VPC-Connected instance, you must specify its VSwitch ID.
* `zone_id` - (Optional) The zone ID of the instance.
* `http_endpoint` - (Optional, ForceNew) Whether to enable access to instance metadata. Valid values:
  - enabled: Enabled.
  - disabled: Disabled.
* `http_tokens` - (Optional, ForceNew) Whether to use the hardened mode (IMDSv2) when accessing instance metadata. Valid values:
  - optional: Not mandatory.
  - required: Mandatory. After this value is set, the normal mode cannot access instance metadata.
* `http_put_response_hop_limit` - (Optional, ForceNew) The HTTP PUT response hop limit required for instance metadata requests.
* `tags` - (Optional) A mapping of tags to assign to instance, block storage, and elastic network.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `template_tags` - (Optional) A mapping of tags to assign to the launch template. 
* `default_version_number` - (Optional, Int, Computed, Available since v1.241.0) The version number of the default launch template version. Default to 1. It is conflict with `update_default_version_number`.
* `update_default_version_number` - (Optional, Bool, Available since v1.241.0) Whether to update the default version of the launch template to the latest version automatically. It is conflict with `default_version_number`.
* `name` - (Deprecated) It has been deprecated from version 1.120.0, and use field `launch_template_name` instead.
* `userdata` - (Deprecated) It has been deprecated from version 1.120.0, and use field `user_data` instead.
* `system_disk_name` - (Deprecated) It has been deprecated from version 1.120.0, and use field `system_disk` instead.
* `system_disk_category` - (Deprecated) It has been deprecated from version 1.120.0, and use field `system_disk` instead.
* `system_disk_size` - (Deprecated) It has been deprecated from version 1.120.0, and use field `system_disk` instead.
* `system_disk_description` - (Deprecated) It has been deprecated from version 1.120.0, and use field `system_disk` instead.

### `system_disk`

The system_disk supports the following: 

* `category` - (Optional, Computed) The category of the system disk. System disk type. Valid values: `all`, `cloud`, `ephemeral_ssd`, `cloud_essd`, `cloud_efficiency`, `cloud_ssd`, `local_disk`.
* `delete_with_instance` - (Optional) Specifies whether to release the system disk when the instance is released. Default to `true`.
* `description` - (Optional, Computed) System disk description. It cannot begin with http:// or https://.
* `iops` - (Optional) The Iops.
* `name` - (Optional, Computed) System disk name. The name is a string of 2 to 128 characters. It must begin with an English or a Chinese character. It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).
* `performance_level` - (Optional) The performance level of the ESSD used as the system disk. Valid Values: `PL0`, `PL1`, `PL2`, and `PL3`. Default to: `PL0`.
* `size` - (Optional, Computed) Size of the system disk, measured in GB. Value range: [20, 500].
* `encrypted` - (Optional) Specifies whether the system disk is encrypted.

### `network_interfaces`

The network_interfaces supports the following: 

* `description` - (Optional) The ENI description.
* `name` - (Optional) The ENI name.
* `primary_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID must be one in the same VPC.
* `vswitch_id` - (Optional) The VSwitch ID for ENI. The instance must be in the same zone of the same VPC network as the ENI, but they may belong to different VSwitches.
* `delete_on_release` - (Optional, Bool, Available since v1.245.0) Specifies whether to release ENI N when the instance is released. Valid values: `true`, `false`.

### `data_disks`

The data_disks supports the following: 

* `category` - (Optional) The category of the disk.
* `delete_with_instance` - (Optional) Indicates whether the data disk is released with the instance.
* `description` - (Optional) The description of the data disk.
* `encrypted` - (Optional) Encrypted the data in this disk.
* `name` - (Optional) The name of the data disk.
* `performance_level` - (Optional) The performance level of the ESSD used as the data disk.
* `size` - (Optional) The size of the data disk.
* `snapshot_id` - (Optional) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
* `device` - (Optional, Available since v1.230.1) The mount point of the data disk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Launch Template.
* `latest_version_number` - The latest version number of the launch template.

## Import

ECS Launch Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_launch_template.example <id>
```
