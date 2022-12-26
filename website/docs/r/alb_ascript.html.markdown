---
subcategory: "Alb"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_ascript"
sidebar_current: "docs-alicloud-resource-alb-ascript"
description: |-
  Provides a Alicloud Alb Ascript resource.
---

# alicloud_alb_ascript

Provides a Alb Ascript resource.

For information about Alb Ascript and how to use it, see [What is AScript](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alb_ascript" "default" {
  script_content = "time()"
  position       = "RequestHead"
  ascript_name   = "test"
  enabled        = true
  listener_id    = var.listenerId
}
```

## Argument Reference

The following arguments are supported:
* `listener_id` - (Required,ForceNew) Listener ID of script attribution
* `position` - (Required,ForceNew) Execution location of AScript.
* `ascript_name` - (Required) The name of AScript.
* `script_content` - (Required) The content of AScript.
* `enabled` - (Required) Whether scripts are enabled.
* `ext_attribute_enabled` - (Optional) Whether extension parameters are enabled.
* `ext_attributes` - (Optional) Extended attribute list. See the following `Block ExtAttributes`.

#### Block ExtAttributes

The ExtAttributes supports the following:
* `attribute_key` - (Optional) The key of the extended attribute.
* `attribute_value` - (Optional) The value of the extended attribute.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `load_balancer_id` - The ID of load balancer instance.
* `status` - The status of AScript.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ascript.
* `delete` - (Defaults to 5 mins) Used when delete the Ascript.
* `update` - (Defaults to 5 mins) Used when update the Ascript.

## Import

Alb AScript can be imported using the id, e.g.

```shell
$terraform import alicloud_alb_ascript.example <id>
```