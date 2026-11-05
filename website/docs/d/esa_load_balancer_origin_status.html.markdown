---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_load_balancer_origin_status"
sidebar_current: "docs-alicloud-datasource-esa-load-balancer-origin-status"
description: |-
  Provides the origin health status of the specified ESA load balancers.
---

# alicloud_esa_load_balancer_origin_status

Provides the origin health status of the specified ESA load balancers. The status reflects the probe results of origins under the default address pool; only the `default_pool` origins affect the load balancer status itself.

For information about ESA Load Balancer origin status and how to use it, see [ListLoadBalancerOriginStatus](https://next.api.alibabacloud.com/document/ESA/2024-09-10/ListLoadBalancerOriginStatus).

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

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "default" {
  site_name   = "${var.name}${random_integer.default.result}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_origin_pool" "default" {
  origins {
    type    = "ip_domain"
    address = "www.example.com"
    header  = "{\"Host\":[\"www.example.com\"]}"
    enabled = true
    weight  = "30"
    name    = "origin1"
  }
  site_id          = alicloud_esa_site.default.id
  origin_pool_name = "originpool1"
  enabled          = true
}

resource "alicloud_esa_load_balancer" "default" {
  load_balancer_name = "lb.exampleloadbalancer.top"
  fallback_pool      = alicloud_esa_origin_pool.default.origin_pool_id
  site_id            = alicloud_esa_site.default.id
  description        = var.name
  default_pools      = [alicloud_esa_origin_pool.default.origin_pool_id]
  steering_policy    = "geo"
  monitor {
    type              = "ICMP Ping"
    timeout           = 5
    monitoring_region = "Global"
    consecutive_up    = 3
    consecutive_down  = 5
    interval          = 60
    port              = 80
    method            = "GET"
    expected_codes    = "2xx"
    follow_redirects  = true
  }
}

data "alicloud_esa_load_balancer_origin_status" "default" {
  site_id           = alicloud_esa_site.default.id
  load_balancer_ids = [alicloud_esa_load_balancer.default.load_balancer_id]
  pool_type         = "default_pool"
}

output "alicloud_esa_load_balancer_origin_status_first_status" {
  value = data.alicloud_esa_load_balancer_origin_status.default.origin_status.0.status
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) The site ID. The origin status is queried within the scope of the site.
* `load_balancer_ids` - (Required) One or more load balancer IDs. The API accepts up to 100 IDs per request. Load balancer IDs can be obtained from the `alicloud_esa_load_balancer` resource or the `ListLoadBalancers` API.
* `pool_type` - (Optional) The origin address pool type. Only the `default_pool` origins affect the load balancer status. Pass `default_pool` to query the origins under the default address pool only.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id` - The ID of the data source, derived from the query conditions.
* `origin_status` - The origin status list of the load balancers. Each element contains the following attributes:
  * `load_balancer_id` - The ID of the load balancer.
  * `origin_id` - The origin ID.
  * `pool_id` - The origin address pool ID.
  * `pool_type` - The type of the origin address pool the origin belongs to. Only `default_pool` is reported; other types return an empty string.
  * `reason` - The probe failure reason.
  * `status` - The origin status. Valid values: `healthy`, `unhealthy`, `unknown` and `undetected`.
