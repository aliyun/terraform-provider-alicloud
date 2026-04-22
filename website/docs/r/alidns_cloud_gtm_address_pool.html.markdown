---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_address_pool"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Address Pool resource.
---

# alicloud_alidns_cloud_gtm_address_pool

Provides a Alidns Cloud Gtm Address Pool resource.

CloudGtm Address Pool  .

For information about Alidns Cloud Gtm Address Pool and how to use it, see [What is Cloud Gtm Address Pool](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmAddressPool).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_alidns_cloud_gtm_address_pool" "default" {
  address_pool_name         = "pool-example-1"
  health_judgement          = "all_ok"
  address_pool_type         = "IPv4"
  enable_status             = "enable"
  address_lb_strategy       = "sequence"
  sequence_lb_strategy_mode = "preemptive"
  remark                    = "remark"
}
```

## Argument Reference

The following arguments are supported:
* `address_lb_strategy` - (Optional) Load balancing strategy among addresses in the address pool:
  - round_robin: Round robin. For any incoming DNS query, all addresses are returned, and their order is rotated with each request.
  - sequence: Sequential. For any incoming DNS query, addresses with lower sequence numbers are returned first (the sequence number indicates the priority of an address, where a lower number means higher priority). If an address with a lower sequence number is unavailable, the next available address with the next lowest sequence number is returned.
  - weight: Weighted. Each address can be assigned a different weight, allowing DNS responses to return addresses according to their weight ratios.
  - source_nearest: Source proximity. This intelligent resolution feature allows GTM to return different addresses based on the geographic location of the DNS query source, enabling users to access the nearest available endpoint.
* `address_pool_name` - (Required) The name of the address pool. Fuzzy search is supported for the entered address pool name.  
* `address_pool_type` - (Required, ForceNew) Address pool type:
  - IPv4: Indicates that the service address to be resolved is an IPv4 address.
  - IPv6: Indicates that the service address to be resolved is an IPv6 address.
  - domain: Indicates that the service address to be resolved is a domain name.
* `enable_status` - (Required) Enable status of the address pool:
  - enable: Enabled. The address pool participates in DNS resolution when its health check is normal.
  - disable: Disabled. The address pool does not participate in DNS resolution regardless of its health check status.
* `health_judgement` - (Required) Conditions for determining the health status of the address pool:  
  - any_ok: At least one address in the address pool is available.  
  - p30_ok: At least 30% of the addresses in the address pool are available.  
  - p50_ok: At least 50% of the addresses in the address pool are available.  
  - p70_ok: At least 70% of the addresses in the address pool are available.  
  - all_ok: All addresses in the address pool are available.  
* `remark` - (Optional) A remark for the address pool to help users distinguish its usage scenario.  
* `sequence_lb_strategy_mode` - (Optional) Service recovery mode for preceding resources when the load balancing strategy among addresses is set to sequential mode:  
  - preemptive: Preemptive mode. When a preceding resource recovers, the address with the smaller sequence number is prioritized.  
  - non_preemptive: Non-preemptive mode. When a preceding resource recovers, the current address continues to be used.  

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - Creation time of the address pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cloud Gtm Address Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Cloud Gtm Address Pool.
* `update` - (Defaults to 5 mins) Used when update the Cloud Gtm Address Pool.

## Import

Alidns Cloud Gtm Address Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_address_pool.example <address_pool_id>
```