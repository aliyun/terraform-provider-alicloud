---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_hd_monitor_region_config"
description: |-
  Provides a Alicloud Network Load Balancer (NLB) Hd Monitor Region Config resource.
---

# alicloud_nlb_hd_monitor_region_config

Provides a Network Load Balancer (NLB) Hd Monitor Region Config resource.

HD monitor config.

For information about Network Load Balancer (NLB) Hd Monitor Region Config and how to use it, see [What is Hd Monitor Region Config](https://next.api.alibabacloud.com/document/Nlb/2022-04-30/SetHdMonitorRegionConfig).

-> **NOTE:** Available since v1.273.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_hd_monitor_region_config&exampleId=056bbd2f-6ca3-1628-474a-8b001701e1289bcdc623&activeTab=example&spm=docs.r.nlb_hd_monitor_region_config.0.056bbd2f6c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}


resource "alicloud_nlb_hd_monitor_region_config" "default" {
  metric_store = "example"
  log_project  = "example"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nlb_hd_monitor_region_config&spm=docs.r.nlb_hd_monitor_region_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `log_project` - (Required) The name of the LogProject.
* `metric_store` - (Required) The name of the MetricStore.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `region_id` - The ID of the region in which the resource resides.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hd Monitor Region Config.
* `delete` - (Defaults to 5 mins) Used when delete the Hd Monitor Region Config.
* `update` - (Defaults to 5 mins) Used when update the Hd Monitor Region Config.

## Import

Network Load Balancer (NLB) Hd Monitor Region Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_hd_monitor_region_config.example <region_id>
```