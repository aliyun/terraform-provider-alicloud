---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key_versions"
sidebar_current: "docs-alicloud-datasource-kms-key-versions"
description: |-
    Provides a list of available KMS KeyVersions.
---

# alicloud\_kms\_key\_versions

This data source provides a list of KMS KeyVersions in an Alibaba Cloud account according to the specified filters.

-> NOTE: Available in v1.85.0+

## Example Usage

```
# Declare the data source
data "alicloud_kms_key_versions" "alicloud_kms_key_versions_ds" {
  key_id = "08438c-b4d5-4d05-928c-07b7xxxx"
  ids    = ["d89e8a53-b708-41aa-8c67-6873axxx"]
}

output "all_versions" {
  value = "${data.alicloud_kms_key_versions.alicloud_kms_key_versions_ds.versions}"
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required) The id of kms key.
* `ids` - (Optional) A list of KMS KeyVersion IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of KMS KeyVersion IDs.
* `versions` - A list of KMS KeyVersions. Each element contains the following attributes:
  * `creation_date` - (Removed from v1.124.4) It has been removed and using `create_time` instead.
  * `create_time` - Date and time when the key version was created (UTC time).
  * `key_id` - ID of the key.
  * `id` - ID of the KMS KeyVersion resource.
  * `key_version_id` - ID of the key version.
