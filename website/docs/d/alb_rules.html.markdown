---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_rules"
sidebar_current: "docs-alicloud-datasource-alb-rules"
description: |-
  Provides a list of Alb Rules to the user.
---

# alicloud\_alb\_rules

This data source provides the Alb Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_rules" "ids" {
  ids = ["example_id"]
}
output "alb_rule_id_1" {
  value = data.alicloud_alb_rules.ids.rules.0.id
}

data "alicloud_alb_rules" "nameRegex" {
  name_regex = "^my-Rule"
}
output "alb_rule_id_2" {
  value = data.alicloud_alb_rules.nameRegex.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Rule IDs.
* `listener_ids` - (Optional, ForceNew) The listener ids.
* `load_balancer_ids` - (Optional, ForceNew) The load balancer ids.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `rule_ids` - (Optional, ForceNew) The rule ids.
* `status` - (Optional, ForceNew) The status of the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of Alb Rules. Each element contains the following attributes:
    * `id` - The ID of the Rule.
    * `listener_id` - The ID of the listener to which the forwarding rule belongs.
    * `load_balancer_id` - The ID of the Application Load Balancer (ALB) instance to which the forwarding rule belongs.
    * `priority` - The priority of the rule. Valid values: 1 to 10000. A smaller value indicates a higher priority.  Note The priority of each rule within the same listener must be unique.
    * `rule_actions` - The actions of the forwarding rules.
        * `type` - The action type.
        * `fixed_response_config` - The configuration of the fixed response.
            * `content` - The fixed response. The response cannot exceed 1 KB in size and can contain only ASCII characters.
            * `content_type` - The format of the fixed response.  Valid values: text/plain, text/css, text/html, application/javascript, and application/json.
            * `http_code` - The HTTP status code of the response. The code must be an HTTP_2xx,HTTP_4xx or HTTP_5xx.x is a digit.
        * `forward_group_config` - The configurations of the destination server groups.
            * `server_group_tuples` - The destination server group to which requests are forwarded.
                * `server_group_id` - The ID of the destination server group to which requests are forwarded.
                * `weight` - The Weight of server group.
        * `insert_header_config` - The configuration of the inserted header field.
            * `key` - The name of the inserted header field. The name must be 1 to 40 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). You cannot use the same name in InsertHeader.  Note You cannot use Cookie or Host in the name.
            * `value` - The content of the inserted header field:  If the ValueType parameter is set to SystemDefined, the following values are used:  ClientSrcPort: the port of the client ClientSrcIp: the IP address of the client Protocol: the protocol used by client requests (HTTP or HTTPS) SLBId: the ID of the ALB instance SLBPort: the listener port of the ALB instance If the ValueType parameter is set to UserDefined: The header value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters whose ASCII value is ch >= 32 && ch < 127, and wildcards such as asterisks (*) and question marks (?). The header value cannot start or end with a space.  If the ValueType parameter is set to ReferenceHeader: The header value must be 1 to 128 characters in length, and can contain lowercase letters, digits, underscores (_), and hyphens (-).
            * `value_type` - Valid values:  UserDefined: a custom value ReferenceHeader: uses a field of the user request header. SystemDefined: a system value.
        * `order` - The order of the forwarding rule actions. Valid values:1 to 50000. The actions are performed in ascending order. You cannot leave this parameter empty. Each value must be unique.
        * `redirect_config` - The configuration of the external redirect action.
            * `host` - The host name of the destination to which requests are directed.  The host name must meet the following rules:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
            * `http_code` - The redirect method. Valid values:301, 302, 303, 307, and 308.
            * `path` - The path of the destination to which requests are directed.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?) and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them with a valid string.
            * `port` - The port of the destination to which requests are redirected.  Valid values: 1 to 63335.  Default value: ${port}. You cannot use this value together with other characters at the same time.
            * `protocol` - The protocol of the requests to be redirected.  Valid values: HTTP and HTTPS.  Default value: ${protocol}. You cannot use this value together with other characters at the same time.  Note HTTPS listeners can redirect only HTTPS requests.
            * `query` - The query string of the request to be redirected.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. You can also reference ${host}, ${protocol}, and ${port}. Each variable can appear at most once. You can use the preceding variables at the same time, or use them together with a valid string.
        * `rewrite_config` - The redirect action within ALB.
            * `host` - The host name of the destination to which requests are redirected within ALB.  Valid values:  The host name must be 3 to 128 characters in length, and can contain letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost domain label can contain only letters, asterisks (*) and question marks (?) and cannot contain digits or hyphens (-). Other domain labels cannot start or end with a hyphen (-). You can include asterisks (*) and question marks (?) anywhere in a domain label. Default value: ${host}. You cannot use this value with other characters at the same time.
            * `path` - The path to which requests are to be redirected within ALB.  Valid values: The path must be 1 to 128 characters in length, and start with a forward slash (/). The path can contain letters, digits, asterisks (*), question marks (?)and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ”. The path is case-sensitive.  Default value: ${path}. This value can be used only once. You can use it with a valid string.
            * `query` - The query string of the request to be redirected within ALB.  The query string must be 1 to 128 characters in length, can contain letters and printable characters. It cannot contain the following special characters: # [ ] { } \ | < > &.  Default value: ${query}. This value can be used only once. You can use it with a valid string.
        * `traffic_limit_config` - The Flow speed limit.
            * `qps` - The Number of requests per second.
        * `traffic_mirror_config` - The Traffic mirroring.
            * `target_type` - The Mirror target type.
            * `mirror_group_config` - The Traffic is mirrored to the server group.
                * `server_group_tuples` - The destination server group to which requests are forwarded.
                    * `server_group_id` - The ID of the destination server group to which requests are forwarded.
    * `rule_conditions` - The conditions of the forwarding rule.
        * `type` - The type of the forwarding rule.
        * `query_string_config` - The configuration of the query string.
            * `values` - The query string.
                * `key` - The key must be 1 to 100 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The key cannot contain spaces or the following special characters: # [ ] { } \ | < > &.
                * `value` - The value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters, asterisks (*), and question marks (?). The value cannot contain spaces or the following special characters: # [ ] { } \ | < > &.
        * `cookie_config` - The configuration of the cookie.
            * `values` - The configuration of the cookie.
        * `header_config` - The configuration of the header field.
            * `key` - The key of the header field. The key must be 1 to 40 characters in length, and can contain letters, digits, hyphens (-) and underscores (_). The key does not support Cookie or Host.
            * `values` - The value of the header field. The value must be 1 to 128 characters in length, and can contain lowercase letters, printable characters whose ASCII value is ch >= 32 && ch < 127, and wildcards such as asterisks (*) and question marks (?). The value cannot start or end with a space.
        * `host_config` - The configuration of the host.
            * `values` - The name of the host. **Note: ** The host name must meet the following rules: The hostname must be 3 to 128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?). The host name must contain at least one period (.), and cannot start or end with a period (.). The rightmost field can contain only letters and wildcards, and cannot contain digits or hyphens (-). Other fields cannot start or end with a hyphen (-). You can enter asterisks (*) and question marks (?) anywhere in a field.
        * `method_config` - The configuration of the request method.
            * `values` - The request method. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, and `DELETE`.
        * `path_config` - The configuration of the path for the request to be forwarded.
            * `values` - The path of the request to be forwarded. The path must be 1 to 128 characters in length and must start with a forward slash (/). The path can contain letters, digits, and the following special characters: $ - _ . + / & ~ @ :. It cannot contain the following special characters: " % # ; ! ( ) [ ] ^ , ". The value is case-sensitive, and can contain asterisks (*) and question marks (?).
        * `source_ip_config` - The Based on source IP traffic matching.
            * `values` - Add one or more IP addresses or IP address segments.
    * `rule_id` - The first ID of the resource.
    * `rule_name` - The name of the forwarding rule. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
    * `status` - The status of the resource.
