---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_znodes"
sidebar_current: "docs-alicloud-datasource-mse-znodes"
description: |-
  Provides a list of Mse Znodes to the user.
---

# alicloud\_mse\_znodes

This data source provides the Mse Znodes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.162.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mse_znodes" "ids" {
  cluster_id = "example_value"
  path       = "/"
  ids        = ["example_value-1", "example_value-2"]
}
output "mse_znode_id_1" {
  value = data.alicloud_mse_znodes.ids.znodes.0.id
}

data "alicloud_mse_znodes" "nameRegex" {
  path       = "/"
  cluster_id = "example_value"
  name_regex = "^my-Znode"
}
output "mse_znode_id_2" {
  value = data.alicloud_mse_znodes.nameRegex.znodes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional, ForceNew) The language type of the returned information. Valid values: `zh` or `en`.
* `ids` - (Optional, ForceNew, Computed)  A list of Znode IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Znode name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `cluster_id` - (Required, ForceNew) The ID of the Cluster.
* `path` - (Required, ForceNew) The Node path.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Znode names.
* `znodes` - A list of Mse Znodes. Each element contains the following attributes:
	* `cluster_id` - The ID of the Cluster.
	* `data` - The Node data.
	* `dir` - Node list information, the value is as follows:
	* `id` - The ID of the Znode. The value formats as `<cluster_id>:<path>`.
	* `path` - The Node path.
	* `znode_name` - The Node name.