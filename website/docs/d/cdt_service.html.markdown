---
subcategory: "CDT"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdt_service"
description: |-
  Provides a Alicloud CDT Internet Service resource.
---

# alicloud_cdt_service

Provides a CDT Internet Service resource. CDT public network service activation.

For information about CDT Internet Service and how to use it, see [What is Internet Service](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cdt_service" "open" {
}
```

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Open status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `read` - (Defaults to 5 mins) Used when create the Internet Service.

## Import

CDT Internet Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdt_service.example <id>
```