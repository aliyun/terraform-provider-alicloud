---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_instance_config"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Instance Config resource.
---

# alicloud_alidns_cloud_gtm_instance_config

Provides a Alidns Cloud Gtm Instance Config resource.

CloudGtm Instance Configuration  .

For information about Alidns Cloud Gtm Instance Config and how to use it, see [What is Cloud Gtm Instance Config](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmInstanceConfig).

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

resource "alicloud_alidns_domain" "default" {
  domain_name = "${var.name}.abc"
}

resource "alicloud_alidns_cloud_gtm_instance_config" "default" {
  address_pool_lb_strategy  = "sequence"
  schedule_rr_type          = "A"
  schedule_zone_name        = alicloud_alidns_domain.default.domain_name
  enable_status             = "disable"
  schedule_host_name        = "example"
  schedule_zone_mode        = "custom"
  ttl                       = "600"
  sequence_lb_strategy_mode = "preemptive"
  remark                    = "remark"
}
```

## Argument Reference

The following arguments are supported:
* `address_pool_lb_strategy` - (Optional) The load balancing strategy among address pools:
  - round_robin: Round-robin. For any incoming DNS resolution request, all address pools are returned, and the order of these address pools is rotated with each request.
  - sequence: Sequential. For any incoming DNS resolution request, the address pool with the smallest sequence number is returned (the sequence number indicates the priority of the address pool, where a smaller number means higher priority). If the address pool with the smallest sequence number is unavailable, the next address pool with the next smallest sequence number is returned.
  - weight: Weighted. You can assign different weights to each address pool so that DNS queries return address pools in proportion to their weights.
  - source_nearest: Source proximity. This is the intelligent DNS resolution feature. GTM returns different address pools based on the geographic location of the DNS query source, enabling users to access the nearest available endpoint.
* `enable_status` - (Required) The enable status of the domain instance:
  - enable: Enabled. The GTM instance's intelligent scheduling policy is active.
  - disable: Disabled. The GTM instance's intelligent scheduling policy is inactive.
* `remark` - (Optional) Remarks. The provided parameter value becomes the updated remark content.  
* `schedule_host_name` - (Optional) The host record of the GTM access domain.
* `schedule_rr_type` - (Required, ForceNew) Record type for the access domain name:  
  - A: IPv4 address  
  - AAAA: IPv6 address  
  - CNAME: domain name.  
* `schedule_zone_mode` - (Required, ForceNew) Access domain name assignment mode:  
  - custom: Custom access domain name. You define the host record and associate it with a primary or subdomain under the same account as the GTM instance to generate the access domain name.  
  - sys_assign: System-assigned access domain name. This mode is no longer supported. Do not select it.
* `schedule_zone_name` - (Optional) Zone name, which is the parent zone of the GTM access domain name. It is typically a hosted domain under the same account as the GTM instance in the Alibaba Cloud DNS console, and can be either a primary domain or a subdomain.  
* `sequence_lb_strategy_mode` - (Optional) When the load balancing policy between address pools is set to sequential mode, the service restoration mode for preceding resources after an anomaly is resolved is as follows:  
  - preemptive: Preemptive mode. When a preceding resource recovers, the address pool with the smaller sequence number is prioritized.  
  - non_preemptive: Non-preemptive mode. When a preceding resource recovers, the current address pool continues to be used.  
* `ttl` - (Required, Int) Global TTL (in seconds). This value determines how long the DNS records resolving the access domain name to addresses in the address pool are cached by the ISP's LocalDNS. Custom TTL values are supported.  

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<config_id>:<instance_id>`.
* `config_id` - The domain name instance configuration ID.
* `instance_id` - The instance ID associated with the GTM 3.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cloud Gtm Instance Config.
* `delete` - (Defaults to 5 mins) Used when delete the Cloud Gtm Instance Config.
* `update` - (Defaults to 5 mins) Used when update the Cloud Gtm Instance Config.

## Import

Alidns Cloud Gtm Instance Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_instance_config.example <config_id>:<instance_id>
```