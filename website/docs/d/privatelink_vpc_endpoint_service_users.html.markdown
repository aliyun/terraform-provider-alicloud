---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_users"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-service-users"
description: |-
  Provides a list of Privatelink Vpc Endpoint Service Users to the user.
---

# alicloud_privatelink_vpc_endpoint_service_users

This data source provides the Privatelink Vpc Endpoint Service Users of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.110.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoint_service_users" "example" {
  service_id = "epsrv-gw81c6vxxxxxx"
}

output "first_privatelink_vpc_endpoint_service_user_id" {
  value = data.alicloud_privatelink_vpc_endpoint_service_users.example.users.0.id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `service_id` - (Required) The Id of Vpc Endpoint Service.
* `user_id` - (Optional) The Id of Ram User. When it is specified, only account ID whitelist entries are queried.
* `user_list_type` - (Optional) The type of the user whitelist to query. Valid values: `Users`, `UserARNs`. `Users` returns account ID whitelist entries; `UserARNs` returns all whitelist entries in ARN form (an account ID entry appears as `acs:ram:*:<account_id>:*`). If it is not set and `user_id` is not specified, both kinds are returned, and each whitelist entry appears exactly once.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Vpc Endpoint Service User IDs. For account ID whitelist entries, the element is formulated as `<service_id>:<user_id>`. For ARN whitelist entries, the element is formulated as `<service_id>::<user_arn>`.
* `users` - A list of Privatelink Vpc Endpoint Service Users. Each element contains the following attributes:
  * `id` - The ID of the Vpc Endpoint Service User.
  * `user_id` - The Id of Ram User.
  * `user_arn` - The whitelist in the format of ARN.
