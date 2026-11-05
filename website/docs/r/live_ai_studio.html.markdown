---
subcategory: "Live"
layout: "alicloud"
page_title: "Alicloud: alicloud_live_ai_studio"
description: |-
  Provides a Alicloud Live AI Studio resource.
---

# alicloud_live_ai_studio

Provides a Live AI Studio resource. An AI Studio is a virtual broadcast template that composes a background, a multimedia layer and a matting (keying) configuration for cloud-based live production.

For information about Live AI Studio and how to use it, see [What is LiveAIStudio](https://www.alibabacloud.com/help/en/live/developer-reference/api-live-2016-11-01-createliveaistudio).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_live_ai_studio&activeTab=example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_live_ai_studio" "default" {
  studio_name  = var.name
  description  = "terraform-example"
  width        = 1920
  height       = 1080
  matting_type = "complex"
}
```

## Argument Reference

The following arguments are supported:

* `studio_name` - (Required, ForceNew) The name of the virtual broadcast template. It must be unique within the region. `ModifyLiveAIStudio` locates the studio by `studio_name` and exposes no rename API, so changing this field forces replacement.
* `matting_type` - (Required) The matting (keying) type. Valid values: `green` (green-screen keying), `blue` (blue-screen keying), `complex` (real-scene keying).
* `description` - (Optional) The custom description of the template.
* `width` - (Optional) The preview canvas width, in pixels (px).
* `height` - (Optional) The preview canvas height, in pixels (px).
* `background_type` - (Optional) The background material type. Valid values: `VOD` (on-demand video), `PIC` (image), `LIVE` (live stream).
* `background_resource_id` - (Optional) The on-demand resource id of the background material. Mutually exclusive with `background_resource_url`.
* `background_resource_url` - (Optional) The access URL of the background material. Mutually exclusive with `background_resource_id`.
* `media_type` - (Optional) The multimedia material type. Valid values: `VOD` (on-demand video), `PIC` (image), `LIVE` (live stream).
* `media_resource_id` - (Optional) The on-demand resource id of the multimedia material. Mutually exclusive with `media_resource_url`.
* `media_resource_url` - (Optional) The access URL of the multimedia material. Mutually exclusive with `media_resource_id`.
* `media_layout` - (Optional) The layout of the multimedia material, a JSON object serialized as a string. The object supports `positionX`, `positionY` and `heightNormalized`, each a normalized value in `[0, 1]`. Example: `"{\"positionX\":\"0\",\"positionY\":\"0\",\"heightNormalized\":\"0.5\"}"`.
* `matting_layout` - (Optional) The layout of the matting layer, a JSON object serialized as a string with the same shape as `media_layout`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id in Terraform, equal to `studio_id`.
* `studio_id` - The id of the AI Studio assigned by the cloud (`StudioId` returned by `CreateLiveAIStudio`, echoed as `TemplateId` by `DescribeLiveAIStudio`).
* `rule_ids` - The ids of the rules attached to the studio.

## Import

Live AI Studio can be imported using the `id` (the cloud `StudioId`), e.g.

```shell
$ terraform import alicloud_live_ai_studio.example <studio_id>
```
