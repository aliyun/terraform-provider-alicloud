---
subcategory: "Ocean Base"
layout: "alicloud"
page_title: "Alicloud: alicloud_ocean_base_instance"
sidebar_current: "docs-alicloud-resource-ocean-base-instance"
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
resource "alicloud_ocean_base_instance" "default" {
  instance_name  = var.name
  series         = "normal"
  disk_size      = 200
  instance_class = "14C70GB"
  zones          = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
  payment_type   = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew` - (ForceNew, Optional) Whether to automatically renew.It takes effect when the parameter ChargeType is PrePaid. Value range:
  - true: automatic renewal.
  - false (default): no automatic renewal.
* `auto_renew_period` - (Optional) The duration of each auto-renewal. When the value of the AutoRenew parameter is True, this parameter is required.-PeriodUnit is Week, AutoRenewPeriod is {"1", "2", "3"}.-PeriodUnit is Month, AutoRenewPeriod is {"1", "2", "3", "6", "12"}.
* `backup_retain_mode` - (Optional) The backup retain mode.
* `payment_type` - (Required, ForceNew) The payment method of the instance. Valid values: `PayAsYouGo`, `Subscription`.
* `disk_size` - (Required) The size of the storage space, in GB.The limits of storage space vary according to the cluster specifications, as follows:
  - 8C32GB:100GB ~ 10000GB
  - 14C70GB:200GB ~ 10000GB
  - 30C180GB:400GB ~ 10000GB
  - 62C400G:800GB ~ 10000GB.
  - The default value of each package is its minimum value.
* `instance_class` - (Required) Cluster specification information. Valid values: `14C70GB` (default), `30C180GB`, `62C400GB`, `8C32GB`, `16C70GB`, `24C120GB`, `32C160GB`, `64C380GB`, `20C32GB`, `40C64GB`, `4C16GB`.
* `instance_name` - (Optional) OceanBase cluster name. The length is 1 to 20 English or Chinese characters. If this parameter is not specified, the default value is the InstanceId of the cluster.
* `period` - (Optional) The duration of the resource purchase. The unit is specified by the PeriodUnit. The parameter `payment_type` takes effect only when the value is `Subscription` and is required. Once the DedicatedHostId is specified, the value cannot exceed the subscription duration of the dedicated host. When `period_unit` = Year, Period values: {"1", "2", "3"}. When `period_unit` = Month, Period values: {"1", "2", "3", "4", "5", "6", "7", "8", "9"}.
* `period_unit` - (Optional) The period unit. Valid values: `Month`,`Year`.
* `resource_group_id` - (Optional, ForceNew) The ID of the enterprise resource group to which the instance resides.
* `series` - (Required, ForceNew) Series of OceanBase clusters. Valid values: `normal`(default), `history`, `normal_ssd`.
* `zones` - (Required, ForceNew) Information about the zone where the cluster is deployed.
* `node_num` - (Optional) The number of nodes in the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of the resource.
* `cpu` - The number of CPU cores of the cluster.
* `commodity_code` - The product code of the OceanBase cluster.
  - oceanbase_oceanbasepre_public_cn: Domestic station cloud database package Year-to-month package.
  - oceanbase_oceanbasepost_public_cn: The domestic station cloud database is paid by the hour.
  - oceanbase_obpre_public_intl: International Station Cloud Database Package Monthly Package.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 41 mins) Used when create the Instance.
* `update` - (Defaults to 61 mins) Used when update the Instance.
* `delete` - (Defaults to 6 mins) Used when delete the Instance.

## Import

Ocean Base Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ocean_base_instance.example <id>
```