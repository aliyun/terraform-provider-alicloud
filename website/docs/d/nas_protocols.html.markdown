---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_protocols"
sidebar_current: "docs-alicloud-datasource-nas-protocols"
description: |-
    Provides a list of FileType owned by an Alibaba Cloud account.
---

# alicloud\_nas_protocols

Provide  a data source to retrieve the type of protocol used to create NAS file system.

-> **NOTE:** Available in 1.42.0

## Example Usage

```terraform
data "alicloud_nas_protocols" "default" {
  type        = "Performance"
  zone_id     = "cn-beijing-e"
  output_file = "protocols.txt"
}

output "nas_protocols_protocol" {
  value = "${data.alicloud_nas_protocols.default.protocols.0}"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The file system type. Valid Values: `Performance` and `Capacity`.  
* `zone_id` - (Optional) String to filter results by zone id. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `protocols` - A list of supported protocol type..
