---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_load_balancer"
description: |-
  Provides a Alicloud ESA Load Balancer resource.
---

# alicloud_esa_load_balancer

Provides a ESA Load Balancer resource.



For information about ESA Load Balancer and how to use it, see [What is Load Balancer](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateLoadBalancer).

-> **NOTE:** Available since v1.262.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_load_balancer&exampleId=0d361cc7-092e-1eef-b8a7-b296dded48797b4484d6&activeTab=example&spm=docs.r.esa_load_balancer.0.0d361cc709&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
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

resource "alicloud_esa_site" "resource_Site_OriginPool" {
  site_name   = "${var.name}${random_integer.default.result}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_origin_pool" "resource_OriginPool_LoadBalancer_1_1" {
  origins {
    type    = "ip_domain"
    address = "www.example.com"
    header  = "{\"Host\":[\"www.example.com\"]}"
    enabled = true
    weight  = "30"
    name    = "origin1"
  }
  site_id          = alicloud_esa_site.resource_Site_OriginPool.id
  origin_pool_name = "originpool1"
  enabled          = true
}

resource "alicloud_esa_load_balancer" "default" {
  load_balancer_name = "lb.exampleloadbalancer.top"
  fallback_pool      = alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_1_1.origin_pool_id
  site_id            = alicloud_esa_site.resource_Site_OriginPool.id
  description        = var.name
  default_pools      = [alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_1_1.origin_pool_id]
  steering_policy    = "geo"
  monitor {
    type              = "ICMP Ping"
    timeout           = 5
    monitoring_region = "ChineseMainland"
    consecutive_up    = 3
    consecutive_down  = 5
    interval          = 60
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_load_balancer&spm=docs.r.esa_load_balancer.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `adaptive_routing` - (Optional, List) Cross-pool origin configuration. See [`adaptive_routing`](#adaptive_routing) below.
* `default_pools` - (Required, List) List of default pool IDs.
* `description` - (Optional) The detailed description of the load balancer for easy management and identification.
* `enabled` - (Optional) Whether the load balancer is enabled.
  - `true`: Enabled.
  - `false`: Not enabled.
* `fallback_pool` - (Required, Int) The fallback pool ID, to which traffic will be redirected if all other pools are unavailable.
* `load_balancer_name` - (Required, ForceNew) The name of the load balancer must meet the domain name format verification and be a subdomain name under the site.
* `monitor` - (Required, List) Monitor configuration for health check. See [`monitor`](#monitor) below.
* `random_steering` - (Optional, List) Weighted round-robin configuration, used to control the traffic distribution weights among different pools. See [`random_steering`](#random_steering) below.
* `region_pools` - (Optional) Address pools corresponding to primary regions.
* `rules` - (Optional, List) Rule configuration list, used to define behavior under specific conditions. See [`rules`](#rules) below.
* `session_affinity` - (Optional) Session persistence. Valid values:
  - `off`: Not enabled.
  - `ip`: Session persistence by IP.
  - `cookie`: Session persistence by Cookie.
* `site_id` - (Required, ForceNew, Int) The site ID.
* `steering_policy` - (Required) Load balancing policy.
* `sub_region_pools` - (Optional) Address pools corresponding to secondary regions. When multiple secondary regions share a set of address pools, the keys can be concatenated with commas.
* `ttl` - (Optional, Int) TTL value, the time-to-live for DNS records. The default value is 30. The value range is 10-600.

### `adaptive_routing`

The adaptive_routing supports the following:
* `failover_across_pools` - (Optional) Whether to failover across pools.
  - `true`: Yes.
  - `false`: No.

### `monitor`

The monitor supports the following:
* `consecutive_down` - (Optional, Int) The number of consecutive failed health checks before the backend is considered down, for example, 5.
* `consecutive_up` - (Optional, Int) The number of consecutive successful probes required to consider the target as up, e.g., 3.
* `expected_codes` - (Optional) Expected status code, such as 200,202, successful HTTP response.
* `follow_redirects` - (Optional) Whether to follow the redirect.
  - `true`: Yes.
  - `false`: No.
* `header` - (Optional) The HTTP headers to be included in the health check request.
* `interval` - (Optional, Int) The monitoring interval, such as 60 seconds, checks the frequency.
* `method` - (Optional) Monitor request methods, such as GET, methods in the HTTP protocol.
* `monitoring_region` - (Optional, Computed) Probe Point Region, default to Global
  - `Global`: Global.
  - `ChineseMainland`: Chinese mainland.
  - `OutsideChineseMainland`: Global (excluding the Chinese mainland).
* `path` - (Optional) The monitor checks the path, such as/healthcheck, the HTTP request path.
* `port` - (Optional, Int) The target port.
* `timeout` - (Optional, Int) The timeout for the health check, in seconds. The value range is 1-10.
* `type` - (Optional) The type of monitor protocol, such as HTTP, used for health checks. When the value is off, it indicates that no check is performed.

### `random_steering`

The random_steering supports the following:
* `default_weight` - (Optional, Int) The default round-robin weight, used for all pools that do not have individually specified weights. The value range is 0-100.
* `pool_weights` - (Optional, Map) Weight configuration for each backend server pool, where the key is the pool ID and the value is the weight coefficient. The weight coefficient represents the proportion of relative traffic distribution.

### `rules`

The rules supports the following:
* `fixed_response` - (Optional, List) Executes a specified response after matching the rule. See [`fixed_response`](#rules-fixed_response) below.
* `overrides` - (Optional) Modifies the load balancer configuration for the corresponding request after matching the rule. The fields in this configuration will override the corresponding fields in the load balancer configuration.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int) Order of rule execution. The smaller the value, the higher the priority for execution.
* `terminates` - (Optional) Whether to terminate the execution of subsequent rules.
  - `true`: Yes.
  - `false`: No.

### `rules-fixed_response`

The rules-fixed_response supports the following:
* `content_type` - (Optional) The Content-Type field in the HTTP Header.
* `location` - (Optional) The location field in the http return.
* `message_body` - (Optional) The body value of the response.
* `status_code` - (Optional, Int) Status Code.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<load_balancer_id>`.
* `load_balancer_id` - The unique identifier ID of the load balancer.
* `status` - The status of the load balancer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

ESA Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_load_balancer.example <site_id>:<load_balancer_id>
```
