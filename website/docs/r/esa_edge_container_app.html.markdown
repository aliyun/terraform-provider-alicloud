---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_edge_container_app"
description: |-
  Provides a Alicloud ESA Edge Container App resource.
---

# alicloud_esa_edge_container_app

Provides a ESA Edge Container App resource.



For information about ESA Edge Container App and how to use it, see [What is Edge Container App](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateEdgeContainerApp).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tfexample"
}

resource "alicloud_esa_edge_container_app" "default" {
  target_port             = "3000"
  health_check_host       = "example.com"
  remarks                 = var.name
  health_check_port       = "80"
  health_check_uri        = "/"
  health_check_timeout    = "3"
  health_check_method     = "HEAD"
  health_check_http_code  = "http_2xx"
  health_check_fail_times = "5"
  service_port            = "80"
  health_check_interval   = "5"
  health_check_succ_times = "2"
  edge_container_app_name = var.name
  health_check_type       = "l7"
}
```

## Argument Reference

The following arguments are supported:
* `edge_container_app_name` - (Required, ForceNew) The application name must start with a lowercase letter. Lowercase letters, numbers, and bars are supported. The length is limited to 6 to 128 characters.
* `health_check_fail_times` - (Optional, ForceNew, Int) The number of consecutive successful health checks required for an application to be considered as healthy. Valid values: 1 to 10. Default value: 2.
* `health_check_host` - (Optional, ForceNew) The health check type. By default, this parameter is left empty.

Valid values:

  - `l4`: Layer 4 health check.
  - `l7`: Layer 7 health check.
* `health_check_http_code` - (Optional, ForceNew, Computed) The domain name that is used for health checks. This parameter is empty by default.
* `health_check_interval` - (Optional, ForceNew, Computed, Int) The timeout period of a health check response. If a backend ECS instance does not respond within the specified timeout period, the ECS instance fails the health check. Unit: seconds.
Valid values: `1` to `100`.
Default value: `3`.
* `health_check_method` - (Optional, ForceNew, Computed) The HTTP status code returned for a successful health check. Valid values:

  - **http\_2xx** (default)
  - **http\_3xx**
* `health_check_port` - (Optional, ForceNew, Computed, Int) The URI used for health checks. The URI must be `1` to `80` characters in length. Default value: "/".
* `health_check_succ_times` - (Optional, ForceNew, Computed, Int) The interval between two consecutive health checks. Unit: seconds. Valid values: `1` to `50`. Default value: `5`.
* `health_check_timeout` - (Optional, ForceNew, Computed, Int) The port used for health checks. Valid values: 1 to 65535. Default value: 80.
* `health_check_type` - (Optional, ForceNew, Computed) The remarks. This parameter is empty by default.
* `health_check_uri` - (Optional, ForceNew, Computed) The HTTP request method for health checks. Valid values:

  - `HEAD` (default): requests the headers of the resource.
  - `GET`: requests the specified resource and returns both the headers and entity body.
* `remarks` - (Optional, ForceNew) The backend port, which is also the service port of the application. Valid values: 1 to 65535.
* `service_port` - (Required, ForceNew, Int) The name of the application. The name must start with a lowercase letter and can contain lowercase letters, digits, and hyphens (-). The name must be 6 to 128 characters in length.
* `target_port` - (Required, ForceNew, Int) The server port. Valid values: 1 to 65535.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the application was created.
* `status` - The status of the application. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 6 mins) Used when create the Edge Container App.
* `delete` - (Defaults to 5 mins) Used when delete the Edge Container App.

## Import

ESA Edge Container App can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_edge_container_app.example <id>
```