---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emr_main_versions"
sidebar_current: "docs-alicloud-datasource-emr-main-versions"
description: |-
    Provides a collection of emr main versions when create emr cluster according to the specified filters.
---

# alicloud\_emr\_main\_versions

The `alicloud_emr_main_versions` data source provides a collection of emr 
main versions available in Alibaba Cloud account when create a emr cluster.

-> **NOTE:** Available in 1.59.0+

## Example Usage

```
data "alicloud_emr_main_versions" "default" {
  emr_version  = "EMR-3.22.0"
  cluster_type = ["HADOOP", "ZOOKEEPER"]
}

output "first_main_version" {
  value = "${data.alicloud_emr_main_versions.default.main_versions.0.emr_version}"
}

output "this_cluster_types" {
  value = "${data.alicloud_emr_main_versions.default.main_versions.0.cluster_types}"
}
```

## Argument Reference

The following arguments are supported:

* `emr_version` - (Optional) The version of the emr cluster instance. Possible values: `EMR-4.0.0`, `EMR-3.23.0`, `EMR-3.22.0`.
* `cluster_type` - (Optional, Available in 1.70.1+) The supported clusterType of this emr version.
Possible values may be any one or combination of these: ["HADOOP", "DRUID", "KAFKA", "ZOOKEEPER", "FLINK", "CLICKHOUSE"]
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of emr instance types IDs. 
* `main_versions` - A list of versions of the emr cluster instance. Each element contains the following attributes:
  * `emr_version` - The version of the emr cluster instance.
  * `image_id` - The image id of the emr cluster instance.
  * `cluster_types` - A list of cluster types the emr cluster supported. Possible values: `HADOOP`, `ZOOKEEPER`, `KAFKA`, `DRUID`.
