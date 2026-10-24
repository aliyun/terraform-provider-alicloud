---
subcategory: "Bailian (SFM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_bailian_connector"
sidebar_current: "docs-alicloud-resource-bailian-connector"
description: |-
  Provides a Bailian Connector resource.
---

# alicloud_bailian_connector

Provides a Bailian Connector resource.

For information about Bailian Connector and how to use it, see [What is Connector](https://www.alibabacloud.com/help/en/model-studio/developer-api/api-reference-12).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bailian_connector" "default" {
  workspace_id   = "llm-xxx"
  connector_name = "tf-example-connector"
  connector_type = "FILE"
  description    = "tf-example-connector-description"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID of Bailian.
* `connector_name` - (Required) The name of the connector.
* `connector_type` - (Required, ForceNew) The type of the connector. Valid value: `FILE`.
* `description` - (Optional) The description of the connector.
* `file_connector_config` - (Optional) The file connector config, a JSON string with fields `bucket_name`, `region_id`, `storage_type`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, format as `<workspace_id>:<connector_id>`.
* `connector_id` - The ID of the connector.

## Import

Bailian Connector can be imported using the id, e.g.

```shell
$ terraform import alicloud_bailian_connector.example <workspace_id>:<connector_id>
```
