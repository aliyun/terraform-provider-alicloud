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

For information about Threat Detection Honeypot Preset and how to use it, see [What is Honeypot Preset](https://help.aliyun.com/document_detail/468960.html).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name                    = var.name
  available_probe_num          = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

resource "alicloud_threat_detection_honeypot_preset" "default" {
  honeypot_image_name = "shiro"
  meta {
    portrait_option = true
    burp            = "open"
  }
  node_id     = alicloud_threat_detection_honeypot_node.default.id
  preset_name = "apiapec_test"
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_image_name` - (Required,ForceNew) Honeypot mirror name
* `meta` - (Required,ForceNew) Honeypot template custom parameters. See the following `Block meta`.
* `node_id` - (Required,ForceNew) Unique id of management node
* `preset_name` - (Required) Honeypot template custom name

#### Block meta

The meta supports the following:

* `portrait_option` - (Optional) Social traceability.
* `burp` - (Required) Burp counter.
* `trojan_git` - (Optional) Git countered.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `honeypot_preset_id` - Unique ID of honeypot Template

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honeypot Preset.
* `delete` - (Defaults to 5 mins) Used when delete the Honeypot Preset.
* `update` - (Defaults to 5 mins) Used when update the Honeypot Preset.

## Import

Threat Detection Honeypot Preset can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honeypot_preset.example <id>
```