---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancers"
sidebar_current: "docs-alicloud-datasource-nlb-load-balancers"
description: |-
  Provides a list of Nlb Load Balancers to the user.
---

# alicloud\_nlb\_load\_balancers 

This data source provides the Nlb Load Balancers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.191.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_load_balancers" "ids" {
  ids = ["example_id"]
}
output "nlb_load_balancer_id_1" {
  value = data.alicloud_nlb_load_balancers.ids.balancers.0.id
}

data "alicloud_nlb_load_balancers" "nameRegex" {
  name_regex = "^my-LoadBalancer"
}
output "nlb_load_balancer_id_2" {
  value = data.alicloud_nlb_load_balancers.nameRegex.balancers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `address_ip_version` - (Optional, ForceNew) The IP version. Valid values: `ipv4`, `DualStack`.
* `address_type` - (Optional, ForceNew) The type of IPv4 address used by the NLB instance. Valid values: `Internet`, `Intranet`.
* `ipv6_address_type` - (Optional, ForceNew) The type of IPv6 address used by the NLB instance. Valid values: `Internet`, `Intranet`.
* `dns_name` - (Optional, ForceNew) The domain name of the NLB instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Load Balancer IDs.
* `load_balancer_business_status` - (Optional, ForceNew) The business status of the NLB instance. Valid values: `Abnormal`, `Normal`.
* `load_balancer_ids` - (Optional, ForceNew) The ID of the NLB instance. You can specify at most 20 IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Load Balancer name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The status of the NLB instance. Valid values: `Inactive`, `Active`, `Provisioning`, `Configuring`, `Deleting`, `Deleted`.
* `vpc_ids` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) where the NLB instance is deployed. You can specify at most 10 IDs.
* `zone_id` - (Optional, ForceNew) The name of the zone.
* `load_balancer_names` - (Optional, ForceNew) The name of the NLB instance. You can specify at most 10 names.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Load Balancer names.
* `balancers` - A list of Nlb Load Balancers. Each element contains the following attributes:
	* `id` - The ID of the NLB instance.
	* `address_ip_version` - The IP version.
	* `address_type` - The type of IPv4 address used by the NLB instance.
	* `create_time` - The time when the resource was created. The time is displayed in UTC in `yyyy-MM-ddTHH:mm:ssZ` format.
	* `cross_zone_enabled` - Indicates whether cross-zone load balancing is enabled for the NLB instance.
	* `dns_name` - The domain name of the NLB instance.
	* `load_balancer_business_status` - The business status of the NLB instance.
	* `load_balancer_name` - The name of the NLB instance.
	* `load_balancer_id` - The ID of the NLB instance.
	* `ipv6_address_type` - The type of IPv6 address used by the NLB instance.
	* `load_balancer_type` - The type of the SLB instance. Only Network is supported, which indicates NLB.
	* `resource_group_id` - The ID of the resource group.
	* `bandwidth_package_id` - The ID of the EIP bandwidth plan that is associated with the NLB instance if the NLB instance uses a public IP address.
	* `security_group_ids` - The security group to which the NLB instance belongs.
	* `status` - The status of the NLB instance.
	* `vpc_id` - The ID of the VPC where the NLB instance is deployed.
	* `tags` - The tag of the resource.
	* `operation_locks` - The configuration of the operation lock. This parameter takes effect if LoadBalancerBussinessStatus is Abnormal.
		* `lock_type` - The type of lock.
		* `lock_reason` - The reason why the NLB instance is locked.
	* `zone_mappings` - The zones and the vSwitches in the zones. An NLB instance can be deployed across 2 to 10 zones.
		* `allocation_id` - The ID of the elastic IP address (EIP).
		* `eni_id` - The ID of the elastic network interface (ENI) attached to the NLB instance.
		* `private_ipv4_address` - The private IPv4 address used by the NLB instance.
		* `public_ipv4_address` - The public IPv4 address used by the NLB instance.
		* `vswitch_id` - The ID of the vSwitch. By default, you can specify one vSwitch (subnet) in each zone of the NLB instance.
		* `zone_id` - The name of the zone. 
		* `ipv6_address` - The IPv6 address of the NLB instance.