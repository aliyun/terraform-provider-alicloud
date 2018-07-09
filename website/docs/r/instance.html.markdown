---
layout: "alicloud"
page_title: "Alicloud: alicloud_instance"
sidebar_current: "docs-alicloud-resource-instance"
description: |-
  Provides an ECS instance resource.
---

# alicloud\_instance

Provides a ECS instance resource.

~> **NOTE:** You can launch an ECS instance for a VPC network via specifying parameter `vswitch_id`. One instance can only belong to one VSwitch.

~> **NOTE:** If a VSwitchId is specified for creating an instance, SecurityGroupId and VSwitchId must belong to one VPC.

~> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

~> **NOTE:** At present, 'PrePaid' instance cannot be deleted and must wait it to be outdated and release it automatically.

~> **NOTE:** The resource supports modifying instance charge type from 'PrePaid' to 'PostPaid' from version 1.9.6.
 However, at present, this modification has some limitation about CPU core count in one month, so strongly recommand that `Don't modify instance charge type frequentlly in one month`.

## Example Usage

```
# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  name        = "tf_test_foo"
  description = "foo"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  security_groups = ["${alicloud_security_group.group.*.id}"]

  # series III
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  image_id             = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_name        = "test_foo"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
  internet_max_bandwidth_out = 10
}

# Create a new ECS instance for VPC
resource "alicloud_vpc" "vpc" {
  # Other parameters...
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  # Other parameters...
}

resource "alicloud_slb" "slb" {
  name       = "test-slb-tf"
  vpc_id     = "${alicloud_vpc.vpc.id}"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The Image to use for the instance. ECS instance's image can be replaced via changing 'image_id'. When it is changed, the instance will reboot to make the change take effect.
* `instance_type` - (Required) The type of instance to start.
* `io_optimized` - (Deprecated) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `security_groups` - (Required)  A list of security group ids to associate with.
* `availability_zone` - (Optional) The Zone to start the instance in. It is ignored and will be computed when set `vswitch_id`.
* `instance_name` - (Optional) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. If not specified, 
Terraform will autogenerate a default name is `ECS-Instance`.
* `allocate_public_ip` - (Deprecated) It has been deprecated from version "1.7.0". Setting "internet_max_bandwidth_out" larger than 0 can allocate a public ip address for an instance.
* `system_disk_category` - (Optional) Valid values are `cloud_efficiency`, `cloud_ssd` and `cloud`. `cloud` only is used to some none I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of the system disk, value range: 40GB ~ 500GB. Default is 40GB. ECS instance's system disk can be reset when replacing system disk.
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `PayByBandwidth`, `PayByTraffic`. Default is `PayByTraffic`. At present, 'PrePaid' instance cannot change the value to "PayByBandwidth" from "PayByTraffic".
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100]. Default to 0 Mbps.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters.
On other OSs such as Linux, the host name can contain a maximum of 30 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“. When it is changed, the instance will reboot to make the change take effect.
* `password` - (Optional) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. When it is changed, the instance will reboot to make the change take effect.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. If you want to create instances in VPC network, this parameter must be set.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, The default is `PostPaid`.
* `period_unit` - (Optional) The duration unit that you will buy the resource. It is valid when `instance_charge_type` is 'PrePaid'. Valid value: ["Week", "Month"]. Default to "Month".
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Default to 1. Valid values:
    - [1-9, 12, 24, 36, 48, 60] when `period_unit` in "Month"
    - [1-3] when `period_unit` in "Week"

* `renewal_status` - (Optional) Whether to renew an ECS instance automatically or not. It is valid when `instance_charge_type` is `PrePaid`. Default to "Normal". Valid values:
    - `AutoRenewal`: Enable auto renewal.
    - `Normal`: Disable auto renewal.
    - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.

* `auto_renew_period` - (Optional) Auto renewal period of an instance, in the unit of month. It is valid when `instance_charge_type` is `PrePaid`. Default to 1. Valid value:
    - [1, 2, 3, 6, 12] when `period_unit` in "Month"
    - [1, 2, 3] when `period_unit` in "Week"

* `tags` - (Optional) A mapping of tags to assign to the resource.
* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance and to pass data into an ECS instance.
* `key_name` - (Optional, Force new resource) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional, Force new resource) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `include_data_disks` - (Optional) Whether to change instance disks charge type when changing instance charge type.
* `dry_run` - (Optional) Whether to pre-detection. When it is true, only pre-detection and not actually modify the payment type operation. It is valid when `instance_charge_type` is 'PrePaid'. Default to false.
* `private_ip` - (Optional) Instance private IP address can be specified when you creating new instance. It is valid when `vswitch_id` is specified.
* `spot_strategy` - (Optional, Force New) The spot strategy of a Pay-As-You-Go instance, and it takes effect only when parameter `instance_charge_type` is 'PostPaid'. Value range:
    - NoSpot: A regular Pay-As-You-Go instance.
    - SpotWithPriceLimit: A price threshold for a spot instance
    - SpotAsPriceGo: A price that is based on the highest Pay-As-You-Go instance

    Default to NoSpot.
* `spot_price_limit` - (Optional, Float, Force New) The hourly price threshold of a instance, and it takes effect only when parameter 'spot_strategy' is 'SpotWithPriceLimit'. Three decimals is allowed at most.


~> **NOTE:** System disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

~> **NOTE:** From version 1.5.0, instance's charge type can be changed to "PrePaid" by specifying `period` and `period_unit`, but it is irreversible.

~> **NOTE:** From version 1.5.0, instance's private IP address can be specified when creating VPC network instance.

~> **NOTE:** From version 1.5.0, instance's vswitch and private IP can be changed in the same availability zone. When they are changed, the instance will reboot to make the change take effect.

~> **NOTE:** From version 1.7.0, setting "internet_max_bandwidth_out" larger than 0 can allocate a public IP for an instance.
 Setting "internet_max_bandwidth_out" to 0 can release allocated public IP for VPC instance(For Classic instnace, its public IP cannot be release once it allocated, even thougth its bandwidth out is 0).
 However, at present, 'PrePaid' instance cannot narrow its max bandwidth out when its 'internet_charge_type' is "PayByBandwidth".

~> **NOTE:** From version 1.7.0, instance's type can be changed. When it is changed, the instance will reboot to make the change take effect.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `availability_zone` - The Zone to start the instance in.
* `instance_name` - The instance name.
* `host_name` - The instance host name.
* `description` - The instance description.
* `status` - The instance status.
* `image_id` - The instance Image Id.
* `instance_type` - The instance type.
* `private_ip` - The instance private ip.
* `public_ip` - The instance public ip.
* `vswitch_id` - If the instance created in VPC, then this value is  virtual switch ID.
* `tags` - The instance tags, use jsonencode(item) to display the value.
* `key_name` - The name of key pair that has been bound in ECS instance.
* `role_name` - The name of RAM role that has been bound in ECS instance.
* `user_data` - The hash value of the user data.
* `period` - The ECS instance using duration.
* `period_unit` - The ECS instance using duration unit.
* `renewal_status` - The ECS instance automatically renew status.
* `auto_renew_period` - Auto renewal period of an instance.
* `dry_run` - Whether to pre-detection.
* `spot_strategy` - The spot strategy of a Pay-As-You-Go instance
* `spot_price_limit` - The hourly price threshold of a instance.


## Import

Instance can be imported using the id, e.g.

```
$ terraform import alicloud_instance.example i-abc12345678
```
