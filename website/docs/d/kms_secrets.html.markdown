---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secrets"
sidebar_current: "docs-alicloud-datasource-kms-secrets"
description: |-
  Provides a list of available KMS Secrets.
---

# alicloud_kms_secrets

This data source provides a list of KMS Secrets in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available since v1.86.0.

## Example Usage

```terraform
# Declare the data source
data "alicloud_kms_secrets" "kms_secrets_ds" {
  fetch_tags = true
  name_regex = "name_regex"
  tags = {
    "k-aa" = "v-aa",
    "k-bb" = "v-bb"
  }
}

output "first_secret_id" {
  value = data.alicloud_kms_secrets.kms_secrets_ds.secrets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `fetch_tags` - (Optional, ForceNew) Whether to include the predetermined resource tag in the return value. Default to `false`.
* `ids` - (Optional, ForceNew) A list of KMS Secret ids. The value is same as KMS secret_name.
* `name_regex` - (Optional, ForceNew) A regex string to filter the results by the KMS secret_name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A mapping of tags to assign to the resource, and can be used to filter secrets.
* `filters` - (Optional, ForceNew, Available since v1.124.0) The secret filter. The filter consists of one or more key-value pairs. 
  More details see API [ListSecrets](https://www.alibabacloud.com/help/en/key-management-service/latest/listsecrets). 
* `enable_details` - (Optional, Available since v1.124.0) Default to `false`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of Kms Secret ids. The value is same as KMS secret_name. 
* `names` -  A list of KMS Secret names.
* `secrets` - A list of KMS Secrets. Each element contains the following attributes:
  * `id` - ID of the Kms Secret. The value is same as KMS secret_name.
  * `secret_name` - Name of the KMS Secret.
  * `planned_delete_time` - Schedule deletion time.
  * `arn` - (Available since v1.124.0) A mapping of tags to assign to the resource.
  * `description` - (Available since v1.124.0)  The description of the secret.
  * `encryption_key_id` - (Available since v1.124.0)  The ID of the KMS CMK that is used to encrypt the secret value.
  * `secret_data` - (Available since v1.124.0)  The value of the secret that you want to create.
  * `secret_data_type` - (Available since v1.124.0)  The type of the secret data value.
  * `secret_type` - (Available since v1.124.0)  The type of the secret.
  * `version_id` - (Available since v1.124.0)  The version number of the initial version.
  * `version_stages` - (Available since v1.124.0)  The stage labels that mark the new secret version.
  * `tags` - A mapping of tags to assign to the resource.

