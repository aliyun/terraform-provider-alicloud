---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_context_store_api_key"
sidebar_current: "docs-alicloud-resource-cms-context-store-api-key"
description: |-
  Provides a Cms Context Store API Key resource.
---

# alicloud_cms_context_store_api_key

Provides a Cms Context Store API Key resource.

For information about Cms Context Store API Key and how to use it, see [CreateContextStoreAPIKey](https://www.alibabacloud.com/help/en/cloudmonitor/developer-reference/api-cms-2024-03-30-createcontextstoreapikey).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_context_store_api_key" "example" {
  workspace          = "example_workspace"
  context_store_name = "example_context_store"
  name               = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required, ForceNew) The name of the workspace to which the context store belongs.
* `context_store_name` - (Required, ForceNew) The name of the context store.
* `name` - (Required, ForceNew) The display name of the API key, which is used to identify the purpose of the key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value is formulated as `<workspace>:<context_store_name>:<name>`.
* `api_key` - The value of the API key. Sensitive: this is not printed in plan or apply output.
* `create_time` - The creation time of the API key.
* `region_id` - The region ID of the API key.

## Timeouts

-> **NOTE:** The Cms API offers no Update endpoint, so `workspace`, `context_store_name` and `name` are all `ForceNew`. The resource retries transient `ServiceUnavailable` responses on Create, Read and Delete for a bounded window; there is no configurable `timeouts` block.

## Import

Cms Context Store API Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_context_store_api_key.example <workspace>:<context_store_name>:<name>
```
