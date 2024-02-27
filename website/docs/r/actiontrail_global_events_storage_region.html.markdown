---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_global_events_storage_region"
sidebar_current: "docs-alicloud-resource-actiontrail-global-events-storage-region"
description: |-
  Provides Alibaba Cloud Actiontrail Global Events Storage Region Resource
---

# alicloud_actiontrail_global_events_storage_region

Provides a Global events storage region resource.

For information about global events storage region and how to use it, see [What is Global Events Storage Region](https://help.aliyun.com/zh/actiontrail/developer-reference/api-actiontrail-2020-07-06-updateglobaleventsstorageregion).

-> **NOTE:** Available since v1.201.0.

## Example Usage

```terraform
resource "alicloud_actiontrail_global_events_storage_region" "foo" {
  storage_region = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `storage_region` - (Optional) Global Events Storage Region.

## Attributes Reference

The following attributes are exported:


## Import

Global events storage region not can be imported.

