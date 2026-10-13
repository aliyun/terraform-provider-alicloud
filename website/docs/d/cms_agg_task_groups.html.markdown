---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_agg_task_groups"
description: |-
  Provides a list of Cms Agg Task Groups to the user.
---

# alicloud_cms_agg_task_groups

This data source provides the Cms Agg Task Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_prometheus_instance" "default" {
  count                    = 2
  prometheus_instance_name = "${var.name}_${count.index}"
  workspace                = alicloud_cms_workspace.default.id
}

resource "alicloud_cms_agg_task_group" "default" {
  source_prometheus_id  = alicloud_cms_prometheus_instance.default.0.id
  target_prometheus_id  = alicloud_cms_prometheus_instance.default.1.id
  agg_task_group_name   = var.name
  agg_task_group_config = <<EOF
groups:
- name: "node.rules"
  interval: "60s"
  rules:
  - record: "node_namespace_pod:kube_pod_info:"
    expr: "max(label_replace(kube_pod_info{job=\"kubernetes-pods-kube-state-metrics\" }, \"pod\", \"$1\", \"pod\", \"(.*)\")) by (node, namespace, pod, cluster)"
EOF
}

data "alicloud_cms_agg_task_groups" "default" {
  source_prometheus_id = alicloud_cms_agg_task_group.default.source_prometheus_id
  ids                  = [alicloud_cms_agg_task_group.default.agg_task_group_id]
}

output "first_agg_task_group_name" {
  value = data.alicloud_cms_agg_task_groups.default.groups.0.agg_task_group_name
}
```

## Argument Reference

The following arguments are supported:

* `source_prometheus_id` - (Required) The ID of the source Prometheus instance of the aggregation task groups.
* `ids` - (Optional) A list of aggregation task group IDs. Both the plain `<agg_task_group_id>` and the resource ID format `<source_prometheus_id>:<agg_task_group_id>` are supported.
* `name_regex` - (Optional) A regex string to filter results by aggregation task group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of aggregation task group names.
* `groups` - A list of aggregation task groups. Each element contains the following attributes:
  * `id` - The resource ID in terraform of the aggregation task group. It formats as `<source_prometheus_id>:<agg_task_group_id>`.
  * `agg_task_group_id` - The ID of the aggregation task group.
  * `agg_task_group_name` - The name of the aggregation task group.
  * `source_prometheus_id` - The ID of the source Prometheus instance.
  * `target_prometheus_id` - The ID of the target Prometheus instance.
  * `cron_expr` - The cron expression for scheduling when `schedule_mode` is set to `Cron`.
  * `delay` - The fixed delay for scheduling.
  * `description` - The description of the aggregation task group.
  * `from_time` - The UNIX timestamp for the scheduling start time.
  * `interval` - The scheduling interval of the aggregation task group.
  * `max_retries` - The maximum number of retries for an aggregation task.
  * `max_run_time_in_seconds` - The maximum run time of an aggregation task.
  * `region_id` - The region ID.
  * `schedule_mode` - The scheduling mode.
  * `schedule_time_expr` - The scheduling time expression.
  * `status` - The status of the aggregation task group.
  * `to_time` - The UNIX timestamp for the scheduling end time.
  * `update_time` - The last update time of the aggregation task group.
  * `tags` - The tags of the aggregation task group.
