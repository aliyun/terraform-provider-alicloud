---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_application_groups"
sidebar_current: "docs-alicloud-datasource-oos-application-groups"
description: |-
  Provides a list of Oos Application Groups to the user.
---

# alicloud\_oos\_application\_groups

This data source provides the Oos Application Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_application_groups" "ids" {
  application_name = "example_value"
  ids              = ["my-ApplicationGroup-1", "my-ApplicationGroup-2"]
}
output "oos_application_group_id_1" {
  value = data.alicloud_oos_application_groups.ids.groups.0.id
}

data "alicloud_oos_application_groups" "nameRegex" {
  application_name = "example_value"
  name_regex       = "^my-ApplicationGroup"
}
output "oos_application_group_id_2" {
  value = data.alicloud_oos_application_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, ForceNew) The name of the Application.
* `deploy_region_id` - (Optional, ForceNew) The region ID of the deployment.
* `ids` - (Optional, ForceNew, Computed)  A list of Application Group IDs. Its element value is same as Application Group Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Application Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Application Group names.
* `groups` - A list of Oos Application Groups. Each element contains the following attributes:
  * `application_group_name` - The name of the Application group.
  * `application_name` - The name of the Application.
  * `cms_group_id` - The ID of the cloud monitor group.
  * `create_time` - The Creation time of the resource.
  * `update_time` - The Update time of the resource.
  * `deploy_region_id` - The region ID of the deployment.
  * `description` - Application group description information.
  * `id` - The ID of the Application Group. Its value is same as Queue Name.
  * `import_tag_key` - Label key.
  * `import_tag_value` - Label value.
