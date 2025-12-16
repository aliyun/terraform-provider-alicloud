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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_url_observation&exampleId=2684f741-8d30-ce29-9784-432d82e6f746e7bad433&activeTab=example&spm=docs.r.esa_url_observation.0.2684f7418d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_url_observation&spm=docs.r.esa_url_observation.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `sdk_type` - (Required) SDK integration mode. Value:
  - `automatic`: automatic integration.
  - `manual`: manual integration.
* `site_id` - (Required, ForceNew) The site ID.
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