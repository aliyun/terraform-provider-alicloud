---
subcategory: "ECD"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop_group"
sidebar_current: "docs-alicloud-resource-ecd-desktop-group"
description: |-
  Provides a Alicloud ECD Desktop Group resource.
---

# alicloud\_ecd\_desktop\_group

Provides a ECD Desktop Group resource.

For information about ECD Desktop Group and how to use it, see [What is Desktop Group](https://help.aliyun.com/).

-> **NOTE:** Available in v1.170.0+.

## Example Usage

Basic Usage

```terraform

## Argument Reference

The following arguments are supported:

* `allow_auto_setup` - (Optional) Whether Or Not to Allow Automatic Creating Desktop: 0 Does Not Allow 1 Allows.
* `allow_buffer_count` - (Optional) Allow You to Leave Your Desktop of the Buffer Number 0-Don't Keep N Are Allowed to Remain in the N.
* `auto_pay` - (Optional) The auto pay.
* `bind_amount` - (Optional) The bind amount.
* `bundle_id` - (Required) The Template ID.
* `charge_type` - (Optional) The charge type.
* `comments` - (Optional) Note.
* `default_init_desktop_count` - (Optional) The default init desktop count.
* `desktop_group_name` - (Optional) Desktop Group Name.
* `directory_id` - (Optional, ForceNew) The Directory ID.
* `end_user_ids` - (Required, ForceNew) To Authorize the Use of the Cloud Desktop Group of User ID.
* `image_id` - (Optional) The image id.
* `keep_duration` - (Optional) The User Connection to the Original Desktop Expiration Time (MS).
* `load_policy` - (Optional) The load policy.
* `max_desktops_count` - (Optional) Desktop Groups Added to a Maximum Desktop, the Default Maximum Number of Children's Cots/100 Desktop.
* `min_desktops_count` - (Optional) Desktop Groups That Must Be Maintained as the Minimum Desktop Number Default Minimum 1 Desktop.
* `office_site_id` - (Required, ForceNew) The Workspace ID.
* `own_type` - (Optional) The own type.
* `period` - (Optional) The period.
* `policy_group_id` - (Optional) Policy Group ID.
* `reset_type` - (Optional) The reset type.
* `scale_strategy_id` - (Optional) The Scaling Policy Group ID.
* `vpc_id` - (Optional) The vpc id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Desktop Group.
* `status` - The status.

## Import

ECD Desktop Group can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_desktop_group.example <id>
```