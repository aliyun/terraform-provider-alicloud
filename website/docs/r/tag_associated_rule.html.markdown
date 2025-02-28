---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_associated_rule"
description: |-
  Provides a Alicloud TAG Associated Rule resource.
---

# alicloud_tag_associated_rule

Provides a TAG Associated Rule resource.



For information about TAG Associated Rule and how to use it, see [What is Associated Rule](https://www.alibabacloud.com/help/en/resource-management/tag/developer-reference/api-tag-2018-08-28-createassociatedresourcerules).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_tag_associated_rule&exampleId=cdc68cc6-56b1-7e90-7852-e80d16ec1d6026e17dfa&activeTab=example&spm=docs.r.tag_associated_rule.0.cdc68cc656&intl_lang=EN_US" target="_blank">
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


resource "alicloud_tag_associated_rule" "default" {
  status                  = "Enable"
  associated_setting_name = "rule:AttachEni-DetachEni-TagInstance:Ecs-Instance:Ecs-Eni"
  tag_keys                = ["user"]
}
```

## Argument Reference

The following arguments are supported:
* `associated_setting_name` - (Required, ForceNew) The setting name of the associated resource tag rule. For specific values, see the Rule Setting Name column in [Resources that Support Associated Resource Tag Settings](https://www.alibabacloud.com/help/en/resource-management/tag/user-guide/associated-resource-label-settings)
* `status` - (Required) Whether to enable the associated resource tag rule. Valid values: `Enable`, `Disable`.
* `tag_keys` - (Optional, List) List of tag keys for the associated resource tag rule.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Associated Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Associated Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Associated Rule.
* `update` - (Defaults to 5 mins) Used when update the Associated Rule.

## Import

TAG Associated Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_associated_rule.example <id>
```
