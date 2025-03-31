---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster_public_endpoint"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster Public Endpoint resource.
---

# alicloud_click_house_enterprise_db_cluster_public_endpoint

Provides a Click House Enterprise Db Cluster Public Endpoint resource.

ClickHouse enterprise instance public network endpoint.

For information about Click House Enterprise Db Cluster Public Endpoint and how to use it, see [What is Enterprise Db Cluster Public Endpoint](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateEndpoint).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_enterprise_db_cluster_public_endpoint&exampleId=bb2c0cf0-e0ba-620b-21c1-8fd26bcbf09cdd897168&activeTab=example&spm=docs.r.click_house_enterprise_db_cluster_public_endpoint.0.bb2c0cf0e0&intl_lang=EN_US" target="_blank">
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

resource "alicloud_click_house_enterprise_db_cluster" "defaultaqnt22" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


resource "alicloud_click_house_enterprise_db_cluster_public_endpoint" "default" {
  db_instance_id           = alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id
  net_type                 = "Public"
  connection_string_prefix = alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id
}
```

## Argument Reference

The following arguments are supported:
* `connection_string_prefix` - (Required) The public network connection string prefix of the instance.
* `db_instance_id` - (Required, ForceNew) The cluster ID.
* `net_type` - (Required, ForceNew) Network type of the connection address. Valid values:
  - Public: Public network.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<net_type>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Db Cluster Public Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Db Cluster Public Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Db Cluster Public Endpoint.

## Import

Click House Enterprise Db Cluster Public Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster_public_endpoint.example <db_instance_id>:<net_type>
```