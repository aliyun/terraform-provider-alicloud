---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugin"
description: |-
  Provides a Alicloud APIG Plugin resource.
---

# alicloud_apig_plugin

Provides a APIG Plugin resource.



For information about APIG Plugin and how to use it, see [What is Plugin](https://next.api.alibabacloud.com/document/APIG/2024-03-27/InstallPlugin).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "plugin_vpc_pre" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "plugin-example-vpc"
}

resource "alicloud_vswitch" "plugin_vswitch_pre" {
  is_default   = false
  vpc_id       = alicloud_vpc.plugin_vpc_pre.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "plugin-example-vswitch"
}

resource "alicloud_apig_gateway" "plugin_gateway_pre" {
  network_access_config {
    type = "Internet"
  }
  vswitch {
    vswitch_id = alicloud_vswitch.plugin_vswitch_pre.id
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = alicloud_vpc.plugin_vpc_pre.id
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = "plugin-example-gateway"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}


resource "alicloud_apig_plugin" "default" {
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
  gateway_id      = alicloud_apig_gateway.plugin_gateway_pre.id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Required, ForceNew) The filter parameter for the gateway instance ID.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `plugin_class_id` - (Required, ForceNew) The plugin class ID.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `gateway_name` - The gateway name.
* `plugin_class_name` - The filter parameter for the plugin class name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Plugin.
* `delete` - (Defaults to 5 mins) Used when delete the Plugin.

## Import

APIG Plugin can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_plugin.example <plugin_id>
```