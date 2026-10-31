---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_rc_cluster"
sidebar_current: "docs-alicloud-resource-rds-rc-cluster"
description: |-
  Provides a Alicloud RDS Custom (RC) Container Cluster resource.
---

# alicloud_rds_rc_cluster

Provides a RDS Custom (RC) Container Cluster resource.

The resource imports an existing ACK managed cluster into RDS Custom so that RDS Custom instances and node pools can run on it. There is no update operation: every argument is `ForceNew`, so a change deletes the import and imports the cluster again.

For information about RDS Custom Container Cluster and how to use it, see [What is RC Container Cluster](https://www.alibabacloud.com/help/en/rds/apsaradb-rds-for-mysql/managing-rds-custom-container-clusters).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
variable "ack_cluster_id" {
  description = "The ID of the ACK managed cluster to import into RDS Custom."
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC that the ACK cluster belongs to."
  type        = string
}

variable "v_switch_id" {
  description = "The ID of the vSwitch that the ACK cluster belongs to."
  type        = string
}

data "alicloud_cs_cluster_credential" "default" {
  cluster_id = var.ack_cluster_id
}

resource "alicloud_rds_rc_cluster" "default" {
  cluster_id         = var.ack_cluster_id
  vpc_id             = var.vpc_id
  v_switch_id        = var.v_switch_id
  client_key         = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_key
  client_certificate = data.alicloud_cs_cluster_credential.default.certificate_authority.0.client_cert
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the ACK managed cluster to import into RDS Custom.
* `vpc_id` - (Required, ForceNew) The ID of the VPC that the ACK cluster belongs to.
* `v_switch_id` - (Required, ForceNew) The ID of the vSwitch that the ACK cluster belongs to.
* `client_key` - (Required, ForceNew, Sensitive) The `client-key-data` of the ACK cluster KubeConfig.
* `client_certificate` - (Required, ForceNew, Sensitive) The `client-certificate-data` of the ACK cluster KubeConfig.
* `region_id` - (Optional, ForceNew) The region ID of the cluster. Defaults to the region of the provider.
* `support_disk_performance_level` - (Optional) The disk performance levels (for example `PL1`) that the cluster supports.
* `support_disk_category` - (Optional) The disk categories (for example `cloud_essd`) that the cluster supports.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the ACK cluster ID.
* `cluster_name` - The name of the cluster.
* `vcluster_status` - The status of the cluster. `UNINITIALIZED` means the cluster is not imported into RDS Custom.

## Timeouts

The resource uses fixed internal polling windows of 30 minutes for create and 10 minutes for delete; no configurable timeouts are exposed.

## Import

RDS Custom Container Cluster can be imported using the ACK cluster ID, e.g.

```shell
$ terraform import alicloud_rds_rc_cluster.example <cluster_id>
```
