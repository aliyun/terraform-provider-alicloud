---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instance"
sidebar_current: "docs-alicloud-resource-ddoscoo-instance"
description: |-
  Provides a Alicloud BGP-line Anti-DDoS Pro(Ddoscoo) Instance Resource.
---

# alicloud_ddoscoo_instance

BGP-Line Anti-DDoS instance resource. "Ddoscoo" is the short term of this product. See [What is Anti-DDoS Pro](https://www.alibabacloud.com/help/doc-detail/69319.htm).

-> **NOTE:** The product region only support cn-hangzhou.

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.37.0+ .

## Example Usage

Basic Usage

```
provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

resource "alicloud_ddoscoo_instance" "newInstance" {
  name              = "yourDdoscooInstanceName"
  bandwidth         = "30"
  base_bandwidth    = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
  period            = "1"
  product_type      = "ddoscoo"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the instance. This name can have a string of 1 to 63 characters.
* `base_bandwidth` - (Required) Base defend bandwidth of the instance. Valid values: 30, 60, 100, 300, 400, 500, 600. The unit is Gbps. Only support upgrade.
* `bandwidth` - (Required) Elastic defend bandwidth of the instance. This value must be larger than the base defend bandwidth. Valid values: 30, 60, 100, 300, 400, 500, 600. The unit is Gbps. Only support upgrade.
* `service_bandwidth` - (Required) Business bandwidth of the instance. At leaset 100. Increased 100 per step, such as 100, 200, 300. The unit is Mbps. Only support upgrade.
* `port_count` - (Required) Port retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65. Only support upgrade.
* `domain_count` - (Required) Domain retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65. Only support upgrade.
* `period` - (Optional, ForceNew) The duration that you will buy Ddoscoo instance (in month). Valid values: [1~9], 12, 24, 36. Default to 1. At present, the provider does not support modify "period".
* `product_type` - (Optional,Available in 1.125.0+ ) The product type for purchasing DDOSCOO instances used to differ different account type. Valid values:
  - ddoscoo: Only supports domestic account.
  - ddoscoo_intl: Only supports to international account.
  Default to ddoscoo.
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddoscoo.

## Import

Ddoscoo instance can be imported using the id, e.g.

```
$ terraform import alicloud_ddoscoo_instance.example ddoscoo-cn-123456
```
