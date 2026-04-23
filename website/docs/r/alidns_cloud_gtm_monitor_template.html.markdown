---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_monitor_template"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Monitor Template resource.
---

# alicloud_alidns_cloud_gtm_monitor_template

Provides a Alidns Cloud Gtm Monitor Template resource.

A Cloud GTM monitor template defines reusable health-check probe configurations (protocol, interval, failure thresholds, probe nodes) that can be attached to Cloud GTM address pools.

For information about Alidns Cloud Gtm Monitor Template and how to use it, see [What is Cloud Gtm Monitor Template](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmMonitorTemplate).

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

resource "alicloud_alidns_cloud_gtm_monitor_template" "default" {
  name             = var.name
  protocol         = "http"
  ip_version       = "IPv4"
  interval         = "60"
  timeout          = "2000"
  evaluation_count = 2
  failure_rate     = 50
  extend_info = jsonencode({
    code           = 500
    followRedirect = true
    path           = "/"
  })
  remark = "terraform-example-remark"

  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }

  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the monitor template. It is recommended to use a name that reflects the health-check protocol for easier identification.
* `protocol` - (Required, ForceNew) The probing protocol of the template. Valid values: `ping`, `tcp`, `http`, `https`.
* `ip_version` - (Required, ForceNew) The IP version of the probing node. Valid values: `IPv4`, `IPv6`.
* `interval` - (Required) The interval between consecutive probes, in seconds. Valid values: `15`, `60`, `300`, `900`, `1800`, `3600`. The `15` seconds interval is only available for Flagship Edition instances.
* `timeout` - (Required) Probe request timeout, in milliseconds. Probe packets that do not return within this duration are treated as timeouts. Valid values: `2000`, `3000`, `5000`, `10000`.
* `evaluation_count` - (Required, Int) The number of retries after a probe failure. A service is marked abnormal only after this many consecutive failures, preventing transient network fluctuations from triggering false alarms. Valid values: `0`, `1`, `2`, `3`.
* `failure_rate` - (Required, Int) The failure-rate threshold (%) among selected probe nodes. If the percentage of failing nodes exceeds this value, the service address is marked as abnormal. Valid values: `0`, `20`, `50`, `80`, `100`.
* `isp_city_nodes` - (Required, Set) The set of monitoring nodes that this template will probe from. Use the [ListCloudGtmMonitorNodes](https://help.aliyun.com/document_detail/2797349.html) API to look up available `city_code` / `isp_code` combinations. See [`isp_city_nodes`](#isp_city_nodes) below.
* `extend_info` - (Optional, Computed) A JSON string containing protocol-specific probe configuration. The supported keys depend on `protocol`. See [`extend_info`](#extend_info) below.
* `remark` - (Optional) The remark of the monitor template. Passing an empty value clears the existing remark.

### `isp_city_nodes`

The `isp_city_nodes` block supports:

* `city_code` - (Optional) The city code of the monitoring node.
* `isp_code` - (Optional) The ISP (Internet Service Provider) code of the monitoring node.

### `extend_info`

The `extend_info` argument takes a JSON-encoded string. The supported keys depend on `protocol`.

Keys for `http` / `https`:

| Key              | Description                                                                                                                                                                                                                         |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `host`           | The value of the `Host` header carried in the HTTP(S) request. Defaults to the primary domain of the target. Set this when the target site requires a specific Host header.                                                        |
| `path`           | The URL path to probe. Defaults to `/`.                                                                                                                                                                                             |
| `code`           | HTTP response-code threshold used to classify the service as abnormal. Valid values: `400` (Bad Request — use when the probe URL includes parameters and you want to detect invalid-request responses), `500` (Server Error, default). |
| `sni`            | Whether to enable Server Name Indication during the TLS handshake (HTTPS only). Valid values: `true`, `false`.                                                                                                                      |
| `followRedirect` | Whether to follow HTTP 3XX redirects (`301`, `302`, `303`, `307`, `308`). Valid values: `true`, `false`.                                                                                                                             |

Keys for `ping`:

| Key              | Description                                                                                                                                                             |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `packetNum`      | The number of ICMP packets sent in each probe cycle. Valid values: `20`, `50`, `100`.                                                                                   |
| `packetLossRate` | The packet-loss rate (%) that triggers an alert. Computed as `lost_packets / total_packets * 100`. Valid values: `10`, `30`, `40`, `80`, `90`, `100`.                   |

No additional keys are required for `tcp` probes.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the monitor template. It is the same as the `template_id` returned by the API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Cloud Gtm Monitor Template.
* `update` - (Defaults to 5 mins) Used when updating the Cloud Gtm Monitor Template.
* `delete` - (Defaults to 5 mins) Used when deleting the Cloud Gtm Monitor Template.

## Import

Alidns Cloud Gtm Monitor Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_monitor_template.example <template_id>
```
