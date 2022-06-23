---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_namespace"
sidebar_current: "docs-alicloud-resource-edas-namespace"
description: |-
  Provides a Alicloud EDAS Namespace resource.
---

# alicloud\_edas\_namespace

Provides a EDAS Namespace resource.

For information about EDAS Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/enterprise-distributed-application-service/latest/insertorupdateregion).

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_edas_namespace" "example" {
  namespace_logical_id = "example_value"
  namespace_name       = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `debug_enable` - (Optional, Computed) Specifies whether to enable remote debugging.
* `description` - (Optional, Computed) The description of the namespace, The description can be up to `128` characters in length.
* `namespace_logical_id` - (Required, ForceNew) The ID of the namespace.
  - The ID of a custom namespace is in the `region ID:namespace identifier` format. An example is `cn-beijing:tdy218`.
  - The ID of the default namespace is in the `region ID` format. An example is cn-beijing.
* `namespace_name` - (Required, Computed) The name of the namespace, The name can be up to `63` characters in length.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Namespace.
* `delete` - (Defaults to 1 mins) Used when delete the Namespace.
* `update` - (Defaults to 1 mins) Used when update the Namespace.

## Import

EDAS Namespace can be imported using the id, e.g.

```
$ terraform import alicloud_edas_namespace.example <id>
```