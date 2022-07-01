---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_bundles"
sidebar_current: "docs-alicloud-datasource-ecd-bundles"
description: |-
  Provides a list of Ecd bundles to the user.
---

# alicloud\_ecd\_bundles

This data source provides the Ecd bundles of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Bundle IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Bundle name.
* `bundle_type` - (Optional, ForceNew) The bundle type of  the bundle. Valid values: `SYSTEM`,`CUSTOM`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Bundle names.
* `bundles` - A list of Ecd Bundle. Each element contains the following attributes:
    * `bundle_type` - The bundle type of  the bundle. Valid values: `SYSTEM`,`CUSTOM`.
    * `desktop_type` - The desktop type of the bundle.
    * `description` - The description of the bundle.
    * `bundle_id` - The bundle id of the bundle.
    * `desktop_type_attribute` - The desktop type attribute of the bundle.
        * `cpu_count` - The cpu count attribute of the bundle.
        * `gpu_count` - The gpu count attribute of the bundle.
        * `gpu_spec` - The gpu spec attribute of the bundle.
        * `memory_size` - The memory size attribute of the bundle.
    * `disks` - The disks of the bundle.
      * `disk_size` - The disk size attribute of the bundle.
      * `disk_type` - The disk type attribute of the bundle.
    * `id` - The ID of the bundle.
    * `image_id` - The image id attribute of the bundle.
    * `os_type` - The os type attribute of the bundle.
    * `bundle_name` - The name of the bundle.
   
