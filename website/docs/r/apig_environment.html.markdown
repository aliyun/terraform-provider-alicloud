---
subcategory: "APIG"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_environment"
description: |-
  Provides a Alicloud APIG Environment resource.
---

# alicloud_apig_environment

Provides a APIG Environment resource.



For information about APIG Environment and how to use it, see [What is Environment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = format("%s2", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
    }
  }
}

resource "alicloud_apig_environment" "default" {
  description       = var.name
  environment_name  = var.name
  gateway_id        = alicloud_apig_gateway.defaultgateway.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description
* `environment_name` - (Required, ForceNew) The name of the resource
* `gateway_id` - (Required, ForceNew) Gateway id
* `resource_group_id` - (Optional, Computed) The ID of the resource group

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Environment.
* `update` - (Defaults to 5 mins) Used when update the Environment.

## Import

APIG Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_environment.example <id>
```