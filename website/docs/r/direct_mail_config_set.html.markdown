---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_config_set"
description: |-
  Provides a Alicloud Direct Mail Config Set resource.
---

# alicloud_direct_mail_config_set

Provides a Direct Mail Config Set resource.

A configuration set groups sending-related configurations. It can be associated with a dedicated IP pool and with sender addresses, and supports address validation options. Up to 100 configuration sets can be created.

For information about Direct Mail Config Set and how to use it, see [What is Config Set](https://next.api.alibabacloud.com/document/Dm/2015-11-23/ConfigSetCreate).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_direct_mail_config_set" "default" {
  name = var.name
}
```

Associate with a dedicated IP pool and enable address validation:

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "default" {
  name = var.name
}

resource "alicloud_direct_mail_config_set" "with_pool" {
  name                      = var.name
  description               = "example configuration set"
  ip_pool_id                = alicloud_direct_mail_dedicated_ip_pool.default.id
  is_public_channel_backoff = true
  validation_option {
    enabled               = true
    forbidden_status_list = ["INVALID", "DO_NOT_MAIL"]
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the configuration set. The description can be up to 50 characters long.
* `ip_pool_id` - (Optional) The ID of the associated IP pool.
* `is_force` - (Optional) Whether to force deletion of the sender address associations.

-> **NOTE:** This parameter configures deletion behavior and is only evaluated when Terraform attempts to destroy the resource. Changes to this parameter during updates are stored but have no immediate effect.

* `is_public_channel_backoff` - (Optional) Whether to fall back to the public channel when the associated dedicated IP pool has problems. true: The system removes the dedicated IPs that are throttled or blacklisted by the target mailbox service provider during each delivery and preferentially uses the remaining available IPs in the pool; when all IPs in the pool are throttled or blacklisted, it automatically switches to the public channel to ensure delivery success, and restores dedicated IPs after the throttling or blacklisting is lifted. false: Only the associated dedicated IP pool is used for sending; throttled or blacklisted IPs continue to be used, which may cause delivery failures or delays, but ensures that the sending IP is always your dedicated IP.
* `name` - (Optional) The name of the configuration set. The name can be up to 50 characters long and must be unique.
* `validation_option` - (Optional, List) Address validation options. See [`validation_option`](#validation_option) below.

### `validation_option`

The validation_option supports the following:
* `enabled` - (Optional) Whether the address validation feature is enabled. true: enabled (the BASIC_ONLY mode of the ValidateEmail API). false: disabled.
* `forbidden_status_list` - (Optional, List) The list of address validation main statuses to block. Before sending, the system pre-validates the recipient address; messages whose validation main status matches this list are rejected. Valid values: INVALID, CATCH_ALL, DO_NOT_MAIL, UNKNOWN.
* `forbidden_sub_status_list` - (Optional, List) The list of address validation sub-statuses to block, providing finer-grained interception control than the main statuses. Takes effect as a union with ForbiddenStatusList (a hit in either list blocks the message).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `ip_pool_name` - The name of the associated IP pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Config Set.
* `delete` - (Defaults to 5 mins) Used when delete the Config Set.
* `update` - (Defaults to 5 mins) Used when update the Config Set.

## Import

Direct Mail Config Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_config_set.example <id>
```
