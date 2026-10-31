---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_rc_cluster_node_pool"
sidebar_current: "docs-alicloud-resource-rds-rc-cluster-node-pool"
description: |-
  Provides a Alicloud RDS Custom (RC) Cluster Node Pool resource.
---

# alicloud_rds_rc_cluster_node_pool

Provides a RDS Custom (RC) Cluster Node Pool resource.

A node pool is the scaling group that provides the ACK worker nodes which carry RDS Custom instances of a RDS Custom container cluster. The cluster must be an ACK managed cluster already imported into RDS Custom; this resource cannot create the cluster itself.

For information about RDS Custom Cluster Node Pool and how to use it, see [What is RC Cluster Node Pool](https://www.alibabacloud.com/help/en/rds/apsaradb-rds-for-mysql/managing-rds-custom-container-clusters).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
variable "cluster_id" {
  description = "The ID of an existing RDS Custom container cluster, that is, an ACK cluster imported into RDS Custom."
  type        = string
}

resource "alicloud_rds_rc_cluster_node_pool" "default" {
  cluster_id = var.cluster_id
  name       = "terraform-example"

  scaling_group = [
    {
      v_switch_ids         = ["vsw-xxxxxxxx"]
      instance_types       = ["mysql.x4.xlarge.6cm"]
      instance_charge_type = "PostPaid"
      period               = 1
      period_unit          = "Month"
      auto_renew           = false
      security_group_ids   = ["sg-xxxxxxxx"]
      desired_size         = 0
    }
  ]

  kubernetes_config = [
    {
      labels = [
        {
          key   = "env"
          value = "prod"
        }
      ]
      taints = [
        {
          key    = "dedicated"
          value  = "rds"
          effect = "NoSchedule"
        }
      ]
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the RDS Custom container cluster that owns the node pool.
* `name` - (Optional) The name of the node pool.
* `scaling_group` - (Optional) The scaling group configuration of the node pool. See [`scaling_group`](#scaling_group) below. If omitted, the value returned by the API is used.
* `kubernetes_config` - (Optional) The Kubernetes configuration of the node pool. See [`kubernetes_config`](#kubernetes_config) below.

### `scaling_group`

The `scaling_group` block supports:

* `v_switch_ids` - (Optional) The IDs of the vSwitches used by the nodes.
* `instance_types` - (Optional) The instance types available to the node pool.
* `instance_charge_type` - (Optional) The billing method of the nodes. Valid values: `PrePaid` and `PostPaid`.
* `period` - (Optional) The subscription duration. Default value: `1`. Only valid when `instance_charge_type` is `PrePaid`.
* `period_unit` - (Optional) The unit of the subscription duration. Valid values: `Month` and `Year`. Default value: `Month`. Only valid when `instance_charge_type` is `PrePaid`.
* `auto_renew` - (Optional) Whether the nodes are renewed automatically. Default value: `true`. Only valid when `instance_charge_type` is `PrePaid`.
* `security_group_ids` - (Optional) The IDs of the security groups used by the nodes.
* `desired_size` - (Optional) The expected number of nodes. Default value: `0`.

### `kubernetes_config`

The `kubernetes_config` block supports:

* `labels` - (Optional) The labels attached to the nodes. See [`labels`](#labels) below.
* `taints` - (Optional) The taints attached to the nodes. See [`taints`](#taints) below.

### `labels`

The `labels` block supports:

* `key` - (Optional) The key of the label.
* `value` - (Optional) The value of the label.

### `taints`

The `taints` block supports:

* `key` - (Optional) The key of the taint.
* `value` - (Optional) The value of the taint.
* `effect` - (Optional) The scheduling effect of the taint. Valid values: `NoSchedule`, `NoExecute` and `PreferNoSchedule`. Default value: `NoSchedule`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, in the format `<cluster_id>:<node_pool_id>`.
* `node_pool_id` - The ID of the node pool.
* `region_id` - The region ID of the node pool.
* `created` - The time when the node pool was created.
* `updated` - The time when the node pool was last updated.
* `status` - The status of the node pool. See [`status`](#status) below.

### `status`

The `status` block exports:

* `state` - The state of the node pool.
* `total_nodes` - The total number of nodes.
* `retry_times` - The number of scaling retries.
* `scaling_error_description` - The description of the last scaling error.

## Timeouts

The resource uses fixed internal polling windows of 30 minutes for create, update and delete; no configurable timeouts are exposed.

## Import

RDS Custom Cluster Node Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_rc_cluster_node_pool.example <cluster_id>:<node_pool_id>
```
