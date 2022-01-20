---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_data_flows"
sidebar_current: "docs-alicloud-datasource-nas-data-flows"
description: |-
  Provides a list of Nas Data Flows to the user.
---

# alicloud\_nas\_data\_flows

This data source provides the Nas Data Flows of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_data_flows" "ids" {
  file_system_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "nas_data_flow_id_1" {
  value = data.alicloud_nas_data_flows.ids.flows.0.id
}

data "alicloud_nas_data_flows" "status" {
  file_system_id = "example_value"
  status         = "Running"
}
output "nas_data_flow_id_2" {
  value = data.alicloud_nas_data_flows.status.flows.0.id
}

```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `ids` - (Optional, ForceNew, Computed)  A list of Data Flow IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the Data flow. Including: `Starting`, `Running`, `Updating`, `Deleting`, `Stopping`, `Stopped`, `Misconfigured`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `flows` - A list of Nas Data Flows. Each element contains the following attributes:
 * `create_time` - The time when Fileset was created. Executing the ISO8601 standard means that the return format is: 'yyyy-MM-ddTHH:mm:ssZ'.
 * `data_flow_id` - The ID of the Data Flow.
 * `description` - The Description of data flow.
 * `error_message` - Error message.
 * `file_system_id` - The ID of the file system.
 * `file_system_path` - The path of Fileset in the CPFS file system.
 * `fset_description` - Description of automatic update.
 * `fset_id` - The ID of the Fileset.
 * `id` - The resource ID of the data flow. The value formats as `<file_system_id>:<data_flow_id>`.
 * `source_security_type` - The security protection type of the source storage.
 * `source_storage` - The access path of the source store. Format: `<storage type>://<path>`.
 * `status` - The status of the Data flow.
 * `throughput` - The maximum transmission bandwidth of data flow, unit: `MB/s`.