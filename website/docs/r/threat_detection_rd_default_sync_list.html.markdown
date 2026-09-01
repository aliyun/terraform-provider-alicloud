---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_rd_default_sync_list"
description: |-
  Provides a Alicloud Threat Detection Rd Default Sync List resource.
---

# alicloud_threat_detection_rd_default_sync_list

Provides a Threat Detection Rd Default Sync List resource. The default synchronization list of resource directory folders for Threat Detection. The resource is account-level (singleton): each Alibaba Cloud account holds at most one such list, and applying a new configuration replaces the previous folder list in full (set/replace semantics). Setting an empty folder list clears the synchronized folders, which is equivalent to disabling the default synchronization.

For information about Threat Detection Rd Default Sync List and how to use it, see [CreateRdDefaultSyncList](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateRdDefaultSyncList) and [ListRdDefaultSyncList](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListRdDefaultSyncList).

-> **NOTE:** Available since v1.292.0.

-> **NOTE:** The resource folder ids can be obtained via the `GetRdTree` API. A resource directory management account or a Threat Detection delegated administrator account is required to call these APIs.

-> **NOTE:** This resource is account-level (singleton): each account holds at most one synchronization list. When adopting an account that already holds a list (for example after `terraform import`), declare `folder_ids` in the configuration with the current folder ids. Omitting `folder_ids` on create will not send the `FolderIds` parameter, so the existing synchronized list is not wiped out; to clear the list intentionally, set `folder_ids = []`.

-> **NOTE:** Declaring multiple `alicloud_threat_detection_rd_default_sync_list` resources for the same account is not recommended. Each apply replaces the entire folder list (set/replace semantics), so the last applied configuration wins and resources created earlier will show drift on the next plan. Always pass the full set of `folder_ids` on every update; omitting them clears the existing list.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_rd_default_sync_list" "default" {
  folder_ids = ["fd-xxxxx", "fd-yyyyy"]
}
```

## Argument Reference

The following arguments are supported:

* `folder_ids` - (Optional) The list of resource directory folder ids to be synchronized by Threat Detection. Setting an empty list clears the synchronized folders (disables the default synchronization). Folder ids can be obtained via the `GetRdTree` API.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id. It is the id of the Alibaba Cloud account that owns the synchronization list.
* `folder_ids` - The current list of synchronized resource directory folder ids.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Rd Default Sync List.
* `update` - (Defaults to 5 mins) Used when update the Rd Default Sync List.
* `delete` - (Defaults to 5 mins) Used when delete the Rd Default Sync List.

## Import

Threat Detection Rd Default Sync List can be imported using the id, e.g. the Alibaba Cloud account id, or an empty string to import the current account's list:

```shell
$ terraform import alicloud_threat_detection_rd_default_sync_list.example
```
