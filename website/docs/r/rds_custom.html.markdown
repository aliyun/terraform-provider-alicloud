---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom"
description: |-
  Provides a Alicloud RDS Custom resource.
---

# alicloud_rds_custom

Provides a RDS Custom resource.

Dedicated RDS User host.

For information about RDS Custom and how to use it, see [What is Custom](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_custom&exampleId=077cbbd5-07d0-20a2-a525-2b5c4eb752adf141c4d6&activeTab=example&spm=docs.r.rds_custom.0.077cbbd507&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

variable "example_zone_id" {
  default = "cn-chengdu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = var.example_zone_id
}

resource "alicloud_vpc" "vpcId" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vSwitchId" {
  vpc_id       = alicloud_vpc.vpcId.id
  cidr_block   = "172.16.5.0/24"
  zone_id      = var.example_zone_id
  vswitch_name = format("%s1", var.name)
}

resource "alicloud_security_group" "securityGroupId" {
  vpc_id = alicloud_vpc.vpcId.id
}

resource "alicloud_ecs_deployment_set" "deploymentSet" {
  domain      = "Default"
  granularity = "Host"
  strategy    = "Availability"
}

resource "alicloud_ecs_key_pair" "KeyPairName" {
  key_pair_name = format("%s4", var.name)
}

resource "alicloud_rds_custom" "default" {
  data_disk {
    category          = "cloud_essd"
    size              = "50"
    performance_level = "PL1"
  }

  host_name         = "1731641300"
  create_mode       = "0"
  description       = var.name
  instance_type     = "mysql.x2.xlarge.6cm"
  password          = "example@12356"
  amount            = "1"
  io_optimized      = "optimized"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  deployment_set_id = alicloud_ecs_deployment_set.deploymentSet.id
  status            = "Running"
  system_disk {
    category = "cloud_essd"
    size     = "40"
  }

  auto_pay                   = "true"
  internet_max_bandwidth_out = "0"
  internet_charge_type       = "PayByTraffic"
  security_group_ids = [
    alicloud_security_group.securityGroupId.id
  ]
  instance_charge_type          = "Prepaid"
  vswitch_id                    = alicloud_vswitch.vSwitchId.id
  key_pair_name                 = alicloud_ecs_key_pair.KeyPairName.key_pair_name
  zone_id                       = var.example_zone_id
  auto_renew                    = "false"
  period                        = "1"
  image_id                      = "aliyun_2_1903_x64_20G_alibase_20240628.vhd"
  security_enhancement_strategy = "Active"
  period_unit                   = "Month"
}
```

## Argument Reference

The following arguments are supported:
* `amount` - (Required, Int) Represents the number of instances created
* `auto_pay` - (Optional) Whether to pay automatically. Value range:
  - `true` (default): automatic payment. You need to ensure that your account balance is sufficient.
  - `false`: only orders are generated without deduction.

-> **NOTE:**  If the balance of your payment method is insufficient, you can set the parameter AutoPay to false, and an unpaid order will be generated. You can log on to the RDS management console to pay by yourself.

-> **NOTE:** >

* `auto_renew` - (Optional) Whether the instance is automatically renewed. Valid values: true/false. The default is false.
* `create_mode` - (Optional) Whether to allow joining the ACK cluster. When this parameter is set to `1`, the created instance can be added to the ACK cluster through The `AttachRCInstances` API to efficiently manage container applications.
  - `1`: Yes.
  - `0` (default): No.
* `data_disk` - (Optional, ForceNew, List) Data disk See [`data_disk`](#data_disk) below.
* `deployment_set_id` - (Optional, ForceNew) The ID of the deployment set.
* `description` - (Optional, ForceNew) Instance description. It must be 2 to 256 characters in length and cannot start with http:// or https.
* `direction` - (Optional) Instance configuration type, value range:

-> **NOTE:**  This parameter does not need to be uploaded, and the system can automatically determine whether to upgrade or downgrade. If you want to upload, please follow the following logic rules.
  - `Up` (default): upgrade the instance specification. Please ensure that your account balance is sufficient.
  - `Down`: Downgrade instance specifications. When the instance type set to InstanceType is lower than the current instance type, set Direction = down.
* `dry_run` - (Optional) Whether to pre-check the operation of creating an instance. Valid values:
  - `true`: The PreCheck operation is performed without creating an instance. Check items include request parameters, request formats, business restrictions, and inventory.
  - `false` (default): Sends a normal request and directly creates an instance after the check is passed.
* `force` - (Optional) Whether to forcibly release the running instance. Value: true/false
* `force_stop` - (Optional) Whether to force shutdown. Value range:
  - `true`: Force shutdown.
  - `false` (default): Normal shutdown.
* `host_name` - (Optional) The instance host name.
* `image_id` - (Optional) The ID of the image used by the instance.
* `instance_charge_type` - (Optional) The Payment type. Currently, only `Prepaid` (package year and month) types are supported.
* `instance_type` - (Required) The type of the created RDS Custom dedicated host instance.
* `internet_charge_type` - (Optional) Reserved parameters are not supported.
* `internet_max_bandwidth_out` - (Optional, Int) Reserved parameters are not supported.
* `io_optimized` - (Optional) Reserved parameters are not supported.
* `key_pair_name` - (Optional) The key pair name. Only flyer names are supported.
* `password` - (Optional) The account and password of the instance.
* `period` - (Optional, Int) Prepaid renewal duration, unit: Month/Year. 
* `period_unit` - (Optional) The unit of duration of the year-to-month billing method. Value range:
  - `Year`: Year
  - `Month` (default): Month
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `security_enhancement_strategy` - (Optional) Reserved parameters are not supported.
* `security_group_ids` - (Optional, ForceNew, List) Security group list
* `status` - (Optional, Computed) The status of the resource
* `system_disk` - (Optional, List) System disk specifications. See [`system_disk`](#system_disk) below.
* `tags` - (Optional, Map) The tag of the resource
* `vswitch_id` - (Required, ForceNew) The ID of the virtual switch. The zone in which the vSwitch is located must correspond to the zone ID entered in ZoneId.

  The network type InstanceNetworkType must be VPC.
* `zone_id` - (Optional, ForceNew) The zone ID  of the resource

### `data_disk`

The data_disk supports the following:
* `category` - (Optional, ForceNew) Instance storage type

  local_ssd: local SSD disk

  cloud_essd:ESSD PL1 cloud disk
* `performance_level` - (Optional, ForceNew) Cloud Disk Performance

  Currently only supports PL1
* `size` - (Optional, ForceNew, Int) Instance storage space. Unit: GB.

### `system_disk`

The system_disk supports the following:
* `category` - (Optional) The cloud disk type of the system disk. Currently, only `cloud_essd`(ESSD cloud disk) is supported.
* `size` - (Optional) System disk size, unit: GiB. Only ESSD PL1 is supported. Valid values range from 20 to 2048.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID. Callable DescribeRegions to get.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom.
* `delete` - (Defaults to 5 mins) Used when delete the Custom.
* `update` - (Defaults to 5 mins) Used when update the Custom.

## Import

RDS Custom can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom.example <id>
```