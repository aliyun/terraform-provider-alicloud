---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_domain"
sidebar_current: "docs-alicloud-resource-wafv3-domain"
description: |-
  Provides a Alicloud Wafv3 Domain resource.
---

# alicloud_wafv3_domain

Provides a Wafv3 Domain resource.

For information about Wafv3 Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/web-application-firewall/latest/api-waf-openapi-2021-10-01-createdomain).

-> **NOTE:** Available since v1.200.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_wafv3_domain&exampleId=6883cb3a-92e1-c272-f108-a9584d515b4ce2e50787&activeTab=example&spm=docs.r.wafv3_domain.0.6883cb3a92&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `access_type` - (Optional) The access type of the WAF instance. Value: **share** (default): CNAME access.
* `domain` - (Required, ForceNew) The name of the domain name to query.
* `instance_id` - (Required, ForceNew) WAF instance ID
* `listen` - (Required) Configure listening information. See [`listen`](#listen) below.
* `redirect` - (Required) Configure forwarding information. See [`redirect`](#redirect) below.

### `listen`

The listen supports the following:
* `cert_id` - (Optional) The ID of the certificate to be added. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
* `cipher_suite` - (Optional) The type of encryption suite to add. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value:
		- **1**: indicates that all encryption suites are added.
		- **2**: indicates that a strong encryption package is added. You can select this value only if the value of **tls_version** is `tlsv1.2`.
		- **99**: indicates that a custom encryption suite is added.
* `custom_ciphers` - (Optional) The specific custom encryption suite to add.
* `enable_tlsv3` - (Optional) Whether TSL1.3 version is supported. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value:
		- **true**: indicates that TSL1.3 is supported.
		- **false**: indicates that TSL1.3 is not supported.
* `exclusive_ip` - (Optional) Whether to enable exclusive IP address. This parameter is used only when the value of **ipv6_enabled** is **false** (indicating that IPv6 is not enabled) and the value of **protection_resource** is **share** (indicating that a shared cluster is used). Value:
		- **true**: indicates that the exclusive IP address is enabled.
		- **false** (default): indicates that exclusive IP address is not enabled.
* `focus_https` - (Optional) Whether to enable the forced jump of HTTPS. This parameter is used only when the value of `https_ports` is not empty (indicating that the domain name uses HTTPS protocol) and the value of httports is empty (indicating that the domain name does not use HTTP protocol). Value:
		- **true**: indicates that HTTPS forced redirection is enabled.
		- **false**: indicates that HTTPS forced redirection is not enabled.
* `http2_enabled` - (Optional) Whether to turn on http2. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value:
		- **true:** indicates that HTTP2 is enabled.
		- **false** (default): indicates that HTTP2 is not enabled.
* `http_ports` - (Optional) The listening port of the HTTP protocol.
* `https_ports` - (Optional) The listening port of the HTTPS protocol.
* `ipv6_enabled` - (Optional) Whether IPv6 is turned on. Value:
		- **true**: indicates that IPv6 is enabled.
		- **false** (default): indicates that IPv6 is not enabled.
* `protection_resource` - (Optional) The type of protection resource to use. Value:
		- **share** (default): indicates that a shared cluster is used.
		- **gslb**: indicates that the shared cluster intelligent load balancing is used.
* `tls_version` - (Optional) The version of TLS to add. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value: **tlsv1**, **tlsv1.1**, **tlsv1.2**.
* `xff_header_mode` - (Optional) WAF obtains the real IP address of the client. Value:
		- **0** (default): indicates that the client has not forwarded the traffic to WAF through other layer -7 agents.
		- **1**: indicates that the first value of the X-Forwarded-For(XFF) field in the WAF read request header is used as the client IP address.
		- **2**: indicates that the custom field value set by you in the WAF read request header is used as the client IP address.
* `xff_headers` - (Optional) Set the list of custom fields used to obtain the client IP address.

### `redirect`

The Redirect supports the following:
* `backends` - (Optional) The IP address of the origin server corresponding to the domain name or the back-to-origin domain name of the server.
* `connect_timeout` - (Optional) Connection timeout. Unit: seconds, value range: 5~120.
* `focus_http_backend` - (Optional) Whether to enable forced HTTP back-to-origin. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value:
		- **true**: indicates that forced HTTP back-to-origin is enabled.
		- **false**: indicates that forced HTTP back-to-origin is not enabled.
* `keepalive` - (Optional) Open long connection, default true.
* `keepalive_requests` - (Optional) Number of long connections,  default: `60`. range :60-1000.
* `keepalive_timeout` - (Optional) Long connection over time, default: `15`. Range: 1-60.
* `loadbalance` - (Required) The load balancing algorithm used when returning to the source. Value:
		- **iphash**: indicates the IPHash algorithm.
		- **roundRobin**: indicates the polling algorithm.
		- **leastTime**: indicates the Least Time algorithm.
		- This value can be selected only if the value of **protection_resource** is **gslb** (indicating that the protected resource type uses shared cluster intelligent load balancing).
* `read_timeout` - (Optional) Read timeout duration. **Unit**: seconds, **Value range**: 5~1800.
* `request_headers` - (Optional) The traffic tag field and value of the domain name which used to mark the traffic processed by WAF. 
  It formats as `[{" k ":"_key_"," v ":"_value_"}]`. Where the `k` represents the specified custom request header field, 
  and the `v` represents the value set for this field. By specifying the custom request header field and the corresponding value, 
  when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value
  to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the
  custom header field already exists in the request, the system will overwrite the value of the custom field in the
  request with the set traffic tag value. See [`request_headers`](#redirect-request_headers) below.
* `retry` - (Optional) Back to Source Retry. default: true, retry 3 times by default.
* `sni_enabled` - (Optional) Whether to enable back-to-source SNI. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol). Value:
		- **true**: indicates that the back-to-source SNI is enabled.
		- **false** (default) indicates that the back-to-source SNI is not enabled.
* `sni_host` - (Optional) Sets the value of the custom SNI extension field. If this parameter is not set, the value of the **Host** field in the request header is used as the value of the SNI extension field by default.In general, you do not need to customize SNI unless your business has special configuration requirements. You want WAF to use SNI that is inconsistent with the actual request Host in the back-to-origin request (that is, the custom SNI set here).> This parameter is required only when **sni_enalbed** is set to **true** (indicating that back-to-source SNI is enabled).
* `write_timeout` - (Optional) Write timeout duration> **Unit**: seconds, **Value range**: 5~1800.

### `redirect-request_headers`

The request headers supports the following:
* `key` - (Optional) The traffic tag field and value of the domain name, which is used to mark the traffic processed by WAF. the format of this parameter value is **[{" k ":"_key_"," v ":"_value_"}]**. where_key_represents the specified custom request header field, and_value_represents the value set for this field.By specifying the custom request header field and the corresponding value, when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the custom header field already exists in the request, the system will overwrite the value of the custom field in the request with the set traffic tag value.
* `value` - (Optional) The traffic tag field and value of the domain name, which is used to mark the traffic processed by WAF. the format of this parameter value is **[{" k ":"_key_"," v ":"_value_"}]**. where_key_represents the specified custom request header field, and_value_represents the value set for this field.By specifying the custom request header field and the corresponding value, when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the custom header field already exists in the request, the system will overwrite the value of the custom field in the request with the set traffic tag value.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<domain>`.
* `resource_manager_resource_group_id` - The ID of the resource group.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.
* `update` - (Defaults to 5 mins) Used when update the Domain.

## Import

Wafv3 Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_domain.example <instance_id>:<domain>
```