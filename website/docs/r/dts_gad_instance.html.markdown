---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_gad_instance"
description: |-
  Provides a Alicloud DTS Gad Instance resource.
---

# alicloud_dts_gad_instance

Provides a DTS Gad Instance resource.

Global Active Database Cluster.

For information about DTS Gad Instance and how to use it, see [What is Gad Instance](https://next.api.alibabacloud.com/document/Dts/2020-01-01/CreateGadInstance).

-> **NOTE:** Available since v2.2.0.

## Example Usage

The GAD instance is a DTS-managed global active database cluster that attaches to an existing
master RDS instance. The master RDS instance provides four values the GAD instance needs:
`master_db_instance_id`, `db_engine_type`, `master_db_instance_region` and `master_db_zone_id`.
Create or reference the master RDS first, then reference its attributes so all four stay consistent.

Two regions are involved:

- **Provider region** (the `provider "alicloud" { region = ... }` block, or `ALICLOUD_REGION`): the
  region where the GAD instance itself is created and lives (`region_id` after apply). This is the
  DTS-side region.
- **master_db_instance_region**: the region where the master RDS instance lives. For a cross-region
  DR (disaster recovery) topology the master is in a different region from the GAD instance, so the
  master RDS is provisioned through a provider alias pinned to the master region.

The example below creates the master RDS in `cn-beijing` (via the `bj` alias) and the GAD instance in
`cn-hangzhou` (default provider), a typical cross-region DR topology. If your master RDS already
exists, replace the `alicloud_db_instance.master` resource with a `data.alicloud_db_instances` /
`data.alicloud_db_instance` data source and pass its `id`, `engine`, `zone_id` and region to the GAD
instance the same way.

```terraform
variable "name" {
  default = "terraform-example"
}

# Default provider: the GAD instance is created in this region (the DTS-side region).
provider "alicloud" {
  region = "cn-hangzhou"
}

# Alias pinned to the master RDS region, so the master RDS and its vswitch land in cn-beijing.
provider "alicloud" {
  alias  = "bj"
  region = "cn-beijing"
}

# Master RDS in cn-beijing, created through the bj alias. GAD requires a dedicated
# (mysql.x* / mssql.x* / pg.x*) RDS class; general-purpose classes are rejected.
data "alicloud_db_zones" "master" {
  provider                 = alicloud.bj
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "master" {
  provider                 = alicloud.bj
  zone_id                  = data.alicloud_db_zones.master.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

locals {
  gad_master_class = [
    for c in data.alicloud_db_instance_classes.master.instance_classes : c.instance_class
    if startswith(c.instance_class, "mysql.x")
  ][0]
}

resource "alicloud_vpc" "master" {
  provider   = alicloud.bj
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "master" {
  provider     = alicloud.bj
  vpc_id       = alicloud_vpc.master.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.master.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "master" {
  provider                 = alicloud.bj
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = local.gad_master_class
  instance_storage         = "100"
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = alicloud_vswitch.master.id
  instance_name            = var.name
}

# GAD instance in cn-hangzhou (default provider), attaching the cn-beijing master as its DR primary.
# All four master_* values come from the same RDS instance so they cannot drift apart.
resource "alicloud_dts_gad_instance" "default" {
  gad_instance_name         = var.name
  db_engine_type            = alicloud_db_instance.master.engine
  gad_instance_type         = "DR"
  master_db_instance_id     = alicloud_db_instance.master.id
  master_db_instance_region = "cn-beijing"
  master_db_zone_id         = data.alicloud_db_zones.master.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `db_engine_type` - (Optional, ForceNew) Database engine type of the master RDS instance. Valid values: `MySQL`, `SQLServer`, `PostgreSQL`. Must match the engine of `master_db_instance_id`. Changing this forces a replacement.
* `gad_instance_name` - (Optional) GAD instance name. Globally unique across the Alibaba Cloud account.
* `gad_instance_type` - (Optional, ForceNew) GAD instance type. Valid values: `DR` (disaster recovery) and `GAD` (global active database, dual-active). Changing this forces a replacement.
* `master_db_instance_id` - (Optional, ForceNew) ID of the master RDS instance the GAD instance attaches to. Use the ID of an existing RDS instance (e.g. `alicloud_db_instance` resource or `data.alicloud_db_instances` data source). Changing this forces a replacement.
* `master_db_instance_region` - (Optional, ForceNew) Region of the master RDS instance (e.g. `cn-beijing`). This is where the master RDS lives and may differ from the provider region where the GAD instance is created. Changing this forces a replacement.
* `master_db_zone_id` - (Optional, ForceNew) Zone ID of the master RDS instance (e.g. `cn-beijing-l`). Must be a zone in `master_db_instance_region`.

  -> **NOTE:** Changing this forces a replacement.

* `resource_group_id` - (Optional, ForceNew) The ID of the resource group. Create-only: CreateGadInstance accepts ResourceGroupId at creation, but DTS exposes no GAD-specific resource-group relocation API and ConvertInstanceResourceGroup (the DTS instance relocation API) returns Forbidden.InstanceNotFound for GAD instance IDs (`gad-*`), so a resource_group_id change forces replacement.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource (the provider region where the GAD instance lives).
* `status` - GAD instance status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gad Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Gad Instance.
* `update` - (Defaults to 5 mins) Used when update the Gad Instance.

## Import

DTS Gad Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_gad_instance.example <gad_instance_id>
```
