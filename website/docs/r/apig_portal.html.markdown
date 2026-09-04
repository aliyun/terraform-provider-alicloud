---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_portal"
description: |-
  Provides a Alicloud APIG Portal resource.
---

# alicloud_apig_portal

Provides a APIG Portal resource.

HiMarket Developer Portal.

For information about APIG Portal and how to use it, see [What is Portal](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreatePortal).

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_apig_portal" "default" {
  name        = var.name
  description = "example portal"

  portal_setting_config {
    builtin_auth_enabled       = true
    auto_approve_developers    = false
    auto_approve_subscriptions = false
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Portal description
* `name` - (Optional) The name of the resource
* `portal_setting_config` - (Optional, Computed, Set) Portal configuration See [`portal_setting_config`](#portal_setting_config) below.

### `portal_setting_config`

The portal_setting_config supports the following:
* `auto_approve_developers` - (Optional) Whether to automatically approve Developer Registration
* `auto_approve_subscriptions` - (Optional) Automatically approve subscription requests
* `builtin_auth_enabled` - (Optional) Whether account and password authentication is enabled

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `portal_domain_config` - The domain information of the Portal.
  * `domain` - The domain name bound to the Portal.
  * `type` - The type of the domain name.
  * `protocol` - The request protocol of the domain name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Portal.
* `delete` - (Defaults to 5 mins) Used when delete the Portal.
* `update` - (Defaults to 5 mins) Used when update the Portal.

## Import

APIG Portal can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_portal.example <portal_id>
```
