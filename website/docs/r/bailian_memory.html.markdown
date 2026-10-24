---
subcategory: "Bailian (SFM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_bailian_memory"
sidebar_current: "docs-alicloud-resource-bailian-memory"
description: |-
  Provides a Bailian Memory resource.
---

# alicloud_bailian_memory

Provides a Bailian Memory resource.

For information about Bailian Memory and how to use it, see [What is Memory](https://www.alibabacloud.com/help/en/model-studio/developer-api/api-reference-37).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bailian_memory" "default" {
  workspace_id = "llm-xxx"
  description  = "tf-example-memory"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID of Bailian.
* `description` - (Required) The description of the memory.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, format as `<workspace_id>:<memory_id>`.
* `memory_id` - The ID of the memory.

## Import

Bailian Memory can be imported using the id, e.g.

```shell
$ terraform import alicloud_bailian_memory.example <workspace_id>:<memory_id>
```
