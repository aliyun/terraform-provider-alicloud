---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instances"
sidebar_current: "docs-alicloud-datasource-cen-instances"
description: |-
    Provides a list of CENs(Cloud Enterprise Network) which owned by an Alicloud account.
---


# alicloud\_cen\_instances

The CENs data source lists a number of CENs resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_cen_instances" "cen"{
	cen_ids = ["cen-id1"]
	cen_bandwidth_package_ids = ["cen_bwp_id1"]
	name_regex="^foo"
}

```

## Argument Reference

The following arguments are supported:

* `cen_ids` - (Optional) Limit search to a list of specific CEN IDs, like ["cen-id1","cen-id2"], which supports at most 5 items.
* `cen_bandwidth_package_ids` - (Optional) Limit search to a list of specific CEN Bandwidth Package IDs, like ["cen_bwp_id1", "cen_bwp_id2"], which supports at most 5 items.
* `name_regex` - (Optional) A regex string of CEN name.
* `output_file` - (Optional) The name of file that can save CENs data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the CEN.
* `name` - Name of the CEN.
* `status` - Status of the CEN, including "Creating", "Active" and "Deleting".
* `cen_bandwidth_package_ids` - List of CEN Bandwidth Package IDs in the specified CEN.
* `instance_ids` - List of child instance IDs in the specified CEN.
* `description` - Description of the CEN.