---
subcategory: "Live"
layout: "alicloud"
page_title: "Alicloud: alicloud_live_ai_subtitle"
description: |-
  Provides a Alicloud Live AI Subtitle resource.
---

# alicloud_live_ai_subtitle

Provides a Live AI Subtitle resource.

For information about Live AI Subtitle and how to use it, see [What is LiveAISubtitle](https://next.api.alibabacloud.com/document/live/2016-11-01/AddLiveAISubtitle).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_live_ai_subtitle" "default" {
  subtitle_name        = var.name
  src_language         = "zh-CN"
  dst_language         = "en-US"
  font_name            = "KaiTi"
  position_normalized  = [0.32, 0.27]
  word_per_line        = 20
  font_size_normalized = 0.037
  max_lines            = 2
  description          = "terraform-example live ai subtitle"
}
```

## Argument Reference

The following arguments are supported:

* `subtitle_name` - (Required, String) The name of the AI subtitle template. It supports only digits, letters and hyphens (-), and the hyphen cannot be the first character.
* `src_language` - (Required, String) The source language. Valid values: `en-US`, `zh-CN`, `ru-RU`.
* `position_normalized` - (Required, List) The subtitle position, a pair of normalized `[x, y]` coordinates measured from the lower-left corner of the screen. The actual pixel offset is the value multiplied by the preview height, e.g. `[0.32, 0.27]`. At least two elements are required.
* `word_per_line` - (Required, Int) The number of characters displayed per line. Valid range: `1` to `500`.
* `font_size_normalized` - (Required, Float) The font size. Valid range: `0` to `1`.
* `dst_language` - (Optional, String) The translation target language. Valid values: `en-US`, `zh-CN`, `es-ES`, `ru-RU`.
* `font_color` - (Optional, String) The font color, in RGBA format, e.g. `0xFFFFFF`.
* `font_name` - (Optional, String) The font. Valid values: `KaiTi`, `AlibabaPuHuiTi-Regular`, `AlibabaPuHuiTi-Bold`, `AlibabaPuHuiTi-Light`, `NotoSansHans-Regular`, `NotoSansHans-Bold`, `NotoSansHans-Light`.
* `bg_color` - (Optional, String) The subtitle background color, in RGBA format, e.g. `0xFF0000`.
* `show_source_lan` - (Optional, Bool) Whether to display the source language. Default value: `false`.
* `max_lines` - (Optional, Int) The number of lines to display.
* `border_width_normalized` - (Optional, Float) The font weight (border width). Valid range: `0` to `1`.
* `bg_width_normalized` - (Optional, Float) The subtitle background size. Valid range: `0` to `1`.
* `description` - (Optional, String) The description of the AI subtitle template.
* `width` - (Optional, String) The preview screen width, in pixels (px).
* `height` - (Optional, String) The preview screen height, in pixels (px).
* `copy_from` - (Optional, String, ForceNew) The name of an existing subtitle template to copy settings from when creating. Changing this forces a new resource.

  -> **NOTE:** `copy_from` is only evaluated during resource creation. The API does not return this parameter, so it cannot be verified after import.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. It is the `SubtitleId` of the AI subtitle template.
* `rules_refer` - (Computed, List) The IDs of subtitle rules that reference this template.

## Import

Live AI Subtitle can be imported using the id, e.g.

```shell
$ terraform import alicloud_live_ai_subtitle.example <subtitle_id>
```
