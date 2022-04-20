---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_gws_cluster"
sidebar_current: "docs-alicloud-resource-ehpc-gws-cluster"
description: |-
  Provides a Alicloud Ehpc Gws Cluster resource.
---

# alicloud\_ehpc\_gws\_cluster

Provides a Ehpc Gws Cluster resource.

For information about Ehpc Gws Cluster and how to use it, see [What is Gws Cluster](https://www.alibabacloud.com/help/product/57664.html).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_ehpc_gws_cluster" "example" {
  name         = "example_value"
  cluster_type = "gws.s1.standard"
  vswitch_id   = data.alicloud_vswitches.default.ids.0
  vpc_id       = data.alicloud_vpcs.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Required, ForceNew) The type of the Gws Cluster. Valid values: `gws.s1.standard`.
* `name` - (Optional, ForceNew) The name of the Gws Cluster.
* `vswitch_id` - (Optional, ForceNew) The ID of the vswitch.
* `vpc_id` - (Required, ForceNew) The ID of the vswitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gws Cluster.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Gws Cluster.
* `delete` - (Defaults to 1 mins) Used when delete the Gws Cluster.

## Import

Ehpc Gws Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_ehpc_gws_cluster.example <id>
```