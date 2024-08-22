---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_health_check_template"
sidebar_current: "docs-alicloud-resource-alb-health-check-template"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Health Check Template resource.
---

# alicloud_alb_health_check_template

Provides a Application Load Balancer (ALB) Health Check Template resource.

For information about Application Load Balancer (ALB) Health Check Template and how to use it, see [What is Health Check Template](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createhealthchecktemplate).

-> **NOTE:** Available since v1.134.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_alb_health_check_template&exampleId=daace97f-7616-a656-2c6a-cc7099240c29fb7b8add&activeTab=example&spm=docs.r.alb_health_check_template.0.daace97f76" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
resource "alicloud_alb_health_check_template" "example" {
  health_check_template_name = "example_name"
}
```

## Argument Reference

The following arguments are supported:
* `health_check_template_name` - (Required) The name of the health check template.  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `dry_run` - (Optional) Whether to precheck the API request.
* `health_check_codes` - (Optional) The HTTP status code that indicates a successful health check. **NOTE:** The attribute `HealthCheckProtocol` is valid when the attribute is  `HTTP` .
* `health_check_connect_port` - (Optional) The number of the port that is used for health checks.  Valid values: `0` to `65535`.  Default value: `0`. This default value indicates that the backend server is used for health checks.
* `health_check_host` - (Optional) The domain name that is used for health checks. Default value:  `$SERVER_IP`. The domain name must be 1 to 80 characters in length.  **NOTE:** The attribute `HealthCheckProtocol` is valid when the attribute is  `HTTP` .
* `health_check_http_version` - (Optional) The version of the HTTP protocol.  Valid values: `HTTP1.0` and `HTTP1.1`.  Default value: `HTTP1.1`. **NOTE:** The attribute `HealthCheckProtocol` is valid when the attribute is  `HTTP` .
* `health_check_interval` - (Optional) The time interval between two consecutive health checks.  Valid values: `1` to `50`. Unit: seconds.  Default value: `2`.
* `health_check_method` - (Optional) The health check method.  Valid values: GET and HEAD.  Default value: HEAD. **NOTE:** The attribute `HealthCheckProtocol` is valid when the attribute is  `HTTP` .
* `health_check_path` - (Optional) The URL that is used for health checks.  The URL must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URL can also contain the following extended characters: _ ; ~ ! ( )* [ ] @ $ ^ : ' , +. The URL must start with a forward slash (/). **NOTE:** The attribute `HealthCheckProtocol` is valid when the attribute is  `HTTP` .
* `health_check_protocol` - (Optional) The protocol that is used for health checks.  Valid values: `HTTP` and `TCP`.  Default value: `HTTP`.
* `health_check_timeout` - (Optional) The timeout period of a health check response. If the backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the health check fails.  Valid values: `1` to `300`. Unit: seconds.  Default value: `5`.
* `healthy_threshold` - (Optional) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy (from fail to success).  Valid values: `2` to `10`.  Default value: `3`. Unit: seconds.
* `unhealthy_threshold` - (Optional)The number of times that an healthy backend server must consecutively fail health checks before it is declared unhealthy (from success to fail). Valid values: `2` to `10`.  Default value: `3`. Unit: seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Health Check Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Health Check Template.
* `delete` - (Defaults to 2 mins) Used when delete the Health Check Template.
* `update` - (Defaults to 2 mins) Used when update the Health Check Template.

## Import

Application Load Balancer (ALB) Health Check Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_health_check_template.example <id>
```
