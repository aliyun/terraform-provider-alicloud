---
subcategory: "Database Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_database_gateway_gateways"
sidebar_current: "docs-alicloud-datasource-database-gateway-gateways"
description: |-
  Provides a list of Database Gateway Gateways to the user.
---

# alicloud\_database\_gateway\_gateways

This data source provides the Database Gateway Gateways of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_database_gateway_gateways" "ids" {
  ids = ["example_id"]
}
output "database_gateway_gateway_id_1" {
  value = data.alicloud_database_gateway_gateways.ids.gateways.0.id
}

data "alicloud_database_gateway_gateways" "nameRegex" {
  name_regex = "^my-Gateway"
}
output "database_gateway_gateway_id_2" {
  value = data.alicloud_database_gateway_gateways.nameRegex.gateways.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `search_key` - (Optional, ForceNew) The search key.
* `status` - (Optional, ForceNew) The status of gateway. Valid values: `EXCEPTION`, `NEW`, `RUNNING`, `STOPPED`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Gateway names.
* `gateways` - A list of Database Gateway Gateways. Each element contains the following attributes:
    * `create_time` - The creation time of Gateway.
    * `gateway_desc` - The description of Gateway.
    * `gateway_instance` - Gateway instance List.
        * `current_version` - The version of Gateway instance.
        * `local_ip` -  The Local IP ADDRESS of Gateway instance.
        * `output_ip` - The host of Gateway instance.
        * `connect_endpoint_type` - The connection type of Gateway instance.
        * `current_daemon_version` - The process of version number of Gateway instance.
        * `end_point` - The endpoint address of Gateway instance.
        * `gateway_instance_id` - The id of Gateway instance.
        * `gateway_instance_status` - The status of Gateway instance. Valid values: `EXCEPTION`, `NEW`, `RUNNING`, `STOPPED`.
        * `last_update_time` - The last Updated time stamp of Gateway instance.
        * `message` - The prompt information of Gateway instance.
    * `gateway_name` - The name of the Gateway.
    * `hosts` - A host of information.
    * `id` - The ID of Gateway.
    * `modified_time` - The Modify time of Gateway.
    * `parent_id` - The parent node Id of Gateway.
    * `status` - The status of gateway. Valid values: `EXCEPTION`, `NEW`, `RUNNING`, `STOPPED`.
    * `user_id` - The user's id.
