---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package Attachment resource.
---

# alicloud_ga_bandwidth_package_attachment

Provides a Global Accelerator (GA) Bandwidth Package Attachment resource.

For information about Global Accelerator (GA) Bandwidth Package Attachment and how to use it, see [What is Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-bandwidthpackageaddaccelerator).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance.
* `bandwidth_package_id` - (Required) The ID of the Bandwidth Package. **NOTE:** From version 1.192.0, `bandwidth_package_id` can be modified.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package Attachment. It formats as `<accelerator_id>:<bandwidth_package_id>`.
-> **NOTE:** Before provider version 1.120.0, it formats as `<bandwidth_package_id>`.
* `accelerators` - Accelerators bound with current Bandwidth Package.
* `status` - State of Bandwidth Package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bandwidth Package Attachment.
* `update` - (Defaults to 5 mins) Used when update the Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Bandwidth Package Attachment.

## Import

Ga Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_bandwidth_package_attachment.example <accelerator_id>:<bandwidth_package_id>
```
