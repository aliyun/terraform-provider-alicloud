---
subcategory: "ECD"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop_groups"
sidebar_current: "docs-alicloud-datasource-ecd-desktop-groups"
description: |-
  Provides a list of Ecd Desktop Groups to the user.
---

# alicloud\_ecd\_desktop\_groups

This data source provides the Ecd Desktop Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.170.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_desktop_groups" "ids" {}
output "ecd_desktop_group_id_1" {
  													value = data.alicloud_ecd_desktop_groups.ids.groups.0.id
												}
												
data "alicloud_ecd_desktop_groups" "nameRegex" {
name_regex = "^my-DesktopGroup"
}
output "ecd_desktop_group_id_2" {
  													value = data.alicloud_ecd_desktop_groups.nameRegex.groups.0.id
												}
												
```

## Argument Reference

The following arguments are supported:

* `desktop_group_name` - (Optional, ForceNew) Desktop Group Name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `excluded_end_user_ids` - (Optional, ForceNew) The excluded end user ids.
* `ids` - (Optional, ForceNew, Computed)  A list of Desktop Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Desktop Group name.
* `office_site_id` - (Optional, ForceNew) The Workspace ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `own_type` - (Optional, ForceNew) The own type.
* `period` - (Optional, ForceNew) The period.
* `period_unit` - (Optional, ForceNew) Subscription Billing Method as the ECS of the Time at Which a Unit. Value Range: Month: Month Year: Years Default Value: Month. Valid values: `Month`, `Year`.
* `status` - (Optional, ForceNew) The status.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Desktop Group names.
* `groups` - A list of Ecd Desktop Groups. Each element contains the following attributes:
	* `allow_auto_setup` - Whether Or Not to Allow Automatic Creating Desktop: 0 Does Not Allow 1 Allows.
	* `allow_buffer_count` - Allow You to Leave Your Desktop of the Buffer Number 0-Don't Keep N Are Allowed to Remain in the N.
	* `bundle_id` - The Template ID.
	* `comments` - Note.
	* `cpu` - The Number of CPUs.
	* `create_time` - Creation Time.
	* `creator` - The Creator.
	* `data_disk_category` - The Categories of Data Disks.
	* `data_disk_size` - The Data Disk Size.
	* `desktop_group_id` - Desktop Group ID.
	* `desktop_group_name` - Desktop Group Name.
	* `directory_id` - The Directory ID.
	* `directory_type` - Directory Type.
	* `end_user_count` - Authorized by the Total Number of Users.
	* `end_user_ids` - To Authorize the Use of the Cloud Desktop Group of User ID.
	* `expired_time` - The Expiration Time.
	* `gpu_count` - GPU Number.
	* `gpu_spec` - GPU Specifications.
	* `id` - The ID of the Desktop Group.
	* `keep_duration` - The User Connection to the Original Desktop Expiration Time (MS).
	* `max_desktops_count` - Desktop Groups Added to a Maximum Desktop, the Default Maximum Number of Children's Cots/100 Desktop.
	* `memory` - The Memory Size.
	* `min_desktops_count` - Desktop Groups That Must Be Maintained as the Minimum Desktop Number Default Minimum 1 Desktop.
	* `office_site_id` - The Workspace ID.
	* `office_site_name` - The Workspace Name.
	* `office_site_type` - Workspace Type.
	* `own_bundle_name` - Template Name.
	* `pay_type` - The Pay-as-You-Type, by Default, Are the Pre-Paid Desktop.
	* `policy_group_id` - Policy Group ID.
	* `policy_group_name` - Security Policy Group Name.
	* `res_type` - Resource Type.
	* `system_disk_category` - System Disk Type.
	* `system_disk_size` - System Disk Size.