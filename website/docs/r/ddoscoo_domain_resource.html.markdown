---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_domain_resource"
sidebar_current: "docs-alicloud-resource-ddoscoo-domain-resource"
description: |-
  Provides a Alicloud Anti-DDoS Pro Domain Resource resource.
---

# alicloud_ddoscoo_domain_resource

Provides a Anti-DDoS Pro Domain Resource resource.

For information about Anti-DDoS Pro Domain Resource and how to use it, see [What is Domain Resource](https://www.alibabacloud.com/help/en/ddos-protection/latest/api-ddoscoo-2020-01-01-createwebrule).

-> **NOTE:** Available since v1.123.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_ddoscoo_domain_resource&exampleId=cbaf9f64-f225-12d7-8c6c-154da65567f031b7f251&activeTab=example&spm=docs.r.ddoscoo_domain_resource.0.cbaf9f64f2" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
variable "domain" {
  default = "tf-example.alibaba.com"
}

resource "alicloud_ddoscoo_instance" "default" {
  name              = var.name
  bandwidth         = "30"
  base_bandwidth    = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
  period            = "1"
  product_type      = "ddoscoo"
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = var.domain
  rs_type      = 0
  instance_ids = [alicloud_ddoscoo_instance.default.id]
  real_servers = ["177.167.32.11"]
  https_ext    = "{\"Http2\":1,\"Http2https\":0,\"Https2http\":0}"
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) The domain name of the website that you want to add to the instance.
* `instance_ids` - (Required, Set) A list of instance ID that you want to associate. If this parameter is empty, only the domain name of the website is added but no instance is associated with the website.
-> **NOTE:** There is a potential diff error because of the order of `instance_ids` values indefinite. So, from version 1.161.0, `instance_ids` type has been updated as `set` from `list`, and you can use [tolist](https://www.terraform.io/language/functions/tolist) to convert it to a list.
* `real_servers` - (Required, List) the IP address. This field is required and must be a string array.
* `rs_type` - (Required, Int) The address type of the origin server. Use the domain name of the origin server if you deploy proxies, such as Web Application Firewall (WAF), between the origin server and the Anti-DDoS Pro or Anti-DDoS Premium instance. If you use the domain name, you must enter the address of the proxy, such as the CNAME of WAF. Valid values:
  - `0`: IP address.
  - `1`: domain name.
-> **NOTE:** From version 1.206.0, `rs_type` can be modified.
* `https_ext` - (Optional, Set) The advanced HTTPS settings. This parameter takes effect only when the value of ProxyType includes https. This parameter is a string that contains a JSON struct. The JSON struct includes the following fields:
  - `Http2https`: specifies whether to turn on Enforce HTTPS Routing. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enforce HTTPS Routing is turned off. The value `1` indicates that Enforce HTTPS Routing is turned on. The default value is `0`. If your website supports both HTTP and HTTPS, this feature suits your needs. If you turn on the switch, all HTTP requests are redirected to HTTPS requests on port 443 by default.
  - `Https2http`: specifies whether to turn on Enable HTTP. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enable HTTP is turned off. The value `1` indicates that Enable HTTP is turned on. The default value is `0`. If your website does not support HTTPS, this feature suits your needs. If you turn on the switch, all HTTPS requests are redirected to HTTP requests and forwarded to origin servers. The feature can also redirect WebSockets requests to WebSocket requests. All requests are redirected over port 80.
  - `Http2`: specifies whether to turn on Enable HTTP/2. This field is optional and must be an integer. Valid values: `0` and `1`. The value `0` indicates that Enable HTTP/2 is turned off. The value `1` indicates that Enable HTTP/2 is turned on. The default value is `0`. After you turn on the switch, the protocol type is HTTP/2.
* `ocsp_enabled` - (Optional, Bool, Available since v1.208.0) Specifies whether to enable the OCSP feature. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `proxy_types` - (Required, Set) Protocol type and port number information. See [`proxy_types`](#proxy_types) below.
-> **NOTE:** From version 1.206.0, `proxy_types` can be modified.

### `proxy_types`

The proxy_types supports the following: 

* `proxy_type` - (Optional) the protocol type. This field is required and must be a string. Valid values: `http`, `https`, `websocket`, and `websockets`.
* `proxy_ports` - (Optional, List) the port number. This field is required and must be an integer. **NOTE:** From version 1.206.0, `proxy_ports` can be modified.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain Resource. Value as `domain`.
* `cname` - (Available since v1.207.2) The CNAME assigned to the domain name.

## Import

Anti-DDoS Pro Domain Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_domain_resource.example <domain>
```
