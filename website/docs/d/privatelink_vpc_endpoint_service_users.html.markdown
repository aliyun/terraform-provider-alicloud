---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_users"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-service-users"
description: |-
  Provides a list of Privatelink Vpc Endpoint Service Users to the user.
---

# alicloud\_privatelink\_vpc\_endpoint\_service\_users

This data source provides the Privatelink Vpc Endpoint Service Users of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.110.0+.

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
* `service_id` - (Required, ForceNew) The Id of Vpc Endpoint Service.
* `user_id` - (Optional, ForceNew) The Id of Ram User.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Vpc Endpoint Service User IDs.
* `users` - A list of Privatelink Vpc Endpoint Service Users. Each element contains the following attributes:
	* `id` - The ID of the Vpc Endpoint Service User.
	* `user_id` - The Id of Ram User.