---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoints"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoints"
description: |-
  Provides a list of Privatelink Vpc Endpoints to the user.
---

# alicloud\_privatelink\_vpc\_endpoints

This data source provides the Privatelink Vpc Endpoints of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.109.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoints" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_privatelink_vpc_endpoint_id" {
  value = data.alicloud_privatelink_vpc_endpoints.example.endpoints.0.id
}
```

## Argument Reference

The following arguments are supported:

* `connection_status` - (Optional, ForceNew) The status of Connection.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Vpc Endpoint IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vpc Endpoint name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `service_name` - (Optional, ForceNew) The name of the terminal node service associated with the terminal node.
* `status` - (Optional, ForceNew) The status of Vpc Endpoint.
* `vpc_endpoint_name` - (Optional, ForceNew) The name of Vpc Endpoint.
* `vpc_id` - (Optional, ForceNew) The private network to which the terminal node belongs..

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Vpc Endpoint names.
* `endpoints` - A list of Privatelink Vpc Endpoints. Each element contains the following attributes:
	* `bandwidth` - The Bandwidth.
	* `connection_status` - The status of Connection.
	* `endpoint_business_status` - The status of Endpoint Business.
	* `endpoint_description` - The description of Vpc Endpoint.
	* `endpoint_domain` - The Endpoint Domain.
	* `endpoint_id` - The ID of the Vpc Endpoint.
	* `id` - The ID of the Vpc Endpoint.
	* `security_group_ids` - The security group associated with the terminal node network card.
	* `service_id` - The terminal node service associated with the terminal node.
	* `service_name` - The name of the terminal node service associated with the terminal node.
	* `status` - The status of Vpc Endpoint.
	* `vpc_endpoint_name` - The name of Vpc Endpoint.
	* `vpc_id` - The private network to which the terminal node belongs.
