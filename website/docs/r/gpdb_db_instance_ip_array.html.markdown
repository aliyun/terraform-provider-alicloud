---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_instance_ip_array"
description: |-
  Provides a Alicloud GPDB DB Instance IP Array resource.
---

# alicloud_gpdb_db_instance_ip_array

Provides a GPDB DB Instance IP Array resource.

Whitelist IP Group.

For information about GPDB DB Instance IP Array and how to use it, see [What is DB Instance IP Array](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.231.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_db_instance_ip_array&exampleId=4cc2820f-adfa-7133-4baf-a813634a60f9003f45b5&activeTab=example&spm=docs.r.gpdb_db_instance_ip_array.0.4cc2820fad&intl_lang=EN_US" target="_blank">
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultNpLRa1" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultwLA5v4" {
  vpc_id     = alicloud_vpc.defaultNpLRa1.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultHKdDs3" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultwLA5v4.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultNpLRa1.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
  description           = var.name
}

resource "alicloud_gpdb_db_instance_ip_array" "default" {
  db_instance_ip_array_attribute = "taffyFish"
  security_ip_list               = ["12.34.56.78", "11.45.14.0", "19.19.81.0"]
  db_instance_ip_array_name      = "taffy"
  db_instance_id                 = alicloud_gpdb_instance.defaultHKdDs3.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_ip_array_attribute` - (Optional, ForceNew) The default is empty. To distinguish between different attribute values, the console does not display groups with the 'hidden' attribute.
* `db_instance_ip_array_name` - (Required, ForceNew) The name of the IP address whitelist. If you do not specify this parameter, the default whitelist is queried.

-> **NOTE:**   Each instance supports up to 50 IP address whitelists.

* `db_instance_id` - (Required, ForceNew) The instance ID.

-> **NOTE:**  You can call the [DescribeDBInstances](https://www.alibabacloud.com/help/en/doc-detail/86911.html) operation to query details about all AnalyticDB for PostgreSQL instances within a region, including instance IDs.

* `modify_mode` - (Optional) The method of modification. Valid values:

  - `Cover`: overwrites the whitelist.
  - `Append`: appends data to the whitelist.
  - `Delete`: deletes the whitelist.
* `security_ip_list` - (Required, Set) The IP address whitelist contains a maximum of 1000 IP addresses separated by commas in the following three formats:
  - 0.0.0.0/0
  - 10.23.12.24(IP)
  - 10.23.12.24/24(CIDR mode, Classless Inter-Domain Routing, '/24' indicates the length of the prefix in the address, and the range is '[1,32]')

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<db_instance_ip_array_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the DB Instance IP Array.
* `delete` - (Defaults to 5 mins) Used when delete the DB Instance IP Array.
* `update` - (Defaults to 5 mins) Used when update the DB Instance IP Array.

## Import

GPDB DB Instance IP Array can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_db_instance_ip_array.example <db_instance_id>:<db_instance_ip_array_name>
```