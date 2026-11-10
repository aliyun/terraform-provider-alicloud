---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_auto_repair_policies"
sidebar_current: "docs-alicloud-datasource-cs-auto-repair-policies"
description: |-
  Provides a list of Container Service for Kubernetes (ACK) Auto Repair Policies.
---

# alicloud_cs_auto_repair_policies

This data source provides the Container Service for Kubernetes (ACK) Auto Repair Policies of the current Alibaba Cloud user. Both custom policies and builtin policies are returned.

-> **NOTE:** Available since v2.3.0.

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
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  cluster_spec                 = "ack.pro.small"
  vswitch_ids                  = [alicloud_vswitch.default.id]
  pod_vswitch_ids              = [alicloud_vswitch.default.id]
  new_nat_gateway              = false
  is_enterprise_security_group = true
  slb_internet_enabled         = false
  deletion_protection          = false
}

resource "alicloud_cs_auto_repair_policy" "default" {
  cluster_id        = alicloud_cs_managed_kubernetes.default.id
  name              = var.name
  resource_sub_type = "ess"

  rules {
    incidents {
      name = "NTPProblem"
    }
    repair_procedure {
      name = "QuarantineNode"
    }
    repair_procedure {
      name = "Evict"
    }
    repair_procedure {
      name = "Repair"
    }
    repair_procedure {
      name = "UnquarantineNode"
    }
  }
}

# Query custom policies created by the user
data "alicloud_cs_auto_repair_policies" "custom" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  type       = "custom"
  name_regex = var.name
}

# Query builtin policies that apply to elastic (ESS) node pools
data "alicloud_cs_auto_repair_policies" "builtin" {
  cluster_id        = alicloud_cs_managed_kubernetes.default.id
  type              = "builtin"
  resource_sub_type = "ess"
}

output "custom_policy_ids" {
  value = data.alicloud_cs_auto_repair_policies.custom.ids
}

output "builtin_policy_ids" {
  value = data.alicloud_cs_auto_repair_policies.builtin.ids
}
```

## Builtin policies

Besides the custom policies, every cluster has two read-only builtin policies. They are returned by this data source with `type = "builtin"` and can be bound to a node pool through `auto_repair_policy_id`.

| Policy ID | Incidents | Procedures | Node pool types |
| --- | --- | --- | --- |
| `r-repair-sys-comp-failure` | `NTPProblem`, `RuntimeOffline`, `Ready` | `QuarantineNode`, `Evict` (`evict_timeout = 300`), `Repair` (`restart_node = true`), `UnquarantineNode` | `ess`, `lingjun` |
| `r-repair-sys-comp-failure` | `SystemdOffline`, `ReadonlyFilesystem` | `QuarantineNode`, `Evict` (`evict_timeout = 300`), `Repair`, `UnquarantineNode` | `ess`, `lingjun` |
| `r-repair-sys-comp-failure-without-reboots` | `NTPProblem`, `RuntimeOffline`, `Ready` | `Repair` (`restart_node = false`) | `ess`, `lingjun` |

A builtin policy covers several node pool types at the same time, so it is returned for both `resource_sub_type = "ess"` and `resource_sub_type = "lingjun"`.

The `rules` attribute of an item is the JSON encoding of the rules. It has the same structure as the `rules` block of [`alicloud_cs_auto_repair_policy`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_auto_repair_policy), for example:

```terraform
rules {
  incidents {
    name = "NTPProblem"
  }
  incidents {
    name = "RuntimeOffline"
  }
  incidents {
    name = "Ready"
  }
  repair_procedure {
    name = "QuarantineNode"
  }
  repair_procedure {
    name   = "Evict"
    config = jsonencode({ evict_timeout = 300 })
  }
  repair_procedure {
    name   = "Repair"
    config = jsonencode({ restart_node = true })
  }
  repair_procedure {
    name = "UnquarantineNode"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the ACK cluster.
* `ids` - (Optional) A list of auto repair policy IDs. The ID consists of the cluster ID and the policy ID, in the format of `<cluster_id>:<policy_id>`.
* `name_regex` - (Optional) A regex string to filter results by policy name.
* `type` - (Optional, Available since v2.3.0) The type of the auto repair policy. Valid values: `custom`, `builtin`.
* `resource_sub_type` - (Optional, Available since v2.3.0) The subtype of the node pool that the policy applies to. Valid values: `ess`, `lingjun`. Builtin policies that cover several node pool types match every one of them.
* `nodepool_id` - (Optional, Available since v2.3.0) The ID of the node pool. Only policies bound to this node pool are returned.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source. It is the hash of the returned policy IDs.
* `ids` - A list of auto repair policy IDs.
* `names` - A list of auto repair policy names.
* `policies` - A list of auto repair policies. Each element contains the following attributes:
  * `id` - The ID of the auto repair policy, in the format of `<cluster_id>:<policy_id>`.
  * `policy_id` - The ID of the auto repair policy.
  * `name` - The name of the auto repair policy.
  * `type` - The type of the auto repair policy. Valid values: `custom`, `builtin`.
  * `resource_type` - The type of the resource that the policy applies to.
  * `resource_sub_type` - The subtype of the node pool that the policy applies to. Multiple values are separated by `|`.
  * `resource_ids` - The IDs of the node pools that the policy is bound to.
  * `rules` - The repair rules of the policy, in JSON format.
