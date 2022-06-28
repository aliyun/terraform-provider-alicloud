---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_health_check_templates"
sidebar_current: "docs-alicloud-datasource-alb-health-check-templates"
description: |-
  Provides a list of Alb Health Check Templates to the user.
---

# alicloud\_alb\_health\_check\_templates

This data source provides the Alb Health Check Templates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_health_check_templates" "ids" {
  ids = ["example_id"]
}
output "alb_health_check_template_id_1" {
  value = data.alicloud_alb_health_check_templates.ids.templates.0.id
}

data "alicloud_alb_health_check_templates" "nameRegex" {
  name_regex = "^my-HealthCheckTemplate"
}
output "alb_health_check_template_id_2" {
  value = data.alicloud_alb_health_check_templates.nameRegex.templates.0.id
}

```

## Argument Reference

The following arguments are supported:

* `health_check_template_ids` - (Optional, ForceNew) The health check template ids.
* `health_check_template_name` - (Optional, ForceNew) The name of the health check template.  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter. 
* `ids` - (Optional, ForceNew, Computed)  A list of Health Check Template IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Health Check Template name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Health Check Template names.
* `templates` - A list of Alb Health Check Templates. Each element contains the following attributes:
	* `health_check_codes` - The HTTP status code that indicates a successful health check.
	* `health_check_connect_port` - The number of the port that is used for health checks.  Valid values: `0` to `65535`.  Default value:` 0`. This default value indicates that the backend server is used for health checks.
	* `health_check_host` - The domain name that is used for health checks. Default value:  `$SERVER_IP`. The domain name must be 1 to 80 characters in length. 
	* `health_check_http_version` - The version of the HTTP protocol.  Valid values: `HTTP1.0` and `HTTP1.1`.  Default value: `HTTP1.1`.
	* `health_check_interval` - The time interval between two consecutive health checks.  Valid values: `1` to `50`. Unit: seconds.  Default value: `2`.
	* `health_check_method` - The health check method.  Valid values: `GET` and `HEAD`.  Default value: `HEAD`.
	* `health_check_path` - The URL that is used for health checks.  The URL must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URL can also contain the following extended characters: ` _ ; ~ ! ( )* [ ] @ $ ^ : ' , +. The URL must start with a forward slash (/)`.
	* `health_check_protocol` - The protocol that is used for health checks.  Valid values: HTTP and TCP.  Default value: HTTP.
	* `health_check_template_id` - The ID of the resource.
	* `health_check_template_name` - The name of the health check template.  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
	* `health_check_timeout` - The timeout period of a health check response. If the backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the health check fails.  Valid values: `1` to `300`. Unit: seconds.  Default value: `5`.
	* `healthy_threshold` - The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy (from fail to success). Valid values: `2` to `10`.  Default value: `3`. Unit: seconds.
	* `id` - The ID of the Health Check Template.
	* `unhealthy_threshold` - The number of times that an healthy backend server must consecutively fail health checks before it is declared unhealthy (from success to fail). Valid values: `2` to `10`.  Default value: `3`. Unit: seconds.
