---
subcategory: "Data Security Center (SDDP)"
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sddp_instance&exampleId=b9554da9-7cb5-f4b2-94c6-a3f53b08d51fa7740a78&activeTab=example&spm=docs.r.sddp_instance.0.b9554da97c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sddp_instance&spm=docs.r.sddp_instance.example&intl_lang=EN_US)

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
* `modify_type` - (Optional) The modify type. Valid values: `Upgrade`, `Downgrade`.  **NOTE:** The `modify_type` is required when you execute a update operation.
* `oss_size` - (Optional) The OSS storage capacity.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `authed` - Whether the required RAM authorization is configured.
* `odps_set` - Whether the authorized MaxCompute (ODPS) assets.
* `oss_bucket_set` - Whether the authorized oss assets.
* `rds_set` - Whether the authorized rds assets.
* `status` - The status of the resource.
* `instance_num` - The number of instances.
* `remain_days` -  The remaining days of the protection period of the assets in the current login account.

## Import

Data Security Center Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_sddp_instance.example <id>
```
