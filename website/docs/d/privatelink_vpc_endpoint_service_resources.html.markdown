---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_resources"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-service-resources"
description: |-
  Provides a list of Privatelink Vpc Endpoint Service Resources to the user.
---

# alicloud\_privatelink\_vpc\_endpoint\_service\_resources

This data source provides the Privatelink Vpc Endpoint Service Resources of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.110.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoint_service_resources" "example" {
  service_id = "epsrv-gw8ii1xxxx"
}

output "first_privatelink_vpc_endpoint_service_resource_id" {
  value = data.alicloud_privatelink_vpc_endpoint_service_resources.example.resources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `service_id` - (Required, ForceNew) The ID of Vpc Endpoint Service.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Vpc Endpoint Service Resource IDs.
* `resources` - A list of Privatelink Vpc Endpoint Service Resources. Each element contains the following attributes:
	* `id` - The ID of the Vpc Endpoint Service Resource.
	* `resource_id` - The ID of Resource.
	* `resource_type` - The type of Resource.
