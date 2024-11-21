---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honey_pot"
sidebar_current: "docs-alicloud-resource-threat-detection-honey-pot"
description: |-
  Provides a Alicloud Threat Detection Honey Pot resource.
---

# alicloud_threat_detection_honey_pot

Provides a Threat Detection Honey Pot resource.

For information about Threat Detection Honey Pot and how to use it, see [What is Honey Pot](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypot).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_honey_pot&exampleId=77751521-44b2-13cd-1871-456d39677ef429826421&activeTab=example&spm=docs.r.threat_detection_honey_pot.0.7775152144&intl_lang=EN_US" target="_blank">
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

resource "alicloud_threat_detection_honey_pot" "default" {
  honeypot_image_name = data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_name
  honeypot_image_id   = data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_id
  honeypot_name       = var.name
  node_id             = alicloud_threat_detection_honeypot_node.default.id
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_image_id` - (Required, ForceNew) The image ID of the honeypot.
* `honeypot_image_name` - (Required, ForceNew) Honeypot mirror name.
* `honeypot_name` - (Required) Honeypot custom name.
* `node_id` - (Required, ForceNew) The ID of the honeypot management node.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is the same as `honeypot_id`.
* `honeypot_id` - Honeypot ID.
* `preset_id` - The custom parameter ID of honeypot.
* `state` - Honeypot status.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honey Pot.
* `delete` - (Defaults to 5 mins) Used when delete the Honey Pot.
* `update` - (Defaults to 5 mins) Used when update the Honey Pot.

## Import

Threat Detection Honey Pot can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honey_pot.example <id>
```