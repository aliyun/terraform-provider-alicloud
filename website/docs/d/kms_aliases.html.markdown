---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_aliases"
sidebar_current: "docs-alicloud-datasource-kms-aliases"
description: |-
    Provides a list of available KMS Aliases.
---

# alicloud_kms_aliases

This data source provides a list of KMS aliases in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available since v1.79.0.

## Example Usage

```terraform
# Declare the data source
data "alicloud_kms_aliases" "kms_aliases" {
  ids        = ["d89e8a53-b708-41aa-8c67-6873axxx"]
  name_regex = "alias/tf-example"
}

output "first_key_id" {
  value = "${data.alicloud_kms_aliases.kms_aliases.aliases.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of KMS aliases IDs. The value is same as KMS alias_name.
* `name_regex` - (Optional, ForceNew) A regex string to filter the results by the KMS alias name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of kms aliases IDs. The value is same as KMS alias_name. 
* `names` -  A list of KMS alias name.
* `aliases` - A list of KMS User alias. Each element contains the following attributes:
  * `id` - ID of the alias. The value is same as KMS alias_name.
  * `key_id` - ID of the key.
  * `alias_name` - The unique identifier of the alias.

