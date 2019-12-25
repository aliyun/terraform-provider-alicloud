---
layout: "alicloud"
page_title: "Alicloud: alicloud_file_crc64_checksum"
sidebar_current: "docs-alicloud-datasource-file-crc64-checksum"
description: |-
    Provides compute file crc64 checksum.
---

# alicloud\_file_crc64_checksum

This data source compute file crc64 checksum.

-> **NOTE:** Available in 1.59.0+.

## Example Usage

```
data "alicloud_file_crc64_checksum" "default" {
  filename = "exampleFileName"
}

output "file_crc64_checksum" {
  value = data.alicloud_file_crc64_checksum.defualt.checksum
}
```
## Argument Reference

The following arguments are supported:

* `filename` - (Required) The name of the file to be computed crc64 checksum.

## Attributes Reference

The following attributes are exported:

* `id` - file crc64 ID
* `checksum` - the file checksum of crc64.
