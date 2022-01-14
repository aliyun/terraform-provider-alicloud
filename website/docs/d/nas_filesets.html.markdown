---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_filesets"
sidebar_current: "docs-alicloud-datasource-nas-filesets"
description: |-
  Provides a list of Nas Filesets to the user.
---

# alicloud\_nas\_filesets

This data source provides the Nas Filesets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_filesets" "ids" {
  file_system_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "nas_fileset_id_1" {
  value = data.alicloud_nas_filesets.ids.filesets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `ids` - (Optional, ForceNew, Computed)  A list of Fileset IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the fileset. Valid values: `CREATED`, `CREATING`, `RELEASED`, `RELEASING`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `filesets` - A list of Nas Filesets. Each element contains the following attributes:
	* `create_time` - The time when Fileset was created.
	* `description` - Description of Fileset.
	* `file_system_id` - The ID of the file system.
	* `file_system_path` - The path of Fileset.
	* `fileset_id` - The first ID of the resource.
	* `id` - The ID of the Fileset.
	* `status` - The status of the fileset.
	* `update_time` - The latest update time of Fileset.