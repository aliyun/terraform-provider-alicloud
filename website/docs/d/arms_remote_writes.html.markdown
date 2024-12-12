---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_remote_writes"
sidebar_current: "docs-alicloud-datasource-arms-remote-writes"
description: |-
  Provides a list of Arms Remote Writes to the user.
---

# alicloud_arms_remote_writes

This data source provides the Arms Remote Writes of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.204.0.

-> **DEPRECATED:** This data source has been deprecated since v1.228.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_remote_writes" "ids" {
  ids        = ["example_id"]
  cluster_id = "your_cluster_id"
}

output "arms_remote_writes_id_1" {
  value = data.alicloud_arms_remote_writes.ids.remote_writes.0.id
}

data "alicloud_arms_remote_writes" "nameRegex" {
  name_regex = "tf-example"
  cluster_id = "your_cluster_id"
}

output "arms_remote_writes_id_2" {
  value = data.alicloud_arms_remote_writes.nameRegex.remote_writes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Remote Write IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Remote Write name.
* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Remote Write names.
* `remote_writes` - A list of Remote Writes. Each element contains the following attributes:
  * `id` - The ID of the Remote Write. It formats as `<cluster_id>:<remote_write_name>`.
  * `cluster_id` - The ID of the Prometheus instance.
  * `remote_write_name` - The name of the Remote Write configuration item.
  * `remote_write_yaml` - The details of the Remote Write configuration item. The value is in the YAML format.
  