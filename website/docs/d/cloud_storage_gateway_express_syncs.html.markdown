---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_express_syncs"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-express-syncs"
description: |-
  Provides a list of Cloud Storage Gateway Express Syncs to the user.
---

# alicloud\_cloud\_storage\_gateway\_express\_syncs

This data source provides the Cloud Storage Gateway Express Syncs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_express_syncs" "ids" {}
output "cloud_storage_gateway_express_sync_id_1" {
  value = data.alicloud_cloud_storage_gateway_express_syncs.ids.syncs.0.id
}

data "alicloud_cloud_storage_gateway_express_syncs" "nameRegex" {
  name_regex = "^my-ExpressSync"
}
output "cloud_storage_gateway_express_sync_id_2" {
  value = data.alicloud_cloud_storage_gateway_express_syncs.nameRegex.syncs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Express Sync IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Express Sync name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Express Sync names.
* `syncs` - A list of Cloud Storage Gateway Express Syncs. Each element contains the following attributes:
    * `bucket_name` - The name of the OSS Bucket.
    * `bucket_prefix` - The prefix of the OSS Bucket.
    * `bucket_region` - The region of the OSS Bucket.
    * `description` - The description of the Express Sync.
    * `express_sync_id` - The ID of the Express Sync.
    * `express_sync_name` - The name of the Express Sync.
    * `id` - The resource ID in terraform of Express Sync. The value is formate as <express_sync_id>.
    * `mns_topic` - The name of the message topic (Topic) corresponding to the Express Sync in the Alibaba Cloud Message Service MNS.