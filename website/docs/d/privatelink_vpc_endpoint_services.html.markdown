---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_services"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-services"
description: |-
  Provides a list of Privatelink Vpc Endpoint Services to the user.
---

# alicloud\_privatelink\_vpc\_endpoint\_services

This data source provides the Privatelink Vpc Endpoint Services of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.109.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoint_services" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_privatelink_vpc_endpoint_service_id" {
  value = data.alicloud_privatelink_vpc_endpoint_services.example.services.0.id
}
```

## Argument Reference

The following arguments are supported:

* `auto_accept_connection` - (Optional, ForceNew) Whether to automatically accept terminal node connections.
* `ids` - (Optional, ForceNew, Computed)  A list of Vpc Endpoint Service IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vpc Endpoint Service name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `service_business_status` - (Optional, ForceNew) The business status of the terminal node service. Valid Value: `Normal`, `FinancialLocked` and `SecurityLocked`.
* `status` - (Optional, ForceNew) The Status of Vpc Endpoint Service. Valid Value: `Active`, `Creating`, `Deleted`, `Deleting` and `Pending`.
* `vpc_endpoint_service_name` - (Optional, ForceNew) The name of Vpc Endpoint Service.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Vpc Endpoint Service names.
* `services` - A list of Privatelink Vpc Endpoint Services. Each element contains the following attributes:
	* `auto_accept_connection` - Whether to automatically accept terminal node connections..
	* `connect_bandwidth` - The connection bandwidth.
	* `id` - The ID of the Vpc Endpoint Service.
	* `service_business_status` - The business status of the terminal node service..
	* `service_description` - The description of the terminal node service.
	* `service_domain` - The domain of service.
	* `service_id` - The ID of the Vpc Endpoint Service.
	* `status` - The Status of Vpc Endpoint Service.
	* `vpc_endpoint_service_name` - The name of Vpc Endpoint Service.
