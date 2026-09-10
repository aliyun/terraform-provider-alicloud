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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_smartag_flow_log&exampleId=e8302ffa-3478-9ad3-e708-aa50b56cc6470d8b3308&activeTab=example&spm=docs.r.smartag_flow_log.0.e8302ffa34&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_smartag_flow_log&spm=docs.r.smartag_flow_log.example&intl_lang=EN_US)

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Flow Log.
* `delete` - (Defaults to 1 mins) Used when delete the Flow Log.
* `update` - (Defaults to 1 mins) Used when update the Flow Log.

## Import

Smartag Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_smartag_flow_log.example <id>
```