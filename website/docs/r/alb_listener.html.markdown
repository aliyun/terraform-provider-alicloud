---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listener"
sidebar_current: "docs-alicloud-resource-alb-listener"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Listener resource.
---

# alicloud\_alb\_listener

Provides a Application Load Balancer (ALB) Listener resource.

For information about Application Load Balancer (ALB) Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/doc-detail/214348.htm).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform

variable "name" {
  default = "example-name"
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count        = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_alb_zones.default.zones.0.id
  vswitch_name = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count        = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id      = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_alb_load_balancer" "default_3" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
  access_log_config {
    log_project = alicloud_log_project.default.name
    log_store   = alicloud_log_store.default.name
  }
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.vpcs.0.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
  tags = {
    Created = "TF"
  }
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "test"
  cert             = file("${path.module}/test.crt")
  key              = file("${path.module}/test.key")
}

resource "alicloud_alb_acl" "example" {
  acl_name = var.name
}

resource "alicloud_alb_listener" "example" {
  load_balancer_id     = alicloud_alb_load_balancer.default_3.id
  listener_protocol    = "HTTPS"
  listener_port        = 443
  listener_description = "createdByTerraform"
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
  }
  certificates {
    certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default.id, "-cn-hangzhou"])
  }
  acl_config {
    acl_type = "White"
    acl_relations {
      acl_id = alicloud_alb_acl.example.id
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `access_log_record_customized_headers_enabled` - (Optional, Computed)Indicates whether the access log has a custom header field. Valid values: true and false. Default value: false.

-> **NOTE:** Only Instances outside the Security Group to Access the Log Switch **accesslogenabled** Open, in Order to Set This Parameter to the **True**.
* `access_log_tracing_config` - (Optional) Xtrace Configuration Information. See the following `Block access_log_tracing_config`.
* `certificates` - (Optional) The default certificate of the Listener. See the following `Block certificates`. **NOTE:** When `listener_protocol` is `HTTPS`, The default certificate must be set one。
* `default_actions` - (Optional) The Default Rule Action List. See the following `Block default_actions`.
* `dry_run` - (Optional) The dry run.
* `gzip_enabled` - (Optional, Computed) Whether to Enable Gzip Compression, as a Specific File Type on a Compression. Valid values: `false`, `true`. Default Value: `true`. .
* `http2_enabled` - (Optional, Computed) Whether to Enable HTTP/2 Features. Valid Values: `True` Or `False`. Default Value: `True`.

-> **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.
* `idle_timeout` - (Optional, Computed) Specify the Connection Idle Timeout Value: `1` to `60`. Unit: Seconds.
* `listener_description` - (Optional)The description of the listener. The description must be 2 to 256 characters in length. The name can contain only the characters in the following string: `/^([^\x00-\xff]|[\w.,;/@-]){2,256}$/`.
* `listener_port` - (Required, ForceNew) The ALB Instance Front-End, and Those of the Ports Used. Value: `1` to `65535`.
* `listener_protocol` - (Required, ForceNew) Snooping Protocols. Valid Values: `HTTP`, `HTTPS` Or `QUIC`. 
* `load_balancer_id` - (Required, ForceNew) The ALB Instance Id.
* `quic_config` - (Optional) Configuration Associated with the QuIC Listening. See the following `Block quic_config`.
* `request_timeout` - (Optional, Computed) The Specified Request Timeout Time. Value: `1` to `180`. Unit: Seconds. Default Value: `60`. If the Timeout Time Within the Back-End Server Has Not Answered the ALB Will Give up Waiting, the Client Returns the HTTP 504 Error Code.
* `security_policy_id` - (Optional, Computed) Security Policy.

-> **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.

* `status` - (Optional, Computed, Available in v1.133.0+) The state of the listener. Valid Values: `Running` Or `Stopped`. Valid values: `Running`: The listener is running. `Stopped`: The listener is stopped.
* `xforwarded_for_config` - (Optional, Deprecated from 1.161.0+) xforwardfor Related Attribute Configuration. See the following `Block xforwarded_for_config`. **NOTE:** 'xforwarded_for_config' has been deprecated from provider version 1.161.0+. Use 'x_forwarded_for_config' instead.",
* `x_forwarded_for_config` - (Optional, Available from 1.161.0+) The `x_forward_for` Related Attribute Configuration. See the following `Block x_forwarded_for_config`. **NOTE:** The attribute is valid when the attribute `listener_protocol` is `HTTPS`.
* `acl_config` - (Optional, Available 1.136.0+)The configurations of the access control lists (ACLs). See the following `Block acl_config`. **NOTE:** Field `acl_config` has been deprecated from provider version 1.163.0, and it will be removed in the future version. Please use the new resource `alicloud_alb_listener_acl_attachment`.,


#### Block x_forwarded_for_config

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

#### Block xforwarded_for_config

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

#### Block quic_config

The quic_config supports the following: 

* `quic_listener_id` - (Optional) There Is a Need to Correlate the QuIC Listener ID. The Https Listener, in Effect at the Time. quicupgradeenabled True When Required.
* `quic_upgrade_enabled` - (Optional) Indicates Whether to Enable the QuIC Upgrade.

-> **NOTE:** The attribute is valid when the attribute `ListenerProtocol` is `HTTPS`.

#### Block default_actions

The default_actions supports the following: 

* `type` - (Required) Action Type.
* `forward_group_config` - (Required) The configurations of the actions. This parameter is required if Type is set to FowardGroup.
    *  `server_group_tuples` - (Required) The destination server group to which requests are forwarded.
        * `server_group_id` - (Required) The ID of the destination server group to which requests are forwarded.

#### Block acl_config

The acl_config supports the following:

* `acl_relations` - (Optional, Available 1.136.0+) The ACLs that are associated with the listener.
    * `acl_id` - (Optional, Available 1.136.0+) Snooping Binding of the Access Policy Group ID List.
* `acl_type` - (Optional, Available 1.136.0+) The type of the ACL. Valid values: `White` Or `Black`. `White`: specifies the ACL as a whitelist. Only requests from the IP addresses or CIDR blocks in the ACL are forwarded. Whitelists apply to scenarios where only specific IP addresses are allowed to access an application. Risks may occur if the whitelist is improperly set. After you set a whitelist for an Application Load Balancer (ALB) listener, only requests from IP addresses that are added to the whitelist are distributed by the listener. If the whitelist is enabled without IP addresses specified, the ALB listener does not forward requests. `Black`: All requests from the IP addresses or CIDR blocks in the ACL are denied. The blacklist is used to prevent specified IP addresses from accessing an application. If the blacklist is enabled but the corresponding ACL does not contain IP addresses, the ALB listener forwards all requests.

#### Block access_log_tracing_config

The access_log_tracing_config supports the following: 

* `tracing_enabled` - (Optional) Xtrace Function. Value: `True` Or `False` . Default Value: `False`.

-> **NOTE:** Only Instances outside the Security Group to Access the Log Switch `accesslogenabled` Open, in Order to Set This Parameter to the `True`.
* `tracing_sample` - (Optional) Xtrace Sampling Rate. Value: `1` to `10000`.

-> **NOTE:** This attribute is valid when `tracingenabled` is `true`.
* `tracing_type` - (Optional) Xtrace Type Value Is `Zipkin`.

-> **NOTE:** This attribute is valid when `tracingenabled` is `true`.


#### Block certificates

The certificates supports the following:

* `certificate_id` - (Optional) The ID of the Certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Listener.
* `delete` - (Defaults to 2 mins) Used when delete the Listener.
* `update` - (Defaults to 2 mins) Used when update the Listener.

## Import

Application Load Balancer (ALB) Listener can be imported using the id, e.g.

```
$ terraform import alicloud_alb_listener.example <id>
```
