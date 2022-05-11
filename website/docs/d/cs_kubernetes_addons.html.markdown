---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addons"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-addons"
description: |-
  Provides a list of available addons.
---

# alicloud\_cs\_kubernetes\_addons

This data source provides a list of available addons that the cluster can install.

-> **NOTE:** Available in 1.150.0+.
-> **NOTE:** From version 1.166.0, support for returning custom configuration of kubernetes cluster addon.

## Example Usage

```terraform
data "alicloud_cs_kubernetes_addons" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
}

output "addons" {
  value = data.alicloud_cs_kubernetes_addons.default.addons
}
```

## Argument Reference

The following arguments are supported.
* `cluster_id` - (Required) The id of kubernetes cluster.
* `ids` - (Optional) A list of addon IDs. The id of addon consists of the cluster id and the addon name, with the structure <cluster_ud>:<addon_name>.
* `name_regex` - (Optional) A regex string to filter results by addon name.

## Attributes Reference

* `cluster_id` - The id of kubernetes cluster.
* `names` - A list of addon names.
* `addons` - A list of addons.
  * `name` - The name of addon. 
  * `current_version` - The current version of addon, if this field is an empty string, it means that the addon is not installed.
  * `next_version` - The next version of this addon can be upgraded to.
  * `required` - Whether the addon is a system addon.
  * `current_config` - The current custom configuration of the addon. **Note:** Available in v1.166.0+