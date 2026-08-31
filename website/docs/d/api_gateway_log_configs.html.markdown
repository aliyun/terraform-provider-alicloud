---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_log_configs"
sidebar_current: "docs-alicloud-datasource-api-gateway-log-configs"
description: |-
  Provides a list of Api Gateway Log Configs to the user.
---

# alicloud\_api\_gateway\_log\_configs

This data source provides the Api Gateway Log Configs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_log_configs" "ids" {
  ids = ["example_id"]
}

output "api_gateway_log_config_id_1" {
  value = data.alicloud_api_gateway_log_configs.ids.configs.0.id
}

data "alicloud_api_gateway_log_configs" "logType" {
  log_type = "PROVIDER"
}

output "api_gateway_log_config_id_2" {
  value = data.alicloud_api_gateway_log_configs.logType.configs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Log Config IDs.
* `log_type` - (Optional, ForceNew) The type the of log. Valid values: `PROVIDER`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `configs` - A list of Api Gateway Log Configs. Each element contains the following attributes:
	* `id` - The ID of the Log Config.
	* `log_type` - The type the of log.
	* `sls_project` - The name of the Project.
	* `sls_log_store` - The name of the Log Store.
	* `region_id` - The region ID of the Log Config.