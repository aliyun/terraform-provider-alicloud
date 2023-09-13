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
  instance_class     = "8C32GB"
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
* `auto_renew_period` - (Optional) The duration of each auto-renewal. When the value of the AutoRenew parameter is True, this parameter is required.
  - PeriodUnit is Week, AutoRenewPeriod is {"1", "2", "3"}.
  - PeriodUnit is Month, AutoRenewPeriod is {"1", "2", "3", "6", "12"}.
* `backup_retain_mode` - (Optional) The backup retention policy after the cluster is deleted. The values are as follows:
  - receive_all: Keep all backup sets;
  - delete_all: delete all backup sets;
  - receive_last: Keep the last backup set.
-> **NOTE:**   The default value is delete_all.
* `disk_size` - (Required) The size of the storage space, in GB.
The limits of storage space vary according to the cluster specifications, as follows:
  - 8C32GB:100GB ~ 10000GB
  - 14C70GB:200GB ~ 10000GB
  - 30C180GB:400GB ~ 10000GB
  - 62C400G:800GB ~ 10000GB.
The default value of each package is its minimum value.
* `disk_type` - (Optional, ForceNew, Computed, Available since v1.210.0) The storage type of the cluster. Effective only in the standard cluster version (cloud disk).
Two types are currently supported:
  - cloud_essd_pl1: cloud disk ESSD pl1.
  - cloud_essd_pl0: cloud disk ESSD pl0. The default value is cloud_essd_pl1.
* `instance_class` - (Required) Cluster specification information.
Four packages are currently supported:
  - 8C32GB:8 cores 32GB.
  - 14C70GB (default):14 cores 70GB.
  - 30C180GB:30 cores 180GB.
  - 62C400GB:62 cores 400GB.
* `instance_name` - (Optional, Computed) OceanBase cluster name.The length is 1 to 20 English or Chinese characters.If this parameter is not specified, the default value is the InstanceId of the cluster.
* `node_num` - (Optional, Computed) The number of nodes in the cluster. If the deployment mode is n-n-n, the number of nodes is n * 3.
* `ob_version` - (Optional, ForceNew, Computed, Available since v1.210.0) The OceanBase Server version number.
* `payment_type` - (Required, ForceNew) The payment method of the instance. Value range:
  - Subscription: Package year and month. When you select this type of payment method, you must make sure that your account supports balance payment or credit payment. Otherwise, an InvalidPayMethod error message will be returned. 
  - PayAsYouGo (default): Pay-as-you-go (default hourly billing).
* `period` - (Optional) The duration of the resource purchase. The unit is specified by the PeriodUnit. The parameter InstanceChargeType takes effect only when the value is PrePaid and is required. Once the DedicatedHostId is specified, the value cannot exceed the subscription duration of the dedicated host. When PeriodUnit = Week, Period values: {"1", "2", "3", "4"}. When PeriodUnit = Month, Period values: {"1", "2", "3", "4", "5", "6", "7", "8", "9", "12", "24", "36", "48", "60"}.
* `period_unit` - (Optional) The duration of the purchase of resources.Package year and Month value range: Month.Default value: Month of the package, which is billed by volume. The default period is Hour.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the enterprise resource group to which the instance resides.
* `series` - (Required, ForceNew) Series of OceanBase cluster instances-normal (default): Standard cluster version (cloud disk)-normal_SSD: Standard cluster version (local disk)-history: history Library cluster version.
* `zones` - (Required, ForceNew) Information about the zone where the cluster is deployed.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `commodity_code` - The product code of the OceanBase cluster._oceanbasepre_public_cn: Domestic station cloud database package Year-to-month package._oceanbasepost_public_cn: The domestic station cloud database is paid by the hour._obpre_public_intl: International Station Cloud Database Package Monthly Package.
* `cpu` - The number of CPU cores of the cluster.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

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