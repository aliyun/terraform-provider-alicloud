---
subcategory: "Serverless Workflow (FnF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_flow"
sidebar_current: "docs-alicloud-resource-fnf-flow"
description: |-
  Provides a Alicloud Serverless Workflow Flow resource.
---

# alicloud_fnf_flow

Provides a Serverless Workflow Flow resource.

For information about Serverless Workflow Flow and how to use it, see [What is Flow](https://www.alibabacloud.com/help/en/doc-detail/123079.htm).

-> **NOTE:** Available since v1.105.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fnf_flow&exampleId=62806fad-3df1-0cba-b739-afb816fac0baf7d6992b&activeTab=example&spm=docs.r.fnf_flow.0.62806fad3d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_ram_role" "default" {
  name     = "tf-example-fnfflow"
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
  name        = "tf-example-flow"
  type        = "FDL"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fnf_flow&spm=docs.r.fnf_flow.example&intl_lang=EN_US)

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Flow.
* `delete` - (Defaults to 1 mins) Used when delete the Flow.
* `update` - (Defaults to 1 mins) Used when update the Flow.

## Import

Serverless Workflow Flow can be imported using the id, e.g.

```shell
$ terraform import alicloud_fnf_flow.example <name>
```
