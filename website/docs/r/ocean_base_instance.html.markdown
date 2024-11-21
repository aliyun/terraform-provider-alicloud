---
subcategory: "Ocean Base"
layout: "alicloud"
page_title: "Alicloud: alicloud_ocean_base_instance"
description: |-
  Provides a Alicloud Ocean Base Instance resource.
---

# alicloud_ocean_base_instance

Provides a Ocean Base Instance resource.



For information about Ocean Base Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/apsaradb-for-oceanbase/latest/what-is-oceanbase-database).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ocean_base_instance&exampleId=75010dde-2a44-9c9d-adba-6ff7445c55039eeb54c9&activeTab=example&spm=docs.r.ocean_base_instance.0.75010dde2a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ocean_base_instance" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zones = [
    "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}",
    "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}",
    "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"
  ]
  auto_renew         = "false"
  disk_size          = "100"
  payment_type       = "PayAsYouGo"
  instance_class     = "8C32G"
  backup_retain_mode = "delete_all"
  series             = "normal"
  instance_name      = var.name
}
```

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Whether to automatically renew.

  It takes effect when the parameter ChargeType is PrePaid. Value range:
  - true: automatic renewal.
  - false (default): no automatic renewal.
* `auto_renew_period` - (Optional, Int) The duration of each auto-renewal. When the value of the AutoRenew parameter is True, this parameter is required.
  - PeriodUnit is Week, AutoRenewPeriod is {"1", "2", "3"}.
  - PeriodUnit is Month, AutoRenewPeriod is {"1", "2", "3", "6", "12"}.
* `backup_retain_mode` - (Optional) The backup retention policy after the cluster is deleted. The values are as follows:
  - receive_all: Keep all backup sets;
  - delete_all: delete all backup sets;
  - receive_last: Keep the last backup set.

-> **NOTE:**   The default value is delete_all.
* `cpu` - (Optional, ForceNew, Computed, Int, Available since v1.230.0) The number of CPU cores of the cluster.
* `cpu_arch` - (Optional, ForceNew, Available since v1.230.0) Cpu architecture, x86, arm. If no, the default value is x86

* `disk_size` - (Required, Int) The size of the storage space, in GB.

  The limits of storage space vary according to the cluster specifications, as follows:
  - 8C32GB:100GB ~ 10000GB
  - 14C70GB:200GB ~ 10000GB
  - 30C180GB:400GB ~ 10000GB
  - 62C400G:800GB ~ 10000GB.

  The default value of each package is its minimum value.
* `disk_type` - (Optional, ForceNew, Computed) The storage type of the cluster. Effective only in the standard cluster version (cloud disk).

  Two types are currently supported:
  - cloud_essd_pl1: cloud disk ESSD pl1.
  - cloud_essd_pl0: cloud disk ESSD pl0. The default value is cloud_essd_pl1.
* `instance_class` - (Required) Cluster specification information. Note Please enter the shape as xCxxG, not xCxxGB

  The x86 cluster architecture currently supports the following packages:
  - 4C16G:4 core 16GB
  - 8C32G:8 core 32GB
  - 14C70G:14 core 70GB
  - 24C120G:24 core 120GB
  - 30C180G:30 core 180GB
  - 62C400G:62 core 400GB
  - 104C600G:104 core 600GB
  - 16C70G:16 core 70GB
  - 32C160G:32 core 160GB
  - 64C380G:64 core 380GB
  - 20C32G:20 core 32GB
  - 40C64G:40 core 64GB
  - 16C32G:16 core 32GB
  - 32C70G:32 core 70GB
  - 64C180G:64 core 180GB
  - 32C180G:32 core 180GB
  - 64C400G:64 core 400GB,

  The cluster architecture of arm currently supports the following packages:
  - 8C32G:8 core 32GB
  - 16C70G:16 core 70GB
  - 32C180G:32 core 180GB
* `instance_name` - (Optional, Computed) OceanBase cluster name.

  The length is 1 to 20 English or Chinese characters.

  If this parameter is not specified, the default value is the InstanceId of the cluster.
* `node_num` - (Optional, Computed) The number of nodes in the cluster. If the deployment mode is n-n-n, the number of nodes is n * 3
* `ob_version` - (Optional, ForceNew, Computed) The OceanBase Server version number.
* `payment_type` - (Required, ForceNew) The payment method of the instance. Value range:
  - Subscription: Package year and month. When you select this type of payment method, you must make sure that your account supports balance payment or credit payment. Otherwise, an InvalidPayMethod error message will be returned. 
  - PayAsYouGo (default): Pay-as-you-go (default hourly billing).
* `period` - (Optional, Int) The duration of the resource purchase. The unit is specified by the PeriodUnit. The parameter InstanceChargeType takes effect only when the value is PrePaid and is required. Once the DedicatedHostId is specified, the value cannot exceed the subscription duration of the dedicated host. When PeriodUnit = Week, Period values: {"1", "2", "3", "4"}. When PeriodUnit = Month, Period values: {"1", "2", "3", "4", "5", "6", "7", "8", "9", "12", "24", "36", "48", "60"}.
* `period_unit` - (Optional) The duration of the purchase of resources.

  Package year and Month value range: Month.

  Default value: Month of the package, which is billed by volume. The default period is Hour.
* `primary_instance` - (Optional, ForceNew, Available since v1.230.0) The ID of the primary instance.
* `primary_region` - (Optional, ForceNew, Available since v1.230.0) The primary instance Region.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the enterprise resource group to which the instance resides.
* `series` - (Required, ForceNew) Series of OceanBase cluster instances-normal (default): Standard cluster version (cloud disk)-normal_SSD: Standard cluster version (local disk)-history: history Library cluster version.
* `zones` - (Required, ForceNew, Set) Information about the zone where the cluster is deployed.
* `upgrade_spec_native` - (Optional, Available since v1.230.0) Valid values:
  - false: migration and configuration change.
  - true: in-situ matching

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `commodity_code` - The product code of the OceanBase cluster._oceanbasepre_public_cn: Domestic station cloud database package Year-to-month package._oceanbasepost_public_cn: The domestic station cloud database is paid by the hour._obpre_public_intl: International Station Cloud Database Package Monthly Package.
* `create_time` - The creation time of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.
* `update` - (Defaults to 80 mins) Used when update the Instance.

## Import

Ocean Base Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ocean_base_instance.example <id>
```