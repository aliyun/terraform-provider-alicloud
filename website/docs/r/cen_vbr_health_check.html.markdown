---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_vbr_health_check"
sidebar_current: "docs-alicloud-resource-cen-vbr-health-check"
description: |-
  Provides a Alicloud CEN VBR HealthCheck resource.
---

# alicloud_cen_vbr_health_check

This topic describes how to configure the health check feature for a Cloud Enterprise Network (CEN) instance. 
After you attach a Virtual Border Router (VBR) to the CEN instance and configure the health check feature, you can monitor the network conditions of the on-premises data center connected to the VBR.

For information about CEN VBR HealthCheck and how to use it, see [Manage CEN VBR HealthCheck](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-enablecenvbrhealthcheck).

-> **NOTE:** Available since v1.88.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_vbr_health_check&exampleId=50c05159-1c42-f7c0-633f-47793d90f6cdd2882937&activeTab=example&spm=docs.r.cen_vbr_health_check.0.50c051591c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}
resource "alicloud_cen_instance_attachment" "example" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.example.id
  child_instance_type      = "VBR"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}
resource "alicloud_cen_vbr_health_check" "example" {
  cen_id                 = alicloud_cen_instance.example.id
  health_check_source_ip = "192.168.1.2"
  health_check_target_ip = "10.0.0.2"
  vbr_instance_id        = alicloud_express_connect_virtual_border_router.example.id
  vbr_instance_region_id = alicloud_cen_instance_attachment.example.child_instance_region_id
  health_check_interval  = 2
  healthy_threshold      = 8
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_vbr_health_check&spm=docs.r.cen_vbr_health_check.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `health_check_interval` - (Optional, Default) Specifies the interval at which the health check sends continuous detection packets. Default value: 2. Value range: 2 to 3.
* `health_check_source_ip` - (Optional) The source IP address of health checks.
* `health_check_target_ip` - (Required) The destination IP address of health checks.
* `healthy_threshold` - (Optional, Default) Specifies the number of probe messages sent by the health check. Default value: 8. Value range: 3 to 8.
* `vbr_instance_id` - (Required, ForceNew) The ID of the VBR.
* `vbr_instance_owner_id` - (Optional) The ID of the account to which the VBR belongs.
* `vbr_instance_region_id` - (Required, ForceNew) The ID of the region to which the VBR belongs.

->**NOTE:** The `alicloud_cen_vbr_health_check` resource depends on the related `alicloud_cen_instance_attachment` resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, formatted as `<vbr_instance_id>:<vbr_instance_region_id>`.

## Timeouts

-> **NOTE:** Available since v1.98.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the CEN VBR Health Check. (until it reaches the available status).
* `update` - (Defaults to 6 mins) Used when update the CEN VBR Health Check.
* `delete` - (Defaults to 6 mins) Used when delete the CEN VBR Health Check.

## Import

CEN VBR HealthCheck can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_vbr_health_check.example vbr-xxxxx:cn-hangzhou
```
