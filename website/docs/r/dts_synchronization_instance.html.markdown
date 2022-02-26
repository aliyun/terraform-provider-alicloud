---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_synchronization_instance"
sidebar_current: "docs-alicloud-resource-dts-synchronization-instance"
description: |-
  Provides a Alicloud DTS Synchronization Instance resource.
---

# alicloud\_dts\_synchronization\_instance

Provides a DTS Synchronization Instance resource.

For information about DTS Synchronization Instance and how to use it, see [What is Synchronization Instance](https://www.alibabacloud.com/help/en/doc-detail/130744.html).

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "PolarDB"
  source_endpoint_region           = "cn-hangzhou"
  destination_endpoint_engine_name = "ADB30"
  destination_endpoint_region      = "cn-hangzhou"
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`.
* `payment_duration_unit` - (Optional) The payment duration unit. Valid values: `Month`, `Year`. When `payment_type` is `Subscription`, this parameter is valid and must be passed in.
* `payment_duration` - (Optional) The duration of prepaid instance purchase. this parameter is required When `payment_type` equals `Subscription`.
* `source_endpoint_region` - (Required, ForceNew) The region of source instance.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source endpoint engine. Valid values: `ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`. For the correspondence between the supported source and target libraries, see [Supported Databases, Synchronization Initialization Types and Synchronization Topologies](https://help.aliyun.com/document_detail/130744.html), [Supported Databases and Migration Types](https://help.aliyun.com/document_detail/26618.html).
* `destination_endpoint_region` - (Required, ForceNew) The region of destination instance. List of [supported regions](https://help.aliyun.com/document_detail/141033.html).
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination engine. Valid values: `ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`. For the correspondence between the supported source and target libraries, see [Supported Databases, Synchronization Initialization Types and Synchronization Topologies](https://help.aliyun.com/document_detail/130744.html), [Supported Databases and Migration Types](https://help.aliyun.com/document_detail/26618.html).
* `instance_class` - (Optional) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`. You can only upgrade the configuration, not downgrade the configuration. If you downgrade the instance, you need to [submit a ticket](https://selfservice.console.aliyun.com/ticket/category/dts/today).
* `sync_architecture` - (Optional) The sync architecture. Valid values: `oneway`, `bidirectional`.
* `compute_unit` - (Optional) [ETL specifications](https://help.aliyun.com/document_detail/212324.html). The unit is the computing unit ComputeUnit (CU), 1CU=1vCPU+4 GB memory. The value range is an integer greater than or equal to 2.
* `database_count` - (Optional) The number of private customized RDS instances under PolarDB-X. The default value is 1. This parameter needs to be passed only when `source_endpoint_engine_name` equals `drds`.
* `auto_pay` - (Optional) Whether to automatically renew when it expires. Valid values: `true`, `false`.
* `auto_start` - (Optional) Whether to automatically start the task after the purchase completed. Valid values: `true`, `false`.
* `quantity` - (Optional) The number of instances purchased.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Synchronization Instance.
* `status` - The status.

## Import

DTS Synchronization Instance can be imported using the id, e.g.

```
$ terraform import alicloud_dts_synchronization_instance.example <id>
```
