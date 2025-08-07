---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_domain"
description: |-
  Provides a Alicloud WAFV3 Domain resource.
---

# alicloud_wafv3_domain

Provides a WAFV3 Domain resource.



For information about WAFV3 Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/web-application-firewall/latest/api-waf-openapi-2021-10-01-createdomain).

-> **NOTE:** Available since v1.200.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_wafv3_instances" "default" {}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "example-value"
  cert             = file("${path.module}/test.crt")
  key              = file("${path.module}/test.key")
}

locals {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "cn-hangzhou"])
}

resource "alicloud_wafv3_domain" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    https_ports         = [443]
    http_ports          = [80]
    cert_id             = local.certificate_id
    cipher_suite        = 99
    xff_header_mode     = 2
    protection_resource = "share"
    tls_version         = "tlsv1"
    enable_tlsv3        = true
    http2_enabled       = true
    custom_ciphers      = ["ECDHE-ECDSA-AES128-GCM-SHA256"]
    focus_https         = false
    ipv6_enabled        = false
    exclusive_ip        = false
    xff_headers         = ["A", "B"]
  }
  redirect {
    backends           = ["1.1.1.1"]
    loadbalance        = "iphash"
    sni_enabled        = true
    sni_host           = "www.aliyun.com"
    focus_http_backend = false
    keepalive          = true
    retry              = true
    keepalive_requests = 80
    keepalive_timeout  = 30
    connect_timeout    = 30
    read_timeout       = 30
    write_timeout      = 30
    request_headers {
      key   = "A"
      value = "B"
    }
  }
  domain      = var.name
  access_type = "share"
}
```

## Argument Reference

The following arguments are supported:
* `access_type` - (Optional) The mode in which the domain name is added to WAF. Valid values:
share: CNAME record mode. This is the default value.
* `domain` - (Required, ForceNew) The name of the domain name to query.
* `instance_id` - (Required, ForceNew) The ID of the Web Application Firewall (WAF) instance.
* `listen` - (Required, List) Configure listening information. See [`listen`](#listen) below.
* `redirect` - (Required, List) Configure forwarding information. See [`redirect`](#redirect) below.
* `resource_manager_resource_group_id` - (Optional, ForceNew, Computed) The ID of the Alibaba Cloud resource group.
* `tags` - (Optional, Map, Available since v1.257.0) The tags. You can specify up to 20 tags.

### `listen`

The listen supports the following:
* `cert_id` - (Optional) The ID of the certificate to be added. This parameter is used only if the value of `HttpsPorts` is not empty (indicating that the domain name uses the HTTPS protocol).
* `cipher_suite` - (Optional, Int) The type of the cipher suites that you want to add. This parameter is available only if you specify `HttpsPorts`. Valid values:
  - `1`: all cipher suites.
  - `2`: strong cipher suites. This value is available only if you set `TLSVersion` to **tlsv1.2**.
  - `99`: custom cipher suites.
* `custom_ciphers` - (Optional, List) The specific custom encryption suite to add.
* `enable_tlsv3` - (Optional) Whether TSL1.3 version is supported. This parameter is used only if the value of `HttpsPorts` is not empty (indicating that the domain name uses the HTTPS protocol). Value:
  - `true`: indicates that TSL1.3 is supported.
  - `false`: indicates that TSL1.3 is not supported.
* `exclusive_ip` - (Optional) Specifies whether to enable the exclusive IP address feature. This parameter is available only if you set `IPv6Enabled` to false and `ProtectionResource` to `share`. Valid values:

  - `true`
  - `false` (default)
* `focus_https` - (Optional) Specifies whether to enable force redirect from HTTP to HTTPS for received requests. This parameter is available only if you specify `HttpsPorts` and leave `HttpPorts` empty. Valid values:

  - `true`
  - `false`
* `http2_enabled` - (Optional) Specifies whether to enable HTTP/2. This parameter is available only if you specify `HttpsPorts`. Valid values:

  - `true`
  - `false` (default)
* `http_ports` - (Optional, List) The HTTP listener ports. Specify the value in the \[**port1,port2,...**] format.
* `https_ports` - (Optional, List) The HTTPS listener ports. Specify the value in the \[**port1,port2,...**] format.
* `ipv6_enabled` - (Optional) Specifies whether to enable IPv6 protection. Valid values:

  - `true`
  - `false` (default)
* `protection_resource` - (Optional, Computed) The type of the protection resource. Valid values:

  - `share` (default): a shared cluster.
  - `gslb`: shared cluster-based intelligent load balancing.
* `sm2_access_only` - (Optional, Available since v1.257.0) Specifies whether to allow access only from SM certificate-based clients. This parameter is available only if you set SM2Enabled to true.

  - true
  - false
* `sm2_cert_id` - (Optional, Available since v1.257.0) The ID of the SM certificate that you want to add. This parameter is available only if you set SM2Enabled to true.
* `sm2_enabled` - (Optional, Computed, Available since v1.257.0) Specifies whether to add an SM certificate.
* `tls_version` - (Optional) The version of TLS to add. This parameter is used only if the value of `HttpsPorts` is not empty (indicating that the domain name uses the HTTPS protocol). Value:
  - `tlsv1`
  - **tlsv1.1**
  - **tlsv1.2**
* `xff_header_mode` - (Optional, Int) The method that is used to obtain the originating IP address of a client. Valid values:

  - `0` (default): Client traffic is not filtered by a Layer 7 proxy before the traffic reaches WAF.
  - `1`: WAF reads the first value of the X-Forwarded-For (XFF) header field as the originating IP address of the client.
  - `2`: WAF reads the value of a custom header field as the originating IP address of the client.
* `xff_headers` - (Optional, List) The custom header fields that are used to obtain the originating IP address of a client. Specify the value in the **\["header1","header2",...]** format.

-> **NOTE:**   This parameter is required only if you set `XffHeaderMode` to 2.


### `redirect`

The redirect supports the following:
* `backends` - (Optional, List) The IP addresses or domain names of the origin server. You cannot specify both IP addresses and domain names. If you specify domain names, the domain names can be resolved only to IPv4 addresses.

  - If you specify IP addresses, specify the value in the **\["ip1","ip2",...]** format. You can enter up to 20 IP addresses.
  - If you specify domain names, specify the value in the **\["domain"]** format. You can enter up to 20 domain names.
* `backup_backends` - (Optional, List, Available since v1.257.0) The secondary IP address or domain name of the origin server.
* `connect_timeout` - (Optional, Int) Connection timeout duration. Unit: seconds.
Value range: 1~3600. Default value: 5.
* `focus_http_backend` - (Optional) Specifies whether to enable force redirect from HTTPS to HTTP for back-to-origin requests. This parameter is available only if you specify `HttpsPorts`. Valid values:

  - `true`
  - `false`
* `keepalive` - (Optional) Specifies whether to enable the persistent connection feature. Valid values:

  - `true` (default)
  - `false`
* `keepalive_requests` - (Optional, Int) The number of reused persistent connections. Valid values: 60 to 1000. Default value: 1000


-> **NOTE:**   This parameter specifies the number of persistent connections that can be reused after you enable the persistent connection feature.

* `keepalive_timeout` - (Optional, Computed, Int) Idle long connection timeout, value range: 1~60, default 15, unit: seconds.

-> **NOTE:**  How long the multiplexed long connection is idle and then released.

* `loadbalance` - (Required) The load balancing algorithm that you want to use to forward requests to the origin server. Valid values:

  - `iphash`
  - `roundRobin`
  - `leastTime`: This value is available only if you set `ProtectionResource` to `gslb`.
* `read_timeout` - (Optional, Int) The timeout period of write connections. Unit: seconds. Valid values: 1 to 3600. Default value: 120.
* `request_headers` - (Optional, List) The traffic marking field and value of the domain name, which is used to mark the traffic processed by WAF.
By specifying custom request header fields and corresponding values, when the access traffic of the domain name passes through WAF, WAF automatically adds the set custom field value to the request header as a traffic mark, which facilitates the statistics of back-end services. See [`request_headers`](#redirect-request_headers) below.
* `retry` - (Optional) Specifies whether WAF retries if WAF fails to forward requests to the origin server. Valid values:

  - `true` (default)
  - `false`
* `sni_enabled` - (Optional) Specifies whether to enable the Server Name Indication (SNI) feature for back-to-origin requests. This parameter is available only if you specify `HttpsPorts`. Valid values:

  - `true`
  - `false` (default)
* `sni_host` - (Optional) The custom value of the SNI field. If you do not specify this parameter, the value of the `Host` header field is automatically used. In most cases, you do not need to specify a custom value for the SNI field. However, if you want WAF to use an SNI field whose value is different from the value of the Host header field in back-to-origin requests, you can specify a custom value for the SNI field.

-> **NOTE:**   This parameter is required only if you set `SniEnabled` to true.

* `write_timeout` - (Optional, Int) The timeout period of write connections. Unit: seconds. Valid values: 1 to 3600. Default value: 120.
* `xff_proto` - (Optional, Available since v1.257.0) Specifies whether to use the X-Forward-For-Proto header field to pass the protocol used by WAF to forward requests to the origin server. Valid values:
  - `true`  (default)
  - `false`

### `redirect-request_headers`

The redirect-request_headers supports the following:
* `key` - (Optional) Specified custom request header fields
* `value` - (Optional) Customize the value of the request header field.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<domain>`.
* `domain_id` - The domain ID.
* `status` - The status of the domain name. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.
* `update` - (Defaults to 5 mins) Used when update the Domain.

## Import

WAFV3 Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_domain.example <instance_id>:<domain>
```