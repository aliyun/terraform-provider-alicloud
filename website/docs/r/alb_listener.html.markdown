---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listener"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Listener resource.
---

# alicloud_alb_listener

Provides a Application Load Balancer (ALB) Listener resource.



For information about Application Load Balancer (ALB) Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createlistener).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_alb_zones" "default" {
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_alb_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  count        = 2
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = format("10.4.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_alb_zones.default.zones[count.index].id
  vswitch_name = format("${var.name}_%d", count.index + 1)
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  description         = var.name
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = alicloud_vpc.default.id
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Basic"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_alb_zones.default.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.0.id
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.default.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_connect_port = "46325"
    health_check_enabled      = true
    health_check_host         = "tf-example.com"
    health_check_codes        = ["http_2xx", "http_3xx", "http_4xx"]
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_path         = "/tf-example"
    health_check_protocol     = "HTTP"
    health_check_timeout      = 5
    healthy_threshold         = 3
    unhealthy_threshold       = 3
  }
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-example"
    sticky_session_type    = "Server"
  }
  servers {
    description = var.name
    port        = 80
    server_id   = alicloud_instance.default.id
    server_ip   = alicloud_instance.default.private_ip
    server_type = "Ecs"
    weight      = 10
  }
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default.id
  listener_protocol    = "HTTP"
  listener_port        = 443
  listener_description = var.name
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `access_log_record_customized_headers_enabled` - (Optional) Access Log Whether to Enable Carry Custom Header Field. Valid values: `true`, `false`. Default Value: `false`.

-> **NOTE:**  Only Instances outside the Security Group to Access the Log Switch `accesslogenabled` Open, in Order to Set This Parameter to the `true`.

* `access_log_tracing_config` - (Optional, List) Xtrace Configuration Information. See [`access_log_tracing_config`](#access_log_tracing_config) below.
* `ca_certificates` - (Optional, List, Available since v1.242.0) The list of certificates. See [`ca_certificates`](#ca_certificates) below.
* `ca_enabled` - (Optional, Available since v1.242.0) Whether to turn on two-way authentication. Value:
  - `true`: on.
  - `false` (default): not enabled.
* `certificates` - (Optional, List) The list of certificates. See [`certificates`](#certificates) below.
* `default_actions` - (Required, List) The Default Rule Action List See [`default_actions`](#default_actions) below.
* `dry_run` - (Optional) Whether to PreCheck only this request. Value:
  - `true`: The check request is sent and the listener configuration is not updated. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - `false` (default): Sends a normal request, returns the 'HTTP 2xx' status code after passing the check, and performs the operation directly.
* `gzip_enabled` - (Optional, Computed) Whether to Enable Gzip Compression, as a Specific File Type on a Compression. Valid Values: True Or False. Default Value: TRUE.
* `http2_enabled` - (Optional, Computed) Whether to Enable HTTP/2 Features. Valid Values: True Or False. Default Value: TRUE.
* `idle_timeout` - (Optional, Computed, Int) Specify the Connection Idle Timeout Value: 1 to 60 seconds.
* `listener_description` - (Optional) Set the IP Address of the Listened Description. Length Is from 2 to 256 Characters.
* `listener_port` - (Required, ForceNew, Int) The SLB Instance Front-End, and Those of the Ports Used. Value: 1~65535.
* `listener_protocol` - (Required, ForceNew) Snooping Protocols. Valid Values: HTTP, HTTPS Or QuIC.
* `load_balancer_id` - (Required, ForceNew) The SLB Instance Id.
* `quic_config` - (Optional, Computed, List) Configuration Associated with the QuIC Listening See [`quic_config`](#quic_config) below.
* `request_timeout` - (Optional, Computed, Int) The Specified Request Timeout Time. Value: 1~180 Seconds. Default Value: 60 seconds. If the Timeout Time Within the Back-End Server Has Not Answered the SLB Will Give up Waiting, the Client Returns the HTTP 504 Error Code.
* `security_policy_id` - (Optional, Computed) Security Policy
* `status` - (Optional, Computed) The Current IP Address of the Listened State
* `tags` - (Optional, Map, Available since v1.215.0) The tag of the resource
* `x_forwarded_for_config` - (Optional, Computed, List, Available since v1.161.0) xforwardfor Related Attribute Configuration See [`x_forwarded_for_config`](#x_forwarded_for_config) below.
* `xforwarded_for_config` - (Optional, Deprecated since v1.161.0) xforwardfor Related Attribute Configuration. See [`xforwarded_for_config`](#xforwarded_for_config) below for details.  **NOTE:** 'xforwarded_for_config' has been deprecated since provider version 1.161.0. Use 'x_forwarded_for_config' instead.",
* `acl_config` - (Optional, Deprecated since v1.163.0)The configurations of the access control lists (ACLs). See [`acl_config`](#acl_config) below for details. **NOTE:** Field `acl_config` has been deprecated from provider version 1.163.0, and it will be removed in the future version. Please use the new resource `alicloud_alb_listener_acl_attachment`.,

### `access_log_tracing_config`

The access_log_tracing_config supports the following:
* `tracing_enabled` - (Required) Xtrace Function. Valid values: `true`, `false`. Default Value: `false`.

-> **NOTE:**  Only Instances outside the Security Group to Access the Log Switch `accesslogenabled` Open, in Order to Set This Parameter to the value `true`.

* `tracing_sample` - (Optional, Int) Xtrace Sampling Rate. Value: 1~10000. `tracingenabled` valued True When Effective.
* `tracing_type` - (Optional) Xtrace Type Value Is `Zipkin`.

-> **NOTE:**  `tracingenabled` valued True When Effective.


### `ca_certificates`

The ca_certificates supports the following:
* `certificate_id` - (Optional) The ID of the certificate. Currently, only server certificates are supported.

### `certificates`

The certificates supports the following:
* `certificate_id` - (Optional, Available since v1.242.0) The ID of the certificate. Currently, only server certificates are supported.

### `default_actions`

The default_actions supports the following:
* `forward_group_config` - (Optional, List) Forwarding Action Configurations See [`forward_group_config`](#default_actions-forward_group_config) below.
* `type` - (Required, ForceNew) The action type. Value: ForwardGroup, indicating forwarding to the server group.

### `default_actions-forward_group_config`

The default_actions-forward_group_config supports the following:
* `server_group_tuples` - (Required, List) The Forwarding Destination Server Group See [`server_group_tuples`](#default_actions-forward_group_config-server_group_tuples) below.

### `default_actions-forward_group_config-server_group_tuples`

The default_actions-forward_group_config-server_group_tuples supports the following:
* `server_group_id` - (Required) Forwarded to the Destination Server Group ID

### `quic_config`

The quic_config supports the following:
* `quic_listener_id` - (Optional) There Is a Need to Correlate the QuIC Listener ID. The Https Listener, in Effect at the Time. quicupgradeenabled True When Required.
* `quic_upgrade_enabled` - (Optional, Computed) Indicates Whether to Enable the QuIC Upgrade

### `x_forwarded_for_config`

The x_forwarded_for_config supports the following:
* `x_forwarded_for_client_cert_client_verify_alias` - (Optional) The Custom Header Field Names Only When xforwardedforclientcertclientverifyenabled Has a Value of True, this Value Will Not Take Effect until.
* `x_forwarded_for_client_cert_client_verify_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Clientcert-clientverify Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate to Verify the Results.
* `x_forwarded_for_client_cert_finger_print_alias` - (Optional) The Custom Header Field Names Only When xforwardedforclientcertfingerprintenabled, Which Evaluates to True When the Entry into Force of.
* `x_forwarded_for_client_cert_finger_print_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Clientcert-fingerprint Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate Fingerprint Value.
* `x_forwarded_for_client_cert_issuer_dn_alias` - (Optional) The Custom Header Field Names Only When xforwardedforclientcertsubjectdnenabled, Which Evaluates to True When the Entry into Force of.
* `x_forwarded_for_client_cert_issuer_dn_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Clientcert-issuerdn Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate after the Manifests Are Signed, the Publisher Information.
* `x_forwarded_for_client_cert_subject_dn_alias` - (Optional) The Custom Header Field Name,
* `x_forwarded_for_client_cert_subject_dn_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Clientcert-subjectdn Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate Owner Information.
* `x_forwarded_for_client_source_ips_enabled` - (Optional) Whether to use the X-Forwarded-Client-Ip header to obtain the source IP address of the server load balancer instance. Value:
  - `true`: Yes.
  - `false` (default): No.

-> **NOTE:** HTTP, HTTPS, and QUIC listeners support this parameter. The function corresponding to this parameter is not open by default. Please contact the account manager if you need to use it.

* `x_forwarded_for_client_source_ips_trusted` - (Optional) Specify the trusted proxy IP. Application-oriented load balancing ALB will traverse the X-Forwarded-For from back to front, and select the first IP that is not in the trusted IP list as the real client IP, which will be used for the source IP speed limit.
* `x_forwarded_for_client_src_port_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Client-Port Header Field Is Used to Obtain Access to Server Load Balancer Instances to the Client, and Those of the Ports.
* `x_forwarded_for_enabled` - (Optional, Computed) Whether to Enable by X-Forwarded-For Header Field Is Used to Obtain the Client IP Addresses.
* `x_forwarded_for_host_enabled` - (Optional, Available since v1.242.0) Whether to enable the X-Forwarded-Host header field to obtain the domain name of the client accessing the Application Load Balancer. Value:
  - `true`: Yes.
  - `false` (default): No.

-> **NOTE:** HTTP, HTTPS, and QUIC listeners support this parameter.

* `x_forwarded_for_processing_mode` - (Optional, Computed, Available since v1.242.0) Schema for processing X-Forwarded-For header fields. This value takes effect only when XForwardedForEnabled is true. Value:

  - `append` (default): append.
  - `remove`: Delete.

  Configure append to add the last hop IP address to the X-Forwarded-For header field before sending the request to the backend service.

  Configure remove to delete the X-Forwarded-For header before the request is sent to the backend service, regardless of whether the request carries X-Forwarded-For header fields.

  HTTP and HTTPS listeners support this parameter.
* `x_forwarded_for_proto_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Proto Header Field Is Used to Obtain the Server Load Balancer Instance Snooping Protocols.
* `x_forwarded_for_slb_id_enabled` - (Optional, Computed) Indicates Whether the SLB-ID Header Field Is Used to Obtain the Load Balancing Instance Id
* `x_forwarded_for_slb_port_enabled` - (Optional, Computed) Indicates Whether the X-Forwarded-Port Header Field Is Used to Obtain the Server Load Balancer Instance Listening Port

### `xforwarded_for_config`

The xforwarded_for_config supports the following:

* `xforwardedforclientcert_issuerdnenabled` - (Optional) Indicates Whether the `X-Forwarded-Clientcert-issuerdn` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate after the Manifests Are Signed, the Publisher Information.
* `xforwardedforclientcertclientverifyalias` - (Optional) The Custom Header Field Names Only When `xforwardedforclientcertclientverifyenabled` Has a Value of True, this Value Will Not Take Effect until.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `xforwardedforclientcertclientverifyenabled` - (Optional) Indicates Whether the `X-Forwarded-Clientcert-clientverify` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate to Verify the Results.
* `xforwardedforclientcertfingerprintalias` - (Optional) The Custom Header Field Names Only When `xforwardedforclientcertfingerprintenabled`, Which Evaluates to True When the Entry into Force of.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `xforwardedforclientcertfingerprintenabled` - (Optional) Indicates Whether the `X-Forwarded-Clientcert-fingerprint` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate Fingerprint Value.
* `xforwardedforclientcertsubjectdnalias` - (Optional) The name of the custom header. This parameter is valid only if `xforwardedforclientcertsubjectdnenabled` is set to true. The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `xforwardedforclientcertsubjectdnenabled` - (Optional) Specifies whether to use the `X-Forwarded-Clientcert-subjectdn` header field to obtain information about the owner of the ALB client certificate. Valid values: true and false. Default value: false.
* `xforwardedforclientcert_issuerdnalias` - (Optional) The Custom Header Field Names Only When `xforwardedforclientcert_issuerdnenabled`, Which Evaluates to True When the Entry into Force of.
* `xforwardedforclientsrcportenabled` - (Optional) Indicates Whether the X-Forwarded-Client-Port Header Field Is Used to Obtain Access to Server Load Balancer Instances to the Client, and Those of the Ports.
* `xforwardedforenabled` - (Optional) Whether to Enable by X-Forwarded-For Header Field Is Used to Obtain the Client IP Addresses.
* `xforwardedforprotoenabled` - (Optional) Indicates Whether the X-Forwarded-Proto Header Field Is Used to Obtain the Server Load Balancer Instance Snooping Protocols.
* `xforwardedforslbidenabled` - (Optional) Indicates Whether the SLB-ID Header Field Is Used to Obtain the Load Balancing Instance Id.
* `xforwardedforslbportenabled` - (Optional) Indicates Whether the X-Forwarded-Port Header Field Is Used to Obtain the Server Load Balancer Instance Listening Port.

### `acl_config`

The acl_config supports the following:

* `acl_relations` - (Optional, Available since v1.136.0) The ACLs that are associated with the listener. See [`acl_relations`](#acl_config-acl_relations) below for details.
* `acl_type` - (Optional, Available since v1.136.0) The type of the ACL. Valid values: `White` Or `Black`. `White`: specifies the ACL as a whitelist. Only requests from the IP addresses or CIDR blocks in the ACL are forwarded. Whitelists apply to scenarios where only specific IP addresses are allowed to access an application. Risks may occur if the whitelist is improperly set. After you set a whitelist for an Application Load Balancer (ALB) listener, only requests from IP addresses that are added to the whitelist are distributed by the listener. If the whitelist is enabled without IP addresses specified, the ALB listener does not forward requests. `Black`: All requests from the IP addresses or CIDR blocks in the ACL are denied. The blacklist is used to prevent specified IP addresses from accessing an application. If the blacklist is enabled but the corresponding ACL does not contain IP addresses, the ALB listener forwards all requests.

### `acl_config-acl_relations`

The acl_relations supports the following:

* `acl_id` - (Optional, Available since v1.136.0) Snooping Binding of the Access Policy Group ID List.
* `status` - (Optional) The status of the ACL relation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Listener.
* `delete` - (Defaults to 5 mins) Used when delete the Listener.
* `update` - (Defaults to 5 mins) Used when update the Listener.

## Import

Application Load Balancer (ALB) Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_listener.example <id>
```