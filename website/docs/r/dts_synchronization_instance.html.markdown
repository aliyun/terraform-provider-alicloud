---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_synchronization_instance"
sidebar_current: "docs-alicloud-resource-dts-synchronization-instance"
description: |-
  Provides a Alicloud DTS Synchronization Instance resource.
---

# alicloud_dts_synchronization_instance

Provides a DTS Synchronization Instance resource.

For information about DTS Synchronization Instance and how to use it, see [What is Synchronization Instance](https://www.alibabacloud.com/help/en/doc-detail/130744.html).

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dts_synchronization_instance&exampleId=77a16b5d-60f4-48fc-9661-0552ec0ae37a01bd196e&activeTab=example&spm=docs.r.dts_synchronization_instance.0.77a16b5d60&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dts_synchronization_instance&spm=docs.r.dts_synchronization_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`.
* `payment_duration_unit` - (Optional) The payment duration unit. Valid values: `Month`, `Year`. When `payment_type` is `Subscription`, this parameter is valid and must be passed in.
* `payment_duration` - (Optional) The duration of prepaid instance purchase. this parameter is required When `payment_type` equals `Subscription`.
* `source_endpoint_region` - (Required, ForceNew) The region of source instance.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source endpoint engine. Valid values: `ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`. For the correspondence between the supported source and target libraries, see [Supported Databases, Synchronization Initialization Types and Synchronization Topologies](https://help.aliyun.com/document_detail/130744.html), [Supported Databases and Migration Types](https://help.aliyun.com/document_detail/26618.html).
* `destination_endpoint_region` - (Required, ForceNew) The region of destination instance. List of [supported regions](https://help.aliyun.com/document_detail/141033.html).
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination engine. Valid values: `ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`. For the correspondence between the supported source and target libraries, see [Supported Databases, Synchronization Initialization Types and Synchronization Topologies](https://help.aliyun.com/document_detail/130744.html), [Supported Databases and Migration Types](https://help.aliyun.com/document_detail/26618.html).
* `instance_class` - (Optional) The instance class. Valid values: `xmicro`, `micro`, `small`, `medium`, `large`, `xlarge`, `2xlarge`, `4xlarge`, `6xlarge`, `8xlarge`. **Note: The parameter is immutable after resource creation.** The class is applied only when the instance is created; this resource does not change the class of an existing instance, and an update to this argument is not applied. To upgrade the class of the synchronization instance, change `instance_class` on the corresponding `alicloud_dts_synchronization_job`, and keep this argument consistent with the resulting class.
* `sync_architecture` - (Optional) The sync architecture. Valid values: `oneway`, `bidirectional`.
* `compute_unit` - (Optional) [ETL specifications](https://help.aliyun.com/document_detail/212324.html). The unit is the computing unit ComputeUnit (CU), 1CU=1vCPU+4 GB memory. The value range is an integer greater than or equal to 2.
* `database_count` - (Optional) The number of private customized RDS instances under PolarDB-X. The default value is 1. This parameter needs to be passed only when `source_endpoint_engine_name` equals `drds`.
* `auto_pay` - (Optional, String, Computed) Whether to automatically renew the instance when it expires. Valid values: `true`, `false`. This argument only takes effect when `payment_type` is `Subscription`; for `PayAsYouGo` instances it is ignored and not set by the provider. When `payment_type` is `Subscription`, the provider reads the actual renewal status from BSS on every refresh, including refreshes of instances whose configuration does not set `auto_pay`, so the RAM user needs `bss:DescribeInstances`; if that lookup fails, the resource read fails. Setting it to `true` switches the instance to `AutoRenewal` and, when an existing instance is updated, requires `payment_duration` and `payment_duration_unit` to define the renewal period. In this case `payment_duration` must be `1`, `2`, `3`, `6`, or `12`, and `payment_duration_unit` must be `Month` or `Year`; these values can differ from the purchase duration accepted when the instance is created, where the renewal setting is passed through `AutoPay`. Setting it to `false` switches the instance to `ManualRenewal`, so the instance is retained but is not renewed automatically at expiry. BSS reports three renewal states (`AutoRenewal`, `ManualRenewal` and `NotRenewal`); `auto_pay` is read as `false` for both `ManualRenewal` and `NotRenewal`, so the two states are not distinguishable in the state. The renewal setting is sent to BSS only when `auto_pay` changes: because an instance in `NotRenewal` is already read as `false`, adding `auto_pay = "false"` to its configuration produces no change and does not trigger an update, so the instance is still released at expiry. To retain such an instance, set `auto_pay = "true"` first and then back to `false`, or change the renewal setting in the console. The renewal period is applied when `auto_pay` is set or changed; changing only `payment_duration` or `payment_duration_unit` on an existing instance does not update the renewal setting. Changing `auto_pay` on an existing subscription instance updates its auto-renew setting in place. The RAM user must be granted `bss:ModifyPrepaidInstanceAutoRenew` to change the setting and `bss:DescribeInstances` to read it back. For an existing subscription instance whose configuration already sets `auto_pay`, the next apply reconciles the instance with the configured value, which can overwrite a renewal setting changed in the console.
* `auto_start` - (Optional) Whether to automatically start the task after the purchase completed. Valid values: `true`, `false`.
* `quantity` - (Optional) The number of instances purchased.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Synchronization Instance.
* `status` - The status.

## Import

DTS Synchronization Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_synchronization_instance.example <id>
```
