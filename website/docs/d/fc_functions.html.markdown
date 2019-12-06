---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_functions"
sidebar_current: "docs-alicloud-datasource-fc-functions"
description: |-
    Provides a list of FC functions to the user.
---

# alicloud\_fc_functions

This data source provides the Function Compute functions of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_fc_functions" "functions_ds" {
  service_name = "sample_service"
  name_regex   = "sample_fc_function"
}

output "first_fc_function_name" {
  value = "${data.alicloud_fc_functions.functions_ds.functions.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - Name of the service that contains the functions to find.
* `name_regex` - (Optional) A regex string to filter results by function name.
* `ids` (Optional, Available in 1.53.0+) - A list of functions ids.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of functions ids.
* `names` - A list of functions names.
* `functions` - A list of functions. Each element contains the following attributes:
  * `id` - Function ID.
  * `name` - Function name.
  * `description` - Function description.
  * `runtime` - Function runtime. The list of possible values is [available here](https://www.alibabacloud.com/help/doc-detail/52077.htm).
  * `handler` - Function [entry point](https://www.alibabacloud.com/help/doc-detail/62213.htm) in the code.
  * `timeout` - Maximum amount of time the function can run in seconds.
  * `memory_size` - Amount of memory in MB the function can use at runtime.
  * `code_size` - Function code size in bytes.
  * `code_checksum` - Checksum (crc64) of the function code.
  * `creation_time` - Function creation time.
  * `last_modification_time` - Function last modification time.
  * `environment_variables` -  A map that defines environment variables for the function.
