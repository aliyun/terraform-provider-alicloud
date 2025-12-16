---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_target"
description: |-
  Provides a Alicloud Resource Manager Shared Target resource.
---

# alicloud_resource_manager_shared_target

Provides a Resource Manager Shared Target resource.



For information about Resource Manager Shared Target and how to use it, see [What is Shared Target](https://www.alibabacloud.com/help/en/resource-management/resource-sharing/developer-reference/api-resourcesharing-2020-01-10-associateresourceshare).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_shared_target&exampleId=a9bbd278-a450-ca1f-c87e-42e9bb284fe9ad56e65d&activeTab=example&spm=docs.r.resource_manager_shared_target.0.a9bbd278a4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_accounts" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_resource_share" "default" {
  resource_share_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_shared_target" "default" {
  resource_share_id = alicloud_resource_manager_resource_share.default.id
  target_id         = data.alicloud_resource_manager_accounts.default.ids.0
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_shared_target&spm=docs.r.resource_manager_shared_target.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `resource_share_id` - (Required, ForceNew) The ID of the resource share.
* `target_id` - (Required, ForceNew) The ID of the principal.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<resource_share_id>:<target_id>`.
* `create_time` - (Available since v1.259.0) The time when the association of the entity was created.
* `status` - The status of shared target.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Shared Target.
* `delete` - (Defaults to 10 mins) Used when delete the Shared Target.

## Import

Resource Manager Shared Target can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_shared_target.example <resource_share_id>:<target_id>
```
