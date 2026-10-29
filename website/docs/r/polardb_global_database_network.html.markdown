---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_global_database_network"
sidebar_current: "docs-alicloud-resource-polardb-global-database-network"
description: |-
  Provides a Alicloud PolarDB Global Database Network resource.
---

# alicloud_polardb_global_database_network

Provides a PolarDB Global Database Network resource.

For information about PolarDB Global Database Network and how to use it, see [What is Global Database Network](https://www.alibabacloud.com/help/en/polardb/api-polardb-2017-08-01-createglobaldatabasenetwork).

-> **NOTE:** Available since v1.181.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_global_database_network&exampleId=0673154c-3963-3b2b-096d-ed036a8b01c2eaee75a3&activeTab=example&spm=docs.r.polardb_global_database_network.0.0673154c39&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
  pay_type   = "PostPaid"
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

resource "alicloud_polardb_global_database_network" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  description   = "terraform-example"
  gdn_version   = "1.0"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_global_database_network&spm=docs.r.polardb_global_database_network.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the primary cluster.
* `status` - (Computed) The status of the Global Database Network.
* `description` - (Optional, Computed) The description of the Global Database Network.
* `enable_global_domain_name` - (Optional, Bool, Available since v2.2.0) Specifies whether to enable the global domain name.

  -> **NOTE:** `enable_global_domain_name` is only sent when the Global Database Network is created or updated. The API does not return this parameter, so it cannot be read back from the remote resource.
* `gdn_version` - (Optional, Computed, ForceNew, String, Available since v2.2.0) The version of the Global Database Network. Defaults to `1.0`.
* `resource_group_id` - (Optional, Computed, String, Available since v2.2.0) The ID of the resource group.

  -> **NOTE:** `resource_group_id` is sent when the Global Database Network is created and when it is moved through `ModifyGlobalDatabaseNetwork`. The current Global Database Network response does not return a readable value, so the attribute keeps the configured value, is empty after `terraform import`, and changes made in the console are not detected as drift.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Global Database Network.
* `status` - The status of the Global Database Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the PolarDB Global Database Network.
* `update` - (Defaults to 3 mins) Used when update the PolarDB Global Database Network.
* `delete` - (Defaults to 10 mins) Used when delete the PolarDB Global Database Network.

## Import

PolarDB Global Database Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_global_database_network.example <id>
```