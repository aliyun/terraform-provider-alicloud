---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_policy"
description: |-
  Provides a Alicloud TAG Policy resource.
---

# alicloud_tag_policy

Provides a TAG Policy resource.



For information about TAG Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/en/resource-management/latest/create-policy).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_tag_policy&exampleId=172f7f49-7f62-77be-d057-24d3850af041d72088e7&activeTab=example&spm=docs.r.tag_policy.0.172f7f497f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_tag_policy" "example" {
  policy_name    = var.name
  policy_desc    = var.name
  user_type      = "USER"
  policy_content = <<EOF
		{"tags":{"CostCenter":{"tag_value":{"@@assign":["Beijing","Shanghai"]},"tag_key":{"@@assign":"CostCenter"}}}}
    EOF
}
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required) The name of the tag policy. The name must be 1 to 128 characters in length and can contain letters, digits, and underscores (_).
* `policy_content` - (Required) The document of the tag policy.
* `policy_desc` - (Optional) The description of the policy. The description must be 1 to 512 characters in length.
* `user_type` - (Optional, ForceNew)The mode of the Tag Policy feature. Valid values: `USER`, `RD`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Policy.

## Timeouts

-> **NOTE:** Available since v1.243.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

TAG Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_policy.example <id>
```
