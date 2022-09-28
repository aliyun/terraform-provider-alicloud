---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_sls_groups"
sidebar_current: "docs-alicloud-datasource-cms-sls-groups"
description: |-
  Provides a list of Cms Sls Groups to the user.
---

# alicloud\_cms\_sls\_groups

This data source provides the Cms Sls Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.171.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_sls_groups" "ids" {
  ids = ["example_id"]
}
output "cms_sls_group_id_1" {
  value = data.alicloud_cms_sls_groups.ids.groups.0.id
}

data "alicloud_cms_sls_groups" "nameRegex" {
  name_regex = "^my-SlsGroup"
}
output "cms_sls_group_id_2" {
  value = data.alicloud_cms_sls_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Sls Group IDs. Its element value is same as Sls Group Name.
* `keyword` - (Optional, ForceNew)  The keywords of the `sls_group_name` or `sls_group_description` of the Sls Group.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Sls Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Sls Group names.
* `groups` - A list of Cms Sls Groups. Each element contains the following attributes:
	* `create_time` - The creation time of the resource.
	* `id` - The ID of the Sls Group. Its value is same as Queue Name.
	* `sls_group_config` - The Config of the Sls Group.
		* `sls_logstore` - The name of the Log Store.
		* `sls_project` - The name of the Project.
		* `sls_region` - The Sls Region.
		* `sls_user_id` - The ID of the Sls User.
	* `sls_group_description` - The Description of the Sls Group.
	* `sls_group_name` - The name of the resource.