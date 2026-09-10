---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_applications"
sidebar_current: "docs-alicloud-datasource-oos-applications"
description: |-
  Provides a list of Oos Applications to the user.
---

# alicloud\_oos\_applications

This data source provides the Oos Applications of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_applications" "ids" {}
output "oos_application_id_1" {
  value = data.alicloud_oos_applications.ids.applications.0.id
}

data "alicloud_oos_applications" "nameRegex" {
  name_regex = "^my-Application"
}
output "oos_application_id_2" {
  value = data.alicloud_oos_applications.nameRegex.applications.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Application IDs. Its element value is same as Application Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Application name.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of oss application names.
* `applications` - A list of Oos Applications. Each element contains the following attributes:
    * `application_name` - The name of the application.
    * `create_time` - The Created time of the application.
    * `update_time` - The Updated time of the application.
    * `description` - Application group description information.
    * `id` - The ID of the Application. The value is formate as <application_name>.
    * `resource_group_id` - The ID of the resource group.
    * `tags` - The tag of the resource.