---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_log_analyses"
description: |-
  Provides a list of File Storage (NAS) Log Analysis owned by an Alicloud account.
---

# alicloud_nas_log_analyses

This data source provides File Storage (NAS) Log Analysis available to the user.

-> **NOTE:** Available since v1.286.0.

## Example Usage

```terraform
resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_log_analysis" "default" {
  file_system_id = alicloud_nas_file_system.default.id
}

data "alicloud_nas_log_analyses" "default" {
  ids = [alicloud_nas_log_analysis.default.id]
}

output "alicloud_nas_log_analysis_example_id" {
  value = data.alicloud_nas_log_analyses.default.analyses.0.id
}
```

## Argument Reference

The following arguments are supported:
* `file_system_type` - (Optional, ForceNew) The type of the file system. Valid values: `standard`, `extreme`, `all`.
* `ids` - (Optional, ForceNew, List) A list of Log Analysis IDs. Each element is the ID of the file system for which log delivery is enabled.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `analyses` - A list of Log Analysis. Each element contains the following attributes:
  * `id` - The ID of the Log Analysis. It is the same as the file system ID.
  * `file_system_id` - The ID of the file system for which log delivery is enabled.
  * `logstore` - The name of the Logstore that receives NAS logs.
  * `project` - The name of the project that receives NAS logs.
  * `region` - The Simple Log Service region of the log project.
  * `role_arn` - The ARN of the service role used by NAS to deliver logs to Simple Log Service.
