---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_domain"
sidebar_current: "docs-alicloud-resource-ga-domain"
description: |-
  Provides a Alicloud Ga Domain resource.
---

# alicloud_ga_domain

Provides a Ga Domain resource.

For information about Ga Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/global-accelerator/latest/createdomain).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_domain" "default" {
  domain         = "changes.com.cn"
  accelerator_id = alicloud_ga_accelerator.default.id
}
```

## Argument Reference

The following arguments are supported:
* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `domain` - (Required, ForceNew) The accelerated domain name to be added. only top-level domain names are supported, such as 'example.com'.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is formulated as `<accelerator_id>:<domain>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.

## Import

Ga Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_domain.example <accelerator_id>:<domain>
```