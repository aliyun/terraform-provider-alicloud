---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Cen Tr Firewall resource.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall

Provides a Cloud Firewall Vpc Cen Tr Firewall resource.

VPC Firewall for Cloud Enterprise Network Enterprise Edition.

For information about Cloud Firewall Vpc Cen Tr Firewall and how to use it, see [What is Vpc Cen Tr Firewall](https://www.alibabacloud.com/help/en/cloud-firewall/cloudfirewall/developer-reference/api-cloudfw-2017-12-07-createtrfirewallv2).

-> **NOTE:** Available since v1.228.0.

## Example Usage

### Managed mode (automatic)

In managed mode, Cloud Firewall creates and manages the firewall VPC. Only `firewall_vpc_cidr` is required.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_cen_tr_firewall&exampleId=37a1e48f-5e2b-3aa6-5c58-4d9edb62f71ef8f83c50&activeTab=example&spm=docs.r.cloud_firewall_vpc_cen_tr_firewall.0.37a1e48f5e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shenzhen"
}

variable "firewall_name" {
  default = "tf-example"
}

variable "region" {
  default = "cn-shenzhen"
}

variable "firewall_vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "zone1" {
  default = "cn-shenzhen-d"
}

variable "zone2" {
  default = "cn-shenzhen-e"
}

resource "alicloud_cen_instance" "cen" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = var.zone1
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = var.zone2
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  auto_publish_route_enabled = false
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  cen_id            = alicloud_cen_instance.cen.id
  firewall_name     = var.firewall_name
  region_no         = var.region
  firewall_vpc_cidr = var.firewall_vpc_cidr
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
  route_mode        = "managed"
}
```

### Manual mode

In manual mode, the firewall ENI is deployed in an existing VPC. `firewall_vpc_id`, `firewall_vswitch_id`, `tr_attachment_master_zone` and `tr_attachment_slave_zone` are required, and `firewall_vpc_cidr` is not used. The example below reuses the CEN, TR, VPC, vSwitch and TR VPC attachment resources from the managed example above.

```terraform
resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "manual" {
  cen_id                    = alicloud_cen_instance.cen.id
  firewall_name             = var.firewall_name
  region_no                 = var.region
  firewall_vpc_id           = alicloud_vpc.vpc1.id
  firewall_vswitch_id       = alicloud_vswitch.vpc1vsw1.id
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  tr_attachment_master_zone = var.zone1
  tr_attachment_slave_zone  = var.zone2
  route_mode                = "manual"
}
```

### Semi-managed mode (手动模式托管VPC)

In semi-managed mode, Cloud Firewall creates and manages the firewall VPC, just like managed mode. Only `firewall_vpc_cidr` is required; `firewall_vpc_id` and `firewall_vswitch_id` are auto-assigned by Cloud Firewall and exposed in state. The difference from managed mode is the routing mode flag passed to the service. The example below reuses the CEN and TR resources from the managed example above.

```terraform
resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "semimanaged" {
  cen_id            = alicloud_cen_instance.cen.id
  firewall_name     = var.firewall_name
  region_no         = var.region
  firewall_vpc_cidr = var.firewall_vpc_cidr
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
  route_mode        = "semi-managed"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_vpc_cen_tr_firewall&spm=docs.r.cloud_firewall_vpc_cen_tr_firewall.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `firewall_description` - (Optional, ForceNew) The description of the firewall.
* `firewall_name` - (Required) The name of the Cloud Firewall instance.
* `firewall_subnet_cidr` - (Optional, String, ForceNew) The CIDR block of the subnet in the firewall VPC in automatic (managed) mode.
* `firewall_vpc_cidr` - (Optional, String, ForceNew) The CIDR block of the firewall VPC in automatic (managed) mode. Required when `route_mode` is `managed` or `semi-managed`.
* `firewall_vpc_id` - (Optional, String, Computed, ForceNew, Available since v2.2.0) The ID of the VPC used to create the firewall ENI in manual mode. Required when `route_mode` is `manual`. In `managed` and `semi-managed` modes this attribute is auto-assigned by Cloud Firewall and exposed in state; omit it from the configuration.

  -> **NOTE:** The VPC specified by `firewall_vpc_id` must support subnet routing. If it does not, create a custom route table for that VPC to enable subnet routing before creating the firewall.
* `firewall_vswitch_id` - (Optional, String, Computed, ForceNew, Available since v2.2.0) The ID of the vSwitch used to create the firewall ENI in manual mode. Required when `route_mode` is `manual`. In `managed` and `semi-managed` modes this attribute is auto-assigned by Cloud Firewall and exposed in state; omit it from the configuration.
* `region_no` - (Required, ForceNew) The region ID of transit router instance I.
* `route_mode` - (Required, String, ForceNew) The routing mode. Valid values:
  - `managed`: indicates automatic mode. Cloud Firewall creates and manages the firewall VPC; `firewall_vpc_cidr` is required.
  - `semi-managed`: indicates semi-managed mode (手动模式托管VPC). Cloud Firewall still creates and manages the firewall VPC; `firewall_vpc_cidr` is required, and `firewall_vpc_id`/`firewall_vswitch_id` are auto-assigned.
  - `manual`: indicates manual mode. The firewall ENI is deployed in an existing VPC; `firewall_vpc_id`, `firewall_vswitch_id`, `tr_attachment_master_zone` and `tr_attachment_slave_zone` are required.

  -> **NOTE:** If this parameter is not specified, VPC border firewalls in all routing modes are queried.

* `tr_attachment_master_cidr` - (Optional, String, ForceNew) (Deprecated) The primary CIDR block of the subnet in the firewall VPC used to connect to the transit router (TR) in automatic mode.
* `tr_attachment_master_zone` - (Optional, String) The primary zone of the vSwitch. Required when `route_mode` is `manual`.

  -> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `tr_attachment_slave_cidr` - (Optional, String, ForceNew) (Deprecated) The secondary CIDR block of the subnet in the firewall VPC used to connect to the transit router (TR) in automatic mode.
* `tr_attachment_slave_zone` - (Optional, String) The secondary zone of the vSwitch. Required when `route_mode` is `manual`.

  -> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `transit_router_id` - (Required, ForceNew) The ID of the transit router instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `firewall_eni_id` - The ID of the firewall ENI.
* `firewall_eni_vpc_id` - The ID of the VPC where the firewall ENI resides.
* `firewall_vpc_attachment_id` - The attachment ID used in the firewall VPC to connect to the transit router (TR) in automatic mode.
* `status` - The status of the firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 41 mins) Used when create the Vpc Cen Tr Firewall.
* `delete` - (Defaults to 46 mins) Used when delete the Vpc Cen Tr Firewall.
* `update` - (Defaults to 5 mins) Used when update the Vpc Cen Tr Firewall.

## Import

Cloud Firewall Vpc Cen Tr Firewall can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_cen_tr_firewall.example <firewall_id>
```
