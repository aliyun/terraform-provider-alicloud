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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_gateway&exampleId=6a041ca8-aa0e-8c55-950d-af10c2df89b01a30ef4f&activeTab=example&spm=docs.r.apig_gateway.0.6a041ca8aa&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_gateway&spm=docs.r.apig_gateway.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `gateway_name` - (Optional) The name of the resource
* `gateway_type` - (Optional, ForceNew, Available since v1.260.1) Describes the gateway type, which is categorized into the following two types:
  - API: indicates an API gateway
  - AI: Indicates an AI gateway
* `log_config` - (Optional, List) Log Configuration See [`log_config`](#log_config) below.
* `network_access_config` - (Optional, List) Network Access Configuration See [`network_access_config`](#network_access_config) below.
* `payment_type` - (Required, ForceNew) The payment type of the resource
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `spec` - (Optional, ForceNew) Gateway instance specifications
* `tags` - (Optional, Map) The tag of the resource
* `vswitch` - (Optional, ForceNew, List) The virtual switch associated with the Gateway. See [`vswitch`](#vswitch) below.
* `vpc` - (Optional, ForceNew, List) The VPC associated with the Gateway. See [`vpc`](#vpc) below.
* `zone_config` - (Required, List) Availability Zone Configuration See [`zone_config`](#zone_config) below.
* `zones` - (Optional, ForceNew, List, Available since v1.260.1) The List of zones associated with the Gateway. See [`zones`](#zones) below.

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

### `zones`

The zones supports the following:
* `vswitch_id` - (Optional, ForceNew) The vswitch ID.
* `zone_id` - (Optional, ForceNew) The zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation timestamp. Unit: milliseconds.
* `status` - The status of the resource
* `vswitch` - The virtual switch associated with the Gateway.
  * `name` - The virtual switch name.
* `vpc` - The VPC associated with the Gateway.
  * `name` - The name of the VPC gateway.
* `zones` - The List of zones associated with the Gateway.
  * `name` - The zone name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 11 mins) Used when create the Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway.
* `update` - (Defaults to 5 mins) Used when update the Gateway.

## Import

APIG Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_gateway.example <id>
```