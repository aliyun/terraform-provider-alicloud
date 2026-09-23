---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_domains"
sidebar_current: "docs-alicloud-datasource-waf-domains"
description: |-
  Provides a datasource to retrieve domain names.
---

# alicloud_waf_domains

Provides a WAF datasource to retrieve domains.

For information about WAF and how to use it, see [What is Alibaba Cloud WAF](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available since v1.86.0.

## Example Usage

```terraform
data "alicloud_waf_instances" "default" {}

data "alicloud_waf_domains" "default" {
  instance_id = data.alicloud_waf_instances.default.ids.0
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of WAF domain names. Each item is domain name.
* `name_regex` - (Optional) A regex string to filter results by domain name.
* `instance_id` - (Required) The Id of waf instance to which waf domain belongs.
* `resource_group_id` - (Optional, Available in v1.94.0+) The ID of the resource group to which the queried domain belongs in Resource Management.
* `enable_details` - (Optional, Available in v1.94.0+) Default to false and only output `id`, `domain_name`. Set it to true can output more details.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of WAF domain self ID, value as `domain_name`.
* `names` - A list of WAF domain names. 
* `domains` - A list of Domains. Each element contains the following attributes:
  * `domain` - Field `domain` has been deprecated from version 1.94.0. Use `domain_name` instead.
  * `domain_name` - Name of the domain.
  * `id` - The ID of domain self ID, value as `domain_name`.
  * `cluster_type` - The type of the WAF cluster.
  * `cname` - The CNAME record assigned by the WAF instance to the specified domain.
  * `connection_time` - The connection timeout for WAF exclusive clusters. Valid values: `PhysicalCluster` and `VirtualCluster`. Default to `PhysicalCluster`.
  * `http2_port` - List of the HTTP 2.0 ports.
  * `http_port` - List of the HTTP ports.
  * `http_to_user_ip` - Specifies whether to enable the HTTP back-to-origin feature. After this feature is enabled, the WAF instance can use HTTP to forward HTTPS requests to the origin server.
  * `https_port` - List of the HTTPS ports.
  * `https_redirect` - Specifies whether to redirect HTTP requests as HTTPS requests. Valid values: `On` and `Off`. Default to `Off`.
  * `is_access_product` - Specifies whether to configure a Layer-7 proxy, such as Anti-DDoS Pro or CDN, to filter the inbound traffic before it is forwarded to WAF. Valid values: `On` and "Off". Default to `Off`.
  * `load_balancing` - The load balancing algorithm that is used to forward requests to the origin. Valid values: `IpHash` and `RoundRobin`. Default to `IpHash`.
  * `log_headers` - The key-value pair that is used to mark the traffic that flows through WAF to the domain. Each item contains two field:
     * `key`: The key of label.
     * `value`: The value of label.
  * `read_time` - The read timeout of a WAF exclusive cluster. Unit: seconds.
  * `resource_group_id` - The ID of the resource group to which the queried domain belongs in Resource Management.
  * `source_ips` - List of the IP address or domain of the origin server to which the specified domain points.
  * `version` - The system data identifier that is used to control optimistic locking.
  * `write_time` - The timeout period for a WAF exclusive cluster write connection. Unit: seconds.
