---
subcategory: "Apsara Agile Live (IMP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_imp_app_templates"
sidebar_current: "docs-alicloud-datasource-imp-app-templates"
description: |-
  Provides a list of Imp App Templates to the user.
---

# alicloud\_imp\_app\_templates

This data source provides the Imp App Templates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_imp_app_templates" "ids" {}
output "imp_app_template_id_1" {
  value = data.alicloud_imp_app_templates.ids.templates.0.id
}

data "alicloud_imp_app_templates" "nameRegex" {
  name_regex = "^my_AppTemplate"
}
output "imp_app_template_id_2" {
  value = data.alicloud_imp_app_templates.nameRegex.templates.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of App Template IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by App Template name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Application template usage status. Valid values: ["attached", "unattached"].

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of App Template names.
* `templates` - A list of Imp App Templates. Each element contains the following attributes:
	* `app_template_creator` - Apply template creator.
	* `app_template_id` - The first ID of the resource.
	* `app_template_name` - The name of the resource.
	* `component_list` - List of components.
	* `config_list` - List of config.
	  * `key` - Config key.
	  * `value` - Config Value.
		
	* `create_time` - Creation time.
	* `id` - The ID of the App Template.
	* `integration_mode` - Integration mode (Integrated SDK:paasSDK, Model Room: standardRoom).
	* `scene` - Application Template scenario, e-commerce business, classroom classroom.
	* `sdk_info` - SDK information.
	* `standard_room_info` - Model room information.
	* `status` - Application template usage status.
