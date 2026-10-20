---
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_tasks"
sidebar_current: "docs-alicloud-datasource-pai-feature-store-tasks"
description: |-
    Provides a list of PaiFeatureStore FeatureStoreTask resources.
---

# alicloud_pai_feature_store_tasks

This data source provides a list of PaiFeatureStore FeatureStoreTask resources owned by a PaiFeatureStore instance.

When `task_id` is set, the data source fetches a single task via `GetTask`; otherwise it lists tasks via `ListTasks` with the given filters. The filter arguments are ignored when `task_id` is set. FeatureStoreTask is read-only: it has no Create/Delete/Update lifecycle API, so only query attributes are exposed.

-> **NOTE:** Available since v2.1.0.

## Example Usage

### List all tasks of an instance

```terraform
data "alicloud_pai_feature_store_tasks" "all" {
  instance_id = "fs-xxxxx"
}

output "task_ids" {
  value = [for t in data.alicloud_pai_feature_store_tasks.all.tasks : t.task_id]
}
```

### List tasks filtered by status and object type

```terraform
data "alicloud_pai_feature_store_tasks" "running" {
  instance_id = "fs-xxxxx"
  status      = "Running"
  object_type = "ModelFeature"
}
```

### Fetch a single task by id

```terraform
data "alicloud_pai_feature_store_tasks" "single" {
  instance_id = "fs-xxxxx"
  task_id     = "4"
}

output "task_status" {
  value = data.alicloud_pai_feature_store_tasks.single.tasks[0].status
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, Available since v2.1.0) The ID of the PaiFeatureStore instance that owns the tasks.
* `task_id` - (Optional, Available since v2.1.0) The ID of a single task. When set, the data source returns one task fetched by `GetTask` instead of listing, and the filter arguments below are ignored.
* `project_id` - (Optional, Available since v2.1.0) Filter by the ID of the project.
* `object_type` - (Optional, Available since v2.1.0) Filter by the object type of the task. Valid values: `ModelFeature`, `FeatureView`.
* `object_id` - (Optional, Available since v2.1.0) Filter by the ID of the task object.
* `status` - (Optional, Available since v2.1.0) Filter by the status of the task. Valid values: `Initializing`, `Running`, `Success`, `Failure`.
* `type` - (Optional, Available since v2.1.0) Filter by the type of the task. Valid values: `OfflineToOnline`, `ExportTrainingSet`.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id` - (Available since v2.1.0) The ID of the data source. It is the `instance_id` when listing, or `instance_id/task_id` when fetching a single task.
* `tasks` - (Available since v2.1.0) A list of FeatureStoreTask resources. Each element contains the following attributes:
  * `task_id` - (Available since v2.1.0) The ID of the task.
  * `instance_id` - (Available since v2.1.0) The ID of the PaiFeatureStore instance that owns the task.
  * `project_id` - (Available since v2.1.0) The ID of the project.
  * `project_name` - (Available since v2.1.0) The name of the project.
  * `object_type` - (Available since v2.1.0) The object type of the task.
  * `object_id` - (Available since v2.1.0) The ID of the task object.
  * `status` - (Available since v2.1.0) The status of the task.
  * `type` - (Available since v2.1.0) The type of the task.
  * `create_time` - (Available since v2.1.0) The creation time (`GmtCreateTime`).
  * `executed_time` - (Available since v2.1.0) The execution time (`GmtExecutedTime`).
  * `finished_time` - (Available since v2.1.0) The finish time (`GmtFinishedTime`).
  * `modified_time` - (Available since v2.1.0) The modification time (`GmtModifiedTime`), only populated by `GetTask`.
  * `config` - (Available since v2.1.0) The task config, only populated by `GetTask`.
  * `running_config` - (Available since v2.1.0) The task running config, only populated by `GetTask`.
