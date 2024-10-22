---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ip_endpoint_relations"
sidebar_current: "docs-alicloud-datasource-ga-basic-accelerate-ip-endpoint-relations"
description: |-
  Provides a list of Global Accelerator (GA) Basic Accelerate Ip Endpoint Relations to the user.
---

# alicloud_ga_basic_accelerate_ip_endpoint_relations

This data source provides the Global Accelerator (GA) Basic Accelerate Ip Endpoint Relations of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_basic_accelerate_ip_endpoint_relations" "ids" {
  ids            = ["example_id"]
  accelerator_id = "example_id"
}

output "ga_basic_accelerate_ip_endpoint_relations_id_1" {
  value = data.alicloud_ga_basic_accelerate_ip_endpoint_relations.ids.relations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Global Accelerator Basic Accelerate Ip Endpoint Relations IDs.
* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator Basic Accelerator instance.
* `accelerate_ip_id` - (Optional, ForceNew) The ID of the Basic Accelerate IP.
* `endpoint_id` - (Optional, ForceNew) The ID of the Basic Endpoint.
* `status` - (Optional, ForceNew) The status of the Global Accelerator Basic Accelerate Ip Endpoint Relation. Valid Value: `active`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `relations` - A list of Global Accelerator Basic Accelerate Ip Endpoint Relations. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Basic Accelerate Ip Endpoint Relation. It formats as `<accelerator_id>:<accelerate_ip_id>:<endpoint_id>`.
  * `accelerator_id` - The ID of the Global Accelerator Basic Accelerator instance.
  * `accelerate_ip_id` - The ID of the Basic Accelerate IP.
  * `endpoint_id` - The ID of the Basic Endpoint.
  * `endpoint_type` - The type of the Basic Endpoint.
  * `endpoint_address` - The address of the Basic Endpoint.
  * `endpoint_sub_address_type` - The sub address type of the Basic Endpoint.
  * `endpoint_sub_address` - The sub address of the Basic Endpoint.
  * `endpoint_zone_id` - The zone id of the Basic Endpoint.
  * `ip_address` - The address of the Basic Accelerate IP.
  * `basic_endpoint_name` - The name of the Basic Endpoint.
  * `status` - The status of the Basic Accelerate Ip Endpoint Relation.
	