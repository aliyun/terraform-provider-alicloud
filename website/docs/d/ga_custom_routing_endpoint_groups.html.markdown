---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_groups"
sidebar_current: "docs-alicloud-datasource-ga-custom-routing-endpoint-groups"
description: |-
  Provides a list of Global Accelerator (GA) Custom Routing Endpoint Groups to the user.
---

# alicloud_ga_custom_routing_endpoint_groups

This data source provides the Global Accelerator (GA) Custom Routing Endpoint Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_custom_routing_endpoint_groups" "ids" {
  ids            = ["example_id"]
  accelerator_id = "your_accelerator_id"
}

output "ga_custom_routing_endpoint_groups_id_1" {
  value = data.alicloud_ga_custom_routing_endpoint_groups.ids.groups.0.id
}

data "alicloud_ga_custom_routing_endpoint_groups" "nameRegex" {
  name_regex     = "tf-example"
  accelerator_id = "your_accelerator_id"
}

output "ga_custom_routing_endpoint_groups_id_2" {
  value = data.alicloud_ga_custom_routing_endpoint_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Custom Routing Endpoint Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Custom Routing Endpoint Group name.
* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Optional, ForceNew) The ID of the custom routing listener.
* `endpoint_group_id` - (Optional, ForceNew) The ID of the endpoint group.
* `status` - (Optional, ForceNew) The status of the endpoint group. Valid Values: `init`, `active`, `updating`, `deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Custom Routing Endpoint Group names.
* `groups` - A list of Custom Routing Endpoint Groups. Each element contains the following attributes:
  * `id` - The id of the Custom Routing Endpoint Group.
  * `endpoint_group_id` - The ID of the Custom Routing Endpoint Group.
  * `accelerator_id` - The ID of the GA instance.
  * `listener_id` - The ID of the custom routing listener.
  * `endpoint_group_region` - The ID of the region where the endpoint group is created.
  * `endpoint_group_ip_list` - The list of endpoint group IP addresses.
  * `endpoint_group_unconfirmed_ip_list` - The endpoint group IP addresses to be confirmed after the GA instance is upgraded.
  * `custom_routing_endpoint_group_name` - The name of the endpoint group.
  * `description` - The description of the endpoint group.
  * `status` - The status of the endpoint group.
  