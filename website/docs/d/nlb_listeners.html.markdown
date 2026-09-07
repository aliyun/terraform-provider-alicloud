---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_listeners"
sidebar_current: "docs-alicloud-datasource-nlb-listeners"
description: |-
  Provides a list of Nlb Listeners to the user.
---

# alicloud_nlb_listeners

This data source provides the Nlb Listeners of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_listeners" "ids" {
  ids = ["example_value"]
}
output "alicloud_nlb_listener_id_1" {
  value = data.alicloud_nlb_listeners.ids.listeners.0.id
}
```

## Argument Reference

The following arguments are supported:

* `listener_protocol` - (Optional, ForceNew) The listening protocol. Valid values: `TCP`, `UDP`, or `TCPSSL`.
* `load_balancer_ids` - (Optional, ForceNew) The ID of the NLB instance. You can specify at most 20 IDs.
* `ids` - (Optional, ForceNew, Computed)  A list of Listener IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `listeners` - A list of Nlb Listeners. Each element contains the following attributes:
	* `alpn_enabled` - ndicates whether Application-Layer Protocol Negotiation (ALPN) is enabled.
	* `alpn_policy` - The ALPN policy.
	* `ca_certificate_ids` - CA certificate list information. Currently, only one CA certificate can be added. **NOTE:** This parameter only takes effect for `TCPSSL` listeners.
	* `ca_enabled` - Whether to start two-way authentication.
	* `certificate_ids` - Server certificate list information. Currently, only one server certificate can be added. This parameter only takes effect for `TCPSSL` listeners.
	* `cps` - The new connection speed limit for a network-based load balancing instance per second. Valid values: `0` ~ `1000000`. `0` indicates unlimited speed.
	* `end_port` - Full port listening end port. Valid values: `0` ~ `65535`. The value of the end port is less than the start port.
	* `idle_timeout` - Connection idle timeout time. Unit: seconds. Valid values: `1` ~ `900`.
	* `listener_description` - Custom listener name. The length is limited to 2 to 256 characters, supports Chinese and English letters, and can include numbers, commas (,), half-width periods (.), half-width semicolons (;), forward slashes (/), at(@), underscores (_), and dashes (-).
	* `listener_id` - The ID of the listener.
	* `listener_port` - Listening port. Valid values: `0` ~ `65535`. `0`: indicates that full port listening is used. When set to 0, you must configure `StartPort` and `EndPort`.
	* `listener_protocol` - The listening protocol. Valid values: `TCP`, `UDP`, or `TCPSSL`.
	* `load_balancer_id` - The ID of the network-based server load balancer instance.
	* `mss` - The maximum segment size of the TCP message. Unit: Bytes. Valid values: `0` ~ `1500`. `0` indicates that the MSS value of the TCP message is not modified. only `TCP` and `TCPSSL` listeners support this field value.
	* `proxy_protocol_enabled` - Whether to enable the Proxy Protocol to carry the source address of the client to the backend server.
	* `sec_sensor_enabled` - Indicates whether fine-grained monitoring is enabled.
	* `security_policy_id` - Security policy ID. Support system security policies and custom security policies. Valid values: `tls_cipher_policy_1_0`, `tls_cipher_policy_1_1`, `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`, or `tls_cipher_policy_1_2_strict_with_1_3`. **Note:** This parameter only takes effect for `TCPSSL` listeners.
	* `server_group_id` - The ID of the server group.
	* `start_port` - Full Port listens to the starting port. Valid values: `0` ~ `65535`.
	* `status` - The status of the resource.
	* `id` - The ID of the Nlb Listener.