---
subcategory: "WAF"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_domain"
sidebar_current: "docs-alicloud-resource-waf-domain"
description: |-
  Provides a WAF  Domain resource.
---


## Example Usage

```
# Add a WAF  Domain.
resource "alicloud_cdn_domain" "domain" {
  
  domain            = "www.test123.abc"
  instance_id       = "waf_elasticity-cn-0xldbqtm005"
  is_access_product = 0
  source_ips        = "1.1.1.1"
  cluster_type      = 0
  http2_port        = "433"
  http_port         = "80"
  https_port        = "433"
  http_to_user_ip   = 0
  https_redirect    = 0
  load_balancing    = 0
  resource_group_id = "rg-atstuj3rtoptyui"
  connection_time   = 60
}
```
## Argument Reference

The following arguments are supported:

* `domain` - (Required) Name of the  domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `instance_id` - (Required) The WAF instance ID. You can view the current WAF instance ID by calling the DescribeInstanceInfo interface.
* `cluster_type` - (Optional) The WAF instance cluster type. Values: `0`: physical cluster (default), `1`: virtual cluster, that is, WAF exclusive.
* `http2_port` - (Optional) Port configured by the HTTP 2.0 protocol. When specifying multiple HTTP 2.0 ports, use commas (,) to separate them.
* `http_port` - (Optional) Port configured by the HTTP protocol. When specifying multiple HTTP ports, use commas (,) to separate them.
* `http_to_user_ip` - (Optional) Whether to enable the HTTP back-to-source function. After enabling, the HTTPS access request will be forwarded back to the source site through the HTTP protocol. The default back-to-origin port is port 80. Values: `0`: disable(default), `1`: enable.
* `https_port` - (Optional) The port configured by the HTTPS protocol. When specifying multiple HTTPS ports, use commas (,) to separate them.
* `https_redirect` - (Optional) Whether to enable HTTPS forced redirection. Values: `0`: disable (default), `1`: enable.
* `is_access_product` - (Required) Whether the domain name is configured with a layer 7 proxy (for example, high defense, CDN...) before the WAF, that is, whether the client access traffic is forwarded by other layer 7 proxy before the WAF. The value: `0` means no,`1` means yes.
* `load_balancing` - (Optional) Load balancing algorithm used when returning to the source. Values: `0` indicates the IP Hash method, and `1` indicates the polling method.
* `log_headers` - (Optional) The traffic marking field and value of the domain name are used to mark the traffic passing through the WAF.
* `resource_group_id` - (Optional) The resource group ID to which the domain name belongs in the resource management product. The default is empty, which belongs to the default resource group.
* `source_ips` - (Optional) The source server IP corresponding to the domain name or the server returns the source domain name.

## Attributes Reference

The following attributes are exported:

* `cname` - The CNAME assigned by the WAF instance to this domain name configuration record.
* `request_id` - The request ID.

