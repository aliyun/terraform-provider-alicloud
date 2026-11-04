---
subcategory: "EAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_eas_resource"
sidebar_current: "docs-alicloud-resource-pai-eas-resource"
description: |-
  Provides a PAI EAS Resource Group resource.
---

# alicloud_pai_eas_resource

Provides a PAI EAS Resource Group resource.

An EAS resource group is a container that groups EAS service instances together
with the backing compute resources. For information about EAS Resource Group and
how to use it, see
[CreateResource](https://www.alibabacloud.com/help/en/pai/user-guide/eas/eas-resource-group-management).

-> **NOTE:** Available since v2.1.1.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_eas_resource" "example" {
  resource_name = "terraform-example"
  charge_type   = "PostPaid"
}
```

## Argument Reference

The following arguments are supported:

* `resource_name` - (Required) (String) The name of the resource group. This is
  the only attribute that can be updated after creation.
* `ecs_instance_count` - (Optional, ForceNew) (Int) The number of instances.
  Required when `resource_type` is `Dedicated`.
* `ecs_instance_type` - (Optional, ForceNew) (String) The ECS instance type.
  Required when `resource_type` is `Dedicated`.
* `system_disk_size` - (Optional, ForceNew) (Int) The system disk size in GB.
  Valid values: 200 to 2000. Defaults to 200.
* `charge_type` - (Optional, ForceNew) (String) The charge type. Valid values:
  `PrePaid` and `PostPaid`. Required when `resource_type` is `Dedicated`.
* `auto_renewal` - (Optional, ForceNew) (Bool) Whether to auto renew. Valid only
  when `charge_type` is `PrePaid`.
* `zone` - (Optional, ForceNew) (String) The zone ID, e.g. `cn-shanghai-f`.
* `resource_type` - (Optional, ForceNew) (String) The resource type. Valid
  values: `Dedicated` and `SelfManaged`. `Dedicated` is a resource group whose
  backing machines are provisioned and managed by EAS;
  `SelfManaged` is a resource group whose machines belong to an external cluster
  supplied via the self-managed resource options. The API has no built-in
  default; when omitted, the server selects the type. When set to `Dedicated`,
  `ecs_instance_count`, `ecs_instance_type` and `charge_type` are required to
  provision the backing machines.
* `usage_mode` - (Optional, ForceNew) (String) The usage mode. Valid values:
  `inference` (an inference service resource group) and `sandbox` (a sandbox
  resource group).
* `labels` - (Optional, ForceNew) (Map) The labels of the resource group. A map
  of string key/value pairs.
* `external_cluster_id` - (Optional, ForceNew) (String) SelfManagedResourceOptions:
  the external cluster ID. Used when `resource_type` is `SelfManaged`.
* `role_name` - (Optional, ForceNew) (String) SelfManagedResourceOptions: the
  role name. Used when `resource_type` is `SelfManaged`.
* `node_match_labels` - (Optional, ForceNew) (Map) SelfManagedResourceOptions:
  the node match labels. A map of string key/value pairs. Used when
  `resource_type` is `SelfManaged`.

  -> **NOTE** The self-managed resource options
    (`external_cluster_id`, `role_name` and `node_match_labels`) are create-only:
    the DescribeResource API does not return them, so they cannot be read back
    or updated in place.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource group.
* `region_id` - (String) The region (cluster) ID of the resource group.
* `status` - (String) The status of the resource group.
* `create_time` - (String) The time when the resource group was created.

## Import

PAI EAS Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_eas_resource.example <resource_id>
```
