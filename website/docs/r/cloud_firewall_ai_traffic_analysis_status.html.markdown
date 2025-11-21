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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_ai_traffic_analysis_status&exampleId=0fbf97b4-adde-c3f8-8f37-5d9c88d4a13639525fa3&activeTab=example&spm=docs.r.cloud_firewall_ai_traffic_analysis_status.0.0fbf97b4ad&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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