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

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultlnZlvA" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultdsRCp1" {
  vpc_id     = alicloud_vpc.defaultlnZlvA.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "defaultuBsECI" {
  vpc_id = alicloud_vpc.defaultlnZlvA.id
}

resource "alicloud_ecs_instance" "default9KDlN7" {
  system_disk {
    category = "cloud_efficiency"
  }
  image_id = "aliyun_2_1903_x64_20G_alibase_20230308.vhd"
  vpc_attributes {
    vswitch_id = alicloud_vswitch.defaultdsRCp1.id
  }
  payment_type       = "PayAsYouGo"
  instance_type      = "ecs.g5ne.xlarge"
  spot_strategy      = "NoSpot"
  zone_id            = alicloud_vswitch.defaultdsRCp1.zone_id
  security_group_ids = ["${alicloud_security_group.defaultuBsECI.id}"]
}

resource "alicloud_eipanycast_anycast_eip_address" "defaultXkpFRs" {
  service_location = "ChineseMainland"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "defaultEfYBJY" {
  bind_instance_id        = alicloud_ecs_instance.default9KDlN7.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = alicloud_ecs_instance.default9KDlN7.region_id
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.anycast_id
  association_mode        = "Default"
}

resource "alicloud_vpc" "defaultVpc2" {
  vpc_name   = "${var.name}6"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultdsVsw2" {
  vpc_id     = alicloud_vpc.defaultVpc2.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_security_group" "defaultuBsECI2" {
  vpc_id = alicloud_vpc.defaultVpc2.id
}

resource "alicloud_ecs_instance" "defaultEcs2" {
  system_disk {
    category = "cloud_efficiency"
  }
  image_id = "aliyun_2_1903_x64_20G_alibase_20230308.vhd"
  vpc_attributes {
    vswitch_id = alicloud_vswitch.defaultdsVsw2.id
  }
  payment_type       = "PayAsYouGo"
  instance_type      = "ecs.g5ne.xlarge"
  spot_strategy      = "NoSpot"
  zone_id            = alicloud_vswitch.defaultdsVsw2.zone_id
  security_group_ids = ["${alicloud_security_group.defaultuBsECI2.id}"]
}


resource "alicloud_eipanycast_anycast_eip_address_attachment" "default" {
  bind_instance_id   = alicloud_ecs_instance.defaultEcs2.network_interface_id
  bind_instance_type = "NetworkInterface"
  pop_locations {
    pop_location = "cn-guangzhou-pop"
  }
  pop_locations {
    pop_location = "cn-shanghai-pop"
  }
  pop_locations {
    pop_location = "cn-beijing-pop"
  }
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.anycast_id
  bind_instance_region_id = alicloud_ecs_instance.defaultEcs2.region_id
}
```

Multiple Usage

-> **NOTE:**  Anycast EIP supports binding cloud resource instances in multiple regions. Only one cloud resource instance is supported as the default origin station, and the rest are normal origin stations. When no access point is specified or an access point is added, the access request is forwarded to the default origin by default.  If you are bound for the first time, the Default value of the binding mode is **Default * *. /li> li> If you are not binding for the first time, you can set the binding mode to **Default**, and the new Default origin will take effect. The original Default origin will be changed to a common origin.

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
  provider                    = "alicloud.beijing"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  provider    = "alicloud.beijing"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  provider          = "alicloud.beijing"
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "defaultVpc" {
  provider   = "alicloud.beijing"
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVsw" {
  provider   = "alicloud.beijing"
  vpc_id     = alicloud_vpc.defaultVpc.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "defaultuBsECI" {
  provider = "alicloud.beijing"
  vpc_id   = alicloud_vpc.defaultVpc.id
}

resource "alicloud_instance" "default9KDlN7" {
  provider             = "alicloud.beijing"
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
  provider         = "alicloud.hangzhou"
  service_location = "ChineseMainland"
}

resource "alicloud_vpc" "defaultVpc2" {
  provider   = "alicloud.hangzhou"
  vpc_name   = "${var.name}6"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "default2" {
  provider                    = "alicloud.hangzhou"
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default2" {
  provider    = "alicloud.hangzhou"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default2" {
  provider          = "alicloud.hangzhou"
  availability_zone = data.alicloud_zones.default2.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vswitch" "defaultdsVsw2" {
  provider   = "alicloud.hangzhou"
  vpc_id     = alicloud_vpc.defaultVpc2.id
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default2.zones.1.id
}

resource "alicloud_security_group" "defaultuBsECI2" {
  provider = "alicloud.hangzhou"
  vpc_id   = alicloud_vpc.defaultVpc2.id
}

resource "alicloud_instance" "defaultEcs2" {
  provider             = "alicloud.hangzhou"
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
  provider                = "alicloud.beijing"
  bind_instance_id        = alicloud_instance.default9KDlN7.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-beijing"
  anycast_id              = alicloud_eipanycast_anycast_eip_address.defaultXkpFRs.id
  association_mode        = "Default"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "normal" {
  provider                = "alicloud.hangzhou"
  bind_instance_id        = alicloud_instance.defaultEcs2.network_interface_id
  bind_instance_type      = "NetworkInterface"
  bind_instance_region_id = "cn-hangzhou"
  anycast_id              = alicloud_eipanycast_anycast_eip_address_attachment.defaultEfYBJY.anycast_id
}
```

## Argument Reference

The following arguments are supported:
* `anycast_id` - (Required, ForceNew, Available since v1.113.0) The ID of the Anycast EIP instance.
* `association_mode` - (Optional, Computed) Binding mode, value:
  - **Default**: The Default mode. The cloud resource instance to be bound is set as the Default origin.
  - **Normal**: In Normal mode, the cloud resource instance to be bound is set to the common source station.
* `bind_instance_id` - (Required, ForceNew, Available since v1.113.0) The ID of the cloud resource instance to be bound.
* `bind_instance_region_id` - (Required, ForceNew, Available since v1.113.0) The region ID of the cloud resource instance to be bound.You can only bind cloud resource instances in some regions. You can call the [describeanystserverregions](~~ 171939 ~~) operation to obtain the region ID of the cloud resource instances that can be bound.
* `bind_instance_type` - (Required, ForceNew, Available since v1.113.0) The type of the cloud resource instance to be bound. Value:
  - **SlbInstance**: a private network SLB instance.
  - **NetworkInterface**: ENI.
* `pop_locations` - (Optional, Computed) The access point information of the associated access area when the cloud resource instance is bound.If you are binding for the first time, this parameter does not need to be configured, and the system automatically associates all access areas. See [`pop_locations`](#pop_locations) below.
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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anycast Eip Address Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Anycast Eip Address Attachment.
* `update` - (Defaults to 5 mins) Used when update the Anycast Eip Address Attachment.

## Import

Eipanycast Anycast Eip Address Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eipanycast_anycast_eip_address_attachment.example <anycast_id>:<bind_instance_id>:<bind_instance_region_id>:<bind_instance_type>
```