---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_groups"
sidebar_current: "docs-alicloud-datasource-nlb-server-groups"
description: |-
  Provides a list of Nlb Server Groups to the user.
---

# alicloud\_nlb\_server\_groups

This data source provides the Nlb Server Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_server_groups" "ids" {}
output "nlb_server_group_id_1" {
  value = data.alicloud_nlb_server_groups.ids.groups.0.id
}

data "alicloud_nlb_server_groups" "nameRegex" {
  name_regex = "^my-ServerGroup"
}
output "nlb_server_group_id_2" {
  value = data.alicloud_nlb_server_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Server Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Server Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which the security group belongs.
* `server_group_names` - (Optional, ForceNew) The names of the server groups to be queried.
* `server_group_type` - (Optional, ForceNew) The type of the server group. Valid values: `Instance`, `Ip`.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Configuring`, `Creating`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Server Group names.
* `groups` - A list of Nlb Server Groups. Each element contains the following attributes:
	* `address_ip_version` - The protocol version.
	* `connection_drain` - Indicates whether connection draining is enabled.
	* `connection_drain_timeout` - The timeout period of connection draining. Unit: seconds.
	* `health_check` - The configurations of health checks.
		* `health_check_interval` - The interval between two consecutive health checks.
		* `health_check_type` - The protocol that is used for health checks.
		* `healthy_threshold` - The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy.
		* `unhealthy_threshold` - The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy.
		* `health_check_http_code` - The HTTP status codes returned for health checks.
		* `health_check_url` - The path to which health check requests are sent.
		* `health_check_connect_port` - The backend port that is used for health checks.
		* `health_check_connect_timeout` - The maximum timeout period of a health check response.
		* `health_check_domain` - The domain name that is used for health checks.
		* `health_check_enabled` - Specifies whether to enable health checks.
		* `http_check_method` - The HTTP method that is used for health checks.
	* `protocol` - The protocol used to forward requests to the backend servers.
	* `related_load_balancer_ids` - The NLB instance.
	* `scheduler` - The routing algorithm.
	* `server_count` - The number of server groups associated with the NLB instance.
	* `server_group_name` - The name of the server group.
	* `server_group_type` - The type of the server group.
	* `status` - The status of the server group.
	* `id` - The ID of the Server Group.
	* `vpc_id` - The ID of the VPC to which the server group belongs.
	* `preserve_client_ip_enabled` - Indicates whether client address retention is enabled.
	* `resource_group_id` - The ID of the resource group to which the security group belongs.
	* `tags` - A mapping of tags to assign to the resource.