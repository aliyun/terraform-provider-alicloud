---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster_backup_policy"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster Backup Policy resource.
---

# alicloud_click_house_enterprise_db_cluster_backup_policy

Provides a Click House Enterprise Db Cluster Backup Policy resource.

Enterprise ClickHouse instance backup policy.

For information about Click House Enterprise Db Cluster Backup Policy and how to use it, see [What is Enterprise Db Cluster Backup Policy](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateBackupPolicy).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_enterprise_db_cluster_backup_policy&exampleId=53a58631-9bf5-a4e6-8ac7-4faf50f1d65d122210ff&activeTab=example&spm=docs.r.click_house_enterprise_db_cluster_backup_policy.0.53a586319b&intl_lang=EN_US" target="_blank">
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

variable "region_id" {
  default = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "default1tTLwe" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


resource "alicloud_click_house_enterprise_db_cluster_backup_policy" "default" {
  preferred_backup_period = "Monday"
  preferred_backup_time   = "04:00Z-05:00Z"
  backup_retention_period = "7"
  db_instance_id          = alicloud_click_house_enterprise_db_cluster.default1tTLwe.id
}
```

## Argument Reference

The following arguments are supported:
* `backup_retention_period` - (Required, Int) Backup retention time.
* `db_instance_id` - (Required, ForceNew) The instance ID.
* `preferred_backup_period` - (Required) Backup period.
* `preferred_backup_time` - (Required) Backup time.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Db Cluster Backup Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Db Cluster Backup Policy.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Db Cluster Backup Policy.

## Import

Click House Enterprise Db Cluster Backup Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster_backup_policy.example <id>
```