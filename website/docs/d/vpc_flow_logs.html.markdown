---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_logs"
sidebar_current: "docs-alicloud-datasource-vpc-flow-logs"
description: |-
  Provides a list of Vpc Flow Logs to the user.
---

# alicloud\_vpc\_flow\_logs

This data source provides the Vpc Flow Logs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_flow_logs" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_vpc_flow_log_id" {
  value = data.alicloud_vpc_flow_logs.example.logs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The Description of flow log.
* `flow_log_name` - (Optional, ForceNew) The flow log name.
* `ids` - (Optional, ForceNew, Computed)  A list of Flow Log IDs.
* `log_store_name` - (Optional, ForceNew) The log store name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Flow Log name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_name` - (Optional, ForceNew) The project name.
* `resource_id` - (Optional, ForceNew) The resource id.
* `resource_type` - (Optional, ForceNew) The resource type. Valid values: `NetworkInterface`, `VPC`, `VSwitch`.
* `status` - (Optional, ForceNew) The status of  flow log. Valid values: `Active`, `Inactive`.
* `traffic_type` - (Optional, ForceNew) The traffic type. Valid values: `All`, `Allow`, `Drop`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Flow Log names.
* `logs` - A list of Vpc Flow Logs. Each element contains the following attributes:
	* `description` - The Description of flow log.
	* `flow_log_id` - The flow log ID.
	* `flow_log_name` - The flow log name.
	* `id` - The ID of the Flow Log.
	* `log_store_name` - The log store name.
	* `project_name` - The project name.
	* `resource_id` - The resource id.
	* `resource_type` - The resource type.
	* `status` - The status of flow log.
	* `traffic_type` - The traffic type.
