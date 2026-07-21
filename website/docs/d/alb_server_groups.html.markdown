---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_groups"
sidebar_current: "docs-alicloud-datasource-alb-server-groups"
description: |-
  Provides a list of Alb Server Groups to the user.
---

# alicloud_alb_server_groups

This data source provides the Alb Server Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.default.id
  server_group_name = var.name
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
}

data "alicloud_alb_server_groups" "ids" {
  ids = [alicloud_alb_server_group.default.id]
}

output "alb_server_group_id_0" {
  value = data.alicloud_alb_server_groups.ids.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Server Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Server Group name.
* `server_group_name` - (Optional, ForceNew) The names of the Server Group.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The status of the Server Group. Valid values: `Available`, `Configuring`, `Provisioning`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `server_group_ids` - (Optional, ForceNew, List) The server group IDs.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Server Group names.
* `groups` - A list of Server Groups. Each element contains the following attributes:
  * `id` - The ID of the Server Group.
  * `server_group_id` - The ID of the Server Group.
  * `server_group_name` - The name of the Server Group.
  * `vpc_id` - The ID of the VPC.
  * `protocol` - The backend protocol.
  * `scheduler` - The scheduling algorithm.
  * `status` - The status of the Server Group.
  * `tags` - The tags of the resource. **Note:** `tags` takes effect only if `enable_details` is set to `true`.
  * `sticky_session_config` - The configuration of the sticky session. **Note:** `sticky_session_config` takes effect only if `enable_details` is set to `true`.
    * `cookie` - the cookie that is configured on the server.
    * `cookie_timeout` - The timeout period of a cookie. The timeout period of a cookie.
    * `sticky_session_enabled` - Indicates whether sticky session is enabled.
    * `sticky_session_type` - The method that is used to handle a cookie.
  * `servers` - The backend server. **Note:** `servers` takes effect only if `enable_details` is set to `true`.
    * `server_type` - The type of the server. The type of the server.
    * `status` - The status of the server.
    * `weight` - The weight of the server.
    * `description` - The description of the server.
    * `port` - The port that is used by the server.
    * `server_id` - The ID of the ECS instance, ENI instance or ECI instance.
    * `server_ip` - The IP address of the ENI instance when it is in the inclusive ENI mode.
  * `health_check_config` - The configuration of health checks. **Note:** `health_check_config` takes effect only if `enable_details` is set to `true`.
    * `unhealthy_threshold` - The number of consecutive health checks that a healthy backend server must consecutively fail before it is declared unhealthy. In this case, the health check state is changed from success to fail.
    * `health_check_codes` - The status code for a successful health check. Multiple status codes can be specified as a list.
    * `health_check_path` - The forwarding rule path of health checks.
    * `healthy_threshold` - The number of health checks that an unhealthy backend server must pass consecutively before it is declared healthy. In this case, the health check state is changed from fail to success.
    * `health_check_http_version` - HTTP protocol version.
    * `health_check_interval` - The time interval between two consecutive health checks.
    * `health_check_method` - Health check method.
    * `health_check_protocol` - Health check protocol.
    * `health_check_timeout` - The timeout period of a health check response. If a backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the ECS instance is considered unhealthy.
    * `health_check_connect_port` - The port of the backend server that is used for health checks.
    * `health_check_enabled` - Indicates whether health checks are enabled.
    * `health_check_host` - The domain name that is used for health checks.
