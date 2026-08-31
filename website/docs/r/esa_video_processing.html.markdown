---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_video_processing"
description: |-
  Provides a Alicloud ESA Video Processing resource.
---

# alicloud_esa_video_processing

Provides a ESA Video Processing resource.



For information about ESA Video Processing and how to use it, see [What is Video Processing](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateVideoProcessing).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_video_processing&exampleId=6790bf73-ec92-cb7e-a8c5-8b6cc4e0d48463c00250&activeTab=example&spm=docs.r.esa_video_processing.0.6790bf73ec&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

resource "alicloud_esa_video_processing" "default" {
  video_seek_enable   = "on"
  rule_enable         = "on"
  mp4_seek_end        = "end"
  flv_seek_start      = "start"
  rule                = "(http.host eq \"video.example.com\")"
  flv_video_seek_mode = "by_byte"
  mp4_seek_start      = "start"
  flv_seek_end        = "end"
  site_id             = alicloud_esa_site.default.id
  sequence            = "1"
  site_version        = "0"
  rule_name           = "example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_video_processing&spm=docs.r.esa_video_processing.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `flv_seek_end` - (Optional) Custom FLV end parameters.
* `flv_seek_start` - (Optional) Custom FLV start parameters.
* `flv_video_seek_mode` - (Optional) FLV drag mode. Value range:
  - `by_byte`: Drag by byte.
  - `by_time`: Drag by time.
* `mp4_seek_end` - (Optional) Custom mp4 end parameters.
* `mp4_seek_start` - (Optional) Custom mp4 start parameters.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - `on`: open.
  - `off`: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int) The rule execution order prioritizes lower numerical values. It is only applicable when setting or modifying the order of individual rule configurations.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
* `video_seek_enable` - (Optional) Drag and drop the play function switch. Value range:
  - `on`: open.
  - `off`: close.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Video Processing.
* `delete` - (Defaults to 5 mins) Used when delete the Video Processing.
* `update` - (Defaults to 5 mins) Used when update the Video Processing.

## Import

ESA Video Processing can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_video_processing.example <site_id>:<config_id>
```