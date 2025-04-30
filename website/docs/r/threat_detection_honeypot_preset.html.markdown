---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_preset"
sidebar_current: "docs-alicloud-resource-threat_detection-honeypot-preset"
description: |-
  Provides a Alicloud Threat Detection Honeypot Preset resource.
---

# alicloud_threat_detection_honeypot_preset

Provides a Threat Detection Honeypot Preset resource.

For information about Threat Detection Honeypot Preset and how to use it, see [What is Honeypot Preset](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypotpreset).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_honeypot_preset&exampleId=44c0969e-9b83-c6cd-3f09-cda647d5c8939b86c67f&activeTab=example&spm=docs.r.threat_detection_honeypot_preset.0.44c0969e9b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
data "alicloud_threat_detection_honeypot_images" "default" {
  name_regex = "^ruoyi"
}
resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name                    = var.name
  available_probe_num          = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

resource "alicloud_threat_detection_honeypot_preset" "default" {
  preset_name         = var.name
  node_id             = alicloud_threat_detection_honeypot_node.default.id
  honeypot_image_name = data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_name
  meta {
    portrait_option = true
    burp            = "open"
    trojan_git      = "open"
  }
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_image_name` - (Required, ForceNew) Honeypot mirror name
* `meta` - (Required, ForceNew) Honeypot template custom parameters. See [`meta`](#meta) below.
* `node_id` - (Required, ForceNew) Unique id of management node
* `preset_name` - (Required) Honeypot template custom name

### `meta`

The meta supports the following:

* `portrait_option` - (Optional) Social traceability.
* `burp` - (Required) Burp counter.
* `trojan_git` - (Optional) Git countered.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `honeypot_preset_id` - Unique ID of honeypot Template

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honeypot Preset.
* `delete` - (Defaults to 5 mins) Used when delete the Honeypot Preset.
* `update` - (Defaults to 5 mins) Used when update the Honeypot Preset.

## Import

Threat Detection Honeypot Preset can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honeypot_preset.example <id>
```