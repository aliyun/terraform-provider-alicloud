---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_rd_default_sync_list"
description: |-
  Provides a Alicloud Threat Detection Rd Default Sync List data source.
---

# alicloud_threat_detection_rd_default_sync_list

Provides a Threat Detection Rd Default Sync List data source. The default synchronization list of resource directory folders for Threat Detection. The list is account-level (singleton): each Alibaba Cloud account holds at most one such list, and this data source reads the current account's list.

For information about Threat Detection Rd Default Sync List and how to use it, see [ListRdDefaultSyncList](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListRdDefaultSyncList).

-> **NOTE:** Available since v1.292.0.

-> **NOTE:** The resource folder ids can be obtained via the `GetRdTree` API. A resource directory management account or a Threat Detection delegated administrator account is required to call these APIs.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_rd_default_sync_list" "default" {
}

output "rd_default_sync_folder_ids" {
  value = data.alicloud_threat_detection_rd_default_sync_list.default.folder_ids
}
```

## Argument Reference

The data source does not support any arguments.

## Attributes Reference

The following attributes are exported:

* `id` - The data source id. It is the id of the Alibaba Cloud account that owns the synchronization list.
* `folder_ids` - The current list of synchronized resource directory folder ids.
