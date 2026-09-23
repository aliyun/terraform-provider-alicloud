---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_templates"
sidebar_current: "docs-alicloud-datasource-oos-templates"
description: |-
    Provides a list of OOS Templates.
---

# alicloud\_oos\_templates

This data source provides a list of OOS Templates in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available in v1.92.0+.

## Example Usage

```
# Declare the data source

data "alicloud_oos_templates" "example" {
  name_regex = "test"
  tags={
    "Created" = "TF",
    "For" = "template Test"
  }
  share_type = "Private"
  has_trigger = false
}


output "first_template_name" {
  value = "${data.alicloud_oos_templates.example.templates.0.template_name}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of OOS Template ids. Each element in the list is same as template_name.
* `name_regex` - (Optional) A regex string to filter the results by the template_name.
* `category` - (Optional) The category of template.
* `created_by` - (Optional) The creator of the template.
* `created_date` - (Optional) The template whose creation time is less than or equal to the specified time. The format is: YYYY-MM-DDThh:mm::ssZ.
* `has_trigger` - (Optional) Is it triggered successfully.
* `share_type` - (Optional) The sharing type of the template. Valid values: `Private`, `Public`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `template_format` - (Optional) The format of the template. Valid values: `JSON`, `YAML`.
* `template_type` - (Optional) The type of OOS Template.
* `sort_field` - (Optional) Sort field. Valid values: `TotalExecutionCount`, `Popularity`, `TemplateName` and `CreatedDate`. Default to `TotalExecutionCount`.
* `sort_order` - (Optional) Sort order. Valid values: `Ascending`, `Descending`. Default to `Descending`
* `created_date_after` - (Optional) Create a template whose time is greater than or equal to the specified time. The format is: YYYY-MM-DDThh:mm:ssZ.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of OOS Template ids. Each element in the list is same as template_name.
* `names` -  (Available in v1.114.0+) A list of OOS Template names.
* `templates` - A list of OOS Templates. Each element contains the following attributes:
  * `id` - ID of the OOS Template. The value is same as template_name.
  * `template_name` - Name of the OOS Template.
  * `description` - Description of the OOS Template.
  * `template_id` - ID of the OOS Template resource.
  * `template_version` - Version of the OOS Template.
  * `updated_by` - The user who updated the template.
  * `updated_date` - The time when the template was updated.

