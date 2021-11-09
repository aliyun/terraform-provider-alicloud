---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator"
sidebar_current: "docs-alicloud-resource-ga-accelerator"
description: |-
  Provides a Alicloud Global Accelerator (GA) Accelerator resource.
---

# alicloud\_ga\_accelerator

Provides a Global Accelerator (GA) Accelerator resource.

For information about Global Accelerator (GA) Accelerator and how to use it, see [What is Accelerator](https://help.aliyun.com/document_detail/153235.html).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_name` - (Optional) The Name of the GA instance.
* `auto_use_coupon` - (Optional) Use coupons to pay bills automatically. Default value is `false`. Valid value: `true`: Use, `false`: Not used.
* `description` - (Optional) Descriptive information of the global acceleration instance.
* `duration` - (Required, ForceNew) The duration. The value range is 1-9.
* `spec` - (Required) The instance type of the GA instance. Specification of global acceleration instance, value:
    `1`: Small 1.
    `2`: Small 2.
    `3`: Small 3.
    `5`: Medium 1.
    `8`: Medium 2.
    `10`: Medium 3.
    
### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ga Accelerator.
* `update` - (Defaults to 6 mins) Used when updating the Ga Accelerator.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Accelerator. Value as `accelerator_id`.
* `status` - The status of the GA instance.

## Import

Ga Accelerator can be imported using the id, e.g.

```
$ terraform import alicloud_ga_accelerator.example <accelerator_id>
```
