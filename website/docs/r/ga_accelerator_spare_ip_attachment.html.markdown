---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator_spare_ip_attachment"
sidebar_current: "docs-alicloud-resource-ga-accelerator-spare-ip-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Accelerator Spare Ip Attachment resource.
---

# alicloud\_ga\_accelerator\_spare\_ip\_attachment

Provides a Global Accelerator (GA) Accelerator Spare Ip Attachment resource.

For information about Global Accelerator (GA) Accelerator Spare Ip Attachment and how to use it, see [What is Accelerator Spare Ip Attachment](https://help.aliyun.com/document_detail/262120.html).

-> **NOTE:** Available in v1.167.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration         = 1
  spec             = "1"
  accelerator_name = var.name
  auto_use_coupon  = true
  description      = var.name
}
resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = var.name
  auto_pay               = true
  auto_use_coupon        = true
}
resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_accelerator_spare_ip_attachment" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  spare_ip       = "127.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `dry_run` - (Optional) The dry run.
* `spare_ip` - (Required, ForceNew) The standby IP address of CNAME. When the acceleration area is abnormal, the traffic is switched to the standby IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Accelerator Spare Ip Attachment. The value formats as `<accelerator_id>:<spare_ip>`.
* `status` - The status of the standby CNAME IP address.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Accelerator Spare Ip Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Accelerator Spare Ip Attachment.

## Import

Global Accelerator (GA) Accelerator Spare Ip Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ga_accelerator_spare_ip_attachment.example <accelerator_id>:<spare_ip>
```