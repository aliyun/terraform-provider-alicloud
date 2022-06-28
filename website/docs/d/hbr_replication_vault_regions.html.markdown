---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_replication_vault_regions"
sidebar_current: "docs-alicloud-datasource-hbr-replication-vault-regions"
description: |-
  Provides a list of HBR Replication Vault Regions to the user.
---

# alicloud\_hbr\_replication\_vault\_regions

This data source provides the HBR Replication Vault Regions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_replication_vault_regions" "default" {}

output "hbr_replication_vault_region_region_id_1" {
  value = data.alicloud_hbr_replication_vault_regions.default.regions.0.replication_region_id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of Ros Regions. Each element contains the following attributes:
	* `replication_region_id` - The ID of the replication region.
