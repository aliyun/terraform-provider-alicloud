---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instance"
sidebar_current: "docs-alicloud-resource-instance"
description: |-
  Provides an ECS instance resource.
---

# alicloud\_instance

Provides a ECS instance resource.

-> **NOTE:** You can launch an ECS instance for a VPC network via specifying parameter `vswitch_id`. One instance can only belong to one VSwitch.

-> **NOTE:** If a VSwitchId is specified for creating an instance, SecurityGroupId and VSwitchId must belong to one VPC, VSwitchId Cannot be modified after creation.

-> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

-> **NOTE:** At present, 'PrePaid' instance cannot be deleted and must wait it to be outdated and release it automatically.

-> **NOTE:** The resource supports modifying instance charge type from 'PrePaid' to 'PostPaid' from version 1.9.6.
 However, at present, this modification has some limitation about CPU core count in one month, so strongly recommand that `Don't modify instance charge type frequentlly in one month`.

-> **NOTE:**  There is unsupported 'deletion_protection' attribute when the instance is spot

## Example Usage

```
variable "name" {
  default = "auto_provisioning_group"
}

# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  name        = "tf_test_foo"
  description = "foo"
  vpc_id      = alicloud_vpc.vpc.id
}

resource "alicloud_kms_key" "key" {
  description            = "Hello KMS"
  pending_window_in_days = "7"
  key_state              = "Enabled"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

# Create a new ECS instance for VPC
resource "alicloud_vpc" "vpc" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_slb_load_balancer" "slb" {
  load_balancer_name       = "test-slb-tf"
  vswitch_id = alicloud_vswitch.vswitch.id
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  security_groups   = alicloud_security_group.group.*.id

  # series III
  instance_type              = "ecs.n4.large"
  system_disk_category       = "cloud_efficiency"
  system_disk_name           = "test_foo_system_disk_name"
  system_disk_description    = "test_foo_system_disk_description"
  image_id                   = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_name              = "test_foo"
  vswitch_id                 = alicloud_vswitch.vswitch.id
  internet_max_bandwidth_out = 10
  data_disks {
    name        = "disk2"
    size        = 20
    category    = "cloud_efficiency"
    description = "disk2"
    encrypted   = true
    kms_key_id  = alicloud_kms_key.key.id
  }
}
```

## Module Support

You can use the existing [ecs-instance module](https://registry.terraform.io/modules/alibaba/ecs-instance/alicloud) 
to create several ECS instances one-click.

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The Image to use for the instance. ECS instance's image can be replaced via changing `image_id`. When it is changed, the instance will reboot to make the change take effect.
* `instance_type` - (Required) The type of instance to start. When it is changed, the instance will reboot to make the change take effect.
* `io_optimized` - (Deprecated) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `security_groups` - (Required)  A list of security group ids to associate with.
* `availability_zone` - (Optional) The Zone to start the instance in. It is ignored and will be computed when set `vswitch_id`.
* `instance_name` - (Optional) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen, and must not begin with http:// or https://. If not specified, 
Terraform will autogenerate a default name is `ECS-Instance`.
* `allocate_public_ip` - (Deprecated) It has been deprecated from version "1.7.0". Setting "internet_max_bandwidth_out" larger than 0 can allocate a public ip address for an instance.
* `system_disk_category` - (Optional,ForceNew) Valid values are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`. `cloud` only is used to some none I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_name` - (Optional, Available in 1.101.0+) The name of the system disk. The name must be 2 to 128 characters in length and can contain letters, digits, periods (.), colons (:), underscores (_), and hyphens (-). It must start with a letter and cannot start with http:// or https://.
* `system_disk_description` - (Optional, Available in 1.101.0+) The description of the system disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `system_disk_size` - (Optional) Size of the system disk, measured in GiB. Value range: [20, 500]. The specified value must be equal to or greater than max{20, Imagesize}. Default value: max{40, ImageSize}. 
* `system_disk_performance_level` (Optional) The performance level of the ESSD used as the system disk, Valid values: `PL0`, `PL1`, `PL2`, `PL3`, Default to `PL1`;For more information about ESSD, See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/122389.htm).
* `system_disk_auto_snapshot_policy_id` - (Optional, Available in 1.73.0+, Modifiable in 1.169.0+) The ID of the automatic snapshot policy applied to the system disk.
* `system_disk_storage_cluster_id` - (Optional, ForceNew, Available 1.177.0+) The ID of the dedicated block storage cluster. If you want to use disks in a dedicated block storage cluster as system disks when you create instances, you must specify this parameter. For more information about dedicated block storage clusters, see [What is Dedicated Block Storage Cluster?](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/dedicated-block-storage-clusters-overview).
* `system_disk_encrypted` - (Optional, ForceNew, Available 1.177.0+) Specifies whether to encrypt the system disk. Valid values: `true`,`false`. Default value: `false`. **Note:** The Encrypt System Disk During Instance Creation feature is in public preview. This public preview is provided only in Hongkong Zone B, Hongkong Zone C, Singapore Zone B, and Singapore Zone C.
    - `true`: encrypts the system disk.
    - `false`: does not encrypt the system disk.
* `system_disk_kms_key_id` - (Optional, ForceNew, Available 1.177.0+) The ID of the Key Management Service (KMS) key to be used for the system disk.
* `system_disk_encrypt_algorithm` - (Optional, ForceNew, Available 1.177.0+) The algorithm to be used to encrypt the system disk. Valid values are `aes-256`, `sm4-128`. Default value is `aes-256`.
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `PayByBandwidth`, `PayByTraffic`. Default is `PayByTraffic`. At present, 'PrePaid' instance cannot change the value to "PayByBandwidth" from "PayByTraffic".
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100]. Default to 0 Mbps.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters. When it is changed, the instance will reboot to make the change take effect.
On other OSs such as Linux, the host name can contain a maximum of 64 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“. When it is changed, the instance will reboot to make the change take effect.
* `password` - (Optional, Sensitive) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. When it is changed, the instance will reboot to make the change take effect.
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored. When it is changed, the instance will reboot to make the change take effect.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating an instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set. When it is changed, the instance will reboot to make the change take effect.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. This parameter must be set unless you can create classic network instances. When it is changed, the instance will reboot to make the change take effect.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, The default is `PostPaid`.
* `resource_group_id` - (Optional, Available in 1.57.0+, Modifiable in 1.115.0+) The Id of resource group which the instance belongs.
* `period_unit` - (Optional) The duration unit that you will buy the resource. It is valid when `instance_charge_type` is 'PrePaid'. Valid value: ["Week", "Month"]. Default to "Month".
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Valid values:
    - [1-9, 12, 24, 36, 48, 60] when `period_unit` in "Month"
    - [1-3] when `period_unit` in "Week"
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.

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

* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance and to pass data into an ECS instance. From version 1.60.0, it can be update in-place. If updated, the instance will reboot to make the change take effect. Note: Not all of changes will take effect and it depends on [cloud-init module type](https://cloudinit.readthedocs.io/en/latest/topics/modules.html).
* `key_name` - (Optional, Force new resource) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional, Force new resource) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `include_data_disks` - (Optional) Whether to change instance disks charge type when changing instance charge type.
* `dry_run` - (Optional) Specifies whether to send a dry-run request. Default to false. 
    - true: Only a dry-run request is sent and no instance is created. The system checks whether the required parameters are set, and validates the request format, service permissions, and available ECS instances. If the validation fails, the corresponding error code is returned. If the validation succeeds, the `DryRunOperation` error code is returned.
    - false: A request is sent. If the validation succeeds, the instance is created.
* `private_ip` - (Optional) Instance private IP address can be specified when you creating new instance. It is valid when `vswitch_id` is specified. When it is changed, the instance will reboot to make the change take effect.
* `credit_specification` - (Optional, Available in 1.57.1+) Performance mode of the t5 burstable instance. Valid values: 'Standard', 'Unlimited'.
* `spot_strategy` - (Optional, ForceNew) The spot strategy of a Pay-As-You-Go instance, and it takes effect only when parameter `instance_charge_type` is 'PostPaid'. Value range:
    - NoSpot: A regular Pay-As-You-Go instance.
    - SpotWithPriceLimit: A price threshold for a spot instance
    - SpotAsPriceGo: A price that is based on the highest Pay-As-You-Go instance

    Default to NoSpot. Note: Currently, the spot instance only supports domestic site account.
* `spot_price_limit` - (Optional, Float, ForceNew) The hourly price threshold of a instance, and it takes effect only when parameter 'spot_strategy' is 'SpotWithPriceLimit'. Three decimals is allowed at most.
* `deletion_protection` - (Optional, true) Whether enable the deletion protection or not. Default value: `false`.
    - true: Enable deletion protection.
    - false: Disable deletion protection.
* `force_delete` - (Optional, Available in 1.18.0+) If it is true, the "PrePaid" instance will be change to "PostPaid" and then deleted forcibly.
However, because of changing instance charge type has CPU core count quota limitation, so strongly recommand that "Don't modify instance charge type frequentlly in one month".
* `auto_release_time` - (Optional, Available in 1.70.0+) The automatic release time of the `PostPaid` instance. 
The time follows the ISO 8601 standard and is in UTC time. Format: yyyy-MM-ddTHH:mm:ssZ. It must be at least half an hour later than the current time and less than 3 years since the current time. 
Set it to null can cancel automatic release attribute and the ECS instance will not be released automatically.

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
    * `performance_level` - (Optional, ForceNew) The performance level of the ESSD used as data disk:
        - `PL0`: A single ESSD can deliver up to 10,000 random read/write IOPS.
        - `PL1`: A single ESSD can deliver up to 50,000 random read/write IOPS.
        - `PL2`: A single ESSD can deliver up to 100,000 random read/write IOPS.
        - `PL3`: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
        Default to `PL1`.
    * `encrypted` -(Optional, Bool, ForceNew) Encrypted the data in this disk. Default value: `false`.
    * `kms_key_id` - (Optional, Available in 1.90.1+) The KMS key ID corresponding to the Nth data disk.
    * `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
    * `auto_snapshot_policy_id` - (Optional, ForceNew, Available in 1.73.0+) The ID of the automatic snapshot policy applied to the system disk.
    * `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on cloud, cloud_efficiency, cloud_essd, cloud_ssd disk. If the category of this data disk was ephemeral_ssd, please don't set this param. Default value: `true`.
    * `description` - (Optional, ForceNew) The description of the data disk.
* `status` - (Optional 1.85.0+) The instance status. Valid values: ["Running", "Stopped"]. You can control the instance start and stop through this parameter. Default to `Running`.
* `hpc_cluster_id` - (Optional, ForceNew, Available in 1.144.0+) The ID of the Elastic High Performance Computing (E-HPC) cluster to which to assign the instance.
* `secondary_private_ips` - (Optional, Available in 1.144.0+) A list of Secondary private IP addresses which is selected from within the CIDR block of the VSwitch.
* `secondary_private_ip_address_count` - (Optional, Available in 1.145.0+) The number of private IP addresses to be automatically assigned from within the CIDR block of the vswitch. **NOTE:** To assign secondary private IP addresses, you must specify `secondary_private_ips` or `secondary_private_ip_address_count` but not both.
* `deployment_set_id` - (Optional, Available in 1.176.0+) The ID of the deployment set to which to deploy the instance. **NOTE:** From version 1.176.0, instance's deploymentSetId can be removed when 'deployment_set_id' = "".
* `operator_type` - (Optional, Available in 1.164.0+) The operation type. It is valid when `instance_charge_type` is `PrePaid`. Default value: `upgrade`. Valid values: `upgrade`, `downgrade`. **NOTE:**  When the new instance type specified by the `instance_type` parameter has lower specifications than the current instance type, you must set `operator_type` to `downgrade`.
* `stopped_mode` - (Optional,Available in 1.170.0+ ) The stop mode of the pay-as-you-go instance. Valid values: `StopCharging`,`KeepCharging`. Default value: If the prerequisites required for enabling the economical mode are met, and you have enabled this mode in the ECS console, the default value is `StopCharging`. For more information, see "Enable the economical mode" in [Economical mode](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/economical-mode). Otherwise, the default value is `KeepCharging`. **Note:** `Not-applicable`: Economical mode is not applicable to the instance.`
    * `KeepCharging` : standard mode. Billing of the instance continues after the instance is stopped, and resources are retained for the instance.
    * `StopCharging` : economical mode. Billing of some resources of the instance stops after the instance is stopped. When the instance is stopped, its resources such as vCPUs, memory, and public IP address are released. You may be unable to restart the instance if some types of resources are out of stock in the current region.

-> **NOTE:** System disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

-> **NOTE:** From version 1.5.0, instance's charge type can be changed to "PrePaid" by specifying `period` and `period_unit`, but it is irreversible.

-> **NOTE:** From version 1.5.0, instance's private IP address can be specified when creating VPC network instance.

-> **NOTE:** From version 1.5.0, instance's vswitch and private IP can be changed in the same availability zone. When they are changed, the instance will reboot to make the change take effect.

-> **NOTE:** From version 1.7.0, setting "internet_max_bandwidth_out" larger than 0 can allocate a public IP for an instance.
 Setting "internet_max_bandwidth_out" to 0 can release allocated public IP for VPC instance(For Classic instnace, its public IP cannot be release once it allocated, even thougth its bandwidth out is 0).
 However, at present, 'PrePaid' instance cannot narrow its max bandwidth out when its 'internet_charge_type' is "PayByBandwidth".

-> **NOTE:** From version 1.7.0, instance's type can be changed. When it is changed, the instance will reboot to make the change take effect.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `public_ip` - The instance public ip.
* `deployment_set_group_no` - (Optional, Available in 1.149.0+) The group number of the instance in a deployment set when the deployment set is use.

### Timeouts

-> **NOTE:** Available in 1.46.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the instance (until it reaches the initial `Running` status).
  `Note`: There are extra at most 2 minutes used to retry to avoid some needless API errors, and it is not in the timeouts configure.
* `update` - (Defaults to 10 mins) Used when stopping and starting the instance when necessary during update - e.g. when changing instance type, password, image, vswitch and private IP.
* `delete` - (Defaults to 20 mins) Used when terminating the instance. `Note`: There are extra at most 5 minutes used to retry to avoid some needless API errors, and it is not in the timeouts configure.



## Import

Instance can be imported using the id, e.g.

```
$ terraform import alicloud_instance.example i-abc12345678
```
