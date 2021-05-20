---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instance"
sidebar_current: "docs-alicloud-resource-cr-ee-instance"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise Edition instances.
---

# alicloud\_cr\_ee\_instance

This resource will help you to manager Container Registry Enterprise Edition instances.

For information about Container Registry Enterprise Edition instances and how to use it, see [Create a Instance](https://www.alibabacloud.com/help/en/doc-detail/208144.htm)

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```
resource "alicloud_cr_ee_instance" "my-instance" {
  payment_type        = "Subscription"
  period              = 12
  renew_period        = 12
  renewal_status      = "AutoRenewal"
  instance_type       = "Standard"
  instance_name       = "test"
}
```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Optional, String) Subscription of Container Registry Enterprise Edition instance. Default value: `Subscription`. Valid values: `Subscription`.
* `period` - (Optional, Int) Service time of Container Registry Enterprise Edition instance. Default value: `12`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`, `48`, `60`. Unit: `month`.
* `renew_period` - (Optional, Int) Renewal period of Container Registry Enterprise Edition instance. Unit: `month`.
* `renewal_status` - (Optional, String) Renewal status of Container Registry Enterprise Edition instance. Valid values: `AutoRenewal`, `ManualRenewal`.
* `instance_type` - (Required, String) Type of Container Registry Enterprise Edition instance. Valid values: `Basic`, `Standard`, `Advanced`.
* `instance_name` - (Required, String) Name of Container Registry Enterprise Edition instance.
* `custom_oss_bucket` - (Optional, String) Name of your customized oss bucket. Use this bucket as instance storage if set.

## Attributes Reference

The following attributes are exported:

* `id` - ID of Container Registry Enterprise Edition instance.
* `status` - Status of Container Registry Enterprise Edition instance.
* `created_time` - Time of Container Registry Enterprise Edition instance creation.
* `end_time` - Time of Container Registry Enterprise Edition instance expiration.

## Import

Container Registry Enterprise Edition instance can be imported using the `id`, e.g.

```
$ terraform import alicloud_cr_ee_instance.default cri-test
```
