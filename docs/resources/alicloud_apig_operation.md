# alicloud_apig_operation

Provides an APIG Operation resource.

An operation (API interface) under an HttpApi in Alibaba Cloud API Gateway. It defines a single API endpoint (path + HTTP method) within the parent HttpApi.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_apig_http_api" "default" {
  http_api_name = "tf-test-api"
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = "/test"
}

resource "alicloud_apig_operation" "default" {
  http_api_id    = alicloud_apig_http_api.default.id
  operation_name = "tf-test-operation"
  path           = "/test-op"
  method         = "GET"
  description    = "test operation"
  mock {
    enable           = true
    response_code    = 200
    response_content = "hello"
  }
}
```

## Argument Reference

The following arguments are supported:

* `http_api_id` - (Required, ForceNew) The ID of the parent HttpApi.
* `operation_name` - (Required) The name of the operation.
* `path` - (Required) The path of the operation. Length must be between 1 and 2048.
* `method` - (Required) The HTTP method of the operation. Valid values: `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS`.
* `description` - (Optional) The description of the operation. Length must be between 0 and 255.
* `mock` - (Optional) The mock configuration block. See [`mock`](#mock) below.

### mock

The `mock` block supports:

* `enable` - (Optional) Whether to enable mock.
* `response_code` - (Optional) The response status code.
* `response_content` - (Optional) The response content.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, formatted as `<http_api_id>:<operation_id>`.
* `operation_id` - The ID of the operation.
* `create_time` - The creation timestamp of the operation.

## Import

APIG Operation can be imported using the ID in the format `<http_api_id>:<operation_id>`, e.g.

```shell
terraform import alicloud_apig_operation.example <http_api_id>:<operation_id>
```
