---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer"
description: |-
  Provides a Alicloud NLB Load Balancer resource.
---

# alicloud_nlb_load_balancer

Provides a NLB Load Balancer resource. 

For information about NLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createloadbalancer).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_nlb_load_balancer&exampleId=ada6281a-e49c-7006-e3a0-ae8127c7c4adab8a3626&activeTab=example&spm=docs.r.nlb_load_balancer.0.ada6281ae4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default1" {
  vswitch_name = var.name
  cidr_block   = "10.4.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew) Protocol version. Value:
  - **Ipv4**:IPv4 type.
  - **DualStack**: Double Stack type.
* `address_type` - (Required) The network address type of IPv4 for network load balancing. Value:
  - **Internet**: public network. Load balancer has a public network IP address, and the DNS domain name is resolved to a public network IP address, so it can be accessed in a public network environment.
  - **Intranet**: private network. The server load balancer only has a private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the intranet environment of the VPC where the server load balancer is located.
* `bandwidth_package_id` - (Optional, ForceNew) The ID of the shared bandwidth package associated with the public network instance.
* `cross_zone_enabled` - (Optional) Whether cross-zone is enabled for a network-based load balancing instance. Value:
  - **true**: on.
  - **false**: closed.
* `deletion_protection_config` - (Optional, Available since v1.217.1) Delete protection. See [`deletion_protection_config`](#deletion_protection_config) below.
* `ipv6_address_type` - (Optional) The IPv6 address type of network load balancing. Value:
  - **Internet**: Server Load Balancer has a public IP address, and the DNS domain name is resolved to a public IP address, so it can be accessed in a public network environment.
  - **Intranet**: SLB only has the private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the Intranet environment of the VPC where SLB is located.
* `load_balancer_name` - (Optional) The name of the network-based load balancing instance.  2 to 128 English or Chinese characters in length, which must start with a letter or Chinese, and can contain numbers, half-width periods (.), underscores (_), and dashes (-).
* `load_balancer_type` - (Optional, ForceNew, Computed) Load balancing type. Only value: **network**, which indicates network-based load balancing.
* `modification_protection_config` - (Optional, Available since v1.217.1) Modify protection. See [`modification_protection_config`](#modification_protection_config) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_group_ids` - (Optional, Available since v1.217.1) The security group to which the network-based SLB instance belongs.
* `tags` - (Optional, Map) List of labels.
* `vpc_id` - (Required, ForceNew) The ID of the network-based SLB instance.
* `zone_mappings` - (Required) The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones. See [`zone_mappings`](#zone_mappings) below.
* `deletion_protection_enabled` - (Optional, Available since v1.206.0) Specifies whether to enable deletion protection. Default value: `false`. Valid values:
  - `true`: Enable deletion protection.
  - `false`: Disable deletion protection. You cannot set the `deletion_protection_reason`. If the `deletion_protection_reason` is set, the value is cleared.
* `deletion_protection_reason` - (Optional, Available since v1.206.0) The reason why the deletion protection feature is enabled or disabled. The `deletion_protection_reason` takes effect only when `deletion_protection_enabled` is set to `true`.
* `modification_protection_status` - (Optional, Available since v1.206.0) Specifies whether to enable the configuration read-only mode. Default value: `NonProtection`. Valid values:
  - `NonProtection`: Does not enable the configuration read-only mode. You cannot set the `modification_protection_reason`. If the `modification_protection_reason` is set, the value is cleared.
  - `ConsoleProtection`: Enables the configuration read-only mode. You can set the `modification_protection_reason`.
* `modification_protection_reason` - (Optional, Available in 1.206.0+) The reason why the configuration read-only mode is enabled. The `modification_protection_reason` takes effect only when `modification_protection_status` is set to `ConsoleProtection`.

### `deletion_protection_config`

The deletion_protection_config supports the following:
* `enabled` - (Optional) Delete protection enable.
* `reason` - (Optional) Reason for opening.

### `modification_protection_config`

The modification_protection_config supports the following:
* `reason` - (Optional) Reason for opening.
* `status` - (Optional) ON.

### `zone_mappings`

The zone_mappings supports the following:
* `allocation_id` - (Optional) The ID of the elastic IP address.
* `private_ipv4_address` - (Optional) The private IPv4 address of a network-based server load balancer instance.
* `status` - (Optional) Zone Status.
* `vswitch_id` - (Required) The switch corresponding to the zone. Each zone uses one switch and one subnet by default.
* `zone_id` - (Required) The name of the zone. You can call the [DescribeZones](~~ 443890 ~~) operation to obtain the name of the zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Resource creation time, using Greenwich Mean Time, formating' yyyy-MM-ddTHH:mm:ssZ '.
* `deletion_protection_config` - Delete protection.
  * `enabled_time` - Opening time.
* `modification_protection_config` - Modify protection.
  * `enabled_time` - Opening time.
* `status` - The status of the resource.
* `load_balancer_business_status` - The business status of the NLB instance.
* `dns_name` - The domain name of the NLB instance.
* `zone_mappings` - The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones.
  * `eni_id` - The ID of ENI.
  * `ipv6_address` - The IPv6 address of a network-based server load balancer instance.
  * `public_ipv4_address` - Public IPv4 address of a network-based server load balancer instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

NLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer.example <id>
```