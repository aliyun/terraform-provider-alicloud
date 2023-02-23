---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_global_events_storage_region"
sidebar_current: "docs-alicloud-datasource-actiontrail-global-events-storage-region"
description: |-
  Provides a list of Actiontrail Global Events Storage Region to the user.
---

# alicloud\_actiontrail\_global\_events\_storage\_region

This data source provides the Actiontrail Global Events Storage Region of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.201.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_actiontrail_global_events_storage_region" "default" {
}
output "alicloud_actiontrail_global_events_storage_region_1" {
  value = data.alicloud_actiontrail_global_events_storage_region.default.storage_region
}
```

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `storage_region` - Global Events Storage Region.