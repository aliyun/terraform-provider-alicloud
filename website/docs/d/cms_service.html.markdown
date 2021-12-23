---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_service"
sidebar_current: "docs-alicloud-datasource-cms-service"
description: |-
    Provides a datasource to open the CMS service automatically.
---

# alicloud\_cms\_service

Using this data source can open CMS service automatically. If the service has been opened, it will return opened.

For information about CMS and how to use it, see [What is CMS](https://help.aliyun.com/product/28572.html).

-> **NOTE:** Available in v1.111.0+

## Example Usage

```
data "alicloud_cms_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the CMS service that means you have read and agreed the [CMS Terms of Service](https://help.aliyun.com/document_detail/68185.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
