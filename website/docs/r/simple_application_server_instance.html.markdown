---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_instance"
sidebar_current: "docs-alicloud-resource-simple-application-server-instance"
description: |-
  Provides a Alicloud Simple Application Server Instance resource.
---

# alicloud\_simple\_application\_server\_instance

Provides a Simple Application Server Instance resource.

For information about Simple Application Server Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/doc-detail/190440.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_images" "default" {}
data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

```

## Argument Reference

The following arguments are supported:

* `auto_renew` - (Optional) Specifies whether to enable auto-renewal. Unit: months. Valid values: `true` and `false`.
* `auto_renew_period` - (Optional) The auto renew period. Valid values: `1`,`3`, `6`, `12`, `24`, `36`. **NOTE:** The attribute `auto_renew` is valid when the attribute is `true`.
* `data_disk_size` - (Optional) The size of the data disk. Unit: GB. Valid values: `0` to `16380`.
* `image_id` - (Required) The ID of the image.  You can use the `alicloud_simple_application_server_images` to query the available images in the specified region. The value must be an integral multiple of 20.
* `instance_name` - (Optional) The name of the simple application server.
* `password` - (Optional) The password of the simple application server. The password must be 8 to 30 characters in length. It must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include: `( ) ~ ! @ # $ % ^ & * - + = | { } [ ] : ; < > , . ? /`.
* `payment_type` - (Optional, Computed, ForceNew) The paymen type of the resource. Valid values: `Subscription`.
* `period` - (Required) The period. Unit: months. Valid values: `1`,`3`, `6`, `12`, `24`, `36`.
* `plan_id` - (Required) The ID of the plan. You can use the `alicloud_simple_application_server_plans`  to query all the plans provided by Simple Application Server in the specified region.
* `status` - (Optional, Computed) The status of the simple application server. Valid values: `Resetting`, `Running`, `Stopped`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Simple Application Server Instance can be imported using the id, e.g.

```
$ terraform import alicloud_simple_application_server_instance.example <id>
```
