---
subcategory: "Anycast Elastic IP Address (Eipanycast)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eipanycast_anycast_eip_address_attachment"
description: |-
  Provides a Alicloud Anycast Eip Address Attachment resource.
---

# alicloud_eipanycast_anycast_eip_address_attachment

Provides a Eipanycast Anycast Eip Address Attachment resource.

For information about Eipanycast Anycast Eip Address Attachment and how to use it, see [What is Anycast Eip Address Attachment](https://www.alibabacloud.com/help/en/anycast-eip/latest/api-eipanycast-2020-03-09-associateanycasteipaddress).

-> **NOTE:** Available since v1.113.0.

-> **NOTE:** The following regions support currently while Slb instance support bound. 
[eu-west-1-gb33-a01,cn-hongkong-am4-c04,ap-southeast-os30-a01,us-west-ot7-a01,ap-south-in73-a01,ap-southeast-my88-a01]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eipanycast_anycast_eip_address_attachment&exampleId=60808d6d-90eb-3d4f-6291-9cc92150dba73d9e2abf&activeTab=example&spm=docs.r.eipanycast_anycast_eip_address_attachment.0.60808d6d90&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_slb_zones" "default" {
  available_slb_address_type = "vpc"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
  address_type       = "intranet"
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_name = var.name
  load_balancer_spec = "slb.s1.small"
  master_zone_id     = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_eipanycast_anycast_eip_address" "default" {
  anycast_eip_address_name = var.name
  service_location         = "ChineseMainland"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "default" {
  bind_instance_id        = alicloud_slb_load_balancer.default.id
  bind_instance_type      = "SlbInstance"
  bind_instance_region_id = data.alicloud_regions.default.regions.0.id
  anycast_id              = alicloud_eipanycast_anycast_eip_address.default.id
}
```

Multiple Usage

-> **NOTE:**  Anycast EIP supports binding cloud resource instances in multiple regions. Only one cloud resource instance is supported as the default origin station, and the rest are normal origin stations. When no access point is specified or an access point is added, the access request is forwarded to the default origin by default.  If you are bound for the first time, the Default value of the binding mode is **Default * *. /li> li> If you are not binding for the first time, you can set the binding mode to **Default**, and the new Default origin will take effect. The original Default origin will be changed to a common origin.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eipanycast_anycast_eip_address_attachment&exampleId=aa4fd8d3-f48b-e412-01d0-81401540669093b77a06&activeTab=example&spm=docs.r.eipanycast_anycast_eip_address_attachment.1.aa4fd8d3f4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  alias  = "beijing"
  region = "cn-beijing"
}

provider "alicloud" {
  alias  = "hangzhou"
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  provider                    = alicloud.beijing
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  provider    = alicloud.beijing
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  provider          = alicloud.beijing
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "defaultVpc" {
  provider   = alicloud.beijing
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVsw" {
  provider   = alicloud.beijing
  vpc_id     = alicloud_vpc.defaultVpc.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "defaultuBsECI" {
  provider = alicloud.beijing
  vpc_id   = alicloud_vpc.defaultVpc.id
}

resource "alicloud_instance" "default9KDlN7" {
  provider             = alicloud.beijing
  image_id             = data.alicloud_images.default.images[0].id
  instance_type        = data.alicloud_instance_types.default.instance_types[0].id
  instance_name        = var.name
  security_groups      = ["${alicloud_security_group.defaultuBsECI.id}"]
  availability_zone    = alicloud_vswitch.defaultVsw.zone_id
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id           = alicloud_vswitch.defaultVsw.id
}

resource "alicloud_eipanycast_anycast_eip_address" "defaultXkpFRs" {
  provider         = alicloud.hangzhou
  service_location = "ChineseMainland"
}

resource "alicloud_vpc" "defaultVpc2" {
  provider   = alicloud.hangzhou
  vpc_name   = "${var.name}6"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "default2" {
  provider                    = alicloud.hangzhou
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default2" {
  provider    = alicloud.hangzhou
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default2" {
  provider          = alicloud.hangzhou
  availability_zone = data.alicloud_zones.default2.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vswitch" "defaultdsVsw2" {
  provider   = alicloud.hangzhou
  vpc_id     = alicloud_vpc.defaultVpc2.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default2.zones.1.id
}

resource "alicloud_security_group" "defaultuBsECI2" {
  provider = alicloud.hangzhou
  vpc_id   = alicloud_vpc.defaultVpc2.id
}

resource "alicloud_instance" "defaultEcs2" {
  provider             = alicloud.hangzhou
  image_id             = data.alicloud_images.default2.images[0].id
  instance_type        = data.alicloud_instance_types.default2.instance_types[0].id
  instance_name        = var.name
  security_groups      = ["${alicloud_security_group.defaultuBsECI2.id}"]
  availability_zone    = alicloud_vswitch.defaultdsVsw2.zone_id
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id           = alicloud_vswitch.defaultdsVsw2.id
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "defaultEfYBJY" {
  provider                = alicloud.beijing
  bind_instance_id        = alicloud_instance.default9KDlN7.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-beijing"
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id
  association_mode        = "Default"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "normal" {
  provider                = alicloud.hangzhou
  bind_instance_id        = alicloud_instance.defaultEcs2.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-hangzhou"
  anycast_id              = alicloud_eipanycast_anycast_eip_address_attachment.defaultEfYBJY.anycast_id
}
```

## Argument Reference

The following arguments are supported:
* `anycast_id` - (Required, ForceNew, Available since v1.113.0) The ID of the Anycast EIP instance.
* `association_mode` - (Optional) Binding mode, value:
  - **Default**: The Default mode. The cloud resource instance to be bound is set as the Default origin.
  - **Normal**: In Normal mode, the cloud resource instance to be bound is set to the common source station.
* `bind_instance_id` - (Required, ForceNew, Available since v1.113.0) The ID of the cloud resource instance to be bound.
* `bind_instance_region_id` - (Required, ForceNew, Available since v1.113.0) The region ID of the cloud resource instance to be bound.You can only bind cloud resource instances in some regions. You can call the [describeanystserverregions](~~ 171939 ~~) operation to obtain the region ID of the cloud resource instances that can be bound.
* `bind_instance_type` - (Required, ForceNew, Available since v1.113.0) The type of the cloud resource instance to be bound. Value:
  - **SlbInstance**: a private network SLB instance.
  - **NetworkInterface**: ENI.
* `pop_locations` - (Optional) The access point information of the associated access area when the cloud resource instance is bound.If you are binding for the first time, this parameter does not need to be configured, and the system automatically associates all access areas. See [`pop_locations`](#pop_locations) below.
* `private_ip_address` - (Optional, ForceNew) The secondary private IP address of the elastic network card to be bound.This parameter takes effect only when **BindInstanceType** is set to **NetworkInterface. When you do not enter, this parameter is the primary private IP of the ENI by default.

### `pop_locations`

The pop_locations supports the following:
* `pop_location` - (Optional) The access point information of the associated access area when the cloud resource instance is bound.If you are binding for the first time, this parameter does not need to be configured, and the system automatically associates all access areas.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<anycast_id>:<bind_instance_id>:<bind_instance_region_id>:<bind_instance_type>`.
* `bind_time` - Binding time.Time is expressed according to ISO8601 standard and UTC time is used. The format is: 'YYYY-MM-DDThh:mm:ssZ'.
* `status` - The status of the bound cloud resource instance. Value:BINDING: BINDING.Bound: Bound.UNBINDING: UNBINDING.DELETED: DELETED.MODIFYING: being modified.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anycast Eip Address Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Anycast Eip Address Attachment.
* `update` - (Defaults to 5 mins) Used when update the Anycast Eip Address Attachment.

## Import

Eipanycast Anycast Eip Address Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eipanycast_anycast_eip_address_attachment.example <anycast_id>:<bind_instance_id>:<bind_instance_region_id>:<bind_instance_type>
```