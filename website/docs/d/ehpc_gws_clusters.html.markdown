---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_gws_clusters"
sidebar_current: "docs-alicloud-datasource-ehpc-gws-clusters"
description: |-
  Provides a list of Ehpc Gws Clusters to the user.
---

# alicloud\_ehpc\_gws\_clusters

This data source provides the Ehpc Gws Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ehpc_gws_clusters" "ids" {
  ids = ["example_value"]
}
output "ehpc_gws_cluster_id_1" {
  value = data.alicloud_ehpc_gws_clusters.ids.clusters.0.id
}

data "alicloud_ehpc_gws_clusters" "status" {
  status = "running"
}
output "ehpc_gws_cluster_id_2" {
  value = data.alicloud_ehpc_gws_clusters.status.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Gws Cluster IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `creating`, `deleted`, `running`, `starting`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `clusters` - A list of Ehpc Gws Clusters. Each element contains the following attributes:
	* `create_time` - Visualize cluster creation time.
	* `gws_cluster_id` - The first ID of the resource.
	* `id` - The ID of the Gws Cluster.
	* `status` - The status of the resource.
	* `vpc_id` - The ID of the VPC.