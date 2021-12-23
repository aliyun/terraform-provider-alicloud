---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_domain"
sidebar_current: "docs-alicloud-resource-waf-domain"
description: |-
  Provides a Web Application Firewall Domain resource.
---

# alicloud\_waf\_domain

Provides a WAF Domain resource to create domain in the Web Application Firewall.

For information about WAF and how to use it, see [What is Alibaba Cloud WAF](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available in 1.82.0+ .

## Example Usage

```
resource "alicloud_waf_domain" "domain" {
  domain            = "www.aliyun.com"
  instance_id       = "waf-123455"
  is_access_product = "On"
  source_ips        = ["1.1.1.1"]
  cluster_type      = "PhysicalCluster"
  http2_port        = [443]
  http_port         = [80]
  https_port        = [443]
  http_to_user_ip   = "Off"
  https_redirect    = "Off"
  load_balancing    = "IpHash"
  log_headers {
    key   = "foo"
    value = "http"
  }
}
```
## Argument Reference

The following arguments are supported:

* `cluster_type` - (Optional) The type of the WAF cluster. Valid values: `PhysicalCluster` and `VirtualCluster`. Default to `PhysicalCluster`.
* `connection_time` - (Optional) The connection timeout for WAF exclusive clusters. Unit: seconds.
* `domain` - (Optional, ForceNew, Deprecated in v1.94.0+)  Field `domain` has been deprecated from version 1.94.0. Use `domain_name` instead.
* `domain_name` - (Optional, ForceNew, Available in v1.94.0+) The domain that you want to add to WAF. The `domain_name` is required when the value of the `domain`  is Empty.
* `http2_port` - (Optional) List of the HTTP 2.0 ports.
* `http_port` - (Optional) List of the HTTP ports.
* `http_to_user_ip` - (Optional) Specifies whether to enable the HTTP back-to-origin feature. After this feature is enabled, the WAF instance can use HTTP to forward HTTPS requests to the origin server. 
By default, port 80 is used to forward the requests to the origin server. Valid values: `On` and `Off`. Default to `Off`.
* `https_port` - (Optional) List of the HTTPS ports.
* `https_redirect` - (Optional) Specifies whether to redirect HTTP requests as HTTPS requests. Valid values: "On" and `Off`. Default to `Off`.
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.
* `is_access_product` - (Required) Specifies whether to configure a Layer-7 proxy, such as Anti-DDoS Pro or CDN, to filter the inbound traffic before it is forwarded to WAF. Valid values: `On` and `Off`. Default to `Off`.
* `load_balancing` - (Optional) The load balancing algorithm that is used to forward requests to the origin. Valid values: `IpHash` and `RoundRobin`. Default to `IpHash`.
* `log_headers` - (Optional) The key-value pair that is used to mark the traffic that flows through WAF to the domain. Each item contains two field:
   * key: The key of label
   * value: The value of label
   
* `read_time` - (Optional) The read timeout of a WAF exclusive cluster. Unit: seconds.
* `resource_group_id` - (Optional) The ID of the resource group to which the queried domain belongs in Resource Management. By default, no value is specified, indicating that the domain belongs to the default resource group.
* `source_ips` - (Optional) List of the IP address or domain of the origin server to which the specified domain points.
* `write_time` - (Optional) The timeout period for a WAF exclusive cluster write connection. Unit: seconds.
			
## Attributes Reference

The following attributes are exported:

* `id` - This resource id. It formats as `<instance_id>:<domain>`
* `cname` - The CNAME record assigned by the WAF instance to the specified domain.

## Import

WAF domain can be imported using the id, e.g.

```
$ terraform import alicloud_waf_domain.domain waf-132435:www.domain.com
```
