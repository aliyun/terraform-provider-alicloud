---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_private_srv_network_address"
description: |-
  Provides a Alicloud Mongodb Private Srv Network Address resource.
---

# alicloud_mongodb_private_srv_network_address

Provides a Mongodb Private Srv Network Address resource.

Private network SRV highly available link address.

For information about Mongodb Private Srv Network Address and how to use it, see [What is Private Srv Network Address](https://next.api.alibabacloud.com/document/Dds/2015-12-01/AllocateDBInstanceSrvNetworkAddress).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_private_srv_network_address&exampleId=f7832271-eb3b-840f-aaeb-df71975174cbd730c21d&activeTab=example&spm=docs.r.mongodb_private_srv_network_address.0.f7832271eb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

resource "alicloud_vpc" "defaultie35CW" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultg0DCAR" {
  vpc_id     = alicloud_vpc.defaultie35CW.id
  zone_id    = var.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "defaultHrZmxC" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.defaultg0DCAR.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.defaultie35CW.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = var.zone_id
}


resource "alicloud_mongodb_private_srv_network_address" "default" {
  db_instance_id = alicloud_mongodb_instance.defaultHrZmxC.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `private_srv_connection_string_uri` - Private network SRV highly available connection address

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 6 mins) Used when create the Private Srv Network Address.
* `delete` - (Defaults to 5 mins) Used when delete the Private Srv Network Address.

## Import

Mongodb Private Srv Network Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_private_srv_network_address.example <id>
```