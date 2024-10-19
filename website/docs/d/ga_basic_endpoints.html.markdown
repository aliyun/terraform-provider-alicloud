---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_endpoints"
sidebar_current: "docs-alicloud-datasource-ga-basic-endpoints"
description: |-
  Provides a list of Global Accelerator (GA) Basic Endpoints to the user.
---

# alicloud_ga_basic_endpoints

This data source provides the Global Accelerator (GA) Basic Endpoints of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_basic_endpoints" "ids" {
  ids               = ["example_id"]
  endpoint_group_id = "example_id"
}

output "ga_basic_endpoints_id_1" {
  value = data.alicloud_ga_basic_endpoints.ids.endpoints.0.id
}

data "alicloud_ga_basic_endpoints" "nameRegex" {
  name_regex        = "tf-example"
  endpoint_group_id = "example_id"
}

output "ga_basic_endpoints_id_2" {
  value = data.alicloud_ga_basic_endpoints.nameRegex.endpoints.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Global Accelerator Basic Endpoints IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Global Accelerator Basic Endpoints name.
* `endpoint_group_id` - (Required, ForceNew) The ID of the Basic Endpoint Group.
* `endpoint_id` - (Optional, ForceNew) The ID of the Basic Endpoint.
* `endpoint_type` - (Optional, ForceNew) The type of the Basic Endpoint. Valid values: `ENI`, `SLB`, `ECS` and `NLB`.
* `name` - (Optional, ForceNew) The name of the Basic Endpoint.
* `status` - (Optional, ForceNew) The status of the Global Accelerator Basic Endpoint. Valid Value: `init`, `active`, `updating`, `binding`, `unbinding`, `deleting`, `bound`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Global Accelerator Basic Endpoint names.
* `endpoints` - A list of Global Accelerator Basic Endpoints. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Basic Endpoint. It formats as `<endpoint_group_id>:<endpoint_id>`.
  * `endpoint_group_id` - The ID of the Basic Endpoint Group.
  * `endpoint_id` - The ID of the Basic Endpoint.
  * `accelerator_id` - The ID of the Global Accelerator Basic Accelerator instance.
  * `endpoint_type` - The type of the Basic Endpoint.
  * `endpoint_address` - The address of the Basic Endpoint.
  * `endpoint_sub_address_type` - The sub address type of the Basic Endpoint.
  * `endpoint_sub_address` - The sub address of the Basic Endpoint.
  * `endpoint_zone_id` - The zone id of the Basic Endpoint.
  * `basic_endpoint_name` - The name of the Basic Endpoint.
  * `status` - The status of the Basic Endpoint.
	