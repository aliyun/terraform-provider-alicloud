---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_smartag_flow_log"
sidebar_current: "docs-alicloud-resource-smartag-flow-log"
description: |-
  Provides a Alicloud Smartag Flow Log resource.
---

# alicloud\_smartag\_flow\_log

Provides a Smartag Flow Log resource.

For information about Smartag Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createflowlog).

-> **NOTE:** Available in v1.168.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_smartag_flow_log" "example" {
  flow_log_name       = "example_value"
  logstore_name       = "example_value"
  netflow_server_ip   = "example_value"
  netflow_server_port = 1
  project_name        = "example_value"
  sls_region_id       = "example_value"
  output_type         = "all"
}
```

## Argument Reference

The following arguments are supported:

* `active_aging` - (Optional, Computed) The time interval at which log data of active connections is collected. Valid values: `60` to `6000`. Default value: `300`. Unit: second.
* `description` - (Optional) The description of the flow log.
* `flow_log_name` - (Optional) The name of the flow log.
* `inactive_aging` - (Optional, Computed) The time interval at which log data of inactive connections is connected. Valid values: `10` to `600`. Default value: `15`. Unit: second.
* `logstore_name` - (Optional) The Logstore in Log Service. If `output_type` is set to `sls` or `all`, this parameter is required.
* `netflow_server_ip` - (Optional) The IP address of the NetFlow collector where the flow log is stored. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `netflow_server_port` - (Optional, Computed) The port of the NetFlow collector. Default value: `9995`. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `netflow_version` - (Optional, Computed) The NetFlow version. Default value: `V9`. Valid values: `V10`, `V5`, `V9`. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `output_type` - (Required) The location where the flow log is stored. Valid values:  
  - `sls`: The flow log is stored in Log Service. 
  - `netflow`: The flow log is stored on a NetFlow collector. 
  - `all`: The flow log is stored both in Log Service and on a NetFlow collector.
* `project_name` - (Optional) The project in Log Service. If `output_type` is set to `sls` or `all`, this parameter is required.
* `sls_region_id` - (Optional) The ID of the region where Log Service is deployed. If `output_type` is set to `sls` or `all`, this parameter is required.
* `status` - (Optional, Computed) The status of the flow log. Valid values:  `Active`: The flow log is enabled. `Inactive`: The flow log is disabled.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Flow Log.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Flow Log.
* `delete` - (Defaults to 1 mins) Used when delete the Flow Log.
* `update` - (Defaults to 1 mins) Used when update the Flow Log.

## Import

Smartag Flow Log can be imported using the id, e.g.

```
$ terraform import alicloud_smartag_flow_log.example <id>
```