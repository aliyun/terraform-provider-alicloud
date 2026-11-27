---
layout: "alicloud"
page_title: "Alicloud: alicloud_file_crc64_checksum"
sidebar_current: "docs-alicloud-datasource-file-crc64-checksum"
description: |-
    Provides compute file crc64 checksum.
---

# alicloud_file_crc64_checksum

This data source compute file crc64 checksum.

-> **NOTE:** Available since v1.59.0+.

## Example Usage

```terraform
data "alicloud_file_crc64_checksum" "default" {
  filename = "exampleFileName"
}

output "file_crc64_checksum" {
  value = "${data.alicloud_file_crc64_checksum.defualt.checksum}"
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) The name of the file to be computed crc64 checksum.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported:

* `id` - file crc64 ID
* `checksum` - the file checksum of crc64.
