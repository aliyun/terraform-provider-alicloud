---
subcategory: "Elastic Desktop Service(ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_policy_group"
sidebar_current: "docs-alicloud-resource-ecd-policy-group"
description: |-
  Provides a Alicloud Elastic Desktop Service(ECD) Policy Group resource.
---

# alicloud\_ecd\_policy\_group

Provides a Elastic Desktop Service(ECD) Policy Group resource.

For information about Elastic Desktop Service(ECD) Policy Group and how to use it, see [What is Policy Group](https://help.aliyun.com/).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_policy_group" "example" {
  policy_group_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `authorize_access_policy_rules` - (Optional) The rule of authorize access rule.
* `authorize_security_policy_rules` - (Optional) The policy rule.
* `clipboard` - (Optional, Computed) The clipboard policy. Valid values: `off`, `read`, `readwrite`.
* `domain_list` - (Optional) The list of domain.
* `html_access` - (Optional, Computed) The access of html5. Valid values: `off`, `on`.
* `html_file_transfer` - (Optional, Computed) The html5 file transfer. Valid values: `all`, `download`, `off`, `upload`.
* `local_drive` - (Optional, Computed) Local drive redirect policy. Valid values: ` readwrite`, `off`, `read`.
* `policy_group_name` - (Optional) The name of policy group.
* `preempt_login` - (Optional) The preempt login.
* `preempt_login_users` - (Optional) A list of preempt log users.
* `usb_redirect` - (Optional, Computed) The usb redirect policy. Valid values: `off`, `on`.
* `visual_quality` - (Optional, Computed) The quality of visual. Valid values: `high`, `lossless`, `low`, `medium`.
* `watermark` - (Optional, Computed) The watermark policy. Valid values: `off`, `on`.
* `watermark_custom_text` - (Optional) The custort text of water mark.
* `watermark_transparency` - (Optional, Computed) The watermark transparency. Valid values: `DARK`, `LIGHT`, `MIDDLE`.
* `watermark_type` - (Optional) The type of watemark. Valid values: `EndUserId`, `HostName`.

#### Block revoke_security_policy_rule

The revoke_security_policy_rule supports the following: 

* `cidr_ip` - (Optional) The cidr ip.
* `description` - (Optional) The description.
* `ip_protocol` - (Optional) The ip protocol.
* `policy` - (Optional) The policy.
* `port_range` - (Optional) The port range.
* `priority` - (Optional) The priority.
* `type` - (Optional) The type.

#### Block revoke_access_policy_rule

The revoke_access_policy_rule supports the following: 

* `cidr_ip` - (Optional) The cidr ip.
* `description` - (Optional) The description.

#### Block authorize_security_policy_rules

The authorize_security_policy_rules supports the following: 

* `cidr_ip` - (Optional) The cidrip of security rules.
* `description` - (Optional) The description of security rules.
* `ip_protocol` - (Optional) The ip protocol of security rules.
* `policy` - (Optional) The policy of security rules.
* `port_range` - (Optional) The port range of security rules.
* `priority` - (Optional) The priority of security rules.
* `type` - (Optional) The type of security rules.

#### Block authorize_access_policy_rules

The authorize_access_policy_rules supports the following: 

* `cidr_ip` - (Optional) The cidrip of authorize access rule..
* `description` - (Optional) The description of authorize access rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Group.
* `status` - The status of policy.

## Import

Elastic Desktop Service(ECD) Policy Group can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_policy_group.example <id>
```