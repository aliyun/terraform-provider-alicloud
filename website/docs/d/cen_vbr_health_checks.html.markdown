---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_vbr_health_checks"
sidebar_current: "docs-alicloud-datasource-cen-vbr-health-checks"
description: |-
    Provides a list of CEN VBR Health Checks owned by an Alibaba Cloud account.
---

# alicloud\_cen\_vbr\_health\_checks

This data source provides CEN VBR Health Checks available to the user.

-> **NOTE:** Available in 1.98.0+

## Example Usage

```terraform
data "alicloud_cen_vbr_health_chekcs" "example" {
  cen_id                 = "cen-wxad980***"
  vbr_instance_id        = "vbr-mkahi8h****"
  vbr_instance_owner_id  = "1189203******"
  vbr_instance_region_id = "cn-beijing"
}

output "first_cen_vbr_health_check_id" {
  value = data.alicloud_cen_vbr_health_checks.example.checks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Optional, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `vbr_instance_id` - (Optional, ForceNew) The ID of the VBR instance.
* `vbr_instance_owner_id` - (Optional, ForceNew) The User ID (UID) of the account to which the VBR instance belongs.
* `vbr_instance_region_id` - (Required, ForceNew) The ID of the region where the VBR instance is deployed.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of the CEN VBR Heath Check IDs.
* `checks` - A list of CEN VBR Heath Checks. Each element contains the following attributes:
  * `id` - The ID of the CEN VBR Heath Check.
  * `cen_id` - The ID of the Cloud Enterprise Network (CEN) instance.
  * `health_check_interval` - The time interval at which probe packets are sent during the health check.
  * `health_check_source_ip` - The source IP address of the health check.
  * `health_check_target_ip` - The destination IP address of the health check.
  * `healthy_threshold` - The number of probe packets that are sent during the health check.
  * `vbr_instance_id` - The ID of the VBR instance.
  * `vbr_instance_region_id` - The ID of the region where the VBR instance is deployed.
