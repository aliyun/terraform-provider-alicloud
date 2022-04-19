---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package Attachment resource.
---

# alicloud\_ga\_bandwidth\_package\_attachment

Provides a Global Accelerator (GA) Bandwidth Package Attachment resource.

For information about Global Accelerator (GA) Bandwidth Package Attachment and how to use it, see [What is Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/doc-detail/153241.htm).

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}

```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required) The ID of the Global Accelerator instance from which you want to disassociate the bandwidth plan.
* `bandwidth_package_id` - (Required, ForceNew) The ID of the bandwidth plan to disassociate.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package Attachment. Value as `<accelerator_id>:<bandwidth_package_id>`. Before version 1.120.0, the value is `<bandwidth_package_id>`.
* `accelerators` - Accelerators bound with current Bandwidth Package.
* `status` - State of Bandwidth Package.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Bandwidth Package Attachment.

## Import

Ga Bandwidth Package Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ga_bandwidth_package_attachment.example <accelerator_id>:<bandwidth_package_id>
```
