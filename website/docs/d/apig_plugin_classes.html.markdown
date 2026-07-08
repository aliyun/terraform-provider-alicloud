---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugin_classes"
description: |-
  Provides a list of APIG Plugin Class owned by an Alibaba Cloud account.
---

# alicloud_apig_plugin_classes

This data source provides the APIG Plugin Class of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.285.0.

## Example Usage

```terraform
resource "alicloud_apig_plugin_class" "default" {
  wasm_url            = "https://example.com/plugin.wasm"
  description         = "A example plugin class"
  version_description = "Initial version"
  plugin_class_name   = "example-plugin-class"
  version             = "1.0.2"
  execute_priority    = "1"
  wasm_language       = "TinyGo"
  execute_stage       = "UNSPECIFIED_PHASE"
}

data "alicloud_apig_plugin_classes" "ids" {
  ids = [alicloud_apig_plugin_class.default.id]
}

output "apig_plugin_class_id_0" {
  value = data.alicloud_apig_plugin_classes.ids.classes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `type` - (Optional) The type of the plugin class used to filter results. Valid values: `Auth`, `FlowControl`, `FlowObservation`, `Security`, `TransportProtocol`, `Other`.
* `ids` - (Optional, List) A list of Plugin Class IDs.
* `name_regex` - (Optional) A regex string to filter results by plugin class name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Plugin Classes.
* `classes` - A list of Plugin Class. Each element contains the following attributes:
  * `id` - The ID of the Plugin Class. It is the same as `plugin_class_id`.
  * `plugin_class_id` - The ID of the plugin class.
  * `plugin_class_name` - The name of the plugin class.
  * `alias` - The alias of the plugin class.
  * `description` - The description of the plugin class, which introduces the main functions of the plugin.
  * `type` - The type of the plugin class. Valid values: `Auth` (authentication and authorization plugin), `FlowControl` (traffic control plugin), `FlowObservation` (traffic observation plugin), `Security` (security protection plugin), `TransportProtocol` (transport protocol plugin), `Other` (custom plugin).
  * `version` - The version of the plugin class.
  * `status` - The publish status of the plugin class. Valid values: `Success` (published successfully), `Failed` (failed to be published), `Publishing` (being published). Only a plugin class in the `Success` status can be installed.
  * `document` - The document of the plugin class, which describes the functions and usage of the plugin in detail. It is available when `enable_details` is set to `true`.
  * `wasm_language` - The programming language of the wasm plugin. It is available when `enable_details` is set to `true`.
