---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_database"
description: |-
  Provides a Alicloud GPDB Database resource.
---

# alicloud_gpdb_database

Provides a GPDB Database resource.



For information about GPDB Database and how to use it, see [What is Database](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_database&exampleId=2544aefa-a42b-db0d-e339-ebc46269e760f42465dc&activeTab=example&spm=docs.r.gpdb_database.0.2544aefaa4&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "default35OkxY" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultl8haQ3" {
  vpc_id     = alicloud_vpc.default35OkxY.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultTC08a9" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  engine                = "gpdb"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultl8haQ3.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.default35OkxY.id
  db_instance_mode      = "StorageElastic"
}


resource "alicloud_gpdb_database" "default" {
  character_set_name = "UTF8"
  owner              = "adbpgadmin"
  description        = "go-to-the-docks-for-french-fries"
  database_name      = "seagull"
  collate            = "en_US.utf8"
  ctype              = "en_US.utf8"
  db_instance_id     = alicloud_gpdb_instance.defaultTC08a9.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_gpdb_database&spm=docs.r.gpdb_database.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `character_set_name` - (Optional, ForceNew) Character set, default value is UTF8
* `collate` - (Optional, ForceNew) Database locale parameters, specifying string comparison/collation
* `ctype` - (Optional, ForceNew) Database locale parameters, specifying character classification/case conversion rules
* `database_name` - (Required, ForceNew) Database Name
* `db_instance_id` - (Required, ForceNew) Instance ID
* `description` - (Optional, ForceNew) Database Description
* `owner` - (Required, ForceNew) Data Sheet owner

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<database_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Database.
* `delete` - (Defaults to 5 mins) Used when delete the Database.

## Import

GPDB Database can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_database.example <db_instance_id>:<database_name>
```