---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_domain_resource"
description: |-
  Provides a Alicloud Ddos Coo Domain Resource resource.
---

# alicloud_ddoscoo_domain_resource

Provides a Ddos Coo Domain Resource resource.



For information about Ddos Coo Domain Resource and how to use it, see [What is Domain Resource](https://www.alibabacloud.com/help/en/anti-ddos/anti-ddos-pro-and-premium/developer-reference/api-ddoscoo-2020-01-01-createdomainresource).

-> **NOTE:** Available since v1.123.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddoscoo_domain_resource&exampleId=747dbf60-9e94-3f21-e0a9-2fae035dcaece1b3ac9e&activeTab=example&spm=docs.r.ddoscoo_domain_resource.0.747dbf609e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  https_ext    = <<EOF
    {
    "Http2": 1,
    "Http2https": 0,
    "Https2http": 0
  }
  EOF
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
}
```

## Argument Reference

The following arguments are supported:
* `cert` - (Optional, Available since v1.231.0) The private key of the certificate that you want to associate. This parameter must be used together with the CertName and Cert parameters.

-> **NOTE:**   If you specify a value for the CertName, Cert, and Key parameters, you do not need to specify a value for the CertId parameter.

* `cert_identifier` - (Optional, Available since v1.231.0) The name of the certificate.

-> **NOTE:**   You can specify the name of the certificate that you want to associate. From version 1.249.0, `cert_identifier` is in the "CertificateID-RegionId" format. For example, if the ID of the certificateId is `123`, and the region ID is `cn-hangzhou`, the value of the `cert_identifier` is `123-cn-hangzhou`.

* `cert_name` - (Optional, Computed, Available since v1.231.0) The public key of the certificate that you want to associate. This parameter must be used together with the CertName and Key parameters.

-> **NOTE:**   If you specify a value for the CertName, Cert, and Key parameters, you do not need to specify a value for the CertId parameter.

* `cert_region` - (Optional, Available since v1.231.0) The region of the certificate. `cn-hangzhou` and `ap-southeast-1` are supported. The default value is `cn-hangzhou`. 
* `custom_headers` - (Optional, Available since v1.261.0) The key-value pair of the custom header. The key specifies the header name, and the value specifies the header value. You can specify up to five key-value pairs. The key-value pairs can be up to 200 characters in length.
  Take note of the following items:
  - Do not use the following default HTTP headers:
  - X-Forwarded-ClientSrcPort: This header is used to obtain the source ports of clients that access Anti-DDoS Proxy (a Layer 7 proxy).
  - X-Forwarded-ProxyPort: This header is used to obtain the ports of listeners that access Anti-DDoS Proxy (a Layer 7 proxy).
  - X-Forwarded-For: This header is used to obtain the IP addresses of clients that access Anti-DDoS Proxy (a Layer 7 proxy).
  - Do not use standard HTTP headers or specific widely used custom HTTP headers. The standard HTTP headers include Host, User-Agent, Connection, and Upgrade, and the widely used custom HTTP headers include X-Real-IP, X-True-IP, X-Client-IP, Web-Server-Type, WL-Proxy-Client-IP, eEagleEye-RpcID, EagleEye-TraceID, X-Forwarded-Cluster, and X-Forwarded-Proto. If the preceding headers are used, the original content of the headers is overwritten.
* `domain` - (Required, ForceNew) The domain name for which you want to configure the Static Page Caching policy.

-> **NOTE:**  You can call the [DescribeDomains](https://www.alibabacloud.com/help/en/doc-detail/91724.html) operation to query all the domain names that are added to Anti-DDoS Pro or Anti-DDoS Premium.

* `https_ext` - (Optional, Computed, JsonString) The advanced HTTPS settings. This parameter takes effect only when the value of the `ProxyType` parameter includes `https`. The value is a string that consists of a JSON struct. The JSON struct contains the following fields:

  - `Http2https`: specifies whether to turn on Enforce HTTPS Routing. This field is optional and must be an integer. Valid values: `0` and `1`. The value 0 indicates that Enforce HTTPS Routing is turned off. The value 1 indicates that Enforce HTTPS Routing is turned on. The default value is 0.

    If your website supports both HTTP and HTTPS, this feature meets your business requirements. If you enable this feature, all HTTP requests to access the website are redirected to HTTPS requests on the standard port 443.

  - `Https2http`: specifies whether to turn on Enable HTTP. This field is optional and must be an integer. Valid values: `0` and `1`. The value 0 indicates that Enable HTTP is turned off. The value 1 indicates that Enable HTTP is turned on. The default value is 0.

    If your website does not support HTTPS, this feature meets your business requirements If this feature is enabled, all HTTPS requests are redirected to HTTP requests and forwarded to origin servers. This feature can redirect WebSockets requests to WebSocket requests. Requests are redirected over the standard port 80.

  - `Http2`: specifies whether to turn on Enable HTTP/2. This field is optional. Data type: integer. Valid values: `0` and `1`. The value 0 indicates that Enable HTTP/2 is turned off. The value 1 indicates that Enable HTTP/2 is turned on. The default value is 0.

    After you turn on the switch, HTTP/2 is used.
* `instance_ids` - (Required, Set) InstanceIds
* `key` - (Optional) The globally unique ID of the certificate. The value is in the "Certificate ID-cn-hangzhou" format. For example, if the ID of the certificate is 123, the value of the CertIdentifier parameter is 123-cn-hangzhou.

-> **NOTE:**   You can specify only one of this parameter and the CertId parameter.

* `ocsp_enabled` - (Optional, Bool, Available since v1.208.0) Specifies whether to enable the OCSP feature. Valid values:
  - `true`: Opened
  - `false`: Not enabled
* `proxy_types` - (Required, Set) Protocol type and port number information. See [`proxy_types`](#proxy_types) below.
* `real_servers` - (Required, Set, Available since v1.231.0) Server address information of the source station.
* `rs_type` - (Required, Int) The address type of the origin server. Valid values:

  - `0`: IP address

  - `1`: domain name

    If you deploy proxies, such as a Web Application Firewall (WAF) instance, between the origin server and the Anti-DDoS Pro or Anti-DDoS Premium instance, set the value to 1. If you use the domain name, you must enter the address of the proxy, such as the CNAME of WAF.

### `proxy_types`

The proxy_types supports the following:
* `proxy_ports` - (Required, Set) The port numbers.
* `proxy_type` - (Optional) The type of the protocol. Valid values:

  - `http`
  - `https`
  - `websocket`
  - `websockets`

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cname` - The CNAME address to query.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain Resource.
* `delete` - (Defaults to 5 mins) Used when delete the Domain Resource.
* `update` - (Defaults to 5 mins) Used when update the Domain Resource.

## Import

Ddos Coo Domain Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_domain_resource.example <id>
```