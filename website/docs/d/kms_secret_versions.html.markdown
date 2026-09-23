---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secret_versions"
sidebar_current: "docs-alicloud-datasource-kms-secret-versions"
description: |-
    Provides a list of available KMS Secret Versions.
---

# alicloud\_kms\_secret\_versions

This data source provides a list of KMS Secret Versions in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available in v1.88.0+.

## Example Usage

```
# Declare the data source
data "alicloud_kms_secret_versions" "kms_secret_versions_ds" {
  secret_name = "secret_name"
  enable_details = true
}

output "first_secret_data" {
  value = "${data.alicloud_kms_secret_versions.kms_secret_versions_ds.versions.0.secret_data}"
}
```

## Argument Reference

The following arguments are supported:

* `include_deprecated` - (Optional, ForceNew) Specifies whether to return deprecated secret versions. Default to `false`.
* `ids` - (Optional, ForceNew) A list of KMS Secret Version ids.
* `secret_name` - (Required, ForceNew) The name of the secret.
* `version_stage` - (Optional, ForceNew, Available in 1.89.0+) The stage of the secret version.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to false and only output `secret_name`, `version_id`, `version_stages`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of Kms Secret Version ids. 
* `versions` - A list of KMS Secret Versions. Each element contains the following attributes:
  * `secret_data` - The secret value. Secrets Manager decrypts the stored secret value in ciphertext and returns it. (Returned when `enable_details` is true).
  * `secret_data_type` - The type of the secret value. (Returned when `enable_details` is true).
  * `secret_name` - The name of the secret.
  * `version_id` - The version number of the secret value.
  * `version_stages` - Stage labels that mark the secret version.

