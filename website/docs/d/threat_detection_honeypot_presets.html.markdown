---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_presets"
sidebar_current: "docs-alicloud-datasource-threat_detection-honeypot-presets"
description: |-
  Provides a list of Threat Detection Honeypot Preset owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_honeypot_presets

This data source provides Threat Detection Honeypot Preset available to the user.

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_threat_detection_honeypot_presets" "default" {
  ids                 = ["${alicloud_threat_detection_honeypot_preset.default.id}"]
  honeypot_image_name = "shiro"
  node_id             = "example_value"
  preset_name         = "apiapec_test"
}

output "alicloud_threat_detection_honeypot_preset_example_id" {
  value = data.alicloud_threat_detection_honeypot_presets.default.presets.0.id
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_image_name` - (ForceNew,Optional) Honeypot mirror name
* `node_id` - (ForceNew,Optional) Unique id of management node
* `preset_name` - (ForceNew,Optional) Honeypot template custom name
* `ids` - (Optional, ForceNew, Computed) A list of Honeypot Preset IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Honeypot Preset IDs.
* `presets` - A list of Honeypot Preset Entries. Each element contains the following attributes:
  * `honeypot_image_name` - Honeypot mirror name.
  * `honeypot_preset_id` - Unique ID of honeypot Template.
  * `meta` - Honeypot template custom parameters.
    * `portrait_option` - Social traceability.
    * `burp` - Burp counter.
    * `trojan_git` - Git countered.
  * `node_id` - Unique id of management node.
  * `preset_name` - Honeypot template custom name.
  * `id` - The id of the Honeypot template.
