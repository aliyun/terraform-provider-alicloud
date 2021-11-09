---
subcategory: "Intelligent Media Management (IMM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_imm_project"
sidebar_current: "docs-alicloud-resource-imm-project"
description: |-
  Provides a Alicloud Intelligent Media Management Project resource.
---

# alicloud\_imm\_project

Provides a Intelligent Media Management Project resource.

For information about Intelligent Media Management Project and how to use it, see [What is Project](https://help.aliyun.com/document_detail/63496.html).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_role" "role" {
  name        = "example_role_name"
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "imm.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a role test."
  force       = true
}
resource "alicloud_imm_project" "example" {
  project      = "example_name"
  service_role = alicloud_ram_role.role.name
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The name of Project.
* `service_role` - (Optional) The service role authorized to the Intelligent Media Management service to access other cloud resources. Default value: `AliyunIMMDefaultRole`. You can also create authorization  roles through the `alicloud_ram_role`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Project. Its value is same as `project`.

## Import

Intelligent Media Management Project can be imported using the id, e.g.

```
$ terraform import alicloud_imm_project.example <project>
```
