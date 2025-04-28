---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster_security_ip"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster Security I P resource.
---

# alicloud_click_house_enterprise_db_cluster_security_ip

Provides a Click House Enterprise Db Cluster Security I P resource.

Enterprise Clickhouse instance Security IP.

For information about Click House Enterprise Db Cluster Security I P and how to use it, see [What is Enterprise Db Cluster Security I P](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/ModifySecurityIPList).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_enterprise_db_cluster_security_ip&exampleId=ab765030-4584-f0ae-faba-0bce420c29c5e29d13d9&activeTab=example&spm=docs.r.click_house_enterprise_db_cluster_security_ip.0.ab76503045&intl_lang=EN_US" target="_blank">
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

resource "alicloud_click_house_enterprise_db_cluster" "defaultn0nVrN" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


resource "alicloud_click_house_enterprise_db_cluster_security_ip" "default" {
  group_name       = "example_group"
  security_ip_list = "127.0.0.2"
  db_instance_id   = alicloud_click_house_enterprise_db_cluster.defaultn0nVrN.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The cluster ID.
* `group_name` - (Required, ForceNew) The whitelist name.
* `security_ip_list` - (Required) The IP address list under the whitelist group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<group_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Db Cluster Security I P.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Db Cluster Security I P.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Db Cluster Security I P.

## Import

Click House Enterprise Db Cluster Security I P can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster_security_ip.example <db_instance_id>:<group_name>
```