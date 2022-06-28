---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_service"
sidebar_current: "docs-alicloud-datasource-brain_industrial-service"
description: |-
    Provides a datasource to open the Brain Industrial service automatically.
---

# alicloud\_brain\_industrial\_service

Using this data source can open Brain Industrial service automatically. If the service has been opened, it will return opened.

-> **NOTE:** Available in v1.115.0+

-> **NOTE:** The Brain Industrial service is not support in the international site.

## Example Usage

```terraform
data "alicloud_brain_industrial_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Brain Industrial service. The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
