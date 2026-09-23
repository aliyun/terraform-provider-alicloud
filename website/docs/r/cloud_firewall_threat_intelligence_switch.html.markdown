---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_threat_intelligence_switch"
description: |-
  Provides a Alicloud Cloud Firewall Threat Intelligence Switch resource.
---

# alicloud_cloud_firewall_threat_intelligence_switch

Provides a Cloud Firewall Threat Intelligence Switch resource.

Cloud Firewall Switch Threat Intelligence.

For information about Cloud Firewall Threat Intelligence Switch and how to use it, see [What is Threat Intelligence Switch](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyThreatIntelligenceSwitch).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_threat_intelligence_switch&exampleId=95ce3730-53af-abda-6093-2f86736e8fba622adfd9&activeTab=example&spm=docs.r.cloud_firewall_threat_intelligence_switch.0.95ce373053&intl_lang=EN_US" target="_blank">
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


resource "alicloud_cloud_firewall_threat_intelligence_switch" "default" {
  action        = "alert"
  enable_status = "0"
  category_id   = "IpOutThreatTorExit"
}
```

### Deleting `alicloud_cloud_firewall_threat_intelligence_switch` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_threat_intelligence_switch`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_threat_intelligence_switch&spm=docs.r.cloud_firewall_threat_intelligence_switch.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `action` - (Optional) Rule action. Value:
  - `alert`: Watch
  - `drop`: Intercept
* `category_id` - (Optional, ForceNew, Computed) The threat intelligence classification ID.
* `enable_status` - (Optional, Int) Switch status. Value:
  - `1`: On
  - `0`: Closed

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Threat Intelligence Switch.
* `update` - (Defaults to 5 mins) Used when update the Threat Intelligence Switch.

## Import

Cloud Firewall Threat Intelligence Switch can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_threat_intelligence_switch.example <id>
```