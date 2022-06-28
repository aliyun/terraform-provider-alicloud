---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_connections"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-connections"
description: |-
  Provides a list of Privatelink Vpc Endpoint Connections to the user.
---

# alicloud\_privatelink\_vpc\_endpoint\_connections

This data source provides the Privatelink Vpc Endpoint Connections of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.110.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoint_connections" "example" {
  service_id = "example_value"
  status     = "Connected"
}

output "first_privatelink_vpc_endpoint_connection_id" {
  value = data.alicloud_privatelink_vpc_endpoint_connections.example.connections.0.id
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Optional, ForceNew) The ID of the Vpc Endpoint.
* `endpoint_owner_id` - (Optional, ForceNew) The endpoint owner id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `service_id` - (Required, ForceNew) The ID of the Vpc Endpoint Service.
* `status` - (Optional, ForceNew) The status of Vpc Endpoint Connection. Valid Values: `Connected`, `Connecting`, `Deleted`, `Deleting`, `Disconnected`, `Disconnecting`, `Pending` and `ServiceDeleted`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Vpc Endpoint Connection IDs.
* `connections` - A list of Privatelink Vpc Endpoint Connections. Each element contains the following attributes:
	* `bandwidth` - The Bandwidth.
	* `endpoint_id` - The ID of the Vpc Endpoint.
	* `id` - The ID of the Vpc Endpoint Connection.
	* `status` - The status of Vpc Endpoint Connection.
