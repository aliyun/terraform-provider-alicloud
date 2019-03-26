---
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instance"
sidebar_current: "docs-alicloud-resource-ddoscoo-instance"
description: |-
  Provides a Alicloud BGP-line Anti-DDoS Pro(Ddoscoo) Instance Resource.
---

# alicloud_ddoscoo_instance

Provides a Ddoscoo instance resource.When you create a Ddoscoo instance, you must enter the basic information about the instance, and define the instance request information, the instance backend service and response information.

Ddoscoo is the short term of BGP-line Anti-DDoS Pro and this is the product which provides protection against DDoS attacks.

~> **NOTE:** Terraform will auto build api while it uses `alicloud_ddoscoo_instance` to build instance.

## Example Usage

Basic Usage

```
resource "alicloud_ddoscoo_instance" "newInstance" {
    business_endpoint = "${var.business_endpoint}"
    band_width = "${var.band_width}"
    base_band_width     = "${var.base_band_width}"
    service_band_width       = "${var.service_band_width}"
    port_count  = "${var.port_count}"
    domain_count  = "${var.domain_count}"
}

variable "business_endpoint" {
  default = "business.aliyuncs.com"
}

variable "base_band_width" {
  default = "30"
}

variable "band_width" {
  default = "30"
}

variable "service_band_width" {
  default = "300"
}

variable "port_count" {
  default = "50"
}

variable "domain_count" {
  default = "55"
}
```
## Argument Reference

The following arguments are supported:

* `business_endpoint` - (Required) The endpoint to invoke for creating a new instance. The users in China Mainland should use 'business.aliyuncs.com'.
* `base_band_width` - (Required) The base defend bandwidth of a ddoscoo instance. The value of this argument is among 30, 60, 100, 300, 400, 500, 600. The unit is Gbps.
* `band_width` - (Required) The elastic bandwidth of a ddoscoo instance. This value must be larger than the base bandwidth and the max value is 600. The unit is Gbps.
* `service_band_width` - (Required) The business bandwidth of a ddoscoo instance. This value mustn't be less than 100. The increase of this argument must be 100. The unit is Mbps.
* `port_count` - (Required) Request_config defines how users can send requests to your API.
* `domain_count` - (Required) The type of backend service. Type including HTTP,VPC and MOCK. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddoscoo.