---
subcategory: "PAI LangStudio"
layout: "alicloud"
page_title: "Alicloud: alicloud_pailangstudio_knowledge_base_job"
sidebar_current: "docs-alicloud-resource-pailangstudio-knowledge-base-job"
description: |-
  Provides a PAI LangStudio Knowledge Base Job resource.
---

# alicloud_pailangstudio_knowledge_base_job

Provides a PAI LangStudio Knowledge Base Job resource.

For information about PAI LangStudio Knowledge Base Job and how to use it, see
[CreateKnowledgeBaseJob](https://www.alibabacloud.com/help/en/pai/developer-reference/api-pailangstudio-2024-07-10-createknowledgebasejob).

-> **NOTE:** Available since v1.295.0.

The resource identity is composite (`<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>`)
because a knowledge base job is addressed by all three IDs together.

## Example Usage

Basic usage:

```terraform
resource "alicloud_pailangstudio_knowledge_base_job" "default" {
  workspace_id      = "ws-xxxxxx"
  knowledge_base_id = "kb-xxxxxx"
  description       = "my knowledge base job"
}
```

With an embedding configuration:

```terraform
resource "alicloud_pailangstudio_knowledge_base_job" "default" {
  workspace_id      = "ws-xxxxxx"
  knowledge_base_id = "kb-xxxxxx"
  description       = "my knowledge base job"
  job_action        = "Sync"

  embedding_config {
    batch_size  = 10
    concurrency = 8
  }
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID. Part of the resource identity.
* `knowledge_base_id` - (Required, ForceNew) The knowledge base ID. Part of the resource identity and the Create path.
* `description` - (Optional) The description of the knowledge base job. This is the only field updatable via `UpdateKnowledgeBaseJob`.
* `accessibility` - (Optional, ForceNew) The accessibility of the job. Valid values: `PUBLIC`, `PRIVATE`. Defaults to `PUBLIC`.
* `job_action` - (Optional, ForceNew) The job action.
* `max_running_time_in_seconds` - (Optional, ForceNew) The maximum running time of the job in seconds.
* `resource_id` - (Optional, ForceNew) The prepaid resource ID.
* `ecs_specs` - (Optional, ForceNew) The ECS specifications. See [`ecs_specs`](#ecs_specs) below.
* `user_vpc` - (Optional, ForceNew) The user VPC configuration. See [`user_vpc`](#user_vpc) below.
* `embedding_config` - (Optional, ForceNew) The embedding configuration. See [`embedding_config`](#embedding_config) below.

### ecs_specs

* `type` - (Optional) The node type. Defaults to `Worker`.
* `instance_type` - (Optional) The instance type.
* `pod_count` - (Optional) The pod count.

### user_vpc

* `vpc_id` - (Optional) The VPC ID.
* `v_switch_id` - (Optional) The vSwitch ID.
* `security_group_id` - (Optional) The security group ID.

### embedding_config

* `batch_size` - (Optional) The embedding batch size.
* `concurrency` - (Optional) The embedding concurrency. Valid values: 1 to 200.

## Attributes Reference

The following attributes are exported:

* `id` - The composite ID of the resource, formatted as `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>`.
* `knowledge_base_job_id` - The knowledge base job ID, assigned by the server on creation.
* `status` - The job status. Valid values: `Running`, `Success`, `Failed`.
* `gmt_create_time` - The creation time.
* `gmt_modified_time` - The modification time.
* `gmt_finish_time` - The finish time.
* `creator` - The creator.
* `error_message` - The job error message, populated when the job fails.
* `region_id` - The region ID.

## Import

PAI LangStudio Knowledge Base Job can be imported using the composite id, e.g.

```shell
$ terraform import alicloud_pailangstudio_knowledge_base_job.default <workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>
```
