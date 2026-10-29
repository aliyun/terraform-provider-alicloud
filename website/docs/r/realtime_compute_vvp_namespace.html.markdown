---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_vvp_namespace"
description: |-
  Provides an Alicloud Realtime Compute Vvp Namespace resource.
---

# alicloud_realtime_compute_vvp_namespace

Provides a Realtime Compute Vvp Namespace resource.

For information about Realtime Compute Vvp Namespace and how to use it, see [What is Vvp Namespace](https://next.api.alibabacloud.com/document/foasconsole/2019-06-01/CreateNamespace).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** A namespace belongs to one Realtime Compute instance. `instance_id` and `namespace` are the composite identifier (`<instance_id>:<namespace>`) and cannot be changed in place; changing either one replaces the resource.

-> **NOTE:** `ha` is only accepted while the namespace is created. Changing it afterwards forces a new namespace.

-> **NOTE:** `tags` is exported from the API and cannot be managed through this resource: `CreateNamespace` does not accept tags.

-> **NOTE:** `resource_spec` consumes the capacity purchased on the parent instance. A workspace created with a `resource_spec` also gets a default namespace, named `<workspace_name>-default` for a single-availability-zone workspace and `<workspace_name>-default-ha` for a multi-availability-zone one, and that namespace takes the whole spec. The workspace must therefore own spare capacity before another namespace can be created: scale the workspace up in the console or by calling `ModifyInstanceSpec` first. A namespace that needs more capacity than the workspace has left fails with error code `903040`.

-> **NOTE:** `resource_spec.memory_gb` must be exactly 4 times `resource_spec.cpu` (for example 2/8, 3/12 or 4/16). The provider rejects other ratios while planning, instead of letting the request fail at the API.

-> **NOTE:** Although `resource_spec` is optional in the schema, the API rejects a namespace without an allocated spec (`903038`), so set it in practice.

-> **NOTE:** A namespace that is the last valid one in its workspace cannot be deleted on its own: the API rejects `DeleteNamespace` with `903054`. `terraform destroy` or removing such a namespace from the configuration therefore fails with that error and keeps the resource in state — the namespace still exists. To remove it, delete the parent workspace (the namespace disappears with it) and then drop the resource from state with `terraform state rm`, or create another namespace in the workspace first so this one is no longer the last.

-> **NOTE:** `resource_spec` can be adjusted after the namespace is created, in either direction. The target must not be smaller than the total already used by the namespace's resource queues; otherwise `ModifyNamespaceSpec` rejects it with `903099`. Those queues are not managed by this resource.

## Example Usage

Basic Usage

The namespace is created inside an existing workspace that still owns spare compute capacity, see the note above.

```terraform
variable "vvp_instance_id" {
  description = "ID of a Realtime Compute workspace that owns spare compute capacity"
}

resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = var.vvp_instance_id
  namespace   = "terraform-example-ns"
  ha          = false

  resource_spec = [
    {
      cpu       = 2
      memory_gb = 8
    }
  ]
}
```

A workspace prepared for its own namespaces looks like this; its default namespace consumes the whole `resource_spec`, so scale it up before the namespace above is added.

```terraform
variable "name" {
  default = "terraform-example"
}

variable "zone_id" {
  default = "cn-hangzhou-i"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = var.zone_id
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.bucket
    }
  }
  vvp_instance_name = "${var.name}-${random_integer.default.result}"
  vpc_id            = alicloud_vpc.default.id
  zone_id           = var.zone_id
  vswitch_ids       = [alicloud_vswitch.default.id]
  payment_type      = "PayAsYouGo"

  resource_spec {
    cpu       = 4
    memory_gb = 16
  }
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew, String) The ID of the Realtime Compute instance that owns the namespace.
* `namespace` - (Required, ForceNew, String) The name of the namespace. A workspace gets its default namespace automatically, named `<workspace_name>-default` (or `<workspace_name>-default-ha` when the workspace is multi-availability-zone).
* `ha` - (Optional, ForceNew, Bool) Whether the namespace is highly available. It can only be set when the namespace is created.
* `resource_spec` - (Optional, List) The resource specification (CU) of the namespace. See [`resource_spec`](#resource_spec) below.
* `timeouts` - (Optional) See [Timeouts](#timeouts) below.

### `resource_spec`

The resource_spec supports the following:
* `cpu` - (Required, Int) The number of CPU cores.
* `memory_gb` - (Required, Int) The memory size in GB. It must be exactly 4 times `cpu`: the product rejects any other ratio, so the provider rejects it at plan time.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<namespace>`.
* `gmt_create` - (Computed, Int) The creation time.
* `gmt_modified` - (Computed, Int) The modification time.
* `resource_used` - (Computed, List) The resource usage of the namespace.
  * `cpu` - (Computed, Float) The used CPU cores.
  * `cu` - (Computed, Float) The used compute units.
  * `memory_gb` - (Computed, Float) The used memory size in GB.
* `status` - (Computed, String) The status of the namespace.
* `tags` - (Computed, Map) The tags of the namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vvp Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Vvp Namespace.
* `update` - (Defaults to 5 mins) Used when update the Vvp Namespace.

## Import

Realtime Compute Vvp Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_vvp_namespace.example <instance_id>:<namespace>
```
