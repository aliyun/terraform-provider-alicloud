---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_target"
sidebar_current: "docs-alicloud-resource-resource-manager-shared-target"
description: |-
  Provides a Alicloud Resource Manager Shared Target resource.
---

# alicloud\_resource\_manager\_shared\_target

Provides a Resource Manager Shared Target resource.

For information about Resource Manager Shared Target and how to use it, see [What is Shared Target](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_accounts" "example" {}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = "example_value"
}

resource "alicloud_resource_manager_shared_target" "example" {
  resource_share_id = alicloud_resource_manager_resource_share.example.resource_share_id
  target_id         = data.alicloud_resource_manager_accounts.example.ids.0
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Shared Target.
* `delete` - (Defaults to 11 mins) Used when delete the Shared Target.

## Import

Resource Manager Shared Target can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_shared_target.example <resource_share_id>:<target_id>
```
