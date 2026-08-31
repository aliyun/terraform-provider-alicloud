---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_hybrid_double_writes"
description: |-
  Provides a list of Cloud Monitor Service Hybrid Double Writes to the user.
---

# alicloud_cloud_monitor_service_hybrid_double_writes

This data source provides the Cloud Monitor Service Hybrid Double Writes of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.220.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {
}

resource "alicloud_cms_namespace" "source" {
  namespace = "your-source-namespace"
}

resource "alicloud_cms_namespace" "default" {
  namespace = "your-namespace"
}

resource "alicloud_cloud_monitor_service_hybrid_double_write" "default" {
  source_namespace = alicloud_cms_namespace.source.id
  source_user_id   = data.alicloud_account.default.id
  namespace        = alicloud_cms_namespace.default.id
  user_id          = data.alicloud_account.default.id
}

data "alicloud_cloud_monitor_service_hybrid_double_writes" "ids" {
  ids = [alicloud_cloud_monitor_service_hybrid_double_write.default.id]
}

output "cloud_monitor_service_hybrid_double_writes_id_1" {
  value = data.alicloud_cloud_monitor_service_hybrid_double_writes.ids.hybrid_double_writes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Hybrid Double Write IDs.
* `source_namespace` - (Optional, ForceNew) Source Namespace.
* `source_user_id` - (Optional, ForceNew) Source UserId.
* `namespace` - (Optional, ForceNew) Target Namespace.
* `user_id` - (Optional, ForceNew) Target UserId.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `hybrid_double_writes` - A list of Hybrid Double Writes. Each element contains the following attributes:
  * `id` - The ID of the Hybrid Double Write. It formats as `<source_namespace>:<source_user_id>`.
  * `source_namespace` - Source Namespace.
  * `source_user_id` - Source UserId.
  * `namespace` - Target Namespace.
  * `user_id` - Target UserId.
