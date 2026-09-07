---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_advanced_query_template"
description: |-
  Provides a Alicloud Actiontrail Advanced Query Template resource.
---

# alicloud_actiontrail_advanced_query_template

Provides a Actiontrail Advanced Query Template resource.

sql template of advanced query.

For information about Actiontrail Advanced Query Template and how to use it, see [What is Advanced Query Template](https://next.api.alibabacloud.com/document/Actiontrail/2020-07-06/CreateAdvancedQueryTemplate).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_actiontrail_advanced_query_template&exampleId=97c64ca0-c3cd-08cb-3f15-894240df72a9e4cbe3c0&activeTab=example&spm=docs.r.actiontrail_advanced_query_template.0.97c64ca0c3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_actiontrail_advanced_query_template" "default" {
  simple_query  = true
  template_name = "exampleTemplateName"
  template_sql  = "*"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_actiontrail_advanced_query_template&spm=docs.r.actiontrail_advanced_query_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `simple_query` - (Required) Distinguish whether the current template is a simple query
* `template_name` - (Optional) The name of the resource
* `template_sql` - (Required) SQL content saved on behalf of the current template

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Advanced Query Template.
* `delete` - (Defaults to 5 mins) Used when delete the Advanced Query Template.
* `update` - (Defaults to 5 mins) Used when update the Advanced Query Template.

## Import

Actiontrail Advanced Query Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail_advanced_query_template.example <id>
```