---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project_policy"
sidebar_current: "docs-alicloud-resource-log-project_policy"
description: |-
  Provides a Alicloud Log Policy resource.
---

# alicloud\_log\_project\_policy
The log project policy is used to set a policy for a project in log service. The policy struce is the same as RAM policy.
[Refer to details](https://www.alibabacloud.com/help/en/resource-access-management/latest/policy-overview).

-> **NOTE:** Available in 1.192.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_log_project" "default" {
  name        = "tf-project"
  description = "tf unit test"
}

resource "alicloud_log_project_policy" "example" {
  project = "tf-project"
  policy  = <<EOF
  {
      "Statement": [{
            "Action": [
                "log:Post*"
            ],
            "Effect": "Allow",
            "Resource": "acs:log:*:*:project/test-project-policy-1/*"
      }],
      "Version": "1"
  }
EOF
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of the log project. It is the only in one Alicloud account.
* `policy` - (Required) The policy detail of project, empty string means no policy.

## Attributes Reference

The following attributes are exported:

* `id` - The Name of the Log Project. It sames as its name.

## Import

Log project policy can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_log_project_policy.example tf-project
```
