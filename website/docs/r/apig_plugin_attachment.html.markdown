---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugin_attachment"
description: |-
  Provides a Alicloud APIG Plugin Attachment resource.
---

# alicloud_apig_plugin_attachment

Provides a APIG Plugin Attachment resource.

Plug-in attachment information.

For information about APIG Plugin Attachment and how to use it, see [What is Plugin Attachment](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreatePluginAttachment).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "10.0.0.0/24"
}

resource "alicloud_apig_gateway" "default" {
  gateway_name    = var.name
  gateway_edition = "Professional"
  gateway_type    = "API"
  payment_type    = "PayAsYouGo"
  spec            = "apigw.small.x1"
  vpc {
    vpc_id = alicloud_vpc.default.id
  }
  vswitch {
    vswitch_id = alicloud_vswitch.default.id
  }
  zone_config {
    select_option = "Auto"
  }
  network_access_config {
    type = "Internet"
  }
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_plugin" "default" {
  gateway_id      = alicloud_apig_gateway.default.id
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
}


resource "alicloud_apig_http_api" "default" {
  http_api_name = var.name
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = "/terraform-example"
}

resource "alicloud_apig_plugin_attachment" "default" {
  attach_resource_type = "HttpApi"
  attach_resource_ids  = [alicloud_apig_http_api.default.id]
  environment_id       = alicloud_apig_gateway.default.environments.0.environment_id
  enable               = true
  plugin_info {
    plugin_id     = alicloud_apig_plugin.default.id
    gateway_id    = alicloud_apig_gateway.default.id
    plugin_config = "eyJ0ZXN0IjoiaGVsbG8ifQ=="
  }
}
```

## Argument Reference

The following arguments are supported:
* `attach_resource_id` - (Optional, Computed, ForceNew) The ID of the attached resource.
* `attach_resource_ids` - (Optional, List) The list of mount point IDs.
* `attach_resource_type` - (Optional, ForceNew) The type of the resource to which the plug-in is attached, such as GatewayRoute, Gateway, GatewayDomain, HttpApi, or Operation.
* `enable` - (Optional) Specifies whether to enable the feature. Default value: false.
* `environment_id` - (Optional, ForceNew) The environment ID.
* `plugin_info` - (Optional, List) The association information between the plug-in and the gateway instance. See [`plugin_info`](#plugin_info) below.

### `plugin_info`

The plugin_info supports the following:
* `gateway_id` - (Optional, ForceNew) The gateway instance ID. **NOTE:** The ListPluginAttachments API does not return this field. After importing, the configuration must specify `gateway_id` again to keep the state consistent.
* `plugin_config` - (Optional) The Base64-encoded content of the original plug-in configuration.
* `plugin_id` - (Optional, ForceNew) The plug-in ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `attach_resource_names` - The list of names of resources to which the plug-in is attached.
* `attach_resource_parent_ids` - The list of parent node IDs for the resources to which the plug-in is attached.
* `plugin_class_info` - The type information of the attached plug-in.
  * `direction` - The direction in which the plug-in acts on traffic: InBound, OutBound, or Both.
  * `execute_priority` - The execution priority of the plug-in.
  * `execute_stage` - The execution stage of the plug-in.
  * `name` - The name of the plug-in.
  * `type` - The plug-in type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Plugin Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Plugin Attachment.
* `update` - (Defaults to 5 mins) Used when update the Plugin Attachment.

## Import

APIG Plugin Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_plugin_attachment.example <plugin_attachment_id>
```
