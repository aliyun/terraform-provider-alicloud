---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_lakebase_tenant_token"
sidebar_current: "docs-alicloud-datasource-polardb-lakebase-tenant-token"
description: |-
  Retrieves a Lakebase tenant token for a PolarFS instance.
---

# alicloud_polardb_lakebase_tenant_token

Retrieves a Lakebase tenant token for a PolarFS instance.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_lakebase_tenant_token" "example" {
  polar_fs_instance_id = alicloud_polardb_polar_fs.example.id
  db_cluster_id        = alicloud_polardb_cluster.example.id
  subdir               = "/home/project"
  tenant               = "admin"
}
```

## Argument Reference

* `polar_fs_instance_id` - (Required, Available since v1.294.0) The PolarFS instance ID.
* `db_cluster_id` - (Optional, Available since v1.294.0) The associated PolarDB cluster ID.
* `subdir` - (Required, Available since v1.294.0) The directory authorized by the token.
* `tenant` - (Optional, Available since v1.294.0) The Lakebase tenant.

## Attributes Reference

* `token` - (Sensitive, Available since v1.294.0) The temporary tenant token.
* `status` - (Available since v1.294.0) The token status.
