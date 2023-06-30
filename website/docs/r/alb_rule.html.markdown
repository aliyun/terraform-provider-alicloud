---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_rule"
sidebar_current: "docs-alicloud-resource-alb-rule"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Rule resource.
---

# alicloud_alb_rule

Provides a Application Load Balancer (ALB) Rule resource.

For information about Application Load Balancer (ALB) Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-doc-alb-2020-06-16-api-doc-createrule).

-> **NOTE:** Available since v1.133.0.

-> **NOTE:** This version only supports forwarding rules in the request direction.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_alb_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

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

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = alicloud_vpc.default.id
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
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.default.id
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

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default.id
  listener_protocol    = "HTTP"
  listener_port        = 80
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
* `rule_actions` - (Required) The actions of the forwarding rules. See [`rule_actions`](#rule_actions) below for details.
* `rule_conditions` - (Required) The conditions of the forwarding rule. See [`rule_conditions`](#rule_conditions) below for details.
* `rule_name` - (Required) The name of the forwarding rule. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `direction` - (Optional, ForceNew, Available in 1.205.0+) The direction to which the forwarding rule is applied. Default value: `Request`. Valid values:
  - `Request`: The forwarding rule is applied to the client requests received by ALB.
  - `Response`: The forwarding rule is applied to the responses returned by backend servers.

### `rule_conditions`

The rule_conditions supports the following: 

* `type` - (Required) The type of the forwarding rule. Valid values: `Header`, `Host`, `Path`,  `Cookie`, `QueryString`, `Method` and `SourceIp`.
  * `Host`: Requests are forwarded based on the domain name. 
  * `Path`: Requests are forwarded based on the path. 
  * `Header`: Requests are forwarded based on the HTTP header field. 
  * `QueryString`: Requests are forwarded based on the query string. 
  * `Method`: Request are forwarded based on the request method. 
  * `Cookie`: Requests are forwarded based on the cookie.
  * `SourceIp`: Requests are forwarded based on the source ip. **NOTE:** The `SourceIp` option is available in 1.162.0+.
* `header_config` - (Optional) The configuration of the header field. See [`header_config`](#rule_conditions-header_config) below for details.
* `cookie_config` - (Optional) The configuration of the cookie. See See [`cookie_config`](#rule_conditions-cookie_config) below for details.
* `host_config` - (Optional) The configuration of the host field. See [`host_config`](#rule_conditions-host_config) below for details.
* `method_config` - (Optional) The configuration of the request method. See [`method_config`](#rule_conditions-method_config) below for details.
* `path_config` - (Optional) The configuration of the path for the request to be forwarded. See [`path_config`](#rule_conditions-path_config) below for details.
* `query_string_config` - (Optional) The configuration of the query string. See [`query_string_config`](#rule_conditions-query_string_config) below for details.
* `source_ip_config` - (Optional, Available in 1.162.0+) The Based on source IP traffic matching. Required and valid when Type is SourceIP. See [`source_ip_config`](#rule_conditions-source_ip_config) below for details.

### `rule_conditions-header_config`

The header_config supports the following:

* `key` - (Optional) The key of the header field. The key must be 1 to 40 characters in length, and can contain letters, digits, hyphens (-) and underscores (_). The key does not support Cookie or Host.
* `values` - (Optional, Array) The value of the header field. The value must be 1 to 128 characters in length, and can contain lowercase letters, printable ASCII characters whose values are ch >= 32 && ch < 127, asterisks (*), and question marks (?). The value cannot start or end with a space.

### `rule_conditions-cookie_config`

The cookie_config supports the following: 

* `values` - (Optional, Array) The configuration of the cookie. See [`values`](#rule_conditions-cookie_config-values) below for details.

### `rule_conditions-cookie_config-values`

The values supports the following:

* `key` - (Optional) The key of the values list.
* `value` - (Optional) The value of the values list.

### `rule_conditions-host_config`

The host_config supports the following:

* `values` - (Optional, Array) The name of the host. **Note: ** The host name must meet the following rules: The hostname must be 3 to 128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost field can contain only letters and wildcards, and cannot contain digits or hyphens (-). Other fields cannot start or end with a hyphen (-). You can enter asterisks (*) and question marks (?) anywhere in a field.

### `rule_conditions-method_config`

The method_config supports the following:

* `values` - (Optional, Array) The request method. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, and `DELETE`.

### `rule_conditions-path_config`

The path_config supports the following:

* `values` - (Optional, Array) The path of the request to be forwarded. The path must be 1 to 128 characters in length and must start with a forward slash (/). The path can contain letters, digits, and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ". The value is case-sensitive, and can contain asterisks (*) and question marks (?).

### `rule_conditions-source_ip_config`

The source_ip_config supports the following:

* `values` - (Optional, Array) Add one or more IP addresses or IP address segments. You can add up to 5 forwarding rules in a SourceIp.

### `rule_conditions-query_string_config`

The query_string_config supports the following:

* `values` - (Optional, Array) The query string. See [`values`](#rule_conditions-query_string_config-values) below for details.

### `rule_conditions-query_string_config-values`

The values supports the following:

* `key` - (Optional) The key must be 1 to 100 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The key cannot contain spaces or the following special characters: # [ ] { } \ | < > &.
* `value` - (Optional) The value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The value cannot contain spaces or the following special characters: # [ ] { } \ | < > &.

### `rule_actions`

The rule_actions supports the following: 

* `order` - (Required) The order of the forwarding rule actions. Valid values: 1 to 50000. The actions are performed in ascending order. You cannot leave this parameter empty. Each value must be unique.
* `type` - (Required) The action. Valid values: `ForwardGroup`, `Redirect`, `FixedResponse`, `Rewrite`, `InsertHeader`, `TrafficLimit`, `TrafficMirror` and `Cors`.
**Note:**  The preceding actions can be classified into two types:  `FinalType`: A forwarding rule can contain only one `FinalType` action, which is executed last. This type of action can contain only one `ForwardGroup`, `Redirect` or `FixedResponse` action. `ExtType`: A forwarding rule can contain one or more `ExtType` actions, which are executed before `FinalType` actions and need to coexist with the `FinalType` actions. This type of action can contain multiple `InsertHeader` actions or one `Rewrite` action.
**NOTE:** The `TrafficLimit` and `TrafficMirror` option is available in 1.162.0+.
**NOTE:** From version 1.205.0+, `type` can be set to `Cors`.
* `fixed_response_config` - (Optional) The configuration of the fixed response. See [`fixed_response_config`](#rule_actions-fixed_response_config) below for details.
* `insert_header_config` - (Optional) The configuration of the inserted header field. See [`insert_header_config`](#rule_actions-insert_header_config) below for details.
* `redirect_config` - (Optional) The configuration of the external redirect action. See [`redirect_config`](#rule_actions-redirect_config) below for details.
* `rewrite_config` - (Optional) The redirect action within ALB. See [`rewrite_config`](#rule_actions-rewrite_config) below for details.
* `forward_group_config` - (Optional) The forward response action within ALB. See [`forward_group_config`](#rule_actions-forward_group_config) below for details.
* `traffic_limit_config` - (Optional, Available in 1.162.0+) The Flow speed limit. See [`traffic_limit_config`](#rule_actions-traffic_limit_config) below for details.
* `traffic_mirror_config` - (Optional, Available in 1.162.0+) The Traffic mirroring. See [`traffic_mirror_config`](#rule_actions-traffic_mirror_config) below for details.
* `cors_config` - (Optional, Available in 1.205.0+) Request forwarding based on CORS. See [`cors_config`](#rule_actions-cors_config) below for details.

### `rule_actions-rewrite_config`

The rewrite_config supports the following: 

* `host` - (Optional) The host name of the destination to which requests are redirected within ALB.  Valid values:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
* `path` - (Optional) The path to which requests are to be redirected within ALB.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?)and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. This value can be used only once. You can use it with a valid string.
* `query` - (Optional) The query string of the request to be redirected within ALB.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. This value can be used only once. You can use it with a valid string.

### `rule_actions-redirect_config`
The redirect_config supports the following: 

* `host` - (Optional) The host name of the destination to which requests are directed.  The host name must meet the following rules:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
* `http_code` - (Optional) The redirect method. Valid values:301, 302, 303, 307, and 308.
* `path` - (Optional) The path of the destination to which requests are directed.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?) and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them with a valid string.
* `port` - (Optional) The port of the destination to which requests are redirected.  Valid values: 1 to 63335.  Default value: ${port}. You cannot use this value together with other characters at the same time.
* `protocol` - (Optional) The protocol of the requests to be redirected.  Valid values: HTTP and HTTPS.  Default value: ${protocol}. You cannot use this value together with other characters at the same time.  Note HTTPS listeners can redirect only HTTPS requests.
* `query` - (Optional) The query string of the request to be redirected.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them together with a valid string.

### `rule_actions-insert_header_config`

The insert_header_config supports the following: 

* `key` - (Optional) The name of the inserted header field. The name must be 1 to 40 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). You cannot use the same name in InsertHeader.  Note You cannot use Cookie or Host in the name.
* `value` - (Optional) The content of the inserted header field:  If the ValueType parameter is set to SystemDefined, the following values are used:  ClientSrcPort: the port of the client ClientSrcIp: the IP address of the client Protocol: the protocol used by client requests (HTTP or HTTPS) SLBId: the ID of the ALB instance SLBPort: the listener port of the ALB instance If the ValueType parameter is set to UserDefined: The header value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters whose ASCII value is ch >= 32 && ch < 127, and wildcards such as asterisks (*) and question marks (?). The header value cannot start or end with a space.  If the ValueType parameter is set to ReferenceHeader: The header value must be 1 to 128 characters in length, and can contain lowercase letters, digits, underscores (_), and hyphens (-). Valid values: `ClientSrcPort`, `ClientSrcIp`, `Protocol`, `SLBId`, `SLBPort`, `UserDefined`.
* `value_type` - (Optional) Valid values:  UserDefined: a custom value ReferenceHeader: uses a field of the user request header. SystemDefined: a system value.

### `rule_actions-fixed_response_config`

The fixed_response_config supports the following: 

* `content` - (Required) The fixed response. The response cannot exceed 1 KB in size and can contain only ASCII characters.
* `content_type` - (Optional) The format of the fixed response.  Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, and `application/json`.
* `http_code` - (Optional) The HTTP status code of the response. The code must be an `HTTP_2xx`, `HTTP_4xx` or `HTTP_5xx.x` is a digit.

### `rule_actions-forward_group_config`

The forward_group_config supports the following:

* `server_group_tuples` - (Optional, Array) The destination server group to which requests are forwarded. See [`server_group_tuples`](#rule_actions-forward_group_config-server_group_tuples) below for details.
* `server_group_sticky_session` - (Optional, Available in 1.179.0+) The configuration of session persistence for server groups. See [`server_group_sticky_session`](#rule_actions-forward_group_config-server_group_sticky_session) below for details.

### `rule_actions-forward_group_config-server_group_tuples`

The server_group_tuples supports the following:

* `server_group_id` - (Optional) The ID of the destination server group to which requests are forwarded.
* `weight` - (Optional, Available in 1.162.0+) The Weight of server group. Default value: `100`. **NOTE:** This attribute is required when the number of `server_group_tuples` is greater than 2.

### `rule_actions-forward_group_config-server_group_sticky_session`

The server_group_sticky_session supports the following:

* `enabled` - (Optional, Available in 1.179.0+) Whether to enable session persistence.
* `timeout` - (Optional, Available in 1.179.0+) The timeout period. Unit: seconds. Valid values: `1` to `86400`. Default value: `1000`.

### `rule_actions-traffic_limit_config`

The traffic_limit_config supports the following:

* `qps` - (Optional) The Number of requests per second. Value range: 1~100000.

### `rule_actions-traffic_mirror_config`

The traffic_mirror_config supports the following:

* `target_type` - (Optional) The Mirror target type.
* `mirror_group_config` - (Optional) The Traffic is mirrored to the server group. See [`mirror_group_config`](#rule_actions-traffic_mirror_config-mirror_group_config) below for details.

### `rule_actions-traffic_mirror_config-mirror_group_config`

The mirror_group_config supports the following:

* `server_group_tuples` - (Optional, Array) The destination server group to which requests are forwarded. See [`server_group_tuples`](#rule_actions-traffic_mirror_config-mirror_group_config-server_group_tuples) below for details.

### `rule_actions-traffic_mirror_config-mirror_group_config-server_group_tuples`

The server_group_tuples supports the following:

* `server_group_id` - (Optional) The ID of the destination server group to which requests are forwarded.

### `rule_actions-cors_config`

The cors_config supports the following:

* `allow_origin` - (Optional, Array) The allowed origins of CORS requests.
* `allow_methods` - (Optional, Array) The allowed HTTP methods for CORS requests. Valid values: `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `OPTIONS`, `PATCH`.
* `allow_headers` - (Optional, Array) The allowed headers for CORS requests.
* `expose_headers` - (Optional, Array) The headers that can be exposed.
* `allow_credentials` - (Optional) Specifies whether credentials can be passed during CORS operations. Valid values: `on`, `off`.
* `max_age` - (Optional) The maximum cache time of preflight requests in the browser. Unit: seconds. Valid values: `-1` to `172800`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Rule.
* `update` - (Defaults to 2 mins) Used when update the Rule.
* `delete` - (Defaults to 2 mins) Used when delete the Rule.

## Import

Application Load Balancer (ALB) Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_rule.example <id>
```
