---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_image_transform"
description: |-
  Provides a Alicloud ESA Image Transform resource.
---

# alicloud_esa_image_transform

Provides a ESA Image Transform resource.



For information about ESA Image Transform and how to use it, see [What is Image Transform](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateImageTransform).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "imagetransform.tf.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
}

resource "alicloud_esa_image_transform" "default" {
  rule         = "http.host eq \"video.example.com\""
  site_version = "0"
  rule_name    = "rule_example"
  site_id      = alicloud_esa_site.default.id
  rule_enable  = "off"
  enable       = "off"
}
```

## Argument Reference

The following arguments are supported:
* `enable` - (Optional) Indicates whether the image transformations feature is enabled. Valid values:
  - `on`: on
  - `off`: off
* `rule` - (Optional) The rule content, which is a policy or conditional expression.
* `rule_enable` - (Optional) Indicates whether the rule is enabled. Valid values:
  - `on`: on
  - `off`: off
* `rule_name` - (Optional) Rule name, you can find out the rule whose rule name is the passed field.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the website.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Transform.
* `delete` - (Defaults to 5 mins) Used when delete the Image Transform.
* `update` - (Defaults to 5 mins) Used when update the Image Transform.

## Import

ESA Image Transform can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_image_transform.example <site_id>:<config_id>
```