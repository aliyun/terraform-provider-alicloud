---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_service"
sidebar_current: "docs-alicloud-datasource-hbr-service"
description: |-
  Provides a datasource to open the HBR service automatically.
---

# alicloud_hbr_service

Using this data source can open HBR service automatically. If the service has been opened, it will return opened.

For information about HBR and how to use it, see [What is HBR](https://www.alibabacloud.com/help/en/hybrid-backup-recovery).

-> **NOTE:** Available since v1.184.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the HBR service that means you have read and agreed the [HBR Terms of Service](https://help.aliyun.com/document_detail/62906.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
