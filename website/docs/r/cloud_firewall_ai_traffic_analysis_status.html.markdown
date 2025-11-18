---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_ai_traffic_analysis_status"
description: |-
  Provides a Alicloud Cloud Firewall Ai Traffic Analysis Status resource.
---

# alicloud_cloud_firewall_ai_traffic_analysis_status

Provides a Cloud Firewall Ai Traffic Analysis Status resource.

AI traffic analysis.

For information about Cloud Firewall Ai Traffic Analysis Status and how to use it, see [What is Ai Traffic Analysis Status](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/UpdateAITrafficAnalysisStatus).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_firewall_ai_traffic_analysis_status" "default" {
  status = "Open"
}
```

### Deleting `alicloud_cloud_firewall_ai_traffic_analysis_status` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_ai_traffic_analysis_status`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `status` - (Optional, Computed) Status

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ai Traffic Analysis Status.
* `update` - (Defaults to 5 mins) Used when update the Ai Traffic Analysis Status.

## Import

Cloud Firewall Ai Traffic Analysis Status can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_ai_traffic_analysis_status.example 
```