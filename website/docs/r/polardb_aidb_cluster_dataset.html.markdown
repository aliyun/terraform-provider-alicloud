---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_aidb_cluster_dataset"
sidebar_current: "docs-alicloud-resource-polardb-aidb-cluster-dataset"
description: |-
  Provides a PolarDB AI Cluster Dataset resource.
---

# alicloud_polardb_aidb_cluster_dataset

Provides a PolarDB AI Cluster Dataset resource.

-> **NOTE:** Available since v2.1.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_polardb_aidb_cluster_dataset" "example" {
  db_cluster_id = "pc-xxxxxxxxxxxxxxxx"
  dataset_name  = var.name
  dataset_type  = "OSS"
  train_mode    = "SFT"
  import_mode   = "AUTO"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI cluster that the dataset belongs to.
* `dataset_name` - (Required, ForceNew, Available since v2.1.0) The name of the dataset.
* `dataset_type` - (Optional, ForceNew, Available since v2.1.0) The type of the dataset.
* `train_mode` - (Optional, ForceNew, Available since v2.1.0) The training mode of the dataset.
* `import_mode` - (Optional, ForceNew, Available since v2.1.0) The import mode of the dataset.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the dataset.
* `region_id` - (Available since v2.1.0) The region ID.
* `relative_db_cluster_id` - (Available since v2.1.0) The ID of the associated database cluster.
* `path` - (Available since v2.1.0) The storage path of the dataset.
* `data_service_id` - (Available since v2.1.0) The ID of the data service.
* `file_name` - (Available since v2.1.0) The name of the dataset file.
* `bucket_name` - (Available since v2.1.0) The name of the OSS bucket that stores the dataset.
* `storage_type` - (Available since v2.1.0) The storage type of the dataset.
* `capacity` - (Available since v2.1.0) The capacity of the dataset.
* `last_modified` - (Available since v2.1.0) The time when the dataset was last modified.
* `creation_time` - (Available since v2.1.0) The time when the dataset was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the dataset.
* `delete` - (Defaults to 30 mins) Used when deleting the dataset.

## Import

PolarDB AI Cluster Dataset can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_aidb_cluster_dataset.example <db_cluster_id>:<dataset_id>
```
