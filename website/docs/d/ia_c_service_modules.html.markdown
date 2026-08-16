---
subcategory: "IaC Service (IaCService)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ia_c_service_modules"
sidebar_current: "docs-alicloud-datasource-ia-c-service-modules"
description: |-
  Provides a list of IaC Service Module owned by an Alibaba Cloud account.
---

# alicloud_ia_c_service_modules

This data source provides IaC Service Module available to the user.[What is Module](https://next.api.alibabacloud.com/document/IaCService/2021-08-06/CreateModule)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
data "alicloud_ia_c_service_modules" "default" {

}

output "first_module_id" {
  value = data.alicloud_ia_c_service_modules.default.modules[0].module_id
}
```

## Argument Reference

The following arguments are supported:
* `group_id` - (Optional, ForceNew) The ID of the group to which the modules belong.
* `module_name` - (Optional, ForceNew) The keyword of the module name.
* `page_number` - (Optional, ForceNew) The page number. Default value: `1`.
* `page_size` - (Optional, ForceNew) The number of entries per page. Default value: `20`. Maximum value: `100`.
* `project_id` - (Optional, ForceNew) The ID of the project to which the modules belong.
* `ids` - (Optional, Computed, List) A list of Module IDs.
* `name_regex` - (Optional) A regex string to filter results by the module name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Module Entries.
* `modules` - A list of Module Entries. Each element contains the following attributes:
  * `module_id` - The ID of the module.
  * `create_time` - The time when the module was created.
  * `description` - The description of the module.
  * `latest_version` - The latest version number of the module.
  * `module_name` - The name of the module.
  * `source` - The source of the module.
  * `status` - The status of the module.
  * `group_info` - The group information of the module. Each element contains the following attributes:
    * `group_id` - The ID of the group to which the module belongs.
    * `group_name` - The name of the group to which the module belongs.
    * `project_id` - The ID of the project to which the module belongs.
    * `project_name` - The name of the project to which the module belongs.
  * `tags` - The tags of the module. Each element contains the following attributes:
    * `tag_key` - The key of the tag.
    * `tag_value` - The value of the tag.
  * `id` - The ID of the module.
