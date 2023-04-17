---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_double_writes"
sidebar_current: "docs-alicloud-datasource-cms-hybrid-double-writes"
description: |-
  Provides a list of Cloud Monitor Service Hybrid Double Write owned by an Alibaba Cloud account.
---

# alicloud_cms_hybrid_double_writes

This data source provides Cloud Monitor Service Hybrid Double Write available to the user.[What is Hybrid Double Write](https://next.api.aliyun.com/document/Cms/2018-03-08/CreateHybridDoubleWrite)

-> **NOTE:** Available in 1.204.0+

## Example Usage

```terraform
data "alicloud_cms_hybrid_double_writes" "default" {
  source_namespace = "test-hybrid-pool-1"
}

output "alicloud_cms_hybrid_double_write_example_id" {
  value = data.alicloud_cms_hybrid_double_writes.default.hybrid_double_writes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed)  A list of Hybrid Double Write IDs.
* `source_namespace` - (ForceNew,Optional) Double write source namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `hybrid_double_writes` - A list of Hybrid Double Write Entries. Each element contains the following attributes:
  * `namespace` - Double write target Namespace.
  * `source_namespace` - Double write source namespace.
  * `source_user_id` - The id of the source user.
  * `user_id` - The double written target user ID.
