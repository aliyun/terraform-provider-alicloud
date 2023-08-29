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

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_common_bandwidth_package" "cbwp" {
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "1000"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "${var.name}1"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vsj" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "192.168.10.0/24"
  vswitch_name = "${var.name}2"
}

resource "alicloud_vswitch" "vsk" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  cidr_block   = "192.168.20.0/24"
  vswitch_name = "${var.name}3"
}

resource "alicloud_security_group" "defaultLkkjal" {
  vpc_id = alicloud_vpc.vpc.id
}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "nlb1"
  resource_group_name = "${var.name}5"
}

resource "alicloud_resource_manager_resource_group" "rg2" {
  display_name        = "nlb2"
  resource_group_name = "${var.name}6"
}

resource "alicloud_vswitch" "vsg" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.2.id
  cidr_block   = "192.168.30.0/24"
  vswitch_name = "${var.name}7"
}


resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  zone_mappings {
    vswitch_id           = alicloud_vswitch.vsj.id
    zone_id              = "cn-hangzhou-j"
    private_ipv4_address = "192.168.10.4"
  }
  zone_mappings {
    vswitch_id           = alicloud_vswitch.vsk.id
    zone_id              = "cn-hangzhou-k"
    private_ipv4_address = "192.168.20.4"
  }
  address_type                   = "Intranet"
  address_ip_version             = "Ipv4"
  modification_protection_status = "ConsoleProtection"
  load_balancer_type             = "Network"
  vpc_id                         = alicloud_vpc.vpc.id
  modification_protection_reason = "test"
  resource_group_id              = alicloud_resource_manager_resource_group.rg1.id
  bandwidth_package_id           = alicloud_common_bandwidth_package.cbwp.id
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew, Computed) Protocol version. Value:
  - **Ipv4**:IPv4 type.
  - **DualStack**: Double Stack type.
* `address_type` - (Required) The network address type of IPv4 for network load balancing. Value:
  - **Internet**: public network. Load balancer has a public network IP address, and the DNS domain name is resolved to a public network IP address, so it can be accessed in a public network environment.
  - **Intranet**: private network. The server load balancer only has a private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the intranet environment of the VPC where the server load balancer is located.
* `bandwidth_package_id` - (Optional, Computed) The ID of the EIP bandwidth plan that is associated with the NLB instance if the NLB instance uses a public IP address.
* `cross_zone_enabled` - (Optional, Computed) Whether cross-zone is enabled for a network-based load balancing instance. Value:
  - **true**: on.
  - **false**: closed.
* `deletion_protection_enabled` - (Optional, Computed) Whether to enable deletion protection. Value:
  - **true**: enabled.
  - **false**: Closed.
* `deletion_protection_reason` - (Optional, Computed) Enter the reason why delete protection is turned on. It must be 2 to 128 English or Chinese characters in length. It must start with a letter or a Chinese character and can contain digits, half-width periods (.), underscores (_), and dashes (-).
  -> **NOTE:**  This parameter is valid and legal only when **DeletionProtectionEnabled** is **true.
* `ipv6_address_type` - (Optional, Computed) The IPv6 address type of network load balancing. Value:
  - **Internet**: Server Load Balancer has a public IP address, and the DNS domain name is resolved to a public IP address, so it can be accessed in a public network environment.
  - **Intranet**: SLB only has the private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the Intranet environment of the VPC where SLB is located.
* `load_balancer_name` - (Optional) The name of the network-based load balancing instance. 2 to 128 English or Chinese characters in length, which must start with a letter or Chinese, and can contain numbers, half-width periods (.), underscores (_), and dashes (-).
* `load_balancer_type` - (Optional, ForceNew, Computed) Load balancing type. Only value: **network**, which indicates network-based load balancing.
* `dns_name` - (Optional, Computed) DNS domain name.
* `load_balancer_business_status` - (Optional, Computed) The business status of network-based load balancing. Value:
  - **Abnormal**: Abnormal status.
  - **Normal**: Normal.
* `modification_protection_reason` - (Optional, Computed) Enter the reason for turning on modification protection. It must be 2 to 128 English or Chinese characters in length. It must start with a letter or a Chinese character and can contain digits, half-width periods (.), underscores (_), and dashes (-).
  -> **NOTE:**  This parameter is valid and valid only when **Status** is set to **ConsoleProtection**.
* `modification_protection_status` - (Optional, Computed) Network-based load balancing modifies the protection status. Value:
  - **NonProtection**: indicates that **ModificationProtectionReason** is not allowed * *. If The **ModificationProtectionReason** of the protection configuration is configured, its configuration information is cleared.
  - **ConsoleProtection**: allows the **ModificationProtectionReason** to be passed in the console to modify protection * *.
    -> **NOTE:**  When the value is **ConsoleProtection**, that is, after modification protection is enabled, users cannot modify the instance configuration through the load balancing console, but can modify the instance configuration by calling API.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_group_ids` - (Optional, Available since v1.210.0) The security group to which the network-based SLB instance belongs.
* `tags` - (Optional, Map) List of labels.
* `vpc_id` - (Required, ForceNew) The ID of the network-based SLB instance.
* `zone_mappings` - (Required) The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones. See [`zone_mappings`](#zone_mappings) below.

### `zone_mappings`

The zone_mappings supports the following:
* `allocation_id` - (Optional) The ID of the elastic IP address.
* `private_ipv4_address` - (Optional, Computed) The private IPv4 address of a network-based server load balancer instance.
* `vswitch_id` - (Required) The switch corresponding to the zone. Each zone uses one switch and one subnet by default.
* `zone_id` - (Required) The ID of the zone of the NLB instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Resource creation time, using Greenwich Mean Time, formating' yyyy-MM-ddTHH:mm:ssZ '.
* `status` - The status of the network-based load balancing instance. Value:, indicating that the instance listener will no longer forward traffic.
* `zone_mappings` - The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones.
  * `eni_id` - The ID of the elastic network interface (ENI).
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