---
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_instance"
sidebar_current: "docs-alicloud-resource-ddosbgp-instance"
description: |-
  Provides a Alicloud 
Anti-DDoS Advanced(Ddosbgp) Instance Resource.
---

# alicloud_ddoscoo_instance

Anti-DDoS Advanced instance resource. "Ddosbgp" is the short term of this product.

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

resource "alicloud_ddosbgp_instance" "instance" {
  name              = "yourDdosbgpInstanceName"
  region            = "cn-hangzhou"
  base_bandwidth    = "20"
  bandwidth         = "101"
  ip_count          = "100"
  ip_type           = "1"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the instance. This name can have a string of 1 to 63 characters.
* `region` - (Required) Region of the instance. Valid values: cn-hangzhou,cn-shanghai,cn-qingdao,cn-beijing,cn-zhangjiakou,cn-huhehaote,cn-shenzhen.
* `base_bandwidth` - (Required) Base defend bandwidth of the instance. Valid values: 20. The unit is Gbps.
* `bandwidth` - (Required) Elastic defend bandwidth of the instance. This value must be larger than the base defend bandwidth. Valid values: 51,101. The unit is Gbps.
* `ip_count` - (Required) IP count of the instance. Valid values: 100.
* `ip_type` - (Required) IP version of the instance. Valid values: v4,v6.
* `period` - (Optional, ForceNew) The duration that you will buy Ddosbgp instance (in month). Valid values: [1~9], 12, 24, 36. Default to 1. At present, the provider does not support modify "period".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddosbgp.