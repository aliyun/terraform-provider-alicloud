---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_namespace"
sidebar_current: "docs-alicloud-resource-cms-namespace"
description: |-
  Provides a Alicloud Cloud Monitor Service Namespace resource.
---

# alicloud\_cms\_namespace

Provides a Cloud Monitor Service Namespace resource.

For information about Cloud Monitor Service Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/doc-detail/28608.htm).

-> **NOTE:** Available in v1.171.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_namespace" "example" {
  namespace     = "example_value"
  specification = "cms.s1.large"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of indicator warehouse.
* `namespace` - (Required, ForceNew) Indicator warehouse name. The namespace can contain lowercase letters, digits, and hyphens (-).
* `specification` - (Optional, Computed) Data storage duration. Valid values: `cms.s1.12xlarge`, `cms.s1.2xlarge`, `cms.s1.3xlarge`, `cms.s1.6xlarge`, `cms.s1.large`, `cms.s1.xlarge`. 
  - `cms.s1.large`: Data storage duration is 15 days. 
  - `cms.s1.xlarge`: Data storage duration is 32 days. 
  - `cms.s1.2xlarge`: Data storage duration 63 days.
  - `cms.s1.3xlarge`: (Default) Data storage duration 93 days.
  - `cms.s1.6xlarge`: Data storage duration 185 days.
  - `cms.s1.12xlarge`: Data storage duration 376 days.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Namespace.
* `delete` - (Defaults to 1 mins) Used when delete the Namespace.
* `update` - (Defaults to 1 mins) Used when update the Namespace.

## Import

Cloud Monitor Service Namespace can be imported using the id, e.g.

```
$ terraform import alicloud_cms_namespace.example <namespace>
```