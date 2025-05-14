---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_domains"
sidebar_current: "docs-alicloud-datasource-wafv3-domains"
description: |-
  Provides a list of Wafv3 Domain owned by an Alibaba Cloud account.
---

# alicloud_wafv3_domains

This data source provides the Wafv3 Domains of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.200.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_wafv3_instances" "default" {
}

data "alicloud_wafv3_domains" "ids" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  ids         = ["example_id"]
}

output "wafv3_domains_id_1" {
  value = data.alicloud_wafv3_domains.ids.domains.0.id
}

data "alicloud_wafv3_domains" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  domain      = "zctest12.wafqax.top"
}

output "wafv3_domains_id_2" {
  value = data.alicloud_wafv3_domains.default.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of domain IDs.
* `instance_id` - (Required, ForceNew) The WAF instance ID.
* `domain` - (Optional, ForceNew) The name of the domain name to query.
* `backend` - (Optional, ForceNew) The address type of the origin server. The address can be an IP address or a domain name. You can specify only one type of address.
* `enable_details` - (Optional, Bool) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `domains` - A list of Domain Entries. Each element contains the following attributes:
  * `id` - The ID of the domain. It formats as `<instance_id>:<domain>`.
  * `domain` - The name of the domain.
  * `cname` - The CNAME assigned by WAF to the domain name.
  * `resource_manager_resource_group_id` - The ID of the resource group.
  * `status` - The status of the domain.
  * `listen` - Configure listening information
    * `cert_id` - The ID of the certificate to be added. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `cipher_suite` - The type of encryption suite to add. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `custom_ciphers` - The specific custom encryption suite to add.
    * `enable_tlsv3` - Whether TSL1.3 version is supported. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `exclusive_ip` - Whether to enable exclusive IP address. This parameter is used only when the value of **ipv6_enabled** is **false** (indicating that IPv6 is not enabled) and the value of **protection_resource** is **share** (indicating that a shared cluster is used).
    * `focus_https` - Whether to enable the forced jump of HTTPS. This parameter is used only when the value of `https_ports` is not empty (indicating that the domain name uses HTTPS protocol) and the value of httports is empty (indicating that the domain name does not use HTTP protocol).
    * `http2_enabled` - Whether to turn on http2. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `http_ports` - The listening port of the HTTP protocol.
    * `https_ports` - The listening port of the HTTPS protocol.
    * `ipv6_enabled` - Whether IPv6 is turned on.
    * `protection_resource` - The type of protection resource to use.
    * `tls_version` - The version of TLS to add. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `xff_header_mode` - WAF obtains the real IP address of the client.
    * `xff_headers` - Set the list of custom fields used to obtain the client IP address.
  * `redirect` - Configure forwarding information.
    * `backends` - The IP address of the origin server corresponding to the domain name or the back-to-origin domain name of the server.
    * `connect_timeout` - Connection timeout, Unit: seconds, value range: 5~120.
    * `focus_http_backend` - Whether to enable forced HTTP back-to-origin. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `keepalive` - Open long connection, default true.
    * `keepalive_requests` - Number of long connections, default: `60`. range :60-1000.
    * `keepalive_timeout` - Long connection over time, default: `15`. Range: 1-60.
    * `loadbalance` - The load balancing algorithm used when returning to the source.
    * `read_timeout` - Read timeout duration. Unit: seconds, Value range: 5~1800.
    * `request_headers` - The traffic tag field and value of the domain name, which is used to mark the traffic processed by WAF. the format of this parameter value is **[{" k ":"_key_"," v ":"_value_"}]**. where_key_represents the specified custom request header field, and_value_represents the value set for this field.By specifying the custom request header field and the corresponding value, when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the custom header field already exists in the request, the system will overwrite the value of the custom field in the request with the set traffic tag value.
      * `key` - The traffic tag field and value of the domain name, which is used to mark the traffic processed by WAF. the format of this parameter value is **[{" k ":"_key_"," v ":"_value_"}]**. where_key_represents the specified custom request header field, and_value_represents the value set for this field.By specifying the custom request header field and the corresponding value, when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the custom header field already exists in the request, the system will overwrite the value of the custom field in the request with the set traffic tag value.
      * `value` - The traffic tag field and value of the domain name, which is used to mark the traffic processed by WAF. the format of this parameter value is **[{" k ":"_key_"," v ":"_value_"}]**. where_key_represents the specified custom request header field, and_value_represents the value set for this field.By specifying the custom request header field and the corresponding value, when the access traffic of the domain name passes through WAF, WAF automatically adds the specified custom field value to the request header as the traffic mark, which is convenient for backend service statistics.Explain that if the custom header field already exists in the request, the system will overwrite the value of the custom field in the request with the set traffic tag value.
    * `retry` - Back to Source Retry. default `true`, retry 3 times by default.
    * `sni_enabled` - Whether to enable back-to-source SNI. This parameter is used only if the value of **https_ports** is not empty (indicating that the domain name uses the HTTPS protocol).
    * `sni_host` - Sets the value of the custom SNI extension field. If this parameter is not set, the value of the **Host** field in the request header is used as the value of the SNI extension field by default.In general, you do not need to customize SNI unless your business has special configuration requirements. You want WAF to use SNI that is inconsistent with the actual request Host in the back-to-origin request (that is, the custom SNI set here).> This parameter is required only when **sni_enalbed** is set to **true** (indicating that back-to-source SNI is enabled).
    * `write_timeout` - Write timeout duration. **Unit**: seconds, **Value range**:5~1800.
