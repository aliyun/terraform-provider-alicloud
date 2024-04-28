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

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  name        = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

Project With Policy Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example_policy" {
  name        = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
  policy      = <<EOF
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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.

## Import

SLS Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_project.example <id>
```