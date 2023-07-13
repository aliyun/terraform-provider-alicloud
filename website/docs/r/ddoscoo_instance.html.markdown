---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instance"
sidebar_current: "docs-alicloud-resource-ddoscoo-instance"
description: |-
  Provides a Alicloud BGP-line Anti-DDoS Pro(Ddoscoo) Instance Resource.
---

# alicloud_ddoscoo_instance

BGP-Line Anti-DDoS instance resource. "Ddoscoo" is the short term of this product. See [What is Anti-DDoS Pro](https://www.alibabacloud.com/help/en/ddos-protection/latest/create-an-anti-ddos-pro-or-anti-ddos-premium-instance-by-calling-an-api-operation).

-> **NOTE:** The product region only support cn-hangzhou.

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available since v1.37.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

resource "alicloud_ddoscoo_instance" "default" {
  name              = var.name
  base_bandwidth    = "30"
  bandwidth         = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
  product_type      = "ddoscoo"
  period            = "1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the instance. This name can have a string of 1 to 63 characters.
* `base_bandwidth` - (Required) Base defend bandwidth of the instance. Valid values: `30`, `60`, `100`, `300`, `400`, `500`, `600`. The unit is Gbps. Only support upgrade.
* `bandwidth` - (Required) Elastic defend bandwidth of the instance. This value must be larger than the base defend bandwidth. Valid values: 30, 60, 100, 300, 400, 500, 600. The unit is Gbps. Only support upgrade.
* `service_bandwidth` - (Required) Business bandwidth of the instance. At leaset 100. Increased 100 per step, such as 100, 200, 300. The unit is Mbps. Only support upgrade.
* `port_count` - (Required) Port retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65. Only support upgrade.
* `domain_count` - (Required) Domain retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65. Only support upgrade.
* `edition_sale` - (Optional, ForceNew, Available since v1.212.0) The mitigation plan of the instance. Default value: `coop`. Valid values:
  - `coop`: Anti-DDoS Pro instance of the Profession mitigation plan.
* `address_type` - (Optional, ForceNew, Available since v1.212.0) The IP version of the IP address. Default value: `Ipv4`. Valid values: `Ipv4`, `Ipv6`.
* `bandwidth_mode` - (Optional, Available since v1.212.0) The mitigation plan of the instance. Valid values:
  - `0`: Disables the burstable clean bandwidth feature.
  - `1`: Enables the burstable clean bandwidth feature and uses the daily 95th percentile metering method.
  - `2`: Enables the burstable clean bandwidth feature and uses the monthly 95th percentile metering method.
* `product_type` - (Optional, Available since v1.125.0) The product type for purchasing DDOSCOO instances used to differ different account type. Default value: `ddoscoo`. Valid values:
  - `ddoscoo`: Only supports domestic account.
  - `ddoscoo_intl`: Only supports to international account.
* `period` - (Optional, Int) The duration that you will buy Ddoscoo instance (in month). Valid values: [1~9], `12`, `24`, `36`. Default value: `1`. At present, the provider does not support modify `period`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddoscoo.
* `ip` - (Available since v1.212.0) The IP address of the instance.

## Timeouts

-> **NOTE:** Available since v1.212.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Ddoscoo instance.
* `delete` - (Defaults to 3 mins) Used when delete the Ddoscoo instance.

## Import

Ddoscoo instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_instance.example ddoscoo-cn-123456
```
