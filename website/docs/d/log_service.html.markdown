---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_service"
sidebar_current: "docs-alicloud-datasource-log-service"
description: |-
    Provides a datasource for Log service.
---

# alicloud_log_service

Log service is enabled automatically when a project is created. This data source is retained for compatibility and no longer calls service APIs.

For information about Log service and how to use it, see [What is Log Service](https://www.alibabacloud.com/help/product/28958.htm).

-> **NOTE:** Available since v1.96.0

## Example Usage

```terraform
data "alicloud_log_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` returns `Opened`. Setting it to `Off` returns an empty status. Valid values: "On" or "Off". Default to "Off".

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - Returns `Opened` when `enable` is `On`, or an empty string otherwise. This is a compatibility value, not the actual service status.
