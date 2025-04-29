---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_target"
sidebar_current: "docs-alicloud-resource-resource-manager-shared-target"
description: |-
  Provides a Alicloud Resource Manager Shared Target resource.
---

# alicloud_resource_manager_shared_target

Provides a Resource Manager Shared Target resource.

For information about Resource Manager Shared Target and how to use it, see [What is Shared Target](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_shared_target&exampleId=a10a0f37-3fdb-86a4-26f9-86b934c54cbb7f1afac4&activeTab=example&spm=docs.r.resource_manager_shared_target.0.a10a0f373f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
data "alicloud_resource_manager_accounts" "default" {}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = var.name
}

resource "alicloud_resource_manager_shared_target" "example" {
  resource_share_id = alicloud_resource_manager_resource_share.example.id
  target_id         = data.alicloud_resource_manager_accounts.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `resource_share_id` - (Required, ForceNew) The resource share ID of resource manager.
* `target_id` - (Required, ForceNew) The member account ID in resource directory.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Shared Target. The value is formatted `<resource_share_id>:<target_id>`.
* `status` - The status of shared target.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Shared Target.
* `delete` - (Defaults to 11 mins) Used when delete the Shared Target.

## Import

Resource Manager Shared Target can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_shared_target.example <resource_share_id>:<target_id>
```
