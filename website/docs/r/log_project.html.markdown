---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project"
description: |-
  Provides a Alicloud SLS Project resource.
---

# alicloud_log_project

Provides a SLS Project resource. 

For information about SLS Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/sls/developer-reference/api-createproject).

-> **NOTE:** Available since v1.9.5.

## Example Usage

Basic Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_project&exampleId=d02177bc-72d8-c195-ee67-bd69bfce5817db868e4f&activeTab=example&spm=docs.r.log_project.0.d02177bc72&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

Project With Policy Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_project&exampleId=5fd3dbc6-77b2-82a8-4d99-137bffbadcb2c329105b&activeTab=example&spm=docs.r.log_project.1.5fd3dbc677&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example_policy" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
  policy       = <<EOF
{
  "Statement": [
    {
      "Action": [
        "log:PostLogStoreLogs"
      ],
      "Condition": {
        "StringNotLike": {
          "acs:SourceVpc": [
            "vpc-*"
          ]
        }
      },
      "Effect": "Deny",
      "Resource": "acs:log:*:*:project/tf-log/*"
    }
  ],
  "Version": "1"
}
EOF
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_log_project&spm=docs.r.log_project.example&intl_lang=EN_US)

## Module Support

You can use the existing [sls module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls/alicloud) 
to create SLS project, store and store index one-click, like ECS instances.

## Argument Reference

The following arguments are supported:
* `policy` - (Optional, Available since v1.197.0) Log project policy, used to set a policy for a project.
* `description` - (Optional) Description.
* `project_name` - (Optional, ForceNew, Available since v1.212.0) The name of the log project. It is the only in one Alicloud account. The project name is globally unique in Alibaba Cloud and cannot be modified after it is created. The naming rules are as follows:
  - The project name must be globally unique. 
  - The name can contain only lowercase letters, digits, and hyphens (-). 
  - It must start and end with a lowercase letter or number. 
  - The value contains 3 to 63 characters.
* `resource_group_id` - (Optional, Computed, Available since v1.212.0) The ID of the resource group.
* `tags` - (Optional, Map) Tag.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.223.0). Field 'name' has been deprecated from provider version 1.223.0. New field 'project_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreateTime.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.

## Import

SLS Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_project.example <id>
```