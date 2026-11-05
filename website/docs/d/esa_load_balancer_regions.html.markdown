---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_load_balancer_regions"
sidebar_current: "docs-alicloud-datasource-esa-load-balancer-regions"
description: |-
  Provides the schedulable region catalog of ESA load balancers.
---

# alicloud_esa_load_balancer_regions

Provides the schedulable region catalog of ESA load balancers. The catalog enumerates the first-level regions and their second-level sub-regions and requires no site or load balancer context.

For information about ESA Load Balancer regions and how to use it, see [ListLoadBalancerRegions](https://next.api.alibabacloud.com/document/ESA/2024-09-10/ListLoadBalancerRegions).

-> **NOTE:** Available since v1.262.0.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_load_balancer_regions" "default" {}

output "alicloud_esa_load_balancer_region_codes" {
  value = data.alicloud_esa_load_balancer_regions.default.regions[*].region_code
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source, derived from the returned region codes.
* `regions` - The region list. Each element contains the following attributes:
  * `region_code` - The first-level region code.
  * `region_cn_name` - The first-level region Chinese name.
  * `region_en_name` - The first-level region English name.
  * `sub_regions` - The second-level sub-region list. Each element contains:
    * `sub_region_code` - The second-level sub-region code.
    * `sub_region_cn_name` - The second-level sub-region Chinese name.
    * `sub_region_en_name` - The second-level sub-region English name.
