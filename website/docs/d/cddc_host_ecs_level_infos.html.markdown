---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_host_ecs_level_infos"
sidebar_current: "docs-alicloud-datasource-cddc-host-ecs-level-infos"
description: |-
  Provides a list of Cddc Dedicated Hosts to the user.
---

# alicloud\_cddc\_host_\ecs_\level_\infos

This data source provides the Cddc Host Ecs Level Infos of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cddc_zones" "example" {}

data "alicloud_cddc_host_ecs_level_infos" "ids" {
  zone_id = data.alicloud_cddc_zones.example.ids.0
}
output "cddc_host_ecs_level_infos_id_1" {
  value = data.alicloud_cddc_host_ecs_level_infos.ids.infos.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_type` - (Required, ForceNew) The database engine of the host. Valid values: `mysql`, `mssql`, `pgsql`, `redis`.
* `zone_id` - (Required, ForceNew, Computed) The ID of the zone in the region.
* `storage_type` - (Required, ForceNew) The storage type of the host ecs level info. Valid values: `local_ssd`, `cloud_essd`, `cloud_essd2`, `cloud_essd3`. 
  * `local_ssd`: specifies that the host uses local SSDs. 
  * `cloud_essd`: specifies that the host uses enhanced SSDs (ESSDs) of performance level (PL) 1. 
  * `cloud_essd2`: specifies that the host uses ESSDs of PL2. 
  * `cloud_essd3`: specifies that the host uses ESSDs of PL3.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `image_category` - (Optional, ForceNew) Host image. Valid values: `WindowsWithMssqlEntAlwaysonLicense`, `WindowsWithMssqlStdLicense`, `WindowsWithMssqlEntLicense`, `WindowsWithMssqlWebLicense`, `AliLinux`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `infos` - A list of Cddc Dedicated Hosts. Each element contains the following attributes:
	* `res_class_code` - The ApsaraDB RDS instance type of the host ecs level info.
	* `ecs_class_code` - The Elastic Compute Service (ECS) instance type.
	* `ecs_class` - The instance family of the host ecs level info.
	* `description` - The description of the host ecs level info.
