---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_monitor_template"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Monitor Template resource.
---

# alicloud_alidns_cloud_gtm_monitor_template

Provides a Alidns Cloud Gtm Monitor Template resource.

CloudGtm probing template.

For information about Alidns Cloud Gtm Monitor Template and how to use it, see [What is Cloud Gtm Monitor Template](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmMonitorTemplate).

-> **NOTE:** Available since v1.275.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_alidns_cloud_gtm_monitor_template" "default" {
  ip_version = "IPv4"
  timeout    = "2000"
  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
  evaluation_count = "2"
  protocol         = "http"
  failure_rate     = "50"
  extend_info      = "{\"code\":500,\"followRedirect\":true,\"path\":\"/\"}"
  name             = "template-example-3"
  interval         = "60"
}
```

## Argument Reference

The following arguments are supported:
* `evaluation_count` - (Required, Int) Number of automatic retries after a probe failure
* `extend_info` - (Optional) Different probing protocols require different extended information
* `failure_rate` - (Required, Int) Probe failure rate
* `interval` - (Required) The time interval between probes
* `ip_version` - (Required, ForceNew) IP version of the template
* `isp_city_nodes` - (Required, List) Probe nodes See [`isp_city_nodes`](#isp_city_nodes) below.
* `name` - (Required) Resource property field representing the resource name
* `protocol` - (Required, ForceNew) The probing protocol of this template
* `remark` - (Optional) Remarks for this template
* `timeout` - (Required) Timeout duration for probe requests

### `isp_city_nodes`

The isp_city_nodes supports the following:
* `city_code` - (Optional) City code
* `isp_code` - (Optional) ISP code

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cloud Gtm Monitor Template.
* `delete` - (Defaults to 5 mins) Used when delete the Cloud Gtm Monitor Template.
* `update` - (Defaults to 5 mins) Used when update the Cloud Gtm Monitor Template.

## Import

Alidns Cloud Gtm Monitor Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_monitor_template.example <template_id>
```