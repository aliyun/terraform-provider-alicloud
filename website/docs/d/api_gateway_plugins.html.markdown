---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_plugins"
sidebar_current: "docs-alicloud-datasource-api-gateway-plugins"
description: |-
  Provides a list of Api Gateway Plugins to the user.
---

# alicloud\_api\_gateway\_plugins

This data source provides the Api Gateway Plugins of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_plugins" "ids" {}
output "api_gateway_plugin_id_1" {
  value = data.alicloud_api_gateway_plugins.ids.plugins.0.id
}

data "alicloud_api_gateway_plugins" "nameRegex" {
  name_regex = "^my-Plugin"
}
output "api_gateway_plugin_id_2" {
  value = data.alicloud_api_gateway_plugins.nameRegex.plugins.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Plugin IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Plugin name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `plugin_name` - (Optional, ForceNew) The name of the plug-in that you want to create. It can contain uppercase English letters, lowercase English letters, Chinese characters, numbers, and underscores (_). It must be 4 to 50 characters in length and cannot start with an underscore (_).
* `plugin_type` - (Optional, ForceNew) The type of the plug-in. Valid values: `backendSignature`, `caching`, `cors`, `ipControl`, `jwtAuth`, `trafficControl`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Plugin names.
* `plugins` - A list of Api Gateway Plugins. Each element contains the following attributes:
	* `create_time` - The CreateTime of the resource.
	* `description` - The description of the plug-in, which cannot exceed 200 characters.
	* `id` - The ID of the Plugin.
	* `modified_time` - The ModifiedTime of the resource.
	* `plugin_data` - The definition statement of the plug-in. Plug-in definition statements in the JSON and YAML formats are supported.
	* `plugin_id` - The first ID of the resource.
	* `plugin_name` - The name of the plug-in that you want to create.
	* `plugin_type` - The type of the plug-in.
	* `tags` - The tag of the resource.