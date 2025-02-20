---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_external_data_service"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) External Data Service resource.
---

# alicloud_gpdb_external_data_service

Provides a AnalyticDB for PostgreSQL (GPDB) External Data Service resource.

External Data Services.

For information about AnalyticDB for PostgreSQL (GPDB) External Data Service and how to use it, see [What is External Data Service](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateExternalDataService).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_external_data_service&exampleId=9e880330-16bc-83f6-1631-78a9a94b16f8b5aba08e&activeTab=example&spm=docs.r.gpdb_external_data_service.0.9e88033016&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "defaultrple4a" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultnYWSkl" {
  vpc_id     = alicloud_vpc.defaultrple4a.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultZ7DPgB" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultnYWSkl.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultrple4a.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}


resource "alicloud_gpdb_external_data_service" "default" {
  service_name        = "example6"
  db_instance_id      = alicloud_gpdb_instance.defaultZ7DPgB.id
  service_description = "example"
  service_spec        = "8"
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) Instance ID
* `service_description` - (Optional) Service Description
* `service_name` - (Required, ForceNew) Service Name
* `service_spec` - (Required) Service Specifications

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<service_id>`.
* `create_time` - The creation time of the resource
* `service_id` - Service ID
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the External Data Service.
* `delete` - (Defaults to 5 mins) Used when delete the External Data Service.
* `update` - (Defaults to 5 mins) Used when update the External Data Service.

## Import

AnalyticDB for PostgreSQL (GPDB) External Data Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_external_data_service.example <db_instance_id>:<service_id>
```