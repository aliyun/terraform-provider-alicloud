---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project"
sidebar_current: "docs-alicloud-resource-log-project"
description: |-
  Provides a Alicloud log project resource.
---

# alicloud_log_project

The project is the resource management unit in Log Service and is used to isolate and control resources.
You can manage all the logs and the related log sources of an application by using projects. [Refer to details](https://www.alibabacloud.com/help/doc-detail/48873.htm).

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

* `name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `description` - (Optional) Description of the log project.
* `tags` - (Optional) Log project tags.
  - Key: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". 
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://".

* `policy` - (Optional, Available in 1.197.0+) Log project policy, used to set a policy for a project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It seems as its name.

## Timeouts

-> **NOTE:** Available in 1.126.0+

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the log project.

## Import

Log project can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_log_project.example tf-log
```
