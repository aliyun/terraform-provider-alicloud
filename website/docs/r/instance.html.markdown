---
layout: "alicloud"
page_title: "Alicloud: alicloud_instance"
sidebar_current: "docs-alicloud-resource-instance"
description: |-
  Provides an ECS instance resource.
---

# alicloud\_instance

Provides a ECS instance resource.

-> **NOTE:** You can launch an ECS instance for a VPC network via specifying parameter `vswitch_id`. One instance can only belong to one VSwitch.

-> **NOTE:** If a VSwitchId is specified for creating an instance, SecurityGroupId and VSwitchId must belong to one VPC.

-> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

-> **NOTE:** At present, 'PrePaid' instance cannot be deleted and must wait it to be outdated and release it automatically.

-> **NOTE:** The resource supports modifying instance charge type from 'PrePaid' to 'PostPaid' from version 1.9.6.
 However, at present, this modification has some limitation about CPU core count in one month, so strongly recommand that `Don't modify instance charge type frequentlly in one month`.

-> **NOTE:**  There is unsupported 'deletion_protection' attribute when the instance is spot

## Example Usage

```
# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  name        = "tf_test_foo"
  description = "foo"
  vpc_id      = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  security_groups   = "${alicloud_security_group.group.*.id}"

  # series III
  instance_type              = "ecs.n4.large"
  system_disk_category       = "cloud_efficiency"
  image_id                   = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_name              = "test_foo"
  vswitch_id                 = "${alicloud_vswitch.vswitch.id}"
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
* `system_disk_category` - (Optional) Valid values are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`. `cloud` only is used to some none I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of the system disk, measured in GiB. Value range: [20, 500]. The specified value must be equal to or greater than max{20, Imagesize}. Default value: max{40, ImageSize}. ECS instance's system disk can be reset when replacing system disk.
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `PayByBandwidth`, `PayByTraffic`. Default is `PayByTraffic`. At present, 'PrePaid' instance cannot change the value to "PayByBandwidth" from "PayByTraffic".
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100]. Default to 0 Mbps.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters.
On other OSs such as Linux, the host name can contain a maximum of 30 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“. When it is changed, the instance will reboot to make the change take effect.
* `password` - (Optional, Sensitive) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. When it is changed, the instance will reboot to make the change take effect.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. This parameter must be set unless you can create classic network instances.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, The default is `PostPaid`.
* `resource_group_id` - (ForceNew, Available in 1.57.0+) The Id of resource group which the instance belongs.
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
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `volume_tags` - (Optional) A mapping of tags to assign to the devices created by the instance at launch time.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance and to pass data into an ECS instance.
* `key_name` - (Optional, Force new resource) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional, Force new resource) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `include_data_disks` - (Optional) Whether to change instance disks charge type when changing instance charge type.
* `dry_run` - (Optional) Whether to pre-detection. When it is true, only pre-detection and not actually modify the payment type operation. It is valid when `instance_charge_type` is 'PrePaid'. Default to false.
* `private_ip` - (Optional) Instance private IP address can be specified when you creating new instance. It is valid when `vswitch_id` is specified.
* `spot_strategy` - (Optional, ForceNew) The spot strategy of a Pay-As-You-Go instance, and it takes effect only when parameter `instance_charge_type` is 'PostPaid'. Value range:
    - NoSpot: A regular Pay-As-You-Go instance.
    - SpotWithPriceLimit: A price threshold for a spot instance
    - SpotAsPriceGo: A price that is based on the highest Pay-As-You-Go instance

    Default to NoSpot. Note: Currently, the spot instance only supports domestic site account.
* `spot_price_limit` - (Optional, Float, ForceNew) The hourly price threshold of a instance, and it takes effect only when parameter 'spot_strategy' is 'SpotWithPriceLimit'. Three decimals is allowed at most.
* `deletion_protection` - (Optional, true) Whether enable the deletion protection or not.
    - true: Enable deletion protection.
    - false: Disable deletion protection.
    
    Default to false.
* `force_delete` - (Optional, Available 1.18.0+) If it is true, the "PrePaid" instance will be change to "PostPaid" and then deleted forcibly.
However, because of changing instance charge type has CPU core count quota limitation, so strongly recommand that "Don't modify instance charge type frequentlly in one month".
* `security_enhancement_strategy` - (Optional, ForceNew) The security enhancement strategy.
    - Active: Enable security enhancement strategy, it only works on system images.
    - Deactive: Disable security enhancement strategy, it works on all images.
* `data_disks` - (Optional, ForceNew, Available 1.23.1+) The list of data disks created with instance.
    * `name` - (Optional, ForceNew) The name of the data disk.
    * `size` - (Required, ForceNew) The size of the data disk.
        - cloud：[5, 2000]
        - cloud_efficiency：[20, 32768]
        - cloud_ssd：[20, 32768]
        - cloud_essd：[20, 32768]
        - ephemeral_ssd: [5, 800]
    * `category` - (Optional, ForceNew) The category of the disk:
        - `cloud`: The general cloud disk.
        - `cloud_efficiency`: The efficiency cloud disk.
        - `cloud_ssd`: The SSD cloud disk.
        - `cloud_essd`: The ESSD cloud disk.
        - `ephemeral_ssd`: The local SSD disk.
        Default to `cloud_efficiency`.
    * `encrypted` -(Optional, Bool, ForceNew) Encrypted the data in this disk.

        Default to false
    * `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
    * `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on cloud, cloud_efficiency, cloud_essd, cloud_ssd disk. If the category of this data disk was ephemeral_ssd, please don't set this param.

        Default to true
    * `description` - (Optional, ForceNew) The description of the data disk.

-> **NOTE:** System disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

-> **NOTE:** From version 1.5.0, instance's charge type can be changed to "PrePaid" by specifying `period` and `period_unit`, but it is irreversible.

-> **NOTE:** From version 1.5.0, instance's private IP address can be specified when creating VPC network instance.

-> **NOTE:** From version 1.5.0, instance's vswitch and private IP can be changed in the same availability zone. When they are changed, the instance will reboot to make the change take effect.

-> **NOTE:** From version 1.7.0, setting "internet_max_bandwidth_out" larger than 0 can allocate a public IP for an instance.
 Setting "internet_max_bandwidth_out" to 0 can release allocated public IP for VPC instance(For Classic instnace, its public IP cannot be release once it allocated, even thougth its bandwidth out is 0).
 However, at present, 'PrePaid' instance cannot narrow its max bandwidth out when its 'internet_charge_type' is "PayByBandwidth".

-> **NOTE:** From version 1.7.0, instance's type can be changed. When it is changed, the instance will reboot to make the change take effect.

### Timeouts

-> **NOTE:** Available in 1.46.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the instance (until it reaches the initial `Running` status). 
`Note`: There are extra at most 2 minutes used to retry to aviod some needless API errors and it is not in the timeouts configure.
* `update` - (Defaults to 10 mins) Used when stopping and starting the instance when necessary during update - e.g. when changing instance type, password, image, vswitch and private IP.
* `delete` - (Defaults to 20 mins) Used when terminating the instance. `Note`: There are extra at most 5 minutes used to retry to aviod some needless API errors and it is not in the timeouts configure.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `status` - The instance status.
* `public_ip` - The instance public ip.

## Import

Instance can be imported using the id, e.g.

```
$ terraform import alicloud_instance.example i-abc12345678
```
