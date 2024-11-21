---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_policy_group"
sidebar_current: "docs-alicloud-resource-ecd-policy-group"
description: |-
  Provides a Alicloud Elastic Desktop Service (ECD) Policy Group resource.
---

# alicloud_ecd_policy_group

Provides a Elastic Desktop Service (ECD) Policy Group resource.

For information about Elastic Desktop Service (ECD) Policy Group and how to use it, see [What is Policy Group](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createpolicygroup).

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_policy_group&exampleId=8b6500b3-4da7-1d11-22dc-ef58e8a22b9db4492c0a&activeTab=example&spm=docs.r.ecd_policy_group.0.8b6500b34d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "terraform-example"
  clipboard         = "read"
  local_drive       = "read"
  usb_redirect      = "off"
  watermark         = "off"

  authorize_access_policy_rules {
    description = "terraform-example"
    cidr_ip     = "1.2.3.45/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "terraform-example"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "1.2.3.4/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `authorize_access_policy_rules` - (Optional) The rule of authorize access rule. See [`authorize_access_policy_rules`](#authorize_access_policy_rules) below.
* `authorize_security_policy_rules` - (Optional) The policy rule. See [`authorize_security_policy_rules`](#authorize_security_policy_rules) below.
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
* `recording` - (Optional, Computed, Available in 1.171.0+) Whether to enable screen recording. Valid values: `off`, `all-time`, `period`.
* `recording_start_time` - (Optional, Available in 1.171.0+) The start time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_end_time` - (Optional, Available in 1.171.0+) The end time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_fps` - (Optional, Computed, Available in 1.171.0+) The fps of recording. Valid values: `2`, `5`, `10`, `15`.
* `camera_redirect` - (Optional, Computed, Available in 1.171.0+) Whether to enable local camera redirection. Valid values: `on`, `off`.
* `recording_expires` - (Optional, Available in 1.186.0+) The screen recording video retention. Valid values between 30 and 180. This return value is meaningful only when the value of `recording` is `period` or `all-time`.

### `authorize_security_policy_rules`

The authorize_security_policy_rules supports the following: 

* `cidr_ip` - (Optional) The cidrip of security rules.
* `description` - (Optional) The description of security rules.
* `ip_protocol` - (Optional) The ip protocol of security rules.
* `policy` - (Optional) The policy of security rules.
* `port_range` - (Optional) The port range of security rules.
* `priority` - (Optional) The priority of security rules.
* `type` - (Optional) The type of security rules.

### `authorize_access_policy_rules`

The authorize_access_policy_rules supports the following: 

* `cidr_ip` - (Optional) The cidrip of authorize access rule.
* `description` - (Optional) The description of authorize access rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Group.
* `status` - The status of policy.

## Import

Elastic Desktop Service (ECD) Policy Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_policy_group.example <id>
```
