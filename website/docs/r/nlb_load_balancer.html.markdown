---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer"
description: |-
  Provides a Alicloud NLB Load Balancer resource.
---

# alicloud_nlb_load_balancer

Provides a NLB Load Balancer resource.



For information about NLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-nlb-2022-04-30-createloadbalancer).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_load_balancer&exampleId=ada6281a-e49c-7006-e3a0-ae8127c7c4adab8a3626&activeTab=example&spm=docs.r.nlb_load_balancer.0.ada6281ae4&intl_lang=EN_US" target="_blank">
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

DualStack Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_load_balancer&exampleId=88623514-c3ba-382c-6b47-7af1e5b6755195da00c9&activeTab=example&spm=docs.r.nlb_load_balancer.1.88623514c3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "tf-example"
}

variable "zone" {
  default = ["cn-beijing-i", "cn-beijing-k", "cn-beijing-l"]
}

resource "alicloud_vpc" "vpc" {
  vpc_name    = var.name
  cidr_block  = "10.2.0.0/16"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "vsw" {
  count                = 2
  enable_ipv6          = true
  ipv6_cidr_block_mask = "1${count.index}"
  vswitch_name         = "vsw-${count.index}-for-nlb"
  vpc_id               = alicloud_vpc.vpc.id
  cidr_block           = "10.2.1${count.index}.0/24"
  zone_id              = var.zone[count.index]
}

resource "alicloud_vpc_ipv6_gateway" "default" {
  ipv6_gateway_name = var.name
  vpc_id            = alicloud_vpc.vpc.id
}

resource "alicloud_nlb_load_balancer" "nlb" {
  depends_on         = [alicloud_vpc_ipv6_gateway.default]
  load_balancer_name = var.name
  load_balancer_type = "Network"
  address_type       = "Intranet"
  address_ip_version = "DualStack"
  ipv6_address_type  = "Internet"
  vpc_id             = alicloud_vpc.vpc.id
  cross_zone_enabled = false
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw[0].id
    zone_id    = var.zone[0]
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw[1].id
    zone_id    = var.zone[1]
  }
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew) The protocol version. Valid values:

  - **ipv4:** IPv4. This is the default value.
  - **DualStack:** dual stack.
* `address_type` - (Required) The type of IPv4 address used by the NLB instance. Valid values:
  - `Internet`: The NLB instance uses a public IP address. The domain name of the NLB instance is resolved to the public IP address. Therefore, the NLB instance can be accessed over the Internet.
  - `Intranet`: The NLB instance uses a private IP address. The domain name of the NLB instance is resolved to the private IP address. Therefore, the NLB instance can be accessed over the virtual private cloud (VPC) where the NLB instance is deployed.

-> **NOTE:**   To enable a public IPv6 address for an NLB instance, call the [EnableLoadBalancerIpv6Internet](https://www.alibabacloud.com/help/en/doc-detail/445878.html) operation.

* `bandwidth_package_id` - (Optional, ForceNew) The ID of the EIP bandwidth plan that is associated with the Internet-facing NLB instance.
* `cross_zone_enabled` - (Optional) Specifies whether to enable cross-zone load balancing for the NLB instance. Valid values:

  - `true`
  - `false`
* `deletion_protection_config` - (Optional, List) Specifies whether to enable deletion protection. Default value: `false`. See [`deletion_protection_config`](#deletion_protection_config) below.
* `ipv6_address_type` - (Optional) The type of IPv6 address used by the NLB instance. Valid values:
  - `Internet`: a public IP address. The domain name of the NLB instance is resolved to the public IP address. Therefore, the NLB instance can be accessed over the Internet.
  - `Intranet`: a private IP address. The domain name of the NLB instance is resolved to the private IP address. Therefore, the NLB instance can be accessed over the VPC where the NLB instance is deployed.
* `load_balancer_name` - (Optional) The name of the NLB instance.

  The value must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The value must start with a letter.
* `load_balancer_type` - (Optional, ForceNew, Computed) The type of the Server Load Balancer (SLB) instance. Set the value to `network`, which specifies NLB.
* `modification_protection_config` - (Optional, List) Specifies whether to enable the configuration read-only mode. Default value: `NonProtection`. See [`modification_protection_config`](#modification_protection_config) below.
* `resource_group_id` - (Optional, Computed) The ID of the new resource group.

  You can log on to the [Resource Management console](https://resourcemanager.console.aliyun.com/resource-groups) to view resource group IDs.
* `security_group_ids` - (Optional, Set) The security group to which the network-based SLB instance belongs.
* `tags` - (Optional, Map) List of labels.
* `vpc_id` - (Required, ForceNew) The ID of the VPC where the NLB instance is deployed.
* `zone_mappings` - (Required, Set) Available Area Configuration List. You must add at least two zones. You can add a maximum of 10 zones. See [`zone_mappings`](#zone_mappings) below.

### `deletion_protection_config`

The deletion_protection_config supports the following:
* `enabled` - (Optional) Specifies whether to enable deletion protection. Valid values:
  - `true`: yes
  - `false` (default): no
* `reason` - (Optional) The reason why deletion protection is enabled. The reason must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The reason must start with a letter.


-> **NOTE:**  This parameter takes effect only when `DeletionProtectionEnabled` is set to `true`.


### `modification_protection_config`

The modification_protection_config supports the following:
* `reason` - (Optional) The reason why the configuration read-only mode is enabled. The value must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The value must start with a letter.

-> **NOTE:**   This parameter takes effect only if the `status` parameter is set to `ConsoleProtection`.

* `status` - (Optional) Specifies whether to enable the configuration read-only mode. Valid values:
  - `NonProtection`: disables the configuration read-only mode. In this case, you cannot set the `ModificationProtectionReason` parameter. If you specify `ModificationProtectionReason`, the value is cleared.
  - `ConsoleProtection`: enables the configuration read-only mode. In this case, you can specify `ModificationProtectionReason`.

-> **NOTE:**  If you set this parameter to `ConsoleProtection`, you cannot use the NLB console to modify instance configurations. However, you can call API operations to modify instance configurations.


### `zone_mappings`

The zone_mappings supports the following:
* `allocation_id` - (Optional) The ID of the elastic IP address (EIP) that is associated with the Internet-facing NLB instance. You can specify one EIP for each zone. You must add at least two zones. You can add a maximum of 10 zones.
* `private_ipv4_address` - (Optional, Computed) The private IP address. You must add at least two zones. You can add a maximum of 10 zones.
* `status` - (Optional, Computed) Zone Status
* `vswitch_id` - (Required) The vSwitch in the zone. You can specify only one vSwitch (subnet) in each zone of an NLB instance. You must add at least two zones. You can add a maximum of 10 zones.
* `zone_id` - (Required) The ID of the zone of the NLB instance. You must add at least two zones. You can add a maximum of 10 zones.

  You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/443890.html) operation to query the most recent zone list.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Resource creation time, using Greenwich Mean Time, formating' yyyy-MM-ddTHH:mm:ssZ '.
* `deletion_protection_config` - Specifies whether to enable deletion protection. Default value: `false`.
  * `enabled_time` - Opening time of enable deletion protection.
* `dns_name` - The domain name of the NLB instance.
* `load_balancer_business_status` - The business status of the NLB instance.
* `modification_protection_config` - Specifies whether to enable the configuration read-only mode. Default value: `NonProtection`.
  * `enabled_time` - Opening time of the configuration read-only mode.
* `status` - The status of the NLB instance.
* `zone_mappings` - Available Area Configuration List. You must add at least two zones. You can add a maximum of 10 zones.
  * `eni_id` - The ID of the elastic network interface (ENI).
  * `ipv6_address` - The IPv6 address of the NLB instance.
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