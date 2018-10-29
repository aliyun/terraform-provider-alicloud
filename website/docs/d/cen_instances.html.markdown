---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instances"
sidebar_current: "docs-alicloud-datasource-cen-instances"
description: |-
    Provides a list of CEN(Cloud Enterprise Network) instances owned by an Alibaba Cloud account.
---

# alicloud\_cen\_instances

This data source provides CEN instances available to the user.

## Example Usage

```
data "alicloud_cen_instances" "cen_instances_ds"{
  ids = ["cen-id1"]
  name_regex = "^foo"
}

output "first_cen_instance_id" {
  value = "${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of CEN instances IDs.
* `name_regex` - (Optional) A regex string to filter CEN instances by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of CEN instances. Each element contains the following attributes:
  * `id` - ID of the CEN instance.
  * `name` - Name of the CEN instance.
  * `status` - Status of the CEN instance, including "Creating", "Active" and "Deleting".
  * `bandwidth_package_ids` - List of CEN Bandwidth Package IDs in the specified CEN instance.
  * `child_instance_ids` - List of child instance IDs in the specified CEN instance.
  * `description` - Description of the CEN instance.