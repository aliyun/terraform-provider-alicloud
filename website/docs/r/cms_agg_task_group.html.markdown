---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_agg_task_group"
description: |-
  Provides a Alicloud Cms Agg Task Group resource.
---

# alicloud_cms_agg_task_group

Provides a Cms Agg Task Group resource.

Aggregation Task Group.

For information about Cms Agg Task Group and how to use it, see [What is Agg Task Group](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateAggTaskGroup).

-> **NOTE:** Available since v1.281.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_agg_task_group&exampleId=ccdbb9d6-a06d-a1fb-dfbc-6af285501c5ff0f08f6a&activeTab=example&spm=docs.r.cms_agg_task_group.0.ccdbb9d6a0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
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
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_agg_task_group&spm=docs.r.cms_agg_task_group.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `agg_task_group_config` - (Required) The configuration of the aggregation task group.
* `agg_task_group_config_type` - (Optional) The type of the aggregation task group configuration.
* `agg_task_group_name` - (Required) The name of the aggregation task group.
* `cron_expr` - (Optional) The cron expression for scheduling when `schedule_mode` is set to `Cron`.
* `delay` - (Optional, Int) The fixed delay for scheduling.
* `description` - (Optional) The description of the aggregation task group.
* `max_retries` - (Optional, Int) The maximum number of retries for an aggregation task.
* `max_run_time_in_seconds` - (Optional, Int) The maximum retry time for an aggregation task.
* `override_if_exists` - (Optional, Bool) Specifies whether to overwrite an existing resource with the same name.
* `precheck_string` - (Optional, JsonString) The dry run configuration.
* `schedule_mode` - (Optional) The scheduling mode. Valid values: `Cron` and `FixedRate`.
* `schedule_time_expr` - (Optional) The scheduling time expression.
* `source_prometheus_id` - (Required, ForceNew) The ID of the source Prometheus instance for the aggregation task group.
* `status` - (Optional) The status of the aggregation task group. Valid values: `Running` and `Stopped`.
* `target_prometheus_id` - (Required, ForceNew) The ID of the target Prometheus instance for the aggregation task group.
* `to_time` - (Optional, Int) The UNIX timestamp for the scheduling end time.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Agg Task Group. It formats as `<source_prometheus_id>:<agg_task_group_id>`.
* `agg_task_group_id` - The ID of the aggregation task group.
* `region_id` - The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Agg Task Group.
* `update` - (Defaults to 5 mins) Used when update the Agg Task Group.
* `delete` - (Defaults to 5 mins) Used when delete the Agg Task Group.

## Import

Cms Agg Task Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_agg_task_group.example <source_prometheus_id>:<agg_task_group_id>
```
