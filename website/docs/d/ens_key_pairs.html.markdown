---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_key_pairs"
sidebar_current: "docs-alicloud-datasource-ens-key-pairs"
description: |-
  Provides a list of Ens Key Pairs to the user.
---

# alicloud\_ens\_key\_pairs

This data source provides the Ens Key Pairs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_key_pairs" "nameRegex" {
  version    = "example_value"
  name_regex = "^my-KeyPair"
}
output "ens_key_pair_id_1" {
  value = data.alicloud_ens_key_pairs.nameRegex.pairs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Optional, ForceNew) The name of the key pair.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Key Pair name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `version` - (Required, ForceNew) The version number.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Key Pair IDs.
* `names` - A list of Key Pair names.
* `pairs` - A list of Ens Key Pairs. Each element contains the following attributes:
	* `create_time` - The creation time of the key pair. The date format is in accordance with ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.
	* `id` - The ID of the Key Pair.
	* `key_pair_finger_print` - Fingerprint of the key pair.
	* `key_pair_name` - The name of the key pair.
	* `version` - The version number.
