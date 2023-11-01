---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_instance"
sidebar_current: "docs-alicloud-resource-ddosbgp-instance"
description: |-
  Provides a Alicloud Anti-DDoS Advanced(Ddosbgp) Instance Resource.
---

# alicloud_ddosbgp_instance

Anti-DDoS Advanced instance resource. "Ddosbgp" is the short term of this product.

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "tf-example"
}

resource "alicloud_ddosbgp_instance" "instance" {
  name             = var.name
  base_bandwidth   = 20
  bandwidth        = -1
  ip_count         = 100
  ip_type          = "IPv4"
  normal_bandwidth = 100
  type             = "Enterprise"
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Optional, ForceNew) Type of the instance. Valid values: `Enterprise`, `Professional`. Default to `Enterprise`  
* `name` - (Optional) Name of the instance. This name can have a string of 1 to 63 characters.
* `base_bandwidth` - (Optional, ForceNew) Base defend bandwidth of the instance. Valid values: 20. The unit is Gbps. Default to `20`.
* `bandwidth` - (Required, ForceNew) Elastic defend bandwidth of the instance. This value must be larger than the base defend bandwidth. Valid values: 51,91,101,201,301. The unit is Gbps.
* `ip_count` - (Required, ForceNew) IP count of the instance. Valid values: 100.
* `ip_type` - (Required, ForceNew) IP version of the instance. Valid values: IPv4,IPv6.
* `period` - (Optional) The duration that you will buy Ddosbgp instance (in month). Valid values: [1~9], 12, 24, 36. Default to 12. At present, the provider does not support modify "period".
* `normal_bandwidth` - (Required, ForceNew) Normal defend bandwidth of the instance. The unit is Gbps.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddosbgp.
## Import

Ddosbgp instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddosbgp.example ddosbgp-abc123456
```
