---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_component"
sidebar_current: "docs-alicloud-resource-data-works-component"
description: |-
  Provides a Alicloud Data Works Component resource.
---

# alicloud_data_works_component

Provides a Data Works Component resource.

For information about Data Works Component and how to use it, see [What is Component](https://www.alibabacloud.com/help/en/aliyun-dataworks/developer-reference/).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_data_works_component" "example" {
  project_id     = "320687"
  spec           = "{\"nodeType\":\"NODE_TYPE_DEFAULT\"}"
  component_type = "NODE_TYPE_DEFAULT"
  source         = "MANUAL"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew) The ID of the DataWorks project.
* `spec` - (Required) The specification of the component in JSON format. It is a FlowSpec UDF JSON string that defines the component configuration.
* `component_type` - (Optional, ForceNew) The type of the component.
* `source` - (Optional, ForceNew) The source of the component.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Component. The value formats as `<component_id>:<project_id>`.
* `component_id` - The ID of the component.

## Import

Data Works Component can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_component.example <component_id>:<project_id>
```
