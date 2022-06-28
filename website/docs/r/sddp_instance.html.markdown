---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_instance"
sidebar_current: "docs-alicloud-resource-sddp-instance"
description: |-
  Provides a Alicloud Data Security Center Instance resource.
---

# alicloud\_sddp\_instance

Provides a Data Security Center Instance resource.

For information about Data Security Center Instance and how to use it, see [What is Instance](https://help.aliyun.com/product/88674.html).

-> **NOTE:** Available in v1.136.0+.

-> **NOTE:** The Data Security Center Instance is not support in the international site.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sddp_instance" "default" {
  payment_type = "Subscription"
  sddp_version = "version_company"
  sd_cbool     = "yes"
  period       = "1"
  sdc          = "3"
  ud_cbool     = "yes"
  udc          = "2000"
  dataphin     = "yes"
}

```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `Subscription`.
* `period` - (Required) The Prepaid period. Valid values: `1`, `2`, `3`, `6`,`12`,`24`.
* `renewal_status` - (Optional) Automatic renewal status. Valid values: `AutoRenewal`,`ManualRenewal`. Default Value: `ManualRenewal`.
* `renew_period` - (Optional) Automatic renewal period. **NOTE:** The `renew_period` is required under the condition that renewal_status is `AutoRenewal`,
* `logistics` - (Optional) The logistics.
* `dataphin` - (Optional) The dataphin. Valid values: `yes`,`no`.
* `dataphin_count` - (Optional) The dataphin count. Valid values: 1 to 20.
* `sddp_version` - (Required) The sddp version. Valid values: `version_audit`,`version_company`,`version_dlp`.
* `sdc` - (Required) The number of instances.
* `ud_cbool` - (Required) Whether to use OSS. Valid values: `yes`,`no`.
* `sd_cbool` - (Required) Whether to use the database. Valid values:`yes`,`no`.
* `udc` - (Required) OSS Size.
* `instance_num` - (Optional) The number of instances.
* `modify_type` - (Optional) The modify type. Valid values: `Upgrade`, `Downgrade`.  **NOTE:** The `modify_type` is required when you execute a update operation. 



## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `authed` - Whether the required RAM authorization is configured.
* `odps_set` - Whether the authorized MaxCompute (ODPS) assets.
* `oss_bucket_set` - Whether the authorized oss assets.
* `rds_set` - Whether the authorized rds assets.
* `status` - The status of the resource.

## Import

Data Security Center Instance can be imported using the id, e.g.

```
$ terraform import alicloud_sddp_instance.example <id>
```
