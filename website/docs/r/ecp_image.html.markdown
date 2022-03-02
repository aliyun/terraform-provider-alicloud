---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_image"
sidebar_current: "docs-alicloud-resource-ecp-image"
description: |-
  Provides a Alicloud Elastic Cloud Phone (ECP) Image resource.
---

# alicloud\_ecp\_image

Provides a Elastic Cloud Phone (ECP) Image resource.

For information about Elastic Cloud Phone (ECP) Image and how to use it,
see [What is Image](https://help.aliyun.com/document_detail/258178.html/).

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "%s"
}
data "alicloud_ecp_instances" "default" {
}

locals {
  instance_id = data.alicloud_ecp_instances.default.instances[0].instance_id
}

resource "alicloud_ecp_image" "example" {
  image_name  = var.name
  description = var.name
  instance_id = "${local.instance_id}"
  force       = "true"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the image. 2 to 256 English or Chinese characters in length and cannot start
  with `http://` and `https`.
* `force` - (Optional) The force.
* `image_name` - (Optional) The name of the mirror image.It must be 2 to 128 characters in length and must start with an
  uppercase letter or Chinese. It cannot start with http:// or https. It can contain Chinese, English, numbers,
  half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
* `instance_id` - (Required) The instance id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image.
* `status` - Mirror status.

### Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 301 mins) Used when create the Image.

## Import

Elastic Cloud Phone (ECP) Image can be imported using the id, e.g.

```
$ terraform import alicloud_ecp_image.example <id>
```