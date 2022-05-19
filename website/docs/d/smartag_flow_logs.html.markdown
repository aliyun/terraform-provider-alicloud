---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_smartag_flow_logs"
sidebar_current: "docs-alicloud-datasource-smartag-flow-logs"
description: |-
  Provides a list of Smartag Flow Logs to the user.
---

# alicloud\_smartag\_flow\_logs

This data source provides the Smartag Flow Logs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.168.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_smartag_flow_logs" "ids" {
  ids = ["example_id"]
}
output "smartag_flow_log_id_1" {
  value = data.alicloud_smartag_flow_logs.ids.logs.0.id
}

data "alicloud_smartag_flow_logs" "nameRegex" {
  name_regex = "^my-FlowLog"
}
output "smartag_flow_log_id_2" {
  value = data.alicloud_smartag_flow_logs.nameRegex.logs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the flow log.
* `ids` - (Optional, ForceNew, Computed)  A list of Flow Log IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Flow Log name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the flow log. Valid values:  `Active`: The flow log is enabled. `Inactive`: The flow log is disabled.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Flow Log names.
* `logs` - A list of Smartag Flow Logs. Each element contains the following attributes:
	* `active_aging` - The time interval at which log data of active connections is collected. Valid values: 60 to 6000. Default value: 300. Unit: second.
	* `description` - The description of the flow log.
	* `flow_log_id` - The ID of the flow log.
	* `flow_log_name` - The name of the flow log.
	* `id` - The ID of the Flow Log.
	* `inactive_aging` - The time interval at which log data of inactive connections is connected. Valid values: 10 to 600. Default value: 15. Unit: second.
	* `logstore_name` - The name of the Log Service Logstore.
	* `netflow_server_ip` - The IP address of the NetFlow collector where the flow log is stored.
	* `netflow_server_port` - The port of the NetFlow collector. Default value: 9995.
	* `netflow_version` - The NetFlow version. Default value: V9.
	* `output_type` - The location where the flow log is stored. Valid values:  sls: The flow log is stored in Log Service. netflow: The flow log is stored on a NetFlow collector. all: The flow log is stored both in Log Service and on a NetFlow collector.
	* `project_name` - The name of the Log Service project.
	* `resource_group_id` - The ID of the resource group.
	* `sls_region_id` - The ID of the region where Log Service is deployed.
	* `status` - The status of the flow log. Valid values:  `Active`: The flow log is enabled. `Inactive`: The flow log is disabled.
	* `total_sag_num` - The number of Smart Access gateway (SAG) instances with which the flow log is associated.