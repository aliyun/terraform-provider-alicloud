---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_global_events_storage_region"
sidebar_current: "docs-alicloud-resource-actiontrail-global-events-storage-region"
description: |-
  Provides Alibaba Cloud Actiontrail Global Events Storage Region Resource
---

# alicloud\_actiontrail\_global\_events\_storage\_region

Provides a Global events storage region resource.

For information about global events storage region and how to use it, see [What is Global Events Storage Region](https://help.aliyun.com/document_detail/608293.html).

-> **NOTE:** Available in 1.201.0+

## Example Usage

```terraform
resource "alicloud_actiontrail_global_events_storage_region" "foo" {
  storage_region = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `storage_region` - (Optional, Computed) Global Events Storage Region.

## Attributes Reference

The following attributes are exported:

* `storage_region` - Global Events Storage Region.

## Import

Global events storage region not can be imported.
```
