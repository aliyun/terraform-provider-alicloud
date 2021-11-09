---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_service"
sidebar_current: "docs-alicloud-datasource-dataworks-service"
description: |-
  Provides a datasource to open the DataWorks service automatically.
---

# alicloud\_data\_works\_service

Using this data source can open DataWorks service automatically. If the service has been opened, it will return opened.

For information about DataWorks and how to use it, see [What is DataWorks](https://www.alibabacloud.com/help/en/product/72772.htm).

-> **NOTE:** Available in v1.118.0+. After the version 1.141.0, the data source is renamed as `alicloud_data_works_service`.

## Example Usage

```terraform
data "alicloud_data_works_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the DataWorks service that means you have read and agreed the [DataWorks Terms of Service](https://help.aliyun.com/document_detail/131538.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
