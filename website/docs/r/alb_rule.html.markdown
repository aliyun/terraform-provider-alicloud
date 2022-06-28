---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_rule"
sidebar_current: "docs-alicloud-resource-alb-rule"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Rule resource.
---

# alicloud\_alb\_rule

Provides a Application Load Balancer (ALB) Rule resource.

For information about Application Load Balancer (ALB) Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/doc-detail/214375.htm).

-> **NOTE:** Available in v1.133.0+.

-> **NOTE:** This version only supports forwarding rules in the request direction.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
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

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.vpcs.0.id
  server_group_name = var.name
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default.id
  listener_protocol    = "HTTP"
  listener_port        = 8080
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

resource "alicloud_alb_rule" "default" {
  rule_name   = var.name
  listener_id = alicloud_alb_listener.default.id
  priority    = "555"
  rule_conditions {
    cookie_config {
      values {
        key   = "created"
        value = "tf"
      }
    }
    type = "Cookie"
  }

  rule_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
    order = "9"
    type  = "ForwardGroup"
  }
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) Specifies whether to precheck this request.
* `listener_id` - (Required, ForceNew) The ID of the listener to which the forwarding rule belongs.
* `priority` - (Required) The priority of the rule. Valid values: 1 to 10000. A smaller value indicates a higher priority. **Note*:* The priority of each rule within the same listener must be unique.
* `rule_actions` - (Required) The actions of the forwarding rules. See the following `Block rule_actions`.
* `rule_conditions` - (Required) The conditions of the forwarding rule. See the following `Block rule_conditions`.
* `rule_name` - (Required) The name of the forwarding rule. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.

### Block rule_conditions

The rule_conditions supports the following: 

* `type` - (Required) The type of the forwarding rule. Valid values: `Header`, `Host`, `Path`,  `Cookie`, `QueryString`, `Method` and `SourceIp`.
  * `Host`: Requests are forwarded based on the domain name. 
  * `Path`: Requests are forwarded based on the path. 
  * `Header`: Requests are forwarded based on the HTTP header field. 
  * `QueryString`: Requests are forwarded based on the query string. 
  * `Method`: Request are forwarded based on the request method. 
  * `Cookie`: Requests are forwarded based on the cookie.
  * `SourceIp`: Requests are forwarded based on the source ip. **NOTE:** The `SourceIp` option is available in 1.162.0+.
* `header_config` - (Optional) The configuration of the header field. See the following `Block header_config`.
* `cookie_config` - (Optional) The configuration of the cookie. See the following `Block cookie_config`.
* `host_config` - (Optional) The configuration of the host field. See the following `Block host_config`.
* `method_config` - (Optional) The configuration of the request method. See the following `Block method_config`.
* `path_config` - (Optional) The configuration of the path for the request to be forwarded. See the following `Block path_config`.
* `query_string_config` - (Optional) The configuration of the query string. See the following `Block query_string_config`.
* `source_ip_config` - (Optional, Available in 1.162.0+) The Based on source IP traffic matching. Required and valid when Type is SourceIP. See the following `Block source_ip_config`.

#### Block header_config

The header_config supports the following:

* `key` - (Optional) The key of the header field. The key must be 1 to 40 characters in length, and can contain letters, digits, hyphens (-) and underscores (_). The key does not support Cookie or Host.
* `values` - (Optional, Array) The value of the header field. The value must be 1 to 128 characters in length, and can contain lowercase letters, printable ASCII characters whose values are ch >= 32 && ch < 127, asterisks (*), and question marks (?). The value cannot start or end with a space.

#### Block cookie_config

The cookie_config supports the following: 

* `values` - (Optional, Array) The configuration of the cookie.

#### Block host_config

The host_config supports the following:

* `values` - (Optional, Array) The name of the host. **Note: ** The host name must meet the following rules: The hostname must be 3 to 128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost field can contain only letters and wildcards, and cannot contain digits or hyphens (-). Other fields cannot start or end with a hyphen (-). You can enter asterisks (*) and question marks (?) anywhere in a field.

#### Block method_config

The method_config supports the following:

* `values` - (Optional, Array) The request method. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, and `DELETE`.

#### Block path_config

The path_config supports the following:

* `values` - (Optional, Array) The path of the request to be forwarded. The path must be 1 to 128 characters in length and must start with a forward slash (/). The path can contain letters, digits, and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ". The value is case-sensitive, and can contain asterisks (*) and question marks (?).

#### Block source_ip_config

The source_ip_config supports the following:

* `values` - (Optional, Array) Add one or more IP addresses or IP address segments. You can add up to 5 forwarding rules in a SourceIp.

#### Block query_string_config

The query_string_config supports the following:

* `values` - (Optional, Array) The query string.
  * `key` - (Optional) The key must be 1 to 100 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The key cannot contain spaces or the following special characters: # [ ] { } \ | < > &.
  * `value` - (Optional) The value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The value cannot contain spaces or the following special characters: # [ ] { } \ | < > &.

### Block rule_actions

The rule_actions supports the following: 

* `order` - (Required) The order of the forwarding rule actions. Valid values: 1 to 50000. The actions are performed in ascending order. You cannot leave this parameter empty. Each value must be unique.
* `type` - (Required) The action. Valid values: `ForwardGroup`, `Redirect`, `FixedResponse`, `Rewrite`, `InsertHeader`, `TrafficLimit` and `TrafficMirror`. **Note:**  The preceding actions can be classified into two types:  `FinalType`: A forwarding rule can contain only one `FinalType` action, which is executed last. This type of action can contain only one `ForwardGroup`, `Redirect` or `FixedResponse` action. `ExtType`: A forwarding rule can contain one or more `ExtType` actions, which are executed before `FinalType` actions and need to coexist with the `FinalType` actions. This type of action can contain multiple `InsertHeader` actions or one `Rewrite` action. **NOTE:** The `TrafficLimit` and `TrafficMirror` option is available in 1.162.0+.
* `fixed_response_config` - (Optional) The configuration of the fixed response. See the following `Block fixed_response_config`.
* `insert_header_config` - (Optional) The configuration of the inserted header field. See the following `Block insert_header_config`.
* `redirect_config` - (Optional) The configuration of the external redirect action. See the following `Block redirect_config`.
* `rewrite_config` - (Optional) The redirect action within ALB. See the following `Block rewrite_config`.
* `forward_group_config` - (Optional) The forward response action within ALB. See the following `Block forward_group_config`.
* `traffic_limit_config` - (Optional, Available in 1.162.0+) The Flow speed limit. See the following `Block traffic_limit_config`.
* `traffic_mirror_config` - (Optional, Available in 1.162.0+) The Traffic mirroring. See the following `Block traffic_mirror_config`.

#### Block rewrite_config

The rewrite_config supports the following: 

* `host` - (Optional) The host name of the destination to which requests are redirected within ALB.  Valid values:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
* `path` - (Optional) The path to which requests are to be redirected within ALB.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?)and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. This value can be used only once. You can use it with a valid string.
* `query` - (Optional) The query string of the request to be redirected within ALB.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. This value can be used only once. You can use it with a valid string.

#### Block redirect_config

The redirect_config supports the following: 

* `host` - (Optional) The host name of the destination to which requests are directed.  The host name must meet the following rules:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
* `http_code` - (Optional) The redirect method. Valid values:301, 302, 303, 307, and 308.
* `path` - (Optional) The path of the destination to which requests are directed.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?) and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them with a valid string.
* `port` - (Optional) The port of the destination to which requests are redirected.  Valid values: 1 to 63335.  Default value: ${port}. You cannot use this value together with other characters at the same time.
* `protocol` - (Optional) The protocol of the requests to be redirected.  Valid values: HTTP and HTTPS.  Default value: ${protocol}. You cannot use this value together with other characters at the same time.  Note HTTPS listeners can redirect only HTTPS requests.
* `query` - (Optional) The query string of the request to be redirected.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them together with a valid string.

#### Block insert_header_config

The insert_header_config supports the following: 

* `key` - (Optional) The name of the inserted header field. The name must be 1 to 40 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). You cannot use the same name in InsertHeader.  Note You cannot use Cookie or Host in the name.
* `value` - (Optional) The content of the inserted header field:  If the ValueType parameter is set to SystemDefined, the following values are used:  ClientSrcPort: the port of the client ClientSrcIp: the IP address of the client Protocol: the protocol used by client requests (HTTP or HTTPS) SLBId: the ID of the ALB instance SLBPort: the listener port of the ALB instance If the ValueType parameter is set to UserDefined: The header value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters whose ASCII value is ch >= 32 && ch < 127, and wildcards such as asterisks (*) and question marks (?). The header value cannot start or end with a space.  If the ValueType parameter is set to ReferenceHeader: The header value must be 1 to 128 characters in length, and can contain lowercase letters, digits, underscores (_), and hyphens (-). Valid values: `ClientSrcPort`, `ClientSrcIp`, `Protocol`, `SLBId`, `SLBPort`, `UserDefined`.
* `value_type` - (Optional) Valid values:  UserDefined: a custom value ReferenceHeader: uses a field of the user request header. SystemDefined: a system value.

#### Block fixed_response_config

The fixed_response_config supports the following: 

* `content` - (Optional) The fixed response. The response cannot exceed 1 KB in size and can contain only ASCII characters.
* `content_type` - (Optional) The format of the fixed response.  Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, and `application/json`.
* `http_code` - (Optional) The HTTP status code of the response. The code must be an `HTTP_2xx`, `HTTP_4xx` or `HTTP_5xx.x` is a digit.

#### Block forward_group_config

The forward_group_config supports the following:

* `server_group_tuples` - (Optional, Array) The destination server group to which requests are forwarded.
  * `server_group_id` - (Optional) The ID of the destination server group to which requests are forwarded.
  * `weight` - (Optional, Computed, Available in 1.162.0+) The Weight of server group.

#### Block traffic_limit_config

The traffic_limit_config supports the following:

* `qps` - (Optional) The Number of requests per second. Value range: 1~100000.

#### Block traffic_mirror_config

The traffic_mirror_config supports the following:

* `target_type` - (Optional) The Mirror target type.
* `mirror_group_config` - (Optional) The Traffic is mirrored to the server group. See the following `Block mirror_group_config`.

#### Block mirror_group_config

The mirror_group_config supports the following:

* `server_group_tuples` - (Optional, Array) The destination server group to which requests are forwarded.
  * `server_group_id` - (Optional) The ID of the destination server group to which requests are forwarded.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Rule.
* `delete` - (Defaults to 2 mins) Used when delete the Rule.
* `update` - (Defaults to 2 mins) Used when update the Rule.

## Import

Application Load Balancer (ALB) Rule can be imported using the id, e.g.

```
$ terraform import alicloud_alb_rule.example <id>
```
