---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_url_observation"
description: |-
  Provides a Alicloud ESA Url Observation resource.
---

# alicloud_esa_url_observation

Provides a ESA Url Observation resource.

Web page monitoring.

For information about ESA Url Observation and how to use it, see [What is Url Observation](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateUrlObservation).

-> **NOTE:** Available since v1.259.0.

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
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_url_observation" "default" {
  sdk_type = "automatic"
  site_id  = alicloud_esa_site.default.id
  url      = "terraform.cn/a.html"
}
```

## Argument Reference

The following arguments are supported:
* `sdk_type` - (Required) SDK integration mode. Value:
  - `automatic`: automatic integration.
  - `manual`: manual integration.
* `site_id` - (Required, ForceNew, Int) The site ID.
* `url` - (Required, ForceNew) The URL of the page to monitor.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Url Observation.
* `delete` - (Defaults to 5 mins) Used when delete the Url Observation.
* `update` - (Defaults to 5 mins) Used when update the Url Observation.

## Import

ESA Url Observation can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_url_observation.example <site_id>:<config_id>
```