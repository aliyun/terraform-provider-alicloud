---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_consumer_channel"
sidebar_current: "docs-alicloud-resource-dts-consumer-channel"
description: |-
  Provides a Alicloud DTS Consumer Channel resource.
---

# alicloud_dts_consumer_channel

Provides a DTS Consumer Channel resource.

For information about DTS Consumer Channel and how to use it, see [What is Consumer Channel](https://www.alibabacloud.com/help/en/doc-detail/264593.htm).

-> **NOTE:** Available since v1.146.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dts_consumer_channel&exampleId=7eb9e094-58bc-8836-3b78-8e30349b6d43aaf0a634&activeTab=example&spm=docs.r.dts_consumer_channel.0.7eb9e09458&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_regions" "example" {
  current = true
}
data "alicloud_db_zones" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.example.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.example.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_db_instance" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.example.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
  security_group_ids       = [alicloud_security_group.example.id]
}

resource "alicloud_rds_account" "example" {
  db_instance_id   = alicloud_db_instance.example.id
  account_name     = "example_name"
  account_password = "example_1234"
}

resource "alicloud_db_database" "example" {
  instance_id = alicloud_db_instance.example.id
  name        = var.name
}

resource "alicloud_db_account_privilege" "example" {
  instance_id  = alicloud_db_instance.example.id
  account_name = alicloud_rds_account.example.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.example.name]
}

resource "alicloud_dts_subscription_job" "example" {
  dts_job_name                       = var.name
  payment_type                       = "PayAsYouGo"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.example.regions.0.id
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.example.id
  source_endpoint_database_name      = alicloud_db_database.example.name
  source_endpoint_user_name          = alicloud_rds_account.example.account_name
  source_endpoint_password           = alicloud_rds_account.example.account_password
  db_list                            = "{\"${alicloud_db_database.example.name}\":{\"name\":\"${alicloud_db_database.example.name}\",\"all\":true}}"
  subscription_instance_network_type = "vpc"
  subscription_instance_vpc_id       = alicloud_vpc.example.id
  subscription_instance_vswitch_id   = alicloud_vswitch.example.id
  status                             = "Normal"
}

resource "alicloud_dts_consumer_channel" "example" {
  dts_instance_id          = alicloud_dts_subscription_job.example.dts_instance_id
  consumer_group_name      = var.name
  consumer_group_user_name = "example"
  consumer_group_password  = "example1234"
}
```

## Argument Reference

The following arguments are supported:

* `consumer_group_name` - (Required, ForceNew) The name of the consumer group.
* `consumer_group_password` - (Required) The password of the consumer group account. The length of the `consumer_group_password` is limited to `8` to `32` characters. It can contain two or more of the following characters: uppercase letters, lowercase letters, digits, and special characters.
* `consumer_group_user_name` - (Required, ForceNew) The username of the consumer group. The length of the `consumer_group_user_name` is limited to `1` to `16` characters. It can contain one or more of the following characters: uppercase letters, lowercase letters, digits, and underscores (_).
* `dts_instance_id` - (Required) The ID of the subscription instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Consumer Channel. The value formats as `<dts_instance_id>:<consumer_group_id>`.
* `consumer_group_id` - The ID of the consumer group.

## Import

DTS Consumer Channel can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_consumer_channel.example <dts_instance_id>:<consumer_group_id>
```