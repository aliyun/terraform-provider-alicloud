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

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_dts_gad_instance" "default" {
  gad_instance_name         = "bingyu_example"
  db_engine_type            = "MySQL"
  gad_instance_type         = "DR"
  master_db_zone_id         = "cn-beijing-l"
  master_db_instance_region = "cn-beijing"
  master_db_instance_id     = "rm-2ze3772m0g292g33g"
}
```

## Argument Reference

The following arguments are supported:
* `db_engine_type` - (Optional, ForceNew) Database engine type
* `gad_instance_name` - (Optional) GAD instance name
* `gad_instance_type` - (Optional, ForceNew) GAD instance type
* `master_db_instance_id` - (Optional, ForceNew) Master database instance ID
* `master_db_instance_region` - (Optional, ForceNew) Master database instance region
* `master_db_zone_id` - (Optional, ForceNew) Master database instance zone ID

  -> **NOTE:** Changing this will force a replacement.

* `resource_group_id` - (Optional, ForceNew) The ID of the resource group. Create-only: no GAD resource-group relocation API exists, so changing this forces replacement.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
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
