---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_health_check_template"
sidebar_current: "docs-alicloud-resource-alb-health-check-template"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Health Check Template resource.
---

# alicloud_alb_health_check_template

Provides an Application Load Balancer (ALB) Health Check Template resource.

For information about Application Load Balancer (ALB) Health Check Template and how to use it, see [What is Health Check Template](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createhealthchecktemplate).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_health_check_template&exampleId=5c4b6f28-788d-9a3d-e8e7-b982d52d64e73638f458&activeTab=example&spm=docs.r.alb_health_check_template.0.5c4b6f2878&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_alb_health_check_template" "example" {
  health_check_template_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `health_check_template_name` - (Required) The name of the health check template. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `health_check_connect_port` - (Optional, Int) The port that is used for health checks. Default value: `0`. Valid values: `0` to `65535`.
* `health_check_host` - (Optional) The domain name that is used for health checks. **NOTE:** `health_check_host` takes effect only if `health_check_protocol` is set to `HTTP`.
* `health_check_http_version` - (Optional) The version of the HTTP protocol. Default value: `HTTP1.1`. Valid values: `HTTP1.0`, `HTTP1.1`. **NOTE:** `health_check_http_version` takes effect only if `health_check_protocol` is set to `HTTP`.
* `health_check_interval` - (Optional, Int) The interval at which health checks are performed. Unit: seconds. Default value: `2`. Valid values: `1` to `50`.
* `health_check_method` - (Optional) The HTTP method that is used for health checks. Default value: `HEAD`. Valid values: `HEAD`, `GET`. **NOTE:** `health_check_method` takes effect only if `health_check_protocol` is set to `HTTP`.
* `health_check_path` - (Optional) The URL that is used for health checks. **NOTE:** `health_check_path` takes effect only if `health_check_protocol` is set to `HTTP`.
* `health_check_protocol` - (Optional) The protocol that is used for health checks. Default value: `HTTP`. Valid values: `HTTP`, `TCP`.
* `health_check_timeout` - (Optional, Int) The timeout period of a health check. Default value: `5`. Valid values: `1` to `300`.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`.
* `health_check_codes` - (Optional, List) The HTTP status codes that are used to indicate whether the backend server passes the health check. Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. **NOTE:** `health_check_codes` takes effect only if `health_check_protocol` is set to `HTTP`.
* `dry_run` - (Optional, Bool) Whether to precheck the API request.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Health Check Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Health Check Template.
* `update` - (Defaults to 2 mins) Used when update the Health Check Template.
* `delete` - (Defaults to 2 mins) Used when delete the Health Check Template.

## Import

Application Load Balancer (ALB) Health Check Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_health_check_template.example <id>
```
