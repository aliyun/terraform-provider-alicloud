# alicloud_apig_operations

This data source provides the APIG Operations available in the current Alibaba Cloud account.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_apig_operations" "default" {
  http_api_id = alicloud_apig_http_api.default.id
  name        = "tf-test-operation"
}

output "operation_id" {
  value = data.alicloud_apig_operations.default.operations.0.operation_id
}
```

## Argument Reference

The following arguments are supported:

* `http_api_id` - (Required) The ID of the parent HttpApi.
* `ids` - (Optional, Computed) A list of operation IDs in the format `<http_api_id>:<operation_id>`.
* `name` - (Optional) The exact name of the operation to search by.
* `name_like` - (Optional) A name prefix used to search operations.
* `path_like` - (Optional) A path prefix used to search operations.
* `method` - (Optional) The HTTP method used to filter operations.
* `name_regex` - (Optional) A regex string to filter operations by name.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of operation IDs.
* `names` - A list of operation names.
* `operations` - A list of APIG Operations. Each element contains the following attributes:
  * `id` - The operation ID, formatted as `<http_api_id>:<operation_id>`.
  * `http_api_id` - The ID of the parent HttpApi.
  * `operation_id` - The ID of the operation.
  * `operation_name` - The name of the operation.
  * `path` - The path of the operation.
  * `method` - The HTTP method of the operation.
  * `description` - The description of the operation.
  * `create_time` - The creation timestamp of the operation.
  * `mock` - The mock configuration block.
    * `enable` - Whether mock is enabled.
    * `response_code` - The response status code.
    * `response_content` - The response content.
