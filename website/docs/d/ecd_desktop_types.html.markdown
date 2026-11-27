---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop_types"
sidebar_current: "docs-alicloud-datasource-ecd-desktop-types"
description: |-
  Provides a list of Ecd Desktop Types to the user.
---

# alicloud_ecd_desktop_types

This data source provides the Ecd Desktop Types of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.170.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_desktop_types" "ids" {
  instance_type_family = "eds.hf"
}
output "ecd_desktop_type_id_1" {
  value = data.alicloud_ecd_desktop_types.ids.types.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cpu_count` - (Optional) The CPU cores.
* `gpu_count` - (Optional) The GPU cores.
* `ids` - (Optional, Computed)  A list of Desktop Type IDs.
* `instance_type_family` - (Optional) The Specification family. Valid values: `eds.graphics`, `eds.hf`, `eds.general`, `ecd.graphics`, `ecd.performance`, `ecd.advanced`, `ecd.basic`.
* `memory_size` - (Optional) The Memory size. Unit: MiB.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional) The status of the resource. Valid values: `SUFFICIENT`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Desktop Type IDs.
* `types` - A list of Ecd Desktop Types.

---

### `types`

The types supports the following attributes:

* `cpu_count` - The CPU cores.
* `data_disk_size` - The size of the data disk. Unit: GiB.
* `desktop_type_id` - Specification ID.
* `gpu_count` - The GPU cores.
* `gpu_spec` - The GPU video memory.
* `id` - The ID of the Desktop Type.
* `instance_type_family` - The Specification family.
* `memory_size` - The Memory size. Unit: MiB.
* `status` - The status of the resource.
* `system_disk_size` - The size of the system disk. Unit: GiB.
