---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addons"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-addons"
description: |-
  Provides a list of available addons.
---

# alicloud_cs_kubernetes_addons

This data source provides a list of available addons that the cluster can install.

-> **NOTE:** Available since v1.150.0.
-> **NOTE:** From version v1.166.0, support for returning custom configuration of kubernetes cluster addon.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

data "alicloud_cs_kubernetes_addons" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
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