---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listener"
sidebar_current: "docs-alicloud-resource-alb-listener"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Listener resource.
---

# alicloud_alb_listener

Provides a Application Load Balancer (ALB) Listener resource.

For information about Application Load Balancer (ALB) Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createlistener).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_listener&exampleId=520fcb03-7402-42c2-3a88-7e64f6577970ee50cee7&activeTab=example&spm=docs.r.alb_listener.0.520fcb0374&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

* `access_log_record_customized_headers_enabled` - (Optional)Indicates whether the access log has a custom header field. Valid values: true and false. Default value: false.

-> **NOTE:** Only Instances outside the Security Group to Access the Log Switch **accesslogenabled** Open, in Order to Set This Parameter to the **True**.
* `access_log_tracing_config` - (Optional) Xtrace Configuration Information. See [`access_log_tracing_config`](#access_log_tracing_config) below for details.
* `certificates` - (Optional) The default certificate of the Listener. See [`certificates`](#certificates) below for details. **NOTE:** When `listener_protocol` is `HTTPS`, The default certificate must be set one。
* `default_actions` - (Optional) The Default Rule Action List. See [`default_actions`](#default_actions) below for details.
* `dry_run` - (Optional) The dry run.
* `gzip_enabled` - (Optional) Whether to Enable Gzip Compression, as a Specific File Type on a Compression. Valid values: `false`, `true`. Default Value: `true`. .
* `http2_enabled` - (Optional) Whether to Enable HTTP/2 Features. Valid Values: `True` Or `False`. Default Value: `True`.

-> **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.
* `idle_timeout` - (Optional) Specify the Connection Idle Timeout Value: `1` to `60`. Unit: Seconds.
* `listener_description` - (Optional)The description of the listener. The description must be 2 to 256 characters in length. The name can contain only the characters in the following string: `/^([^\x00-\xff]|[\w.,;/@-]){2,256}$/`.
* `listener_port` - (Required, ForceNew) The ALB Instance Front-End, and Those of the Ports Used. Value: `1` to `65535`.
* `listener_protocol` - (Required, ForceNew) Snooping Protocols. Valid Values: `HTTP`, `HTTPS` Or `QUIC`.
* `load_balancer_id` - (Required, ForceNew) The ALB Instance Id.
* `quic_config` - (Optional) Configuration Associated with the QuIC Listening. See [`quic_config`](#quic_config) below for details.
* `request_timeout` - (Optional) The Specified Request Timeout Time. Value: `1` to `180`. Unit: Seconds. Default Value: `60`. If the Timeout Time Within the Back-End Server Has Not Answered the ALB Will Give up Waiting, the Client Returns the HTTP 504 Error Code.
* `security_policy_id` - (Optional) Security Policy.

-> **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.

* `status` - (Optional, Available since v1.133.0) The state of the listener. Valid Values: `Running` Or `Stopped`. Valid values: `Running`: The listener is running. `Stopped`: The listener is stopped.
* `xforwarded_for_config` - (Optional, Deprecated since v1.161.0) xforwardfor Related Attribute Configuration. See [`xforwarded_for_config`](#xforwarded_for_config) below for details.  **NOTE:** 'xforwarded_for_config' has been deprecated since provider version 1.161.0. Use 'x_forwarded_for_config' instead.",
* `x_forwarded_for_config` - (Optional, Available since v1.161.0) The `x_forward_for` Related Attribute Configuration. See [`x_forwarded_for_config`](#x_forwarded_for_config) below for details. **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.
* `tags` - (Optional, Available since v1.215.0) A mapping of tags to assign to the resource.
* `acl_config` - (Optional, Deprecated since v1.163.0)The configurations of the access control lists (ACLs). See [`acl_config`](#acl_config) below for details. **NOTE:** Field `acl_config` has been deprecated from provider version 1.163.0, and it will be removed in the future version. Please use the new resource `alicloud_alb_listener_acl_attachment`.,

### `x_forwarded_for_config`

The x_forwarded_for_config supports the following:

* `x_forwarded_for_client_cert_issuer_dn_enabled` - (Optional) Indicates Whether the `X-Forwarded-Clientcert-issuerdn` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate after the Manifests Are Signed, the Publisher Information.
* `x_forwarded_for_client_cert_client_verify_alias` - (Optional) The Custom Header Field Names Only When `x_forwarded_for_client_cert_client_verify_enabled` Has a Value of True, this Value Will Not Take Effect until.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `x_forwarded_for_client_cert_client_verify_enabled` - (Optional) Indicates Whether the `X-Forwarded-Clientcert-clientverify` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate to Verify the Results.
* `x_forwarded_for_client_cert_finger_print_alias` - (Optional) The Custom Header Field Names Only When `x_forwarded_for_client_certfingerprint_enabled`, Which Evaluates to True When the Entry into Force of.The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `x_forwarded_for_client_cert_finger_print_enabled` - (Optional) Indicates Whether the `X-Forwarded-client_cert-fingerprint` Header Field Is Used to Obtain Access to the Server Load Balancer Instance of the Client Certificate Fingerprint Value.
* `x_forwarded_for_client_cert_subject_dn_alias` - (Optional) The name of the custom header. This parameter is valid only if `x_forwarded_for_client_certsubjectdn_enabled` is set to true. The name must be 1 to 40 characters in length, and can contain letters, hyphens (-), underscores (_), and digits.
* `x_forwarded_for_client_cert_subject_dn_enabled` - (Optional) Specifies whether to use the `X-Forwarded-client_cert-subjectdn` header field to obtain information about the owner of the ALB client certificate. Valid values: true and false. Default value: false.
* `x_forwarded_for_client_cert_issuer_dn_alias` - (Optional) The Custom Header Field Names Only When `x_forwarded_for_client_cert_issuer_dn_enabled`, Which Evaluates to True When the Entry into Force of.
* `x_forwarded_for_client_src_port_enabled` - (Optional) Indicates Whether the X-Forwarded-Client-Port Header Field Is Used to Obtain Access to Server Load Balancer Instances to the Client, and Those of the Ports.
* `x_forwarded_for_enabled` - (Optional) Whether to Enable by X-Forwarded-For Header Field Is Used to Obtain the Client IP Addresses.
* `x_forwarded_for_proto_enabled` - (Optional) Indicates Whether the X-Forwarded-Proto Header Field Is Used to Obtain the Server Load Balancer Instance Snooping Protocols.
* `x_forwarded_for_slb_id_enabled` - (Optional) Indicates Whether the SLB-ID Header Field Is Used to Obtain the Load Balancing Instance Id.
* `x_forwarded_for_slb_port_enabled` - (Optional) Indicates Whether the X-Forwarded-Port Header Field Is Used to Obtain the Server Load Balancer Instance Listening Port.
* `x_forwarded_for_client_source_ips_enabled` - (Optional) Whether to use the X-Forwarded-Client-Ip header to obtain the source IP address of the server load balancer instance. Value: true, false. Note HTTP, HTTPS, and QUIC listeners support this parameter. The function corresponding to this parameter is not open by default. Please contact the account manager if you need to use it.
* `x_forwarded_for_client_source_ips_trusted` - (Optional) Specify the trusted proxy IP. Application-oriented load balancing ALB will traverse the X-Forwarded-For from back to front, and select the first IP that is not in the trusted IP list as the real client IP, which will be used for the source IP speed limit.

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

### `quic_config`

The quic_config supports the following:

* `quic_listener_id` - (Optional) There Is a Need to Correlate the QuIC Listener ID. The Https Listener, in Effect at the Time. quicupgradeenabled True When Required.
* `quic_upgrade_enabled` - (Optional) Indicates Whether to Enable the QuIC Upgrade.

-> **NOTE:** The attribute is valid when the attribute `ListenerProtocol` is `HTTPS`.

### `default_actions`

The default_actions supports the following:

* `type` - (Required) Action Type.
* `forward_group_config` - (Required) The configurations of the actions. This parameter is required if Type is set to FowardGroup. See [`forward_group_config`](#default_actions-forward_group_config) below for details.

### `default_actions-forward_group_config`

The forward_group_config supports the following:

* `server_group_tuples` - (Required) The destination server group to which requests are forwarded. See [`server_group_tuples`](#default_actions-forward_group_config-server_group_tuples) below for details.

### `default_actions-forward_group_config-server_group_tuples`

The server_group_tuples supports the following:

* `server_group_id` - (Required) The ID of the destination server group to which requests are forwarded.

### `acl_config`

The acl_config supports the following:

* `acl_relations` - (Optional, Available since v1.136.0) The ACLs that are associated with the listener. See [`acl_relations`](#acl_config-acl_relations) below for details.
* `acl_type` - (Optional, Available since v1.136.0) The type of the ACL. Valid values: `White` Or `Black`. `White`: specifies the ACL as a whitelist. Only requests from the IP addresses or CIDR blocks in the ACL are forwarded. Whitelists apply to scenarios where only specific IP addresses are allowed to access an application. Risks may occur if the whitelist is improperly set. After you set a whitelist for an Application Load Balancer (ALB) listener, only requests from IP addresses that are added to the whitelist are distributed by the listener. If the whitelist is enabled without IP addresses specified, the ALB listener does not forward requests. `Black`: All requests from the IP addresses or CIDR blocks in the ACL are denied. The blacklist is used to prevent specified IP addresses from accessing an application. If the blacklist is enabled but the corresponding ACL does not contain IP addresses, the ALB listener forwards all requests.

### `acl_config-acl_relations`

The acl_relations supports the following:

* `acl_id` - (Optional, Available since v1.136.0) Snooping Binding of the Access Policy Group ID List.
* `status` - (Optional) The status of the ACL relation.

### `access_log_tracing_config`

The access_log_tracing_config supports the following:

* `tracing_enabled` - (Optional) Xtrace Function. Value: `True` Or `False` . Default Value: `False`.

-> **NOTE:** Only Instances outside the Security Group to Access the Log Switch `accesslogenabled` Open, in Order to Set This Parameter to the `True`.
* `tracing_sample` - (Optional) Xtrace Sampling Rate. Value: `1` to `10000`.

-> **NOTE:** This attribute is valid when `tracingenabled` is `true`.
* `tracing_type` - (Optional) Xtrace Type Value Is `Zipkin`.

-> **NOTE:** This attribute is valid when `tracingenabled` is `true`.


### `certificates`

The certificates supports the following:

* `certificate_id` - (Optional) The ID of the Certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Listener.
* `update` - (Defaults to 2 mins) Used when update the Listener.
* `delete` - (Defaults to 2 mins) Used when delete the Listener.

## Import

Application Load Balancer (ALB) Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_listener.example <id>
```
