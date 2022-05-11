---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_namespaces"
sidebar_current: "docs-alicloud-datasource-mse-engine-namespaces"
description: |-
  Provides a list of Mse Engine Namespaces to the user.
---

# alicloud\_mse\_engine\_namespaces

This data source provides the Mse Engine Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mse_engine_namespaces" "ids" {
  cluster_id = "example_value"
  ids        = ["example_value"]
}
output "mse_engine_namespace_id_1" {
  value = data.alicloud_mse_engine_namespaces.ids.namespaces.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional, ForceNew) The language type of the returned information. Valid values: `zh`, `en`.
* `ids` - (Optional, ForceNew, Computed)  A list of Engine Namespace IDs. It is formatted to `<cluster_id>:<namespace_id>`.
* `cluster_id` - (Required, ForceNew) The id of the cluster.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `namespaces` - A list of Mse Engine Namespaces. Each element contains the following attributes:
  * `config_count` - The Number of Configuration of the Namespace.
  * `id` - The ID of the Engine Namespace. It is formatted to `<cluster_id>:<namespace_id>`.
  * `namespace_id` - The id of Namespace.
  * `namespace_desc` - The description of the Namespace.
  * `namespace_show_name` - The name of the Namespace.
  * `quota` - The Quota of the Namespace.
  * `service_count` - The number of active services.
  * `type` - The type of the Namespace, the value is as follows:
    - '0': Global Configuration.
    - '1': default namespace.
    - '2': Custom Namespace.