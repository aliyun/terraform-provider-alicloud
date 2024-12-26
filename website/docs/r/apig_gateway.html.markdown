---
subcategory: "APIG"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_gateway"
description: |-
  Provides a Alicloud APIG Gateway resource.
---

# alicloud_apig_gateway

Provides a APIG Gateway resource.



For information about APIG Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/en/).

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

resource "alicloud_apig_gateway" "default" {
  network_access_config {
    type = "Intranet"
  }

  log_config {
    sls {
      enable = "false"
    }
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
  spec              = "apigw.small.x1"
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }

  zone_config {
    select_option = "Auto"
  }

  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = var.name
}
```

### Deleting `alicloud_apig_gateway` or removing it from your configuration

The `alicloud_apig_gateway` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `gateway_name` - (Optional) The name of the resource
* `log_config` - (Optional, List) Log Configuration See [`log_config`](#log_config) below.
* `network_access_config` - (Optional, List) Network Access Configuration See [`network_access_config`](#network_access_config) below.
* `payment_type` - (Required, ForceNew) The payment type of the resource
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `spec` - (Optional, ForceNew) Gateway instance specifications
* `tags` - (Optional, Map) The tag of the resource
* `vswitch` - (Optional, ForceNew, List) The virtual switch associated with the Gateway. See [`vswitch`](#vswitch) below.
* `vpc` - (Optional, ForceNew, List) The VPC associated with the Gateway. See [`vpc`](#vpc) below.
* `zone_config` - (Required, List) Availability Zone Configuration See [`zone_config`](#zone_config) below.

### `log_config`

The log_config supports the following:
* `sls` - (Optional, List) Sls See [`sls`](#log_config-sls) below.

### `log_config-sls`

The log_config-sls supports the following:
* `enable` - (Optional) Enable Log Service

### `network_access_config`

The network_access_config supports the following:
* `type` - (Optional) Network Access Type

### `vswitch`

The vswitch supports the following:
* `vswitch_id` - (Optional, ForceNew) The ID of the virtual switch.

### `vpc`

The vpc supports the following:
* `vpc_id` - (Required, ForceNew) The VPC network ID.

### `zone_config`

The zone_config supports the following:
* `select_option` - (Required) Availability Zone Options

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation timestamp. Unit: milliseconds.
* `status` - The status of the resource
* `vswitch` - The virtual switch associated with the Gateway.
  * `name` - The virtual switch name.
* `vpc` - The VPC associated with the Gateway.
  * `name` - The name of the VPC gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway.
* `update` - (Defaults to 5 mins) Used when update the Gateway.

## Import

APIG Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_gateway.example <id>
```