---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugin_class"
description: |-
  Provides a Alicloud APIG Plugin Class resource.
---

# alicloud_apig_plugin_class

Provides a APIG Plugin Class resource.

plugin class info.

For information about APIG Plugin Class and how to use it, see [What is Plugin Class](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreatePluginClass).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_apig_plugin_class" "default" {
  wasm_url                      = "https://example.com/plugin.wasm"
  description                   = "A example plugin class for CloudSpec coverage"
  version_description           = "Initial version for exampleing"
  plugin_class_name             = "example-plugin-class-cspec-v3"
  version                       = "1.0.2"
  alias                         = "example-plugin-alias-v3"
  execute_priority              = "1"
  wasm_language                 = "TinyGo"
  supported_min_gateway_version = "1.0.0"
  execute_stage                 = "UNSPECIFIED_PHASE"
}
```

## Argument Reference

The following arguments are supported:
* `alias` - (Optional, ForceNew) The alias of the plugin class.
* `description` - (Required, ForceNew) The description of the plugin class, which introduces the main functions of the plugin.
* `execute_priority` - (Required, ForceNew, Int) The execution priority of the plugin. The larger the value, the higher the priority. The value must be greater than `0`.
* `execute_stage` - (Required, ForceNew) The execution stage of the plugin. Valid values:
  - `AUTHN`: The authentication stage.
  - `AUTHZ`: The authorization stage.
  - `STATS`: The statistics stage.
  - `UNSPECIFIED_PHASE`: The default stage.
* `plugin_class_name` - (Required, ForceNew) The name of the plugin class.
* `supported_min_gateway_version` - (Optional, ForceNew) The minimum gateway version supported by the plugin.
* `version` - (Required, ForceNew) The version of the plugin class. The value increments from `1.0.0`.
* `version_description` - (Required, ForceNew) The description of the current version.
* `wasm_language` - (Required, ForceNew) The programming language of the wasm plugin. Valid values: `TinyGo`, `Cpp`, `Rust`, `AssemblyScript`, `Zig`.
* `wasm_url` - (Required, ForceNew) The URL of the wasm plugin.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is the same as `plugin_class_id`.
* `document` - The document of the plugin class, which describes the functions and usage of the plugin in detail.
* `status` - The publish status of the plugin class. Only a plugin class in the `Success` status can be installed. Valid values:
  - `Success`: The plugin class is published successfully.
  - `Failed`: The plugin class fails to be published.
  - `Publishing`: The plugin class is being published.
* `type` - The type of the plugin class. Valid values:
  - `Auth`: The authentication and authorization plugin.
  - `FlowControl`: The traffic control plugin.
  - `FlowObservation`: The traffic observation plugin.
  - `Security`: The security protection plugin.
  - `TransportProtocol`: The transport protocol plugin.
  - `Other`: The custom plugin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Plugin Class.
* `delete` - (Defaults to 5 mins) Used when delete the Plugin Class.

## Import

APIG Plugin Class can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_plugin_class.example <plugin_class_id>
```