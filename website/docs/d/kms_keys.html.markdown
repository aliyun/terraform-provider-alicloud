---
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_keys"
sidebar_current: "docs-alicloud-datasource-kms-keys"
description: |-
    Provides a list of Availability KMS Keys.
---

# alicloud\_kms\_keys

The KMS keys data source provides a list of Alicloud KMS keys in an Alicloud account according to the specified filters.

## Example Usage

```
# Declare the data source
data "alicloud_kms_keys" "keys" {
	description_regex = "Hello KMS"
	output_file = "kms_keys.json"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of KMS key ID.
* `description_regex` - (Optional) A regex string of the KMS key description.
* `status` - (Optional) The status of KMS key. Valid values: "Enabled", "Disabled", "PendingDeletion". Default to nil to get all keys.
* `output_file` - (Optional) The name of file that can save KMS keys data source after running `terraform plan`.

## Attributes Reference

A list of KMS keys will be exported and its every element contains the following attributes:

* `id` - ID of the key.
* `arn` - The Alicloud Resource Name (ARN) of the key.
* `description` - Description of the key.
* `status` - Status of the key, with possible values: "Enabled", "Disabled", "PendingDeletion".
* `creation_date` - Creation date of key.
* `delete_date` - Delete date of key.
* `creator` - The createor to key belongs.