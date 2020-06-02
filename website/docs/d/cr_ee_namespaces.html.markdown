---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_namespaces"
sidebar_current: "docs-alicloud-datasource-cr-ee-namespaces"
description: |-
  Provides a list of Container Registry Enterprise namespaces.
---

# alicloud\_cr_ee\_namespaces

This data source provides a list Container Registry Enterprise namespaces on Alibaba Cloud.

-> **NOTE:** Available in v1.85.0+

## Example Usage

```
# Declare the data source
data "alicloud_cr_ee_namespaces" "my_namespaces" {
  instance_id = "cri-xxx"
  name_regex  = "my-namespace"
  output_file = "my-namespace-json"
}

output "output" {
  value = "${data.alicloud_cr_ee_namespaces.my_namespaces.namespaces}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of CR EE instance.
* `name_regex` - (Optional) A regex string to filter results by namespace name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched CR EE namespaces. Its element is a namespace uuid.
* `names` - A list of namespace names.
* `namespaces` - A list of matched CR EE namespaces. Each element contains the following attributes:
  * `instance_id` - ID of CR EE instance.
  * `id` - ID of CR EE namespace.
  * `name` - Name of CR EE namespace.
  * `auto_create` - Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
  * `default_visibility` - `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

