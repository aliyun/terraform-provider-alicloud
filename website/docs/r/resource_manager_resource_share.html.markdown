---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_share"
description: |-
  Provides a Alicloud Resource Manager Resource Share resource.
---

# alicloud_resource_manager_resource_share

Provides a Resource Manager Resource Share resource.

RS resource sharing.

For information about Resource Manager Resource Share and how to use it, see [What is Resource Share](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_resource_share&exampleId=6be253ac-e5f2-8690-f13b-fe32084c06e39ce31fd6&activeTab=example&spm=docs.r.resource_manager_resource_share.0.6be253ace5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_resource_share&spm=docs.r.resource_manager_resource_share.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `allow_external_targets` - (Optional, Available since v1.261.0) Whether to allow sharing to accounts outside the resource directory. Value:
  - false (default): Only sharing within the resource directory is allowed.
  - true: Allow sharing to any account.
* `permission_names` - (Optional, List, Available since v1.261.0) Share permission name. When it is empty, the system automatically binds the default permissions associated with the resource type. For more information, see [Permission Library](https://www.alibabacloud.com/help/en/resource-management/resource-sharing/user-guide/permissions-for-resource-sharing).
* `resource_group_id` - (Optional, Computed, Available since v1.261.0) The ID of the resource group
* `resource_share_name` - (Required) The name of resource share.
* `resources` - (Optional, List, Available since v1.261.0) List of shared resources. **Note: The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.** See [`resources`](#resources) below.
* `tags` - (Optional, Map, Available since v1.261.0) The tag of the resource
* `targets` - (Optional, List, Available since v1.261.0) Resource user.

### `resources`

The resources supports the following:
* `resource_id` - (Optional, Available since v1.261.0) The ID of the shared resource.
* `resource_type` - (Optional, Available since v1.261.0) Shared resource type. For the types of resources that support sharing, see [Cloud services that support sharing](https://www.alibabacloud.com/help/en/resource-management/resource-sharing/product-overview/services-that-work-with-resource-sharing).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The create time of resource share.
* `resource_share_owner` - The owner of resource share,  `Self` and `OtherAccounts`.
* `status` - The status of resource share.  `Active`,`Deleted` and `Deleting`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Share.
* `delete` - (Defaults to 10 mins) Used when delete the Resource Share.
* `update` - (Defaults to 5 mins) Used when update the Resource Share.

## Import

Resource Manager Resource Share can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_share.example <id>
```