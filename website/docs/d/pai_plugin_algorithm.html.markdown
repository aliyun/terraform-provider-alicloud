---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_algorithm"
description: |-
  This data source provides PAIPlugin Algorithms.
---

# alicloud_pai_plugin_algorithm

This data source provides PAIPlugin Algorithms.

With `algorithm_id` set, the data source reads a single algorithm's full detail
(description, infer/train user config map); otherwise it lists all algorithms,
returning only the id and name of each (the order matches the gateway response).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
# List all algorithms.
data "alicloud_pai_plugin_algorithm" "all" {}

output "first_algorithm_id" {
  value = data.alicloud_pai_plugin_algorithm.all.ids.0
}

# Read a single algorithm's full detail.
data "alicloud_pai_plugin_algorithm" "one" {
  algorithm_id = "alg-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `algorithm_id` - (Optional) The algorithm id. When set, the data source reads the single algorithm via GetAlgorithm and populates the detail fields.
* `algorithm_name` - (Optional) The algorithm name. When set without `algorithm_id`, ListAlgorithms is filtered by name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `algorithm_id_out` - The algorithm id of the first matching algorithm.
* `algorithm_name_out` - The algorithm name of the first matching algorithm.
* `description` - The description. Populated only when `algorithm_id` is set.
* `infer_user_config_map` - The infer user config map (JSON string). Populated only when `algorithm_id` is set.
* `train_user_config_map` - The train user config map (JSON string). Populated only when `algorithm_id` is set.
* `ids` - The ids of all matching algorithms, in the order returned by ListAlgorithms.
* `names` - The names of all matching algorithms, in the order returned by ListAlgorithms.
