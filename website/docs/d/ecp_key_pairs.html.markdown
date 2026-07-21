---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_key_pairs"
sidebar_current: "docs-alicloud-datasource-ecp-key-pairs"
description: |-
  Provides a list of Ecp Key Pairs to the user.
---

# alicloud\_ecp\_key\_pairs

This data source provides the Ecp Key Pairs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecp_key_pairs" "ids" {}
output "ecp_key_pair_id_1" {
  value = data.alicloud_ecp_key_pairs.ids.pairs.0.id
}

data "alicloud_ecp_key_pairs" "nameRegex" {
  name_regex = "^my-KeyPair"
}
output "ecp_key_pair_id_2" {
  value = data.alicloud_ecp_key_pairs.nameRegex.pairs.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Key Pair IDs. Its element value is same as Key Pair Name.
* `key_pair_finger_print` - (Optional, ForceNew) The Private Key of the Fingerprint.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Key Pair name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Key Pair names.
* `pairs` - A list of Ecp Key Pairs. Each element contains the following attributes:
	* `id` - The ID of the Key Pair. Its value is same as Queue Name.
	* `key_pair_finger_print` - The Private Key of the Fingerprint.
	* `key_pair_name` - The Key Name.
