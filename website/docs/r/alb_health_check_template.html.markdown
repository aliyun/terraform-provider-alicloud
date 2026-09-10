---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_health_check_template"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Health Check Template resource.
---

# alicloud_alb_health_check_template

Provides a Application Load Balancer (ALB) Health Check Template resource.

Health check template.

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alb_health_check_template&spm=docs.r.alb_health_check_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Whether to PreCheck only this request, value:
true: sends a check request and does not create a resource. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code DryRunOperation is returned.
false (default): Sends a normal request, returns the http_2xx status code after the check, and directly performs the operation.
* `health_check_codes` - (Optional, Computed, List) The HTTP code of the health check. The default value is http_2xx. The normal HTTP code for health check. Separate multiple codes with commas (,). Valid values: http_2xx, http_3xx, http_4xx, or http_5xx.
* `health_check_connect_port` - (Optional, Computed, Int) The number of the port that is used for health checks.  Valid values: 0 to 65535.  Default value: 0. This value indicates that the backend server is used for health checks.
* `health_check_host` - (Optional, Computed) The domain name that is used for health checks. Valid values:  $SERVER_IP (default value): The private IP addresses of backend servers. If the $_ip parameter is set or the HealthCheckHost parameter is not set, SLB uses the private IP addresses of backend servers as the domain names for health checks.  domain: The domain name must be 1 to 80 characters in length, and can contain only letters, digits, periods (.),and hyphens (-).
* `health_check_http_version` - (Optional, Computed) The version of the HTTP protocol.  Valid values: HTTP 1.0 and HTTP 1.1.  Default value: HTTP 1.1.
* `health_check_interval` - (Optional, Computed, Int) The time interval between two consecutive health checks.  Valid values: 1 to 50. Unit: seconds.  Default value: 2.
* `health_check_method` - (Optional, Computed) The health check method.  Valid values: GET and HEAD.  Default value: HEAD.
* `health_check_path` - (Optional, Computed) The URL that is used for health checks.  The URL must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URL can also contain the following extended characters: _ ; ~ ! ( )* [ ] @ $ ^ : ' , +. The URL must start with a forward slash (/).
* `health_check_protocol` - (Optional, Computed) The protocol used for the health check. Value:
HTTP (default): Sends a HEAD or GET request to simulate the browser's access behavior to check whether the server application is healthy.
HTTPS: Sends a HEAD or GET request to simulate the browser's access behavior to check whether the server application is healthy. (Data encryption is more secure than HTTP.)
TCP: Sends a SYN handshake packet to check whether the server port is alive.
gRPC: Check whether the server application is healthy by sending a POST or GET request.
* `health_check_template_name` - (Required) The name of the health check template.  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `health_check_timeout` - (Optional, Computed, Int) The timeout period of a health check response. If the backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the health check fails.  Valid values: 1 to 300. Unit: seconds.  Default value: 5.
* `healthy_threshold` - (Optional, Computed, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy (from fail to success).
* `resource_group_id` - (Optional, ForceNew, Computed, Available since v1.249.0) The ID of the resource group
* `tags` - (Optional, Map, Available since v1.249.0) The tag of the resource
* `unhealthy_threshold` - (Optional, Computed, Int) Specifies the number of times that an healthy backend server must consecutively fail health checks before it is declared unhealthy (from success to fail).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Health Check Template.
* `delete` - (Defaults to 5 mins) Used when delete the Health Check Template.
* `update` - (Defaults to 5 mins) Used when update the Health Check Template.

## Import

Application Load Balancer (ALB) Health Check Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_health_check_template.example <id>
```