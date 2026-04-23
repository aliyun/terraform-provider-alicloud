---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_address"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Address resource.
---

# alicloud_alidns_cloud_gtm_address

Provides a Alidns Cloud Gtm Address resource.

A Cloud GTM address represents an individual service endpoint (IPv4, IPv6, or domain name) that can be grouped into a Cloud GTM address pool. Each address carries its own health-check configuration — one or more probe tasks that reference Cloud GTM monitor templates — so that Cloud GTM can determine whether the endpoint is available before returning it in DNS responses.

For information about Alidns Cloud Gtm Address and how to use it, see [What is Cloud Gtm Address](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmAddress).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "tcp" {
  name             = "${var.name}-tcp"
  protocol         = "tcp"
  ip_version       = "IPv4"
  interval         = "60"
  timeout          = "3000"
  evaluation_count = 1
  failure_rate     = 50
  extend_info      = "{}"

  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "https" {
  name             = "${var.name}-https"
  protocol         = "https"
  ip_version       = "IPv4"
  interval         = "60"
  timeout          = "2000"
  evaluation_count = 1
  failure_rate     = 50
  extend_info = jsonencode({
    code           = 400
    followRedirect = true
    path           = "/"
    sni            = false
  })

  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "ping" {
  name             = "${var.name}-ping"
  protocol         = "ping"
  ip_version       = "IPv4"
  interval         = "60"
  timeout          = "3000"
  evaluation_count = 1
  failure_rate     = 50
  extend_info = jsonencode({
    packetNum      = 20
    packetLossRate = 10
  })

  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
}

resource "alicloud_alidns_cloud_gtm_address" "default" {
  name                    = var.name
  type                    = "IPv4"
  address                 = "1.1.1.1"
  enable_status           = "enable"
  available_mode          = "manual"
  manual_available_status = "available"
  health_judgement        = "all_ok"
  remark                  = "terraform-example-remark"

  health_tasks {
    template_id = alicloud_alidns_cloud_gtm_monitor_template.ping.id
  }

  health_tasks {
    port        = 53
    template_id = alicloud_alidns_cloud_gtm_monitor_template.tcp.id
  }

  health_tasks {
    port        = 443
    template_id = alicloud_alidns_cloud_gtm_monitor_template.https.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the address. Used to identify the address in the Cloud GTM console.
* `type` - (Required, ForceNew) The address type. Valid values: `IPv4`, `IPv6`, `domain`.
* `address` - (Required) The address value. Must match `type`: an IPv4 address, an IPv6 address, or a domain name.
* `enable_status` - (Required) Whether the address participates in DNS resolution. Valid values:
    * `enable` - The address participates in resolution when its health check is normal.
    * `disable` - The address does not participate in resolution regardless of health status.
* `available_mode` - (Required) How the availability of the address is determined. Valid values:
    * `auto` - Availability is computed from the attached health-check tasks.
    * `manual` - Availability is set explicitly via `manual_available_status`; health-check results are informational only.
* `manual_available_status` - (Optional) The manually-set availability status. Only meaningful when `available_mode` is `manual`. Valid values: `available`, `unavailable`.
* `health_judgement` - (Required) The rule used to judge overall health when the address has multiple health-check tasks. Valid values:
    * `any_ok` - Any task reports healthy.
    * `all_ok` - All tasks report healthy.
    * `p30_ok` - At least 30% of tasks report healthy.
    * `p50_ok` - At least 50% of tasks report healthy.
    * `p70_ok` - At least 70% of tasks report healthy.
* `health_tasks` - (Optional, Set) The health-check tasks attached to this address. Each task references a Cloud GTM monitor template. See [`health_tasks`](#health_tasks) below.
* `remark` - (Optional) The remark of the address. Passing an empty value clears the existing remark.

### `health_tasks`

The `health_tasks` block supports:

* `template_id` - (Optional) The ID of the Cloud GTM monitor template to probe this address with. Usually referenced as `alicloud_alidns_cloud_gtm_monitor_template.<name>.id`.
* `port` - (Optional, Int) The port to probe. If omitted, the default port of the template's protocol is used (for example, 80 for HTTP, 443 for HTTPS).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the address. It is the same as the `address_id` returned by the API.
* `create_time` - The creation time of the address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Cloud Gtm Address.
* `update` - (Defaults to 5 mins) Used when updating the Cloud Gtm Address.
* `delete` - (Defaults to 5 mins) Used when deleting the Cloud Gtm Address.

## Import

Alidns Cloud Gtm Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_address.example <address_id>
```
