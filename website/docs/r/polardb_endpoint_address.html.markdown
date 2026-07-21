---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoint_address"
sidebar_current: "docs-alicloud-resource-poalrdb-endpoint"
description: |-
  Provides a PolarDB instance endpoint resource.
---

# alicloud_polardb_endpoint_address

Provides a PolarDB endpoint address resource to allocate an Internet endpoint address string for PolarDB instance.

-> **NOTE:** Available since v1.68.0. Each PolarDB instance will allocate a intranet connection string automatically and its prefix is Cluster ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_endpoint_address&exampleId=3a3b6e68-c19b-a012-6f4c-c12a229b35a6cbf74ad4&activeTab=example&spm=docs.r.polardb_endpoint_address.0.3a3b6e68c1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

data "alicloud_polardb_endpoints" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
}

resource "alicloud_polardb_endpoint_address" "default" {
  db_cluster_id     = alicloud_polardb_cluster.default.id
  db_endpoint_id    = data.alicloud_polardb_endpoints.default.endpoints[0].db_endpoint_id
  connection_prefix = "polardbexample"
  net_type          = "Public"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_endpoint_address&spm=docs.r.polardb_endpoint_address.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `db_endpoint_id` - (Required, ForceNew) The Id of endpoint that can run database.
* `connection_prefix` - (Optional) Prefix of the specified endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter.
* `net_type` - (Optional, ForceNew) Internet connection net type. Valid value: `Public`. Default to `Public`. Currently supported only `Public`.
* `port` - (Optional, Available since v1.217.0) Port of the specified endpoint. Valid values: 3000 to 5999.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `connection_string` - Connection cluster or endpoint string.
* `ip_address` - The ip address of connection string.

## Import

PolarDB endpoint address can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_endpoint_address.example pc-abc123456:pe-abc123456
```
