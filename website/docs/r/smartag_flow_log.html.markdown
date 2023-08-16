---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_smartag_flow_log"
sidebar_current: "docs-alicloud-resource-smartag-flow-log"
description: |-
  Provides a Alicloud Smartag Flow Log resource.
---

# alicloud_smartag_flow_log

Provides a Smartag Flow Log resource.

For information about Smartag Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createflowlog).

-> **NOTE:** Available since v1.168.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_smartag_flow_log" "example" {
  netflow_server_ip   = "192.168.0.2"
  netflow_server_port = 9995
  netflow_version     = "V9"
  output_type         = "netflow"
}
```

## Argument Reference

The following arguments are supported:

* `active_aging` - (Optional) The time interval at which log data of active connections is collected. Valid values: `60` to `6000`. Default value: `300`. Unit: second.
* `description` - (Optional) The description of the flow log.
* `flow_log_name` - (Optional) The name of the flow log.
* `inactive_aging` - (Optional) The time interval at which log data of inactive connections is connected. Valid values: `10` to `600`. Default value: `15`. Unit: second.
* `logstore_name` - (Optional) The Logstore in Log Service. If `output_type` is set to `sls` or `all`, this parameter is required.
* `netflow_server_ip` - (Optional) The IP address of the NetFlow collector where the flow log is stored. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `netflow_server_port` - (Optional) The port of the NetFlow collector. Default value: `9995`. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `netflow_version` - (Optional) The NetFlow version. Default value: `V9`. Valid values: `V10`, `V5`, `V9`. If `output_type` is set to `netflow` or `all`, this parameter is required.
* `output_type` - (Required) The location where the flow log is stored. Valid values:  
  - `sls`: The flow log is stored in Log Service. 
  - `netflow`: The flow log is stored on a NetFlow collector. 
  - `all`: The flow log is stored both in Log Service and on a NetFlow collector.
* `project_name` - (Optional) The project in Log Service. If `output_type` is set to `sls` or `all`, this parameter is required.
* `sls_region_id` - (Optional) The ID of the region where Log Service is deployed. If `output_type` is set to `sls` or `all`, this parameter is required.
* `status` - (Optional) The status of the flow log. Valid values:  `Active`: The flow log is enabled. `Inactive`: The flow log is disabled.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Flow Log.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Flow Log.
* `delete` - (Defaults to 1 mins) Used when delete the Flow Log.
* `update` - (Defaults to 1 mins) Used when update the Flow Log.

## Import

Smartag Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_smartag_flow_log.example <id>
```