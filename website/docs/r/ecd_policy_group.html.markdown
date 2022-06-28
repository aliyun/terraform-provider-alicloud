---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_policy_group"
sidebar_current: "docs-alicloud-resource-ecd-policy-group"
description: |-
  Provides a Alicloud Elastic Desktop Service(EDS) Policy Group resource.
---

# alicloud\_ecd\_policy\_group

Provides a Elastic Desktop Service(EDS) Policy Group resource.

For information about Elastic Desktop Service(EDS) Policy Group and how to use it, see [What is Policy Group](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "my-policy-group"
  clipboard         = "read"
  local_drive       = "read"
  usb_redirect      = "off"
  watermark         = "off"

  authorize_access_policy_rules {
    description = "my-description1"
    cidr_ip     = "1.2.3.45/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "my-description"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "1.2.3.4/24"
  }
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
* `usb_redirect` - (Optional, Computed) The usb redirect policy. Valid values: `off`, `on`.
* `visual_quality` - (Optional, Computed) The quality of visual. Valid values: `high`, `lossless`, `low`, `medium`.
* `watermark` - (Optional, Computed) The watermark policy. Valid values: `off`, `on`.
* `watermark_transparency` - (Optional, Computed) The watermark transparency. Valid values: `DARK`, `LIGHT`, `MIDDLE`.
* `watermark_type` - (Optional) The type of watemark. Valid values: `EndUserId`, `HostName`.
* `recording` - (Optional, Computed, Available in 1.171.0+) Whether to enable screen recording. Valid values: `off`, `alltime`, `period`.
* `recording_start_time` - (Optional, Available in 1.171.0+) The start time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_end_time` - (Optional, Available in 1.171.0+) The end time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_fps` - (Optional, Computed, Available in 1.171.0+) The fps of recording. Valid values: `2`, `5`, `10`, `15`.
* `camera_redirect` - (Optional, Computed, Available in 1.171.0+) Whether to enable local camera redirection. Valid values: `on`, `off`.

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

* `cidr_ip` - (Optional) The cidrip of authorize access rule.
* `description` - (Optional) The description of authorize access rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Group.
* `status` - The status of policy.

## Import

Elastic Desktop Service(EDS) Policy Group can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_policy_group.example <id>
```
