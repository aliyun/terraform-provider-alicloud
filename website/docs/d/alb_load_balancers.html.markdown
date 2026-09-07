---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancers"
sidebar_current: "docs-alicloud-datasource-alb-load-balancers"
description: |-
  Provides a list of Alb Load Balancers to the user.
---

# alicloud_alb_load_balancers

This data source provides the Alb Load Balancers of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_alb_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = "true"
}

resource "alicloud_eip" "zone_a" {
  bandwidth            = "10"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_vswitch" "zone_a" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "192.168.0.0/18"
  zone_id              = data.alicloud_alb_zones.default.zones.0.id
  ipv6_cidr_block_mask = "6"
}

resource "alicloud_vswitch" "zone_b" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "192.168.128.0/18"
  zone_id              = data.alicloud_alb_zones.default.zones.1.id
  ipv6_cidr_block_mask = "8"
}

resource "alicloud_vpc_ipv6_gateway" "default" {
  ipv6_gateway_name = var.name
  vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth            = 1000
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_alb_load_balancer" "default" {
  load_balancer_edition       = "Basic"
  address_type                = "Internet"
  vpc_id                      = alicloud_vpc_ipv6_gateway.default.vpc_id
  address_allocated_mode      = "Fixed"
  address_ip_version          = "DualStack"
  ipv6_address_type           = "Internet"
  bandwidth_package_id        = alicloud_common_bandwidth_package.default.id
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  load_balancer_name          = var.name
  deletion_protection_enabled = false
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  zone_mappings {
    vswitch_id       = alicloud_vswitch.zone_a.id
    zone_id          = alicloud_vswitch.zone_a.zone_id
    eip_type         = "Common"
    allocation_id    = alicloud_eip.zone_a.id
    intranet_address = "192.168.10.1"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_b.id
    zone_id    = alicloud_vswitch.zone_b.zone_id
  }
  tags = {
    Created = "TF"
  }
}

data "alicloud_alb_load_balancers" "ids" {
  ids = [alicloud_alb_load_balancer.default.id]
}

output "alb_load_balancers_id_0" {
  value = data.alicloud_alb_load_balancers.ids.balancers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `address_type` - (Optional, ForceNew) The type of IP address that the ALB instance uses to provide services. Valid values: `Intranet`, `Internet`.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `ids` - (Optional, ForceNew, List)  A list of Load Balancer IDs.
* `load_balancer_business_status` - (Optional, ForceNew, Available since v1.142.0) Load Balancing of the Service Status. Valid Values: `Abnormal`and `Normal`.
* `load_balancer_ids` - (Optional, ForceNew) The load balancer ids.
* `load_balancer_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The load balancer status. Valid values: `Active`, `Configuring`, `CreateFailed`, `Inactive` and `Provisioning`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Load Balancer name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) where the ALB instance is deployed.
* `vpc_ids` - (Optional, ForceNew) The vpc ids.
* `zone_id` - (Optional, ForceNew) The zone ID of the resource.
* `load_balancer_bussiness_status` - (Deprecated since v1.142.0) Field `load_balancer_bussiness_status` has been deprecated from provider version 1.142.0. New field `load_balancer_business_status` instead.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Load Balancer names.
* `balancers` - A list of Alb Load Balancers. Each element contains the following attributes:
  * `access_log_config` - The Access Logging Configuration Structure.
    * `log_project` -  The log service that access logs are shipped to.
    * `log_store` - The logstore that access logs are shipped to.
  * `address_allocated_mode` - The method in which IP addresses are assigned.
  * `address_type` - The type of IP address that the ALB instance uses to provide services.
  * `bandwidth_package_id` - The ID of the EIP bandwidth plan which is associated with an ALB instance that uses a
    public IP address.
  * `create_time` - The creation time of the resource.
  * `deletion_protection_config` - Remove the Protection Configuration.
    * `enabled` - Remove the Protection Status.
    * `enabled_time` - Deletion Protection Turn-on Time Use Greenwich Mean Time, in the Format of Yyyy-MM-ddTHH: mm:SSZ.
  * `dns_name` - DNS Domain Name.
  * `id` - The ID of the Load Balancer.
  * `load_balancer_billing_config` - The configuration of the billing method.
    * `pay_type` - The billing method of the ALB instance.
  * `load_balancer_bussiness_status` - (Deprecated since v1.142.0) Load Balancing of the Service Status. **NOTE:** Field `load_balancer_bussiness_status` has been deprecated from provider version 1.142.0. New field `load_balancer_business_status` instead.
  * `load_balancer_business_status` - (Available since v1.142.0) Load Balancing of the Service Status.
  * `load_balancer_edition` - The edition of the ALB instance.
  * `load_balancer_id` - The first ID of the resource.
  * `load_balancer_name` - The name of the resource.
  * `load_balancer_operation_locks` - The Load Balancing Operations Lock Configuration.
    * `lock_reason` - The Locking of the Reasons. 
    * `lock_type` - The Locking of the Type.
  * `modification_protection_config` - Modify the Protection Configuration.
    * `status` - Specifies whether to enable the configuration read-only mode for the ALB instance.
    * `reason` - The reason for modification protection.
  * `resource_group_id` - The ID of the resource group.
  * `status` - The The load balancer status.
  * `tags` - The tag of the resource.
  * `vpc_id` - The ID of the virtual private cloud (VPC) where the ALB instance is deployed. 
  * `zone_mappings` - The zones and vSwitches.
    * `vswitch_id` - The ID of the vSwitch that corresponds to the zone.
    * `zone_id` - The ID of the zone to which the ALB instance belongs.
    * `status` - (Available since v1.250.0) The zone status.
    * `load_balancer_addresses` - (Available since v1.250.0) The address of the ALB instance.
      * `allocation_id` - The elastic IP address (EIP).
      * `eip_type` - The type of EIP.
      * `address` - IPv4 address.
      * `intranet_address` - The private IPv4 address.
      * `intranet_address_hc_status` - The health status of the private IPv4 address of the ALB instance.
      * `ipv6_address` - IPv6 address.
      * `ipv6_address_hc_status` - The health status of the private IPv6 address of the ALB instance.
      * `ipv4_local_addresses` - The IPv4 link-local addresses.
      * `ipv6_local_addresses` - The IPv6 link-local addresses.
