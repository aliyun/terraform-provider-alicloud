---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_gateways"
sidebar_current: "docs-alicloud-datasource-apig-gateways"
description: |-
  Provides a list of Apig Gateway owned by an Alibaba Cloud account.
---

# alicloud_apig_gateways

This data source provides Apig Gateway available to the user.[What is Gateway](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateGateway)

-> **NOTE:** Available since v1.284.0.

## Example Usage

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

data "alicloud_apig_gateways" "default" {
  ids          = [alicloud_apig_gateway.default.id]
  name_regex   = alicloud_apig_gateway.default.gateway_name
  gateway_name = var.name
}

output "alicloud_apig_gateway_example_id" {
  value = data.alicloud_apig_gateways.default.gateways.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Gateway IDs.
* `name_regex` - (Optional) A regex string to filter results by Gateway name.
* `gateway_id` - (Optional) Cloud-native API gateway ID.
* `gateway_name` - (Optional) The name of the gateway.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Gateway IDs.
* `names` - A list of name of Gateways.
* `gateways` - A list of Gateway Entries.

---

### `gateways`

The gateways supports the following attributes:

* `create_from` - The source from which the gateway was created.
* `create_time` - Creation timestamp.
* `environments` - The list of environments associated with the gateway.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `expire_time` - Timestamp indicating when the subscription expires.
* `gateway_edition` - Gateway instance edition:.
* `gateway_id` - Cloud-native API gateway ID.
* `gateway_name` - Query by exact match of the gateway name.
* `gateway_type` - The gateway type.
* `load_balancers` - The list of Gateway ingress addresses.
* `payment_type` - Payment type:.
* `resource_group_id` - The ID of the destination resource group.
* `security_group` - Security group of the gateway.
* `spec` - Gateway specification:.
* `status` - Gateway status:.
* `sub_domain_infos` - List of second-level domain names.
* `tags` - The tag of the resource.
* `target_version` - The target version of the gateway instance.
* `update_time` - The timestamp when the resource was last updated.
* `vswitch` - The vSwitch associated with the gateway.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `version` - The current running version of the gateway instance.
* `vpc` - The Virtual Private Cloud (VPC) associated with the gateway.
* `zones` - The list of zones associated with the gateway.
* `id` - The ID of the resource supplied above.

---

### `gateways.environments`

The environments supports the following attributes:

* `alias` - The alias of the environment.
* `environment_id` - The ID of the environment.
* `name` - The name of the environment.

---

### `gateways.load_balancers`

The load_balancers supports the following attributes:

* `address` - The address of the load balancer for the gateway.
* `address_ip_version` - IP version:.
* `address_type` - Load balancer address type:.
* `gateway_default` - Indicates whether this is the default ingress address of the gateway.
* `ipv4_addresses` - The list of IPv4 addresses.
* `ipv6_addresses` - The list of IPv6 addresses.
* `load_balancer_id` - The ID of the load balancer associated with the gateway.
* `mode` - Load balancing provisioning mode for the gateway:.
* `ports` - The list of listening ports.
* `status` - Load balancer status:.
* `type` - Load balancer type:.

---

### `gateways.load_balancers.ports`

The ports supports the following attributes:

* `port` - The port number of the load balancer listener.
* `protocol` - Protocol:.

---

### `gateways.security_group`

The security_group supports the following attributes:

* `name` - The name of the security group.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `security_group_id` - The ID of the security group.

---

### `gateways.sub_domain_infos`

The sub_domain_infos supports the following attributes:

* `domain_id` - The ID of the secondary domain name for the gateway.
* `name` - The name of the secondary domain associated with the gateway.
* `network_type` - Network type:.
* `protocol` - The protocol used by the secondary domain name.

---

### `gateways.vswitch`

The vswitch supports the following attributes:

* `name` - The name of the virtual switch associated with the gateway.
* `vswitch_id` - The ID of the virtual switch.

---

### `gateways.vpc`

The vpc supports the following attributes:

* `name` - The name of the VPC gateway.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `vpc_id` - The ID of the VPC network associated with the gateway.

---

### `gateways.zones`

The zones supports the following attributes:

* `name` - The name of the availability zone for the gateway.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `vswitch_id` - The ID of the virtual switch in the availability zone.
* `zone_id` - The ID of the availability zone for the gateway.
