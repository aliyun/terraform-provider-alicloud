---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_auto_repair_policy"
sidebar_current: "docs-alicloud-resource-cs-auto-repair-policy"
description: |-
  Provides a Alicloud Container Service for Kubernetes (ACK) Auto Repair Policy resource.
---

# alicloud_cs_auto_repair_policy

Provides a Container Service for Kubernetes (ACK) Auto Repair Policy resource. An auto repair policy defines which node incidents are repaired, by which repair procedures, and whether manual approval is required before a repair runs. A policy is bound to node pools through the `auto_repair_policy` block of [alicloud_cs_kubernetes_node_pool](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_node_pool).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** Builtin policies cannot be created, updated or deleted. Use the `alicloud_cs_auto_repair_policies` data source to query them, and reference their IDs in the node pool configuration.

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
      name   = "Evict"
      config = jsonencode({ evict_timeout = 300 })
    }
    repair_procedure {
      name = "Repair"

      intervention {
        enable = true
        type   = "label"
        inquiring_label {
          key   = "intervention.autorepair.node.k8s.aliyun.com/repair"
          value = "inquiring"
        }
        approved_label {
          key   = "intervention.autorepair.node.k8s.aliyun.com/repair"
          value = "approved"
        }
      }
    }
    repair_procedure {
      name = "UnquarantineNode"
    }
  }
}
```

### Repair a GPU fault by rebooting the node

```terraform
resource "alicloud_cs_auto_repair_policy" "gpu" {
  cluster_id        = alicloud_cs_managed_kubernetes.default.id
  name              = "gpu-repair"
  resource_sub_type = "ess"

  rules {
    incidents {
      name = "NvidiaXID48Error"
    }
    repair_procedure {
      name = "Repair"
    }
  }
}
```

## How it works

Auto repair is driven by the `ack-node-problem-detector` (NPD) add-on and `node-lifecycle-controller`:

1. **Fault detection.** ACK watches CloudMonitor events, ECS and Lingjun system events, and node conditions reported by NPD. Each fault is converted into a Node Condition and a Kubernetes Event.
2. **Rule matching.** The policy whose `resource_sub_type` matches the node pool and whose `incidents` contain the fault name is applied.
3. **Repair.** The `repair_procedure` list runs in the order defined by the built-in template of the incident, optionally pausing at any step for manual approval.
4. **Reporting.** Progress and the final result are reported through Kubernetes Events and the ACK alert rule sets.

Each node is repaired individually by default. The parallelism and circuit breaking thresholds are configured on the node pool, see `max_parallel_repairing_nodes` and `max_unhealthy_nodes_threshold` in [alicloud_cs_kubernetes_node_pool](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_node_pool).

-> **NOTE:** Auto repair requires a managed node pool with managed configuration enabled, the `node-lifecycle-controller` add-on at version `v2.3.13-b581b592` or later and the `ack-node-problem-detector` add-on at version `1.2.41` or later, and is supported only for `ess` and `lingjun` node pools (including super nodes) in public cloud regions. GPU and system/component faults additionally rely on NPD fault checks; the service validates the installed add-on versions and rejects requests whose version is below the minimum, so install or upgrade the add-ons with [`alicloud_cs_kubernetes_addon`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_addon) before binding a policy. Do not enable auto repair together with other node auto repair software in the same cluster.

## Incident reference

`incidents.name` must be a built-in system incident. Within one `rules` block all incidents must share the same repair procedure template, an incident can appear only once per policy, and the incident must be supported by the policy's `resource_sub_type`.

### Incidents grouped by repair procedure template

| Template | Incidents | Node pool types | Required procedures |
| --- | --- | --- | --- |
| `SystemComponentsError` | `NTPProblem`, `RuntimeOffline`, `Ready` | `ess`, `lingjun` | `Repair` |
| `SystemComponentsOnlyRebootError` | `SystemdOffline`, `ReadonlyFilesystem` | `ess`, `lingjun` | None, all procedures are optional |
| `NodeNeedRebootError` | `SystemMaintenance.Reboot:Inquiring`, `SystemMaintenance.Reboot:Scheduled` | `ess` | `Repair` |
| `NodeNeedRebootError` | `Node.FaultNeedReboot.HOST:Inquiring` | `lingjun` | `Repair` |
| `NodeNeedRedeployErrorWithPossibleDataLoss` | `SystemMaintenance.Redeploy:Inquiring`, `SystemMaintenance.Redeploy:Scheduled` | `ess` | `Repair` |
| `NodeNeedRemedyErrorWithPossibleDataLoss` | `Node.HardwareFault.HOST:Inquiring` | `lingjun` | `Repair` |
| `NodeNeedRebootErrorDueToGPUFaults` | `NvidiaXID31Error`, `NvidiaXID48Error`, `NvidiaXID62Error`, `NvidiaXID95Error`, `NvidiaXID119Error`, `NvidiaXID120Error`, `NvidiaXID140Error`, `NvidiaPendingRetiredPages`, `NvidiaRemappingRowsRequireReset` | `ess`, `lingjun` | `Repair` |
| `HyperNodeNeedRebootError` | `HyperNode.Node.FaultNeedReboot.HOST:Inquiring` | `lingjun` | `Repair` |
| `HyperNodeNeedRemedyErrorWithPossibleDataLoss` | `HyperNode.Node.HardwareFault.HOST:Inquiring` | `lingjun` | `Repair` |

### Incident meanings

System and component faults, detected by NPD:

| Incident | Meaning | Detection threshold |
| --- | --- | --- |
| `NTPProblem` | The node time synchronization process (`ntpd` or `chronyd`) is abnormal. | 10s |
| `RuntimeOffline` | The container runtime (`containerd` or Docker Engine) stopped and the node is unavailable. | 90s |
| `Ready` | The node is NotReady because kubelet is hung or its heartbeat is abnormal. It also covers PLEG failures and sandbox errors. | 180s |
| `SystemdOffline` | systemd is abnormal, so containers cannot be started or destroyed. | 90s |
| `ReadonlyFilesystem` | The node file system became read-only. | 90s |

Node faults:

| Incident | Meaning |
| --- | --- |
| `SystemMaintenance.Reboot:Inquiring`, `SystemMaintenance.Reboot:Scheduled` | Alibaba Cloud detected a potential software or hardware risk on the underlying host and the instance should be rebooted. Initiate the reboot during the repair; otherwise ECS reboots the instance in 24 to 48 hours. |
| `SystemMaintenance.Redeploy:Inquiring`, `SystemMaintenance.Redeploy:Scheduled` | The instance should be redeployed for system maintenance. Data on the local disk may be lost. Initiate the redeployment during the repair; otherwise ECS redeploys the instance in 24 to 48 hours. |
| `Node.FaultNeedReboot.HOST:Inquiring` | A Lingjun node has an underlying fault and should be rebooted as soon as possible. |
| `Node.HardwareFault.HOST:Inquiring` | A Lingjun node has a hardware fault and requires repair. Data on the local disk may be lost. |
| `HyperNode.Node.FaultNeedReboot.HOST:Inquiring` | A fault is reported for a Lingjun super node and the node should be rebooted. ACK starts a separate repair flow for each child node and requests a single authorization when all child flows reach the repair step. |
| `HyperNode.Node.HardwareFault.HOST:Inquiring` | A hardware fault is reported for a Lingjun super node and the node requires repair. Data on the local disk may be lost. |

GPU faults, detected by NPD, are all repaired by rebooting the node:

| Incident | Meaning |
| --- | --- |
| `NvidiaXID48Error` | Double Bit ECC Error (DBE), the GPU detected an uncorrectable error. |
| `NvidiaXID62Error` | The GPU internal microcontroller stopped. |
| `NvidiaXID140Error`, `NvidiaXID95Error` | The GPU driver detected an uncorrectable error in GPU memory. |
| `NvidiaXID120Error` | An error occurred in the code running on the GPU GSP core. |
| `NvidiaXID119Error` | Timed out waiting for the GSP core to respond to an RPC message. |
| `NvidiaXID31Error` | GPU memory page fault. |
| `NvidiaPendingRetiredPages` | The GPU has retired pages in a pending state. |
| `NvidiaRemappingRowsRequireReset` | GPU row remapping is pending and is performed on the next GPU reset. |

## Repair procedure reference

| Procedure | Description | Dependency |
| --- | --- | --- |
| `QuarantineNode` | Marks the node as unschedulable (`spec.unschedulable=true`). | None |
| `Evict` | Evicts the pods on the node one by one and respects PodDisruptionBudgets. | Requires `QuarantineNode` |
| `Repair` | Runs the repair action of the fault type, for example restarting kubelet or the container runtime, or rebooting or redeploying the node. | None |
| `UnquarantineNode` | Makes the node schedulable again (`spec.unschedulable=false`). | Requires `QuarantineNode` |

The procedures run in the order of the built-in template, must be a subsequence of that order, and cannot be repeated. `Repair` is required by every template except `SystemComponentsOnlyRebootError`, where all procedures are optional.

`QuarantineNode` must be listed before `Evict` and before `UnquarantineNode`, because both procedures depend on it. The provider checks this ordering, the duplicate procedures and the JSON format of `config` before the request is sent, and the API applies the same dependency check together with the required procedures of the built-in template.

Step behavior:

* **Node draining** (`Evict`) respects PodDisruptionBudgets. Pods managed by a DaemonSet or an OpenKruise Advanced DaemonSet and pods labeled `alibabacloud.com/evict-pod=skip` are not evicted. Draining stops after `evict_timeout` and the flow continues with the next step even if some pods were not evicted.
* **Automatic repair** (`Repair`) performs the repair action of the fault type. For system and component faults, if `restart_node` is `true`, ACK reboots the node when the anomaly persists after the repair action.
* **UnquarantineNode** does not restore a node that was already unschedulable before the repair flow started.
* **Manual approval** (`intervention`) pauses the procedure until the approval label is added to the node. While waiting, all subsequent procedures of the same policy are paused as well.

`config` keys accepted by each procedure:

| Procedure | Key | Type | Description |
| --- | --- | --- | --- |
| `Evict` | `evict_timeout` | Integer | Draining timeout in seconds. Value range: `[0, 3600]`. Default value: `300`. `Evict` accepts no other key. |
| `Repair` | `restart_node` | Boolean | Whether to reboot the node when the repair action does not resolve a system or component fault. Default value: `false`. Only accepted by the `SystemComponentsError` template. |

## Behaviour notes

* If the repair flow fails, the node pool stops initiating new repairs until the fault is resolved. Repair the node manually and contact technical support if needed.
* A node that was already unschedulable before the repair flow started is not made schedulable afterwards.
* Nodes in a node pool are repaired one at a time by default, and circuit breaking is disabled by default. Both are configurable through the node pool `auto_repair_policy` block.
* Builtin policies are read-only and cannot be created, updated or deleted. Two builtin policies are provided: `r-repair-sys-comp-failure` (allows rebooting the node on failure) and `r-repair-sys-comp-failure-without-reboots` (does not reboot the node). Bind them to a node pool through `auto_repair_policy_id`.

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the ACK cluster.
* `name` - (Required) The name of the auto repair policy.
* `resource_type` - (Optional, ForceNew, Available since v2.3.0) The type of the resource that the policy applies to. Default value: `nodepool`. Valid values: `nodepool`.
* `resource_sub_type` - (Required, ForceNew) The subtype of the node pool that the policy applies to. Valid values: `ess`, `lingjun`.
* `rules` - (Required, Set) The repair rules. See [`rules`](#rules) below. At least one rule is required.

### `rules`

The rules support the following:

* `incidents` - (Required, Set) The incidents that trigger the repair. See [`incidents`](#rules-incidents) below.
* `repair_procedure` - (Required, List) The repair procedures executed in order when the incidents are detected. See [`repair_procedure`](#rules-repair_procedure) below. The procedures must be a subsequence of the built-in template of the incident, and the procedures marked as required in the template must be included. A procedure cannot be repeated, and `QuarantineNode` must be listed before `Evict` and `UnquarantineNode`.

### `rules-incidents`

The incidents support the following:

* `name` - (Required) The name of the incident. It must be a built-in system incident supported by `resource_sub_type`, see [Incident reference](#incident-reference), for example `NTPProblem`.

### `rules-repair_procedure`

The repair_procedure supports the following:

* `name` - (Required) The name of the repair procedure, one of `QuarantineNode`, `Evict`, `Repair` or `UnquarantineNode`. Each procedure can be used once per rule. See [Repair procedure reference](#repair-procedure-reference) for the dependency and required procedures of each template.
* `config` - (Optional) The configuration of the procedure as a JSON object, for example `jsonencode({ evict_timeout = 300 })` for `Evict`. An empty value and an empty JSON object are equivalent. See the `config` key table in [Repair procedure reference](#repair-procedure-reference) for the accepted keys.
* `intervention` - (Optional, Set, MaxItems: 1) The manual approval configuration of the procedure. See [`intervention`](#rules-repair_procedure-intervention) below.

### `rules-repair_procedure-intervention`

The intervention supports the following:

* `enable` - (Optional) Whether manual approval is required before the procedure runs. Default value: `false`.
* `type` - (Optional) The approval type. Default value: `label`. Valid values: `label`.
* `inquiring_label` - (Optional, Set, MaxItems: 1) The label written to the node to request approval. Default key: `intervention.autorepair.node.k8s.aliyun.com/<procedure_name>`, with the `hypernode` domain for Lingjun super nodes. Default value: `inquiring`. See [`inquiring_label`](#rules-repair_procedure-intervention-inquiring_label) below.
* `approved_label` - (Optional, Set, MaxItems: 1) The label that marks the request as approved. Default value: `approved`. See [`approved_label`](#rules-repair_procedure-intervention-approved_label) below.

When manual approval is enabled, the procedure pauses after the inquiring label is written, and all subsequent procedures of the same policy wait until the approved label is written to the node.

### `rules-repair_procedure-intervention-inquiring_label`

The inquiring_label supports the following:

* `key` - (Required) The key of the label. Both labels must use the same key.
* `value` - (Required) The value of the label.

### `rules-repair_procedure-intervention-approved_label`

The approved_label supports the following:

* `key` - (Required) The key of the label. Both labels must use the same key.
* `value` - (Required) The value of the label. It must differ from the value of the inquiring label.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the auto repair policy, in the format of `<cluster_id>:<policy_id>`.
* `policy_id` - The ID of the auto repair policy.
* `type` - The type of the auto repair policy. Valid values: `custom`, `builtin`.
* `resource_ids` - The IDs of the node pools that the policy is bound to. It is maintained by the server and reflects the bindings configured on the node pools.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the auto repair policy.
* `update` - (Defaults to 5 mins) Used when updating the auto repair policy.
* `delete` - (Defaults to 5 mins) Used when deleting the auto repair policy.

## Import

Auto repair policy can be imported using the cluster ID and the policy ID, separated by a colon (`:`), e.g.

```shell
$ terraform import alicloud_cs_auto_repair_policy.example <cluster_id>:<policy_id>
```
