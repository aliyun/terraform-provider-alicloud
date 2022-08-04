---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_ips"
sidebar_current: "docs-alicloud-datasource-ddos-bgp-ips"
description: |-
  Provides a list of Ddos Bgp Ips to the user.
---

# alicloud\_ddos\_bgp\_ips

This data source provides the Ddos Bgp Ips of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.180.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ddosbgp_ips" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "ddosbgp_ip_id_1" {
  value = data.alicloud_ddosbgp_ips.ids.ips.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ip IDs.
* `instance_id` - (Required, ForceNew) The ID of the native protection enterprise instance to be operated.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_name` - (Optional, ForceNew) The product name. Valid Value:`ECS`, `SLB`, `EIP`, `WAF`.
* `status` - (Optional, ForceNew) The current state of the IP address. Valid Value:
  - normal: indicates normal (not attacked).
  - hole_begin: indicates that you are in a black hole state.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ips` - A list of Ddos Bgp Ips. Each element contains the following attributes:
	* `id` - The ID of the Ip. The value formats as `<instance_id>:<ip>`.
	* `instance_id` - The ID of the native protection enterprise instance to be operated.
	* `ip` - The IP address.
	* `product` - The type of cloud asset to which the IP address belongs.
	* `status` - The current state of the IP address.