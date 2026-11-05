---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_created_by"
sidebar_current: "docs-alicloud-resource-tag-created-by"
description: |-
  Provides a Alicloud TAG Created By resource.
---

# alicloud_tag_created_by

Provides a TAG Created By resource.

The account-level switch that enables the system `CreatedBy` tag on resources created in your account. The existence of the resource indicates that the switch is enabled; removing the resource disables it.

For information about TAG Created By and how to use it, see [What is Created By](https://next.api.alibabacloud.com/document/Tag/2018-08-28/OpenCreatedBy).

-> **NOTE:** Available since v2.1.1.

-> **NOTE:** The switch is account-scoped and only one exists per account. Tag serves `OpenCreatedBy`, `CheckCreatedByEnabled` and `CloseCreatedBy` only from the `cn-shanghai` endpoint, so the resource always calls that control plane and the resource ID is `cn-shanghai`, regardless of the region configured on the provider.

## Example Usage

Basic Usage

```terraform
resource "alicloud_tag_created_by" "example" {}
```

## Argument Reference

The following arguments are supported:

* `timeouts` - (Optional, Object) See [Timeouts](#timeouts) below.

## Attributes Reference

The following attributes are exported:
* `id` - (Computed, String) The CreatedBy control region, which is `cn-shanghai`.
* `open_status` - (Computed, Bool) Indicates whether the createdby tag is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Created By.
* `delete` - (Defaults to 5 mins) Used when delete the Created By.
* `update` - (Defaults to 5 mins) Used when update the Created By.

## Import

TAG Created By can be imported using the control region ID `cn-shanghai`, e.g.

```shell
$ terraform import alicloud_tag_created_by.example cn-shanghai
```
