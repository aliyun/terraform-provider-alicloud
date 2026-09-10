---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_plugins"
sidebar_current: "docs-alicloud-datasource-apig-plugins"
description: |-
  Provides a list of Apig Plugin owned by an Alibaba Cloud account.
---

# alicloud_apig_plugins

This data source provides Apig Plugin available to the user.[What is Plugin](https://next.api.alibabacloud.com/document/APIG/2024-03-27/InstallPlugin)

-> **NOTE:** Available since v1.285.0.

## Example Usage

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

data "alicloud_apig_plugins" "default" {
  ids             = ["${alicloud_apig_plugin.default.id}"]
  gateway_id      = alicloud_apig_gateway.plugin_gateway_pre.id
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
}

output "alicloud_apig_plugin_example_id" {
  value = data.alicloud_apig_plugins.default.plugins.0.id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (ForceNew, Optional) The filter parameter for the gateway instance ID.
* `plugin_class_id` - (ForceNew, Optional) The plugin class ID.
* `plugin_class_name` - (ForceNew, Optional) The filter parameter for the plugin class name.
* `ids` - (Optional, Computed) A list of Plugin IDs. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Plugin IDs.
* `plugins` - A list of Plugin Entries. Each element contains the following attributes:
    * `gateway_id` - The filter parameter for the gateway instance ID.
    * `gateway_name` - The gateway name.
    * `plugin_class_id` - The plugin class ID.
    * `plugin_class_name` - The filter parameter for the plugin class name.
    * `plugin_id` - The plugin ID.
    * `id` - The ID of the resource supplied above.
