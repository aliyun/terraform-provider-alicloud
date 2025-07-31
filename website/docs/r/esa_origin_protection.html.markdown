---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_protection"
description: |-
  Provides a Alicloud ESA Origin Protection resource.
---

# alicloud_esa_origin_protection

Provides a ESA Origin Protection resource.



For information about ESA Origin Protection and how to use it, see [What is Origin Protection](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateOriginProtection).

-> **NOTE:** Available since v1.256.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_origin_protection" "default" {
  origin_converge = "on"
  site_id         = alicloud_esa_site.default.id
}
```

## Argument Reference

The following arguments are supported:
* `origin_converge` - (Optional) The IP convergence status.

  - on
  - off
* `site_id` - (Required, ForceNew, Int) Site Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Protection.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Protection.
* `update` - (Defaults to 5 mins) Used when update the Origin Protection.

## Import

ESA Origin Protection can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_protection.example <id>
```