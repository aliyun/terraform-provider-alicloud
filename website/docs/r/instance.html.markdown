---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instance"
sidebar_current: "docs-alicloud-resource-instance"
description: |-
  Provides an ECS instance resource.
---

# alicloud_instance

Provides a ECS instance resource.

-> **NOTE:** Available since v1.0.0

-> **NOTE:** From version v1.213.0, you can specify `launch_template_id` and `launch_template_version` to use a launch template. This eliminates the need to configure a large number of parameters every time you create instances.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_instance&exampleId=f98d5a2b-3f0d-19d9-4708-bce020dfaf3bc8d21801&activeTab=example&spm=docs.r.instance.0.f98d5a2b3f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "instance_type" {
  default = "ecs.n4.large"
}

variable "image_id" {
  default = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
}

# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  security_group_name = var.name
  description         = "foo"
  vpc_id              = alicloud_vpc.vpc.id
}

resource "alicloud_kms_key" "key" {
  description            = "Hello KMS"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
  available_instance_type     = var.instance_type
}

# Create a new ECS instance for VPC
resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = data.alicloud_zones.default.zones.0.id
  security_groups   = alicloud_security_group.group.*.id

  # series III
  instance_type              = var.instance_type
  system_disk_category       = "cloud_efficiency"
  system_disk_name           = var.name
  system_disk_description    = "test_foo_system_disk_description"
  image_id                   = var.image_id
  instance_name              = var.name
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

* `image_id` - (Optional) The Image to use for the instance. ECS instance's image can be replaced via changing `image_id`. When it is changed, the instance will reboot to make the change take effect. If you do not use `launch_template_id` or `launch_template_name` to specify a launch template, you must specify `image_id`.
* `instance_type` - (Optional) The type of instance to start. When it is changed, the instance will reboot to make the change take effect. If you do not use `launch_template_id` or `launch_template_name` to specify a launch template, you must specify `instance_type`.
* `io_optimized` - (Removed) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type.
* `security_groups` - (Optional)  A list of security group ids to associate with. If you do not use `launch_template_id` or `launch_template_name` to specify a launch template, you must specify `security_groups`.
* `availability_zone` - (Optional, ForceNew) The Zone to start the instance in. It is ignored and will be computed when set `vswitch_id`.
* `instance_name` - (Optional) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen, and must not begin with http:// or https://. **NOTE:** From version 1.243.0, the default value `ECS-Instance` will be removed.
* `allocate_public_ip` - (Deprecated) It has been deprecated from version "1.7.0". Setting "internet_max_bandwidth_out" larger than 0 can allocate a public ip address for an instance.
* `system_disk_category` - (Optional, ForceNew) Valid values are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`, `cloud_auto`, `cloud_essd_entry`. only is used to some none I/O optimized instance. Valid values `cloud_auto` Available since v1.184.0.
* `system_disk_name` - (Optional, Available since v1.101.0) The name of the system disk. The name must be 2 to 128 characters in length and can contain letters, digits, periods (.), colons (:), underscores (_), and hyphens (-). It must start with a letter and cannot start with http:// or https://.
* `system_disk_description` - (Optional, Available since v1.101.0) The description of the system disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `system_disk_size` - (Optional) Size of the system disk, measured in GiB. Value range: [20, 500]. The specified value must be equal to or greater than max{20, Imagesize}. Default value: max{40, ImageSize}.
* `system_disk_performance_level` (Optional) The performance level of the ESSD used as the system disk, Valid values: `PL0`, `PL1`, `PL2`, `PL3`, Default to `PL1`;For more information about ESSD, See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/122389.htm).
* `system_disk_auto_snapshot_policy_id` - (Optional, Available since v1.73.0, Modifiable in 1.169.0) The ID of the automatic snapshot policy applied to the system disk.
* `system_disk_storage_cluster_id` - (Optional, ForceNew, Available since v1.177.0) The ID of the dedicated block storage cluster. If you want to use disks in a dedicated block storage cluster as system disks when you create instances, you must specify this parameter. For more information about dedicated block storage clusters.
* `system_disk_encrypted` - (Optional, ForceNew, Available since v1.177.0) Specifies whether to encrypt the system disk. Valid values: `true`,`false`. Default value: `false`.
  - `true`: encrypts the system disk.
  - `false`: does not encrypt the system disk.
* `system_disk_kms_key_id` - (Optional, ForceNew, Available since v1.177.0) The ID of the Key Management Service (KMS) key to be used for the system disk.
* `system_disk_encrypt_algorithm` - (Optional, ForceNew, Available since v1.177.0) The algorithm to be used to encrypt the system disk. Valid values are `aes-256`, `sm4-128`. Default value is `aes-256`.
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `PayByBandwidth`, `PayByTraffic`. At present, 'PrePaid' instance cannot change the value to "PayByBandwidth" from "PayByTraffic". **NOTE:** From version 1.243.0, the default value `PayByTraffic` will be removed.
* `internet_max_bandwidth_in` - (Optional, Deprecated since v1.121.2) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100]. **NOTE:** From version 1.243.0, the default value `0` will be removed.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters. When it is changed, the instance will reboot to make the change take effect.
  On other OSs such as Linux, the host name can contain a maximum of 64 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“. When it is changed, the instance will reboot to make the change take effect.
* `password` - (Optional, Sensitive) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. When it is changed, the instance will reboot to make the change take effect.
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored. When it is changed, the instance will reboot to make the change take effect.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating an instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set. When it is changed, the instance will reboot to make the change take effect.
* `vpc_id` - (Optional, Available since v1.227.1) The ID of the VPC.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. This parameter must be set unless you can create classic network instances. When it is changed, the instance will reboot to make the change take effect.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`. **NOTE:** From version 1.243.0, the default value `PostPaid` will be removed.
  **NOTE:** Since 1.9.6, it can be changed each other between `PostPaid` and `PrePaid`.
  However, since [some limitation about CPU core count in one month](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/modifyinstancechargetype),
  there strongly recommends that `Don't change instance_charge_type frequentlly in one month`.

* `resource_group_id` - (Optional, Available since v1.57.0, Modifiable in 1.115.0) The Id of resource group which the instance belongs.
* `period_unit` - (Optional) The duration unit that you will buy the resource. It is valid when `instance_charge_type` is 'PrePaid'. Valid value: ["Week", "Month"]. Default to "Month".
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid and required when `instance_charge_type` is `PrePaid`. Valid values:
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
  - Key: It can be up to `128` characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to `128` characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `volume_tags` - (Optional) A mapping of tags to assign to the devices created by the instance at launch time.
  - Key: It can be up to `128` characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to `128` characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance and to pass data into an ECS instance.
  It supports to setting a [base64-encoded value](https://developer.hashicorp.com/terraform/language/functions/base64encode), and it is the recommended usage.
  From version 1.60.0, it can be updated in-place. If updated, the instance will reboot to make the change take effect.
  Note: Not all changes will take effect, and it depends on [cloud-init module type](https://cloudinit.readthedocs.io/en/latest/topics/modules.html).
* `key_name` - (Optional, ForceNew) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional, ForceNew) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `include_data_disks` - (Optional) Whether to change instance disks charge type when changing instance charge type.
* `dry_run` - (Optional) Specifies whether to send a dry-run request. Default to false.
  - true: Only a dry-run request is sent and no instance is created. The system checks whether the required parameters are set, and validates the request format, service permissions, and available ECS instances. If the validation fails, the corresponding error code is returned. If the validation succeeds, the `DryRunOperation` error code is returned.
  - false: A request is sent. If the validation succeeds, the instance is created.
* `private_ip` - (Optional) Instance private IP address can be specified when you creating new instance. It is valid when `vswitch_id` is specified. When it is changed, the instance will reboot to make the change take effect.
* `credit_specification` - (Optional, Available since v1.57.1) Performance mode of the t5 burstable instance. Valid values: 'Standard', 'Unlimited'.
* `spot_strategy` - (Optional, ForceNew) The spot strategy of a Pay-As-You-Go instance, and it takes effect only when parameter `instance_charge_type` is 'PostPaid'. Value range:
  - NoSpot: A regular Pay-As-You-Go instance.
  - SpotWithPriceLimit: A price threshold for a spot instance
  - SpotAsPriceGo: A price that is based on the highest Pay-As-You-Go instance

  Default to NoSpot. Note: Currently, the spot instance only supports domestic site account.
* `spot_price_limit` - (Optional, Float, ForceNew) The hourly price threshold of a instance, and it takes effect only when parameter 'spot_strategy' is 'SpotWithPriceLimit'. Three decimals is allowed at most.
* `deletion_protection` - (Optional, true) Whether enable the deletion protection or not. It does not work when the instance is spot. Default value: `false`.
  - true: Enable deletion protection.
  - false: Disable deletion protection.
* `force_delete` - (Optional, Available since v1.18.0) If it is true, the "PrePaid" instance will be change to "PostPaid" and then deleted forcibly.
  However, because of changing instance charge type has CPU core count quota limitation, so strongly recommand that "Don't modify instance charge type frequentlly in one month".
* `auto_release_time` - (Optional, Available since v1.70.0) The automatic release time of the `PostPaid` instance.
  The time follows the ISO 8601 standard and is in UTC time. Format: yyyy-MM-ddTHH:mm:ssZ. It must be at least half an hour later than the current time and less than 3 years since the current time.
  Setting it to null can cancel automatic release feature, and the ECS instance will not be released automatically.

* `security_enhancement_strategy` - (Optional, ForceNew) The security enhancement strategy.
  - Active: Enable security enhancement strategy, it only works on system images.
  - Deactive: Disable security enhancement strategy, it works on all images.
* `data_disks` - (Optional, ForceNew, Available since v1.23.1) The list of data disks created with instance. See [`data_disks`](#data_disks) below.
* `network_interfaces` - (Optional, ForceNew, Available since v1.212.0) The list of network interfaces created with instance. See [`network_interfaces`](#network_interfaces) below.
* `status` - (Optional 1.85.0) The instance status. Valid values: ["Running", "Stopped"]. You can control the instance start and stop through this parameter. Default to `Running`.
* `hpc_cluster_id` - (Optional, ForceNew, Available since v1.144.0) The ID of the Elastic High Performance Computing (E-HPC) cluster to which to assign the instance.
* `secondary_private_ips` - (Optional, Available since v1.144.0) A list of Secondary private IP addresses which is selected from within the CIDR block of the vSwitch.
* `secondary_private_ip_address_count` - (Optional, Available since v1.145.0) The number of private IP addresses to be automatically assigned from within the CIDR block of the vswitch. **NOTE:** To assign secondary private IP addresses, you must specify `secondary_private_ips` or `secondary_private_ip_address_count` but not both.
* `deployment_set_id` - (Optional, Available since v1.176.0) The ID of the deployment set to which to deploy the instance. **NOTE:** From version 1.176.0, instance's deploymentSetId can be removed when 'deployment_set_id' = "".
* `operator_type` - (Optional, Available since v1.164.0) The operation type. It is valid when `instance_charge_type` is `PrePaid`. Default value: `upgrade`. Valid values: `upgrade`, `downgrade`. **NOTE:**  When the new instance type specified by the `instance_type` parameter has lower specifications than the current instance type, you must set `operator_type` to `downgrade`.
* `stopped_mode` - (Optional, Available since v1.170.0) The stop mode of the pay-as-you-go instance. Valid values: `StopCharging`,`KeepCharging`, `Not-applicable`. Default value: If the prerequisites required for enabling the economical mode are met, and you have enabled this mode in the ECS console, the default value is `StopCharging`. For more information, see "Enable the economical mode" in [Economical mode](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/economical-mode). Otherwise, the default value is `KeepCharging`. **Note:** `Not-applicable`: Economical mode is not applicable to the instance.`
  * `KeepCharging`: standard mode. Billing of the instance continues after the instance is stopped, and resources are retained for the instance.
  * `StopCharging`: economical mode. Billing of some resources of the instance stops after the instance is stopped. When the instance is stopped, its resources such as vCPUs, memory, and public IP address are released. You may be unable to restart the instance if some types of resources are out of stock in the current region.
* `maintenance_time` - (Optional, Available since v1.181.0) The time of maintenance. See [`maintenance_time`](#maintenance_time) below.
* `maintenance_action` - (Optional, Available since v1.181.0) The maintenance action. Valid values: `Stop`, `AutoRecover` and `AutoRedeploy`.
  * `Stop` : stops the instance.
  * `AutoRecover` : automatically recovers the instance.
  * `AutoRedeploy` : fails the instance over, which may cause damage to the data disks attached to the instance.
* `maintenance_notify` - (Optional, Available since v1.181.0) Specifies whether to send an event notification before instance shutdown. Valid values: `true`, `false`. Default value: `false`.
  * `true` : sends an event notification.
  * `false` : does not send an event notification.
* `spot_duration` - (Optional, Available since v1.188.0) The retention time of the preemptive instance in hours. Valid values: `0`, `1`, `2`, `3`, `4`, `5`, `6`. Retention duration 2~6 is under invitation test, please submit a work order if you need to open. If the value is `0`, the mode is no protection period. Default value is `1`.
* `http_tokens` - (Optional, Available since v1.192.0) Specifies whether to forcefully use the security-enhanced mode (IMDSv2) to access instance metadata. Default value: optional. Valid values:
  - optional: does not forcefully use the security-enhanced mode (IMDSv2).
  - required: forcefully uses the security-enhanced mode (IMDSv2). After you set this parameter to required, you cannot access instance metadata in normal mode.
* `http_endpoint` - (Optional, Available since v1.192.0) Specifies whether to enable the access channel for instance metadata. Valid values: `enabled`, `disabled`. Default value: `enabled`.
* `http_put_response_hop_limit` - (Optional, ForceNew) **NOTE:**: This parameter is not available for use yet. The HTTP PUT response hop limit for accessing instance metadata. Valid values: 1 to 64. Default value: 1.
* `ipv6_address_count` - (Optional, ForceNew, Available since v1.193.0) The number of IPv6 addresses to randomly generate for the primary ENI. Valid values: 1 to 10. **NOTE:** You cannot specify both the `ipv6_addresses` and `ipv6_address_count` parameters.
* `ipv6_addresses` - (Optional, Available since v1.193.0) A list of IPv6 address to be assigned to the primary ENI. Support up to 10. **NOTE:** From version 1.241.0, `ipv6_addresses` can be modified.
* `dedicated_host_id` - (Optional, ForceNew, Available since v1.201.0) The ID of the dedicated host on which to create the instance. If you set the DedicatedHostId parameter, the `spot_strategy` and `spot_price_limit` parameters cannot be set. This is because preemptible instances cannot be created on dedicated hosts.
* `subnet_id` - (Removed since v1.210.0) The ID of the subnet. **NOTE:** Field `subnet_id` has been removed from provider version 1.210.0.
* `launch_template_id` - (Optional, ForceNew, Available since v1.213.1) The ID of the launch template. For more information, see [DescribeLaunchTemplates](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-describelaunchtemplates).To use a launch template to create an instance, you must use the `launch_template_id` or `launch_template_name` parameter to specify the launch template.
* `launch_template_name` - (Optional, ForceNew, Available since v1.213.1) The name of the launch template.
* `launch_template_version` - (Optional, ForceNew, Available since v1.213.1) The version of the launch template. If you set `launch_template_id` or `launch_template_name` parameter but do not set the version number of the launch template, the default template version is used.
* `enable_jumbo_frame` - (Optional, Bool, Available since v1.223.2) Specifies whether to enable the Jumbo Frames feature for the instance. Valid values: `true`, `false`.
* `network_interface_traffic_mode` - (Optional, ForceNew, Available since v1.227.1) The communication mode of the Primary ENI. Default value: `Standard`. Valid values:
  - `Standard`: Uses the TCP communication mode.
  - `HighPerformance`: Uses the remote direct memory access (RDMA) communication mode with Elastic RDMA Interface (ERI) enabled.
* `network_card_index` - (Optional, ForceNew, Int, Available since v1.227.1)  The index of the network card for Primary ENI.
* `queue_pair_number` - (Optional, ForceNew, Int, Available since v1.227.1) The number of queues supported by the ERI.
* `password_inherit` - (Optional, Bool, Available since v1.232.0) Specifies whether to use the password preset in the image. Default value: `false`. Valid values:
  - `true`: Uses the preset password.
  - `false`: Does not use the preset password.
-> **NOTE:** If you set `password_inherit` to `true`, make sure that you have not specified `password` or `kms_encrypted_password` and the selected image has a preset password.

* `image_options` - (Optional, Available since v1.237.0) The options of images. See [`image_options`](#image_options) below.

-> **NOTE:** System disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

-> **NOTE:** From version 1.5.0, instance's charge type can be changed to "PrePaid" by specifying `period` and `period_unit`, but it is irreversible.

-> **NOTE:** From version 1.5.0, instance's private IP address can be specified when creating VPC network instance.

-> **NOTE:** From version 1.5.0, instance's vswitch and private IP can be changed in the same availability zone. When they are changed, the instance will reboot to make the change take effect.

-> **NOTE:** From version 1.7.0, setting "internet_max_bandwidth_out" larger than 0 can allocate a public IP for an instance.
Setting "internet_max_bandwidth_out" to 0 can release allocated public IP for VPC instance(For Classic instnace, its public IP cannot be release once it allocated, even thougth its bandwidth out is 0).
However, at present, 'PrePaid' instance cannot narrow its max bandwidth out when its 'internet_charge_type' is "PayByBandwidth".

-> **NOTE:** From version 1.7.0, instance's type can be changed. When it is changed, the instance will reboot to make the change take effect.

### `data_disks`

The data_disks supports the following:

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
  - `cloud_auto`: The AutoPL cloud disk.
    Default to `cloud_efficiency`.
* `performance_level` - (Optional, ForceNew) The performance level of the ESSD used as data disk:
  - `PL0`: A single ESSD can deliver up to 10,000 random read/write IOPS.
  - `PL1`: A single ESSD can deliver up to 50,000 random read/write IOPS.
  - `PL2`: A single ESSD can deliver up to 100,000 random read/write IOPS.
  - `PL3`: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
    Default to `PL1`.
* `encrypted` -(Optional, Bool, ForceNew) Encrypted the data in this disk. Default value: `false`.
* `kms_key_id` - (Optional, Available since v1.90.1) The KMS key ID corresponding to the Nth data disk.
* `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
* `auto_snapshot_policy_id` - (Optional, ForceNew, Available since v1.73.0) The ID of the automatic snapshot policy applied to the system disk.
* `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on cloud, cloud_efficiency, cloud_essd, cloud_ssd disk. If the category of this data disk was ephemeral_ssd, please don't set this param. Default value: `true`.
* `description` - (Optional, ForceNew) The description of the data disk.
* `device` - (Optional, ForceNew, Available since v1.183.0) The mount point of the data disk.

### `network_interfaces`

The network_interfaces supports the following. Currently only one secondary ENI can be specified.

* `network_interface_id` - (Optional, ForceNew) The ID of the Secondary ENI.
* `vswitch_id` - (Optional, ForceNew, Available since v1.223.2) The ID of the vSwitch to which to connect Secondary ENI N.
* `network_interface_traffic_mode` - (Optional, ForceNew, Available since v1.223.2) The communication mode of the Secondary ENI. Default value: `Standard`. Valid values:
  - `Standard`: Uses the TCP communication mode.
  - `HighPerformance`: Uses the remote direct memory access (RDMA) communication mode with Elastic RDMA Interface (ERI) enabled.
* `network_card_index` - (Optional, ForceNew, Int, Available since v1.227.1)  The index of the network card for Secondary ENI.
* `queue_pair_number` - (Optional, ForceNew, Int, Available since v1.227.1) The number of queues supported by the ERI.
* `security_group_ids` - (Optional, ForceNew, List, Available since v1.223.2) The ID of security group N to which to assign Secondary ENI N.

### `maintenance_time`

The maintenance_time supports the following:

* `start_time` - (Optional) The start time of maintenance. The time must be on the hour at exactly 0 minute and 0 second. The `start_time` and `end_time` parameters must be specified at the same time. The `end_time` value must be 1 to 23 hours later than the `start_time` value. Specify the time in the HH:mm:ss format. The time must be in UTC+8.
* `end_time` - (Optional) The end time of maintenance. The time must be on the hour at exactly 0 minute and 0 second. The `start_time` and `end_time` parameters must be specified at the same time. The `end_time` value must be 1 to 23 hours later than the `start_time` value. Specify the time in the HH:mm:ss format. The time must be in UTC+8.

### `image_options`

The image_options supports the following:

* `login_as_non_root` - (Optional, ForceNew) Whether to allow the instance logging in with the ecs-user user.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `public_ip` - The instance public ip.
* `cpu` - The number of vCPUs.
* `memory` - The memory size of the instance. Unit: MiB.
* `os_type` - The type of the operating system of the instance.
* `os_name` - The name of the operating system of the instance.
* `network_interface_id` - The ID of the Primary ENI.
* `system_disk_id` - (Available since v1.210.0) The ID of system disk.
* `primary_ip_address` - The primary private IP address of the ENI.
* `deployment_set_group_no` - The group number of the instance in a deployment set when the deployment set is use.
* `create_time` - (Available since v1.232.0) The time when the instance was created.
* `start_time` - (Available since v1.232.0) The time when the instance was last started.
* `expired_time` - (Available since v1.232.0) The expiration time of the instance.

## Timeouts

-> **NOTE:** Available since v1.46.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the instance (until it reaches the initial `Running` status).
  `Note`: There are extra at most 2 minutes used to retry to avoid some needless API errors, and it is not in the timeouts configure.
* `update` - (Defaults to 10 mins) Used when stopping and starting the instance when necessary during update - e.g. when changing instance type, password, image, vswitch and private IP.
* `delete` - (Defaults to 20 mins) Used when terminating the instance. `Note`: There are extra at most 5 minutes used to retry to avoid some needless API errors, and it is not in the timeouts configure.

## Import

Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_instance.example i-abc12345678
```
