---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_service"
sidebar_current: "docs-alicloud-datasource-api-gateway-service"
description: |-
    Provides a datasource to open the API gateway service automatically.
---

# alicloud\_api\_gateway\_service

Using this data source can open API gateway service automatically. If the service has been opened, it will return opened.

## Example Usage

```
data "alicloud_api_gateway_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been opened, return the result.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 