---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugin_attachments"
sidebar_current: "docs-alicloud-datasource-apig-plugin-attachments"
description: |-
  Provides a list of Apig Plugin Attachment owned by an Alibaba Cloud account.
---

# alicloud_apig_plugin_attachments

This data source provides Apig Plugin Attachment available to the user. [What is Plugin Attachment](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreatePluginAttachment)

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_apig_plugin_attachments" "default" {
  attach_resource_type = "HttpApi"
}

output "alicloud_apig_plugin_attachment_example_id" {
  value = data.alicloud_apig_plugin_attachments.default.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `attach_resource_id` - (Optional) The ID of the attached resource.
* `attach_resource_type` - (Optional) The type of the resource to which the plug-in is attached, such as GatewayRoute, Gateway, GatewayDomain, HttpApi, or Operation.
* `environment_id` - (Optional) The environment ID.
* `plugin_info` - (Optional, List) The association information between the plug-in and the gateway instance. See [`plugin_info`](#plugin_info) below.
* `ids` - (Optional, Computed) A list of Plugin Attachment IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `plugin_info`

The plugin_info supports the following:
* `gateway_id` - (Optional) The gateway instance ID.
* `plugin_config` - (Optional) The Base64-encoded content of the original plug-in configuration.
* `plugin_id` - (Optional) The plug-in ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Plugin Attachment IDs.
* `attachments` - A list of Plugin Attachment Entries. Each element contains the following attributes:
  * `attach_resource_id` - The ID of the attached resource.
  * `attach_resource_ids` - The list of mount point IDs.
  * `attach_resource_names` - The list of names of resources to which the plug-in is attached.
  * `attach_resource_parent_ids` - The list of parent node IDs for the resources to which the plug-in is attached.
  * `attach_resource_type` - The type of the resource to which the plug-in is attached, such as GatewayRoute, Gateway, GatewayDomain, HttpApi, or Operation.
  * `enable` - Specifies whether to enable the feature.
  * `environment_id` - The environment ID.
  * `plugin_attachment_id` - The plug-in attachment ID.
  * `plugin_class_info` - The type information of the attached plug-in.
    * `direction` - The direction in which the plug-in acts on traffic: InBound, OutBound, or Both.
    * `execute_priority` - The execution priority of the plug-in.
    * `execute_stage` - The execution stage of the plug-in.
    * `name` - The name of the plug-in.
    * `type` - The plug-in type.
  * `plugin_info` - The association information between the plug-in and the gateway instance.
    * `gateway_id` - The gateway instance ID.
    * `plugin_config` - The Base64-encoded content of the original plug-in configuration.
    * `plugin_id` - The plug-in ID.
  * `id` - The ID of the resource supplied above.
