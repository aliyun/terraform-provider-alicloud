---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_gateway"
description: |-
  Provides a Alicloud APIG Gateway resource.
---

# alicloud_apig_gateway

Provides a APIG Gateway resource.

Gateway instance  .

For information about APIG Gateway and how to use it, see [What is Gateway](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateGateway).

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
* `gateway_edition` - (Optional, ForceNew, Computed, Available since v1.284.0) Gateway instance edition. Valid values:
  - Professional: Standard instance.
  - Serverless: Serverless instance.
  - MultiTenantServerless: Multi-tenant Serverless instance.
* `gateway_name` - (Optional) Query by exact match of the gateway name.
* `gateway_type` - (Optional, ForceNew, Computed, Available since v1.260.1) The gateway type. Valid values:
  - API: API Gateway
  - AI: AI Gateway
* `log_config` - (Optional, Set) The log configuration for the gateway instance. See [`log_config`](#log_config) below. **Note: The parameter is immutable after resource creation.**
* `network_access_config` - (Optional, Set) The network access type of the gateway instance. See [`network_access_config`](#network_access_config) below. **Note: The parameter is immutable after resource creation.**

* `payment_type` - (Required, ForceNew) Payment type. Valid values:
  - PayAsYouGo: Pay-as-you-go.
  - Subscription: Subscription.
* `resource_group_id` - (Optional, Computed) The ID of the destination resource group.
* `spec` - (Optional, ForceNew) Gateway specification:  
  - apigw.small.x1: Small specification.  
* `tags` - (Optional, Map) The tag of the resource
* `vswitch` - (Optional, ForceNew, Set) The vSwitch associated with the gateway. See [`vswitch`](#vswitch) below.
* `vpc` - (Optional, ForceNew, Set) The Virtual Private Cloud (VPC) associated with the gateway. See [`vpc`](#vpc) below.
* `zone_config` - (Required, Set) The availability zone selection option for the gateway. See [`zone_config`](#zone_config) below. **Note: The parameter is immutable after resource creation.**

* `zones` - (Optional, ForceNew, Computed, List) The list of zones associated with the gateway. See [`zones`](#zones) below.

### `log_config`

The log_config supports the following:
* `sls` - (Optional, Set) The Simple Log Service configuration for the gateway. See [`sls`](#log_config-sls) below.

### `log_config-sls`

The log_config-sls supports the following:
* `enable` - (Optional) The Simple Log Service configuration for the gateway.

### `network_access_config`

The network_access_config supports the following:
* `type` - (Optional) The network access type of the gateway instance.

### `vswitch`

The vswitch supports the following:
* `vswitch_id` - (Optional, ForceNew) The ID of the virtual switch.

### `vpc`

The vpc supports the following:
* `vpc_id` - (Required, ForceNew) The ID of the VPC network associated with the gateway.

### `zone_config`

The zone_config supports the following:
* `select_option` - (Required) Zone selection option.

### `zones`

The zones supports the following:
* `vswitch_id` - (Optional, ForceNew) The ID of the virtual switch in the availability zone.
* `zone_id` - (Optional, ForceNew) The ID of the availability zone for the gateway.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_from` - The source from which the gateway was created.
* `create_time` - The creation timestamp. Unit: milliseconds.
* `environments` - The list of environments associated with the gateway.
    * `alias` - The alias of the environment.
    * `environment_id` - The ID of the environment.
    * `name` - The name of the environment.
* `expire_time` - Timestamp indicating when the subscription expires. Unit: milliseconds.
* `load_balancers` - The list of Gateway ingress addresses.
    * `address` - The address of the load balancer for the gateway.
    * `address_ip_version` - The IP version of the load balancer.
    * `address_type` - The load balancer address type.
    * `gateway_default` - Indicates whether this is the default ingress address of the gateway.
    * `ipv4_addresses` - The list of IPv4 addresses.
    * `ipv6_addresses` - The list of IPv6 addresses.
    * `load_balancer_id` - The ID of the load balancer associated with the gateway.
    * `mode` - The load balancing provisioning mode for the gateway.
    * `ports` - The list of listening ports.
        * `port` - The port number of the load balancer listener.
        * `protocol` - The protocol of the load balancer listener.
    * `status` - The status of the load balancer.
    * `type` - The type of the load balancer.
* `security_group` - The security group of the gateway.
    * `name` - The name of the security group.
    * `security_group_id` - The ID of the security group.
* `status` - The status of the gateway.
* `target_version` - The target version of the gateway instance.
* `update_time` - The timestamp when the gateway was last updated. Unit: milliseconds.
* `vswitch` - The vSwitch associated with the gateway.
    * `name` - The name of the virtual switch associated with the gateway.
* `version` - The current running version of the gateway instance.
* `vpc` - The Virtual Private Cloud (VPC) associated with the gateway.
    * `name` - The name of the VPC gateway.
* `zones` - The list of zones associated with the gateway.
    * `name` - The name of the availability zone for the gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 11 mins) Used when create the Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway.
* `update` - (Defaults to 5 mins) Used when update the Gateway.

## Import

APIG Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_gateway.example <gateway_id>
```