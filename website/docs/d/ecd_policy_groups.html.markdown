---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_policy_groups"
sidebar_current: "docs-alicloud-datasource-ecd-policy-groups"
description: |-
  Provides a list of Ecd Policy Groups to the user.
---

# alicloud\_ecd\_policy\_groups

This data source provides the Ecd Policy Groups of the current Alibaba Cloud user.

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

data "alicloud_ecd_policy_groups" "nameRegex" {
  name_regex = "^my-policy"
}
output "ecd_policy_group_id" {
  value = data.alicloud_ecd_policy_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Policy Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Policy Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of policy.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Policy Group names.
* `groups` - A list of Ecd Policy Groups. Each element contains the following attributes:
	* `authorize_access_policy_rules` - The rule of authorize access rule.
		* `cidr_ip` - The cidrip of authorize access rule..
		* `description` - The description of authorize access rule.
	* `authorize_security_policy_rules` - The policy rule.
		* `cidr_ip` - The cidrip of security rules.
		* `description` - The description of security rules.
		* `ip_protocol` - The ip protocol of security rules.
		* `policy` - The policy of security rules.
		* `port_range` - The port range of security rules.
		* `priority` - The priority of security rules.
		* `type` - The type of security rules.
	* `clipboard` - The clipboard policy.
	* `domain_list` - The list of domain.
	* `eds_count` - The count of eds.
	* `html_access` - The access of html5.
	* `html_file_transfer` - The html5 file transfer.
	* `id` - The ID of the Policy Group.
	* `local_drive` - Local drive redirect policy.
	* `policy_group_id` - The policy group id.
	* `policy_group_name` - The name of policy group.
	* `policy_group_type` - The type of policy group.
	* `preempt_login` - The preempt login.
	* `preempt_login_users` - A list of preempt log users.
	* `status` - The status of policy.
	* `usb_redirect` - The usb redirect policy.
	* `visual_quality` - The quality of visual.sae_ecdsae_nameecd_po
	* `watermark` - The watermark policy.
	* `watermark_custom_text` - The custort text of water mark.
	* `watermark_transparency` - The watermark transparency.
	* `watermark_type` - The type of watemark.
	* `recording` - (Available in 1.171.0+) Whether to enable screen recording. Valid values: `off`, `alltime`, `period`.
	* `recording_start_time` - (Available in 1.171.0+) The start time of recording.
	* `recording_end_time` - (Available in 1.171.0+) The end time of recording.
	* `recording_fps` - (Available in 1.171.0+) The fps of recording. Valid values: `2`, `5`, `10`, `15`.
	* `camera_redirect` - (Available in 1.171.0+) Whether to enable local camera redirection. Valid values: `on`, `off`.
