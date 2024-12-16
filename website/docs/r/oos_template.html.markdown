---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_template"
sidebar_current: "docs-alicloud-resource-oos-template"
description: |-
  Provides a OOS Template resource.
---

# alicloud_oos_template

Provides a OOS Template resource. For information about Alicloud OOS Template and how to use it, see [What is Resource Alicloud OOS Template](https://www.alibabacloud.com/help/doc-detail/120761.htm).

-> **NOTE:** Available since v1.92.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_template&exampleId=10923d26-3220-f2f7-2f86-fa45c48d15e00f0e3c2e&activeTab=example&spm=docs.r.oos_template.0.10923d2632&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oos_template" "example" {
  content       = <<EOF
  {
    "FormatVersion": "OOS-2019-06-01",
    "Description": "Update Describe instances of given status",
    "Parameters":{
      "Status":{
        "Type": "String",
        "Description": "(Required) The status of the Ecs instance."
      }
    },
    "Tasks": [
      {
        "Properties" :{
          "Parameters":{
            "Status": "{{ Status }}"
          },
          "API": "DescribeInstances",
          "Service": "Ecs"
        },
        "Name": "foo",
        "Action": "ACS::ExecuteApi"
      }]
  }
  EOF
  template_name = "tf-example-name-${random_integer.default.result}"
  version_name  = "example"
  tags = {
    "Created" = "TF",
    "For"     = "acceptance Test"
  }
}

```

## Argument Reference

The following arguments are supported:

* `content` - (Required) The content of the template. The template must be in the JSON or YAML format. Maximum size: 64 KB. 
* `auto_delete_executions` - (Optional) When deleting a template, whether to delete its related executions. Default to `false`.
* `template_name` - (Required, ForceNew) The name of the template. The template name can be up to 200 characters in length. The name can contain letters, digits, hyphens (-), and underscores (_). It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, or `ALICLOUD`.
* `version_name` - (Optional) The name of template version.
* `resource_group_id` (Optional, Computed, Available in 1.177.0+) The ID of resource group which the template belongs.  
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the resource. It same with `template_name`.
* `created_by` - The creator of the template.
* `created_date` - The time when the template is created.
* `description` - The description of the template.
* `has_trigger` - Is it triggered successfully.
* `share_type` - The sharing type of the template. The sharing type of templates created by users are set to Private. The sharing type of common templates provided by OOS are set to Public.
* `template_format` - The format of the template. The format can be JSON or YAML. The system automatically identifies the format.
* `template_id` - The id of OOS Template.
* `template_type` - The type of OOS Template. `Automation` means the implementation of Alibaba Cloud API template, `Package` means represents a template for installing software.
* `template_version` - The version of OOS Template.
* `updated_by` - The user who updated the template.
* `updated_date` - The time when the template was updated.

## Import

OOS Template can be imported using the id or template_name, e.g.

```shell
$ terraform import alicloud_oos_template.example template_name
```
