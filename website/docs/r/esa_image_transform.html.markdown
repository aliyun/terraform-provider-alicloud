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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_image_transform&exampleId=bf3e942c-98a6-5fa4-39a5-f000e7836177838fb5df&activeTab=example&spm=docs.r.esa_image_transform.0.bf3e942c98&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  -   `on`: on
  -   `off`: off
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  -   `on`: open
  -   `off`: close
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

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