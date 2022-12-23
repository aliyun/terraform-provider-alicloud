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

For information about Threat Detection Honey Pot and how to use it, see [What is Honey Pot](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-createhoneypot).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_honey_pot" "default" {
  honeypot_image_name = "ruoyi"
  honeypot_image_id   = "sha256:007095d6de9c7a343e9fc1f74a7efc9c5de9d5454789d2fa505a1b3fc623730c"
  honeypot_name       = "huangtiong-test"
  node_id             = "a44e1ab3-6945-444c-889d-5bacee7056e8"
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_image_id` - (Required,ForceNew) The image ID of the honeypot.
* `honeypot_image_name` - (Required,ForceNew) Honeypot mirror name.
* `honeypot_name` - (Required) Honeypot custom name.
* `node_id` - (Required,ForceNew) The ID of the honeypot management node.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is the same as `honeypot_id`.
* `honeypot_id` - Honeypot ID.
* `preset_id` - The custom parameter ID of honeypot.
* `state` - Honeypot status.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honey Pot.
* `delete` - (Defaults to 5 mins) Used when delete the Honey Pot.
* `update` - (Defaults to 5 mins) Used when update the Honey Pot.

## Import

Threat Detection Honey Pot can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honey_pot.example <id>
```