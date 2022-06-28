---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_flow"
sidebar_current: "docs-alicloud-resource-fnf-flow"
description: |-
  Provides a Alicloud Serverless Workflow Flow resource.
---

# alicloud\_fnf\_flow

Provides a Serverless Workflow Flow resource.

For information about Serverless Workflow Flow and how to use it, see [What is Flow](https://www.alibabacloud.com/help/en/doc-detail/123079.htm).

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_role" "default" {
  name     = "tf-testacc-fnfflow"
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "fnf.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_fnf_flow" "example" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: pass
      name: helloworld
  EOF
  role_arn    = alicloud_ram_role.default.arn
  description = "Test for terraform fnf_flow."
  name        = "tf-testacc-flow"
  type        = "FDL"
}
```

## Argument Reference

The following arguments are supported:

* `definition` - (Required) The definition of the flow. It must comply with the Flow Definition Language (FDL) syntax.
* `description` - (Required) The description of the flow.
* `name` - (Required, ForceNew) The name of the flow. The name must be unique in an Alibaba Cloud account.
* `role_arn` - (Optional) The ARN of the specified RAM role that Serverless Workflow uses to assume the role when Serverless Workflow executes a flow.
* `type` - (Required) The type of the flow. Valid values are `FDL` or `DEFAULT`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Flow. The value same as `name`.
* `flow_id` - The unique ID of the flow.
* `last_modified_time` - The time when the flow was last modified.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Flow.
* `delete` - (Defaults to 1 mins) Used when delete the Flow.
* `update` - (Defaults to 1 mins) Used when update the Flow.

## Import

Serverless Workflow Flow can be imported using the id, e.g.

```
$ terraform import alicloud_fnf_flow.example <name>
```
