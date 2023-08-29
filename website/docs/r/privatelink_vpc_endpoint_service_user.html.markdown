---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_user"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-service-user"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service User resource.
---

# alicloud_privatelink_vpc_endpoint_service_user

Provides a Private Link Vpc Endpoint Service User resource.

For information about Private Link Vpc Endpoint Service User and how to use it, see [What is Vpc Endpoint Service User](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-addusertovpcendpointservice).

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
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `service_id` - (Required, ForceNew) The Id of Vpc Endpoint Service.
* `user_id` - (Required, ForceNew) The Id of Ram User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Service User. The value is formatted `<service_id>:<user_id>`.

## Import

Private Link Vpc Endpoint Service User can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service_user.example <service_id>:<user_id>
```
