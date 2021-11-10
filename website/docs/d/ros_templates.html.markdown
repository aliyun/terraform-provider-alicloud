---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_templates"
sidebar_current: "docs-alicloud-datasource-ros-templates"
description: |-
  Provides a list of Ros Templates to the user.
---

# alicloud\_ros\_templates

This data source provides the Ros Templates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.108.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_templates" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_template_id" {
  value = data.alicloud_ros_templates.example.templates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Template IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Template name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `share_type` - (Optional, ForceNew) Share Type. Valid Values: `Private`, `Shared`
* `template_name` - (Optional, ForceNew) The name of the template.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `tags` - (Optional) Query the resource bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Template names.
* `templates` - A list of Ros Templates. Each element contains the following attributes:
	* `change_set_id` - The ID of the change set.
	* `description` - The description of the template. The description can be up to 256 characters in length.
	* `id` - The ID of the Template.
	* `share_type` - Share Type.
	* `stack_group_name` - The name of the stack group. The name must be unique in a region.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
	* `stack_id` - The ID of the stack.
	* `tags` - Tags.
		* `tag_key` - The key of tag N of the resource.
		* `tag_value` - The value of tag N of the resource.
	* `template_body` - The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You must specify one of the TemplateBody and TemplateURL parameters, but you cannot specify both of them.
	* `template_id` - The ID of the template.
	* `template_name` - The name of the template.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
	* `template_version` - Template Version.
