---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_listener"
sidebar_current: "docs-alicloud-resource-ga-listener"
description: |-
  Provides a Alicloud Global Accelerator (GA) Listener resource.
---

# alicloud_ga_listener

Provides a Global Accelerator (GA) Listener resource.

For information about Global Accelerator (GA) Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createlistener).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_listener&exampleId=461725b8-a426-c5c0-47c2-d1dacf19d4d468bd4e1f&activeTab=example&spm=docs.r.ga_listener.0.461725b8a4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  port_ranges {
    from_port = 80
    to_port   = 80
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_listener&spm=docs.r.ga_listener.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The accelerator id.
* `protocol` - (Optional) Type of network transport protocol monitored. Default value: `TCP`. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`.
* `proxy_protocol` - (Optional, Bool) The proxy protocol of the listener. Default value: `false`. Valid values:
  - `true`: Turn on the keep client source IP function. After it is turned on, the back-end service is supported to view the original IP address of the client.
  - `false`: Keep client source IP function is not turned on.
* `security_policy_id` - (Optional, Available since v1.183.0) The ID of the security policy. **NOTE:** Only `HTTPS` listeners support this parameter. Valid values:
  - `tls_cipher_policy_1_0`:
    - Supported TLS versions: TLS 1.0, TLS 1.1, and TLS 1.2.
    - Supported cipher suites: ECDHE-RSA-AES128-GCM-SHA256, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, AES128-GCM-SHA256, AES256-GCM-SHA384, AES128-SHA256, AES256-SHA256, ECDHE-RSA-AES128-SHA, ECDHE-RSA-AES256-SHA, AES128-SHA, AES256-SHA, and DES-CBC3-SHA.
  - `tls_cipher_policy_1_1`:
    - Supported TLS versions: TLS 1.1 and TLS 1.2.
    - Supported cipher suites: ECDHE-RSA-AES128-GCM-SHA256, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, AES128-GCM-SHA256, AES256-GCM-SHA384, AES128-SHA256, AES256-SHA256, ECDHE-RSA-AES128-SHA, ECDHE-RSA-AES256-SHA, AES128-SHA, AES256-SHA, and DES-CBC3-SHA.
  - `tls_cipher_policy_1_2`:
    - Supported TLS version: TLS 1.2.
    - Supported cipher suites: ECDHE-RSA-AES128-GCM-SHA256, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, AES128-GCM-SHA256, AES256-GCM-SHA384, AES128-SHA256, AES256-SHA256, ECDHE-RSA-AES128-SHA, ECDHE-RSA-AES256-SHA, AES128-SHA, AES256-SHA, and DES-CBC3-SHA.
  - `tls_cipher_policy_1_2_strict`:
    - Supported TLS version: TLS 1.2.
    - Supported cipher suites: ECDHE-RSA-AES128-GCM-SHA256, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, ECDHE-RSA-AES128-SHA, and ECDHE-RSA-AES256-SHA.
  - `tls_cipher_policy_1_2_strict_with_1_3`:
    - Supported TLS versions: TLS 1.2 and TLS 1.3.
    - Supported cipher suites: TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_CCM_SHA256, TLS_AES_128_CCM_8_SHA256, ECDHE-ECDSA-AES128-GCM-SHA256, ECDHE-ECDSA-AES256-GCM-SHA384, ECDHE-ECDSA-AES128-SHA256, ECDHE-ECDSA-AES256-SHA384, ECDHE-RSA-AES128-GCM-SHA256, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, ECDHE-ECDSA-AES128-SHA, ECDHE-ECDSA-AES256-SHA, ECDHE-RSA-AES128-SHA, and ECDHE-RSA-AES256-SHA.
* `listener_type` - (Optional, ForceNew, Available since v1.196.0) The routing type of the listener. Default Value: `Standard`. Valid values:
  - `Standard`: intelligent routing.
  - `CustomRouting`: custom routing.
* `http_version` - (Optional, Available since v1.220.0) The maximum version of the HTTP protocol. Default Value: `http2`. Valid values: `http1.1`, `http2`, `http3`.

  -> **NOTE:** `http_version` is only valid when `protocol` is `HTTPS`.
* `idle_timeout` - (Optional, Int, Available since v1.227.0) The timeout period of idle connections. Unit: seconds. Valid values:
  - If you set `protocol` to `TCP`. Default Value: `900`. Valid values: `10` to `900`.
  - If you set `protocol` to `UDP`. Default Value: `20`. Valid values: `10` to `20`.
  - If you set `protocol` to `HTTP` or `HTTPS`. Default Value: `15`. Valid values: `1` to `60`.
* `request_timeout` - (Optional, Int, Available since v1.227.0) The timeout period for HTTP or HTTPS requests. Unit: seconds. Default Value: `60`. Valid values: `1` to `180`.

  -> **NOTE:** `request_timeout` is only valid when `protocol` is `HTTP` or `HTTPS`.
* `client_affinity` - (Optional) The clientAffinity of the listener. Default value: `NONE`. Valid values:
  - `NONE`: client affinity is not maintained, that is, connection requests from the same client cannot always be directed to the same terminal node.
  - `SOURCE_IP`: maintain client affinity. When a client accesses a stateful application, all requests from the same client can be directed to the same terminal node, regardless of the source port and protocol.
* `name` - (Optional) The name of the listener. The length of the name is 2-128 characters. It starts with uppercase and lowercase letters or Chinese characters. It can contain numbers and underscores and dashes.
* `description` - (Optional) The description of the listener.
* `certificates` - (Optional, Set) The certificates of the listener. See [`certificates`](#certificates) below.

  -> **NOTE:** This parameter needs to be configured only for monitoring of the `HTTPS` protocol.
* `port_ranges` - (Required, Set) The portRanges of the listener. See [`port_ranges`](#port_ranges) below.

  -> **NOTE:** For `HTTP` or `HTTPS` protocol monitoring, only one monitoring port can be configured, that is, the start monitoring port and end monitoring port should be the same.
* `forwarded_for_config` - (Optional, Set, Available since v1.207.0) The XForward headers. See [`forwarded_for_config`](#forwarded_for_config) below.
* `backend_ports` - (Optional, Set, Available since v2.3.0) The ports on which the backend server receives requests. See [`backend_ports`](#backend_ports) below.

  -> **NOTE:** This parameter is required only when `protocol` is `HTTP` or `HTTPS` and the listener port is different from the port on which the backend server provides services. Once configured, removing the block keeps the existing backend ports on the service side.
* `endpoint_group_configurations` - (Optional, ForceNew, Set, Available since v2.3.0) The endpoint groups created together with the listener. See [`endpoint_group_configurations`](#endpoint_group_configurations) below.

  -> **NOTE:** This parameter is available only when `listener_type` is `Standard`. Endpoint groups created this way are removed automatically when the listener is destroyed. Do not manage the same endpoint groups through `alicloud_ga_endpoint_group` at the same time. Because `DescribeListener` does not return these settings, they are not restored by `terraform import`; importing a listener loses this configuration from state and a later apply recreates it.
* `custom_routing_endpoint_group_configurations` - (Optional, ForceNew, Set, Available since v2.3.0) The custom routing endpoint groups created together with the listener. See [`custom_routing_endpoint_group_configurations`](#custom_routing_endpoint_group_configurations) below.

  -> **NOTE:** This parameter is available only when `listener_type` is `CustomRouting`. The same import and mixed-management limitations as `endpoint_group_configurations` apply.

### `backend_ports`

The backend_ports supports the following:

* `from_port` - (Required, Int) The starting port on which the backend server receives requests.
* `to_port` - (Required, Int) The ending port on which the backend server receives requests.

### `endpoint_group_configurations`

The endpoint_group_configurations supports the following:

* `endpoint_group_name` - (Optional) The name of the endpoint group. The name must be 2 to 128 characters in length.
* `endpoint_group_description` - (Optional) The description of the endpoint group. The description can be up to 256 characters in length.
* `endpoint_group_region` - (Optional) The ID of the region in which the endpoint group resides.
* `traffic_percentage` - (Optional, Int) The percentage of traffic distributed to the endpoint group. Valid values: `0` to `100`.
* `health_check_enabled` - (Optional, Bool) Specifies whether to enable the health check feature.
* `health_check_interval_seconds` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Valid values: `1` to `60`.
* `health_check_path` - (Optional) The path to which the health check request is sent.
* `health_check_port` - (Optional, Int) The port that is used for health checks. Valid values: `1` to `65535`.
* `health_check_protocol` - (Optional) The protocol over which the health check is performed. Valid values: `TCP`, `HTTP`, `HTTPS`.
* `threshold_count` - (Optional, Int) The number of consecutive failed health checks that must occur before an endpoint is considered unhealthy. Valid values: `2` to `10`.
* `endpoint_configurations` - (Optional, Set) The endpoints in the endpoint group. The endpoint_configurations supports the following:
  * `type` - (Required) The type of the endpoint. Valid values: `Domain`, `Ip`, `IpTarget`, `PublicIp`, `ECS`, `SLB`, `ALB`, `NLB`, `ENI`, `OSS`.
  * `weight` - (Optional, Int) The weight of the endpoint. Valid values: `0` to `255`.
  * `endpoint` - (Required) The address of the endpoint.
  * `api_keys` - (Optional, Set) The keys that are used to authenticate the endpoint.
  * `vpc_id` - (Optional) The ID of the virtual private cloud (VPC) in which the endpoint resides.
  * `vswitch_ids` - (Optional, Set) The IDs of the vSwitches in which the endpoint resides.
  * `sub_address` - (Optional) The sub-address of the endpoint.
  * `enable_proxy_protocol` - (Optional, Bool) Specifies whether to use the proxy protocol to preserve client IP addresses for the endpoint.
  * `provider` - (Optional) The service provider of the endpoint.
  * `enable_client_ip_preservation` - (Optional, Bool) Specifies whether to preserve client IP addresses for the endpoint.
* `endpoint_request_protocol` - (Optional) The protocol over which the endpoint group receives requests. Valid values: `HTTP`, `HTTPS`.
* `endpoint_group_type` - (Optional) The type of the endpoint group. Valid values: `default`, `virtual`.
* `port_overrides` - (Optional, List) The port mapping between the listener and the endpoint group. The port_overrides supports the following:
  * `listener_port` - (Optional, Int) The port of the listener.
  * `endpoint_port` - (Optional, Int) The port of the endpoint.
* `enable_client_ip_preservation_proxy_protocol` - (Optional, Bool) Specifies whether to use the proxy protocol to preserve client IP addresses.
* `enable_client_ip_preservation_toa` - (Optional, Bool) Specifies whether to use the TOA module to preserve client IP addresses.
* `endpoint_protocol_version` - (Optional) The backend service protocol of the endpoint group. Valid values: `HTTP1.1`, `HTTP2`.
* `endpoint_ip_version` - (Optional) The IP version that GA uses to communicate with the endpoint service. Valid values: `IPv4`, `IPv6`, `ProtocolAffinity`.
* `health_check_host` - (Optional) The host header that is used for health checks.

### `custom_routing_endpoint_group_configurations`

The custom_routing_endpoint_group_configurations supports the following:

* `endpoint_group_region` - (Optional) The ID of the region in which the endpoint group resides.
* `name` - (Optional) The name of the endpoint group. The name must be 2 to 128 characters in length.
* `description` - (Optional) The description of the endpoint group. The description can be up to 256 characters in length.
* `destination_configurations` - (Optional, Set) The traffic destinations of the endpoint group. The destination_configurations supports the following:
  * `protocols` - (Optional, Set) The protocols of the traffic destination. Valid values: `TCP`, `UDP`.
  * `from_port` - (Optional, Int) The starting port of the traffic destination.
  * `to_port` - (Optional, Int) The ending port of the traffic destination.
* `custom_routing_endpoint_configurations` - (Optional, Set) The endpoints of the custom routing endpoint group. The custom_routing_endpoint_configurations supports the following:
  * `type` - (Optional) The type of the endpoint. Valid value: `PrivateSubNet`.
  * `endpoint` - (Optional) The address of the endpoint.
  * `traffic_to_endpoint_policy` - (Optional) The traffic policy for the endpoint. Valid values: `AllowAll`, `DenyAll`, `AllowCustom`.
  * `policy_configurations` - (Optional, Set) The traffic policy configurations of the endpoint. The policy_configurations supports the following:
    * `address` - (Optional) The destination address of the traffic policy.
    * `custom_routing_port_ranges` - (Optional, Set) The port ranges of the traffic policy. The custom_routing_port_ranges supports the following:
      * `from_port` - (Optional, Int) The starting port of the port range.
      * `to_port` - (Optional, Int) The ending port of the port range.

### `certificates`

The certificates supports the following:

* `id` - (Optional) The id of the certificate.

### `port_ranges`

The port_ranges supports the following:

* `from_port` - (Required, Int) The initial listening port used to receive requests and forward them to terminal nodes.
* `to_port` - (Required, Int) The end listening port used to receive requests and forward them to terminal nodes.

### `forwarded_for_config`

The forwarded_for_config supports the following:

* `forwarded_for_ga_id_enabled` - (Optional, Bool) Specifies whether to use the GA-ID header to retrieve the ID of the GA instance. Default value: `false`. Valid values:
  - `true`: yes.
  - `false`: no.
* `forwarded_for_ga_ap_enabled` - (Optional, Bool) Specifies whether to use the GA-AP header to retrieve the information about acceleration regions. Default value: `false`. Valid values:
  - `true`: yes.
  - `false`: no.
* `forwarded_for_proto_enabled` - (Optional, Bool) Specifies whether to use the GA-X-Forward-Proto header to retrieve the listener protocol of the GA instance. Default value: `false`. Valid values:
  - `true`: yes.
  - `false`: no.
* `forwarded_for_port_enabled` - (Optional, Bool) Specifies whether to use the GA-X-Forward-Port header to retrieve the listener ports of the GA instance. Default value: `false`. Valid values:
  - `true`: yes.
  - `false`: no.
* `real_ip_enabled` - (Optional, Bool) Specifies whether to use the X-Real-IP header to retrieve client IP addresses. Default value: `false`. Valid values:
  - `true`: yes.
  - `false`: no.

  -> **NOTE:** These `forwarded_for_ga_id_enabled`, `forwarded_for_ga_ap_enabled`, `forwarded_for_proto_enabled`, `forwarded_for_port_enabled`, `real_ip_enabled` are available only when you create an `HTTPS` or `HTTP` listener.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener.
* `status` - The status of the listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Listener.
* `update` - (Defaults to 3 mins) Used when update the Listener.
* `delete` - (Defaults to 10 mins) Used when delete the Listener.

## Import

Ga Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_listener.example <id>
```
