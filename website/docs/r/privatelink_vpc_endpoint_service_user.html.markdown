---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_user"
description: |-
  Provides a Alicloud Privatelink Vpc Endpoint Service User resource.
---

# alicloud_privatelink_vpc_endpoint_service_user

Provides a Privatelink Vpc Endpoint Service User resource.

Endpoint service user whitelist.

For information about Privatelink Vpc Endpoint Service User and how to use it, see [What is Vpc Endpoint Service User](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-addusertovpcendpointservice).

-> **NOTE:** Available since v1.110.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tfexampleuser"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}

resource "alicloud_ram_user" "example" {
  name         = var.name
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "alicloud_privatelink_vpc_endpoint_service_user" "example" {
  service_id = alicloud_privatelink_vpc_endpoint_service.example.id
  user_id    = alicloud_ram_user.example.id
}

resource "alicloud_privatelink_vpc_endpoint_service_user" "arn" {
  service_id = alicloud_privatelink_vpc_endpoint_service.example.id
  user_arn   = "acs:ram:*:11827252xxxxxxxx:*"
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.

-> **NOTE:** This parameter is only evaluated during resource creation and deletion. Modifying it in isolation will not trigger any action.

* `service_id` - (Required, ForceNew) The endpoint service ID.
* `user_arn` - (Optional, Computed, Available since v1.232.0) The whitelist in the format of ARN. At least one of `user_id` and `user_arn` must be specified.
* `user_id` - (Optional, ForceNew, Computed) The ID of the Alibaba Cloud account in the whitelist of the endpoint service. At least one of `user_id` and `user_arn` must be specified.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. For account ID whitelist entries, the value is formulated as `<service_id>:<user_id>`. For ARN whitelist entries, the value is formulated as `<service_id>:<user_id>:<user_arn>`, where `<user_id>` is empty when only `user_arn` is specified.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Service User.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Service User.

## Import

Privatelink Vpc Endpoint Service User can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service_user.example <service_id>:<user_id>
$ terraform import alicloud_privatelink_vpc_endpoint_service_user.arn <service_id>::<user_arn>
```