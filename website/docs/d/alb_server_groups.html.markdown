---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_groups"
sidebar_current: "docs-alicloud-datasource-alb-server-groups"
description: |-
  Provides a list of Alb Server Groups to the user.
---

# alicloud\_alb\_server\_groups

This data source provides the Alb Server Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_server_groups" "ids" {}
output "alb_server_group_id_1" {
  value = data.alicloud_alb_server_groups.ids.groups.0.id
}

data "alicloud_alb_server_groups" "nameRegex" {
  name_regex = "^my-ServerGroup"
}
output "alb_server_group_id_2" {
  value = data.alicloud_alb_server_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Server Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Server Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `server_group_ids` - (Optional, ForceNew) The server group ids.
* `server_group_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The status of the resource.
* `tag` - (Optional, ForceNew) The tag.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC that you want to access.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Server Group names.
* `groups` - A list of Alb Server Groups. Each element contains the following attributes:
	* `health_check_config` - The configuration of health checks.
		* `unhealthy_threshold` - The number of consecutive health checks that a healthy backend server must consecutively fail before it is declared unhealthy. In this case, the health check state is changed from success to fail. Valid values: `2` to `10`. Default value: `3`.
		* `health_check_codes` - The status code for a successful health check. Multiple status codes can be specified as a list. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. Default value: `http_2xx`. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
		* `health_check_path` - The forwarding rule path of health checks. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
		* `healthy_threshold` - The number of health checks that an unhealthy backend server must pass consecutively before it is declared healthy. In this case, the health check state is changed from fail to success. Valid values: `2` to `10`. Default value: `3`.
		* `health_check_http_version` - HTTP protocol version. Valid values: `HTTP1.0` and `HTTP1.1`. Default value: `HTTP1.1`. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
		* `health_check_interval` - The time interval between two consecutive health checks. Unit: seconds. Valid values: `1` to `50`. Default value: `2`.
		* `health_check_method` - Health check method. Valid values: `GET` and `HEAD`. Default: `GET`. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
		* `health_check_protocol` - Health check protocol. Valid values: `HTTP` and `TCP`.
		* `health_check_timeout` - The timeout period of a health check response. If a backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the ECS instance is considered unhealthy. Unit: seconds. Valid values: `1` to `300`. Default value: `5`. **NOTE:** If the value of the `HealthCHeckTimeout` parameter is smaller than that of the `HealthCheckInterval` parameter, the value of the `HealthCHeckTimeout` parameter is ignored and the value of the `HealthCheckInterval` parameter is regarded as the timeout period.
		* `health_check_connect_port` - The port of the backend server that is used for health checks. Valid values: `0` to `65535`. Default value: `0`. A value of `0` indicates that a backend server port is used for health checks.
		* `health_check_enabled` - Indicates whether health checks are enabled. Valid values: `true`, `false`. Default value: `true`.
		* `health_check_host` - The domain name that is used for health checks.
	* `id` - The ID of the Server Group.
	* `protocol` - The server protocol. Valid values: `HTTP` and `HTTPS`. Default value: `HTTP`.
    * `scheduler` - The scheduling algorithm. Valid values: `Wrr`, `Wlc` and `Sch`.
	* `server_group_id` - The first ID of the res ource.
	* `server_group_name` - The name of the resource.
	* `servers` - The backend server.
		* `server_type` - The type of the server. The type of the server. Valid values: `Ecs`, `Eni` and `Eci`.
		* `status` - The status of the server. Valid values: `Adding`, `Available`, `Removing` and `Configuring`.
		* `weight` - The weight of the server.  Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no requests are forwarded to the server.
		* `description` - The description of the server.
		* `port` - The port that is used by the server. Valid values: `1` to `65535`.
		* `server_id` - The ID of the ECS instance, ENI instance or ECI instance.
		* `server_ip` - The IP address of the ENI instance when it is in the inclusive ENI mode.
	* `status` - The status of the resource. Valid values: `Provisioning`, `Available` and `Configuring`.
	* `sticky_session_config` - The configuration of the sticky session.
		* `cookie` - the cookie that is configured on the server. **NOTE:** This parameter exists if the `StickySession` parameter is set to `On` and the `StickySessionType` parameter is set to `server`.
		* `cookie_timeout` - The timeout period of a cookie. The timeout period of a cookie. Unit: seconds. Valid values: `1` to `86400`. Default value: `1000`.
		* `sticky_session_enabled` - Indicates whether sticky session is enabled. Values: `true` and `false`. Default value: `false`.  **NOTE:** This parameter exists if the `StickySession` parameter is set to `On`.
		* `sticky_session_type` - The method that is used to handle a cookie. Values: `Server` and `Insert`. 
	* `vpc_id` - The ID of the VPC that you want to access.
