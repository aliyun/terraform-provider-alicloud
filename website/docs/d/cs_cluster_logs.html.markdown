---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_cluster_logs"
sidebar_current: "docs-alicloud-datasource-cs-cluster-logs"
description: |-
Provide logs, events or tasks of kubernetes cluster.
---

# alicloud\_cs\_cluster\_logs

This data source provides logs, events or tasks of kubernetes cluster.

-> **NOTE:** Available in 1.186.0+.

## Example Usage

```terraform
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex     = "my-cluster"
  enable_details = false
}
data "alicloud_cs_cluster_logs" "default" {
  cluster_id  = data.alicloud_cs_managed_kubernetes_clusters.default.ids.0
  type        = "log"
  entries     = 100
  output_file = "cluster.log"
}
```

## Argument Reference

The following arguments are supported.
* `cluster_id` - (**Required**) The id of kubernetes cluster.
* `type` - (Optional) The type of cluster log you want to export. It's valid value is `log`, `event` or `task`. Default to `log`.
* `entries` - (Optional) The number of log entries you want to export. Please refer to the attribute value of this field for the actual number of logs obtained. Default to `100`.
* `output_file` - (Optional) The path you want to persist the log to. If this value is not set, the log will not be written to the state file.
