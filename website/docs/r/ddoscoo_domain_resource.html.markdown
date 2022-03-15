---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_domain_resource"
sidebar_current: "docs-alicloud-resource-ddoscoo-domain-resource"
description: |-
  Provides a Alicloud Anti-DDoS Pro Domain Resource resource.
---

# alicloud\_ddoscoo\_domain\_resource

Provides a Anti-DDoS Pro Domain Resource resource.

For information about Anti-DDoS Pro Domain Resource and how to use it, see [What is Domain Resource](https://www.alibabacloud.com/help/en/doc-detail/157463.htm).

-> **NOTE:** Available in v1.123.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ddoscoo_domain_resource" "example" {
  domain       = "tftestacc1234.abc"
  rs_type      = 0
  instance_ids = ["ddoscoo-cn-6ja1rl4j****"]
  real_servers = ["177.167.32.11"]
  https_ext    = "{\"Http2\":1,\"Http2https\":0，\"Https2http\":0}"
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
}

```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) The domain name of the website that you want to add to the instance.
* `https_ext` - (Optional) The advanced HTTPS settings. This parameter takes effect only when the value of ProxyType includes https. This parameter is a string that contains a JSON struct. The JSON struct includes the following fields:
    - `Http2https`: specifies whether to turn on Enforce HTTPS Routing. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enforce HTTPS Routing is turned off. The value `1` indicates that Enforce HTTPS Routing is turned on. The default value is `0`. If your website supports both HTTP and HTTPS, this feature suits your needs. If you turn on the switch, all HTTP requests are redirected to HTTPS requests on port 443 by default.
    - `Https2http`: specifies whether to turn on Enable HTTP. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enable HTTP is turned off. The value `1` indicates that Enable HTTP is turned on. The default value is `0`. If your website does not support HTTPS, this feature suits your needs. If you turn on the switch, all HTTPS requests are redirected to HTTP requests and forwarded to origin servers. The feature can also redirect WebSockets requests to WebSocket requests. All requests are redirected over port 80.
    - `Http2`: specifies whether to turn on Enable HTTP/2. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enable HTTP/2 is turned off. The value `1` indicates that Enable HTTP/2 is turned on. The default value is `0`. After you turn on the switch, the protocol type is HTTP/2.
* `instance_ids` - (Required) A list of instance ID that you want to associate. If this parameter is empty, only the domain name of the website is added but no instance is associated with the website.
  **NOTE:** There is a potential diff error because of the order of `instance_ids` values indefinite. 
  So, from version 1.161.0, `instance_ids` type has been updated as `set` from `list`, 
  and you can use [tolist](https://www.terraform.io/language/functions/tolist) to convert it to a list.
* `proxy_types` - (Required, ForceNew) Protocol type and port number information.
* `real_servers` - (Required) the IP address. This field is required and must be a string array.
* `rs_type` - (Required, ForceNew) The address type of the origin server. Valid values: `0`: IP address. `1`: domain name. Use the domain name of the origin server if you deploy proxies, such as Web Application Firewall (WAF), between the origin server and the Anti-DDoS Pro or Anti-DDoS Premium instance. If you use the domain name, you must enter the address of the proxy, such as the CNAME of WAF.

#### Block proxy_types

The proxy_types supports the following: 

* `proxy_ports` - (Optional, ForceNew) the port number. This field is required and must be an integer.
* `proxy_type` - (Optional) the protocol type. This field is required and must be a string. Valid values: `http`, `https`, `websocket`, and `websockets`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain Resource. Value as `domain`.

## Import

Anti-DDoS Pro Domain Resource can be imported using the id, e.g.

```
$ terraform import alicloud_ddoscoo_domain_resource.example <domain>
```
