---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package resource.
---

# alicloud\_ga\_bandwidth\_package

Provides a Global Accelerator (GA) Bandwidth Package resource.

For information about Global Accelerator (GA) Bandwidth Package and how to use it, see [What is Bandwidth Package](https://www.alibabacloud.com/help/en/doc-detail/153241.htm).

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}

```

## Argument Reference

The following arguments are supported:

* `auto_pay` - (Optional) The auto pay. Valid values: `false`, `true`.
* `auto_use_coupon` - (Optional) The auto use coupon. Valid values: `false`, `true`.
* `bandwidth` - (Required, ForceNew) The bandwidth value of bandwidth packet.
* `bandwidth_package_name` - (Optional) The name of the bandwidth packet.
* `bandwidth_type` - (Optional) The bandwidth type of the bandwidth. Valid values: `Advanced`, `Basic`, `Enhanced`.
* `billing_type` - (Optional, ForceNew) The billing type. Valid values: `PayBy95`, `PayByTraffic`.
* `cbn_geographic_region_ida` - (Optional, ForceNew) Interworking area A of cross domain acceleration package. Only international stations support returning this parameter. Default value is `China-mainland`.
* `cbn_geographic_region_idb` - (Optional, ForceNew) Interworking area B of cross domain acceleration package. Only international stations support returning this parameter. Default value is `Global`.
* `description` - (Optional) The description of bandwidth package.
* `duration` - (Optional, ForceNew) The duration.
* `payment_type` - (Optional, ForceNew) The payment type of the bandwidth. Valid values: `PayAsYouGo`, `Subscription`. Default value is `Subscription`.
* `ratio` - (Optional, ForceNew) The ratio.
* `type` - (Required, ForceNew) The type of the bandwidth packet. China station only supports return to basic. Valid values: `Basic`, `CrossDomain`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package.
* `status` - The status of the bandwidth plan.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Bandwidth Package.
* `update` - (Defaults to 2 mins) Used when updating the Bandwidth Package.

## Import

Ga Bandwidth Package can be imported using the id, e.g.

```
$ terraform import alicloud_ga_bandwidth_package.example <id>
```