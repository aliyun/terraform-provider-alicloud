---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_backup_client"
sidebar_current: "docs-alicloud-resource-hbr-hana-backup-client"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Hana Backup Client resource.
---

# alicloud\_hbr\_hana\_backup\_client

Provides a Hybrid Backup Recovery (HBR) Hana Backup Client resource.

For information about Hybrid Backup Recovery (HBR) Hana Backup Client and how to use it, see [What is Hana Backup Client](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-doc-hbr-2017-09-08-api-doc-createclients).

-> **NOTE:** Available in v1.198.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbr_hana_backup_client" "default" {
  vault_id      = data.alicloud_hbr_vaults.default.vaults.0.id
  client_info   = "[ { \"instanceId\": \"i-bp116lr******te9q2\", \"clusterId\": \"cl-000csy09q******9rfz9\", \"sourceTypes\": [ \"HANA\" ]  }]"
  alert_setting = "INHERITED"
  use_https     = true
}
```

## Argument Reference

The following arguments are supported:

* `vault_id` - (Required, ForceNew) The ID of the backup vault.
* `client_info` - (Optional) The installation information of the HBR clients.
* `alert_setting` - (Optional, ForceNew, Computed) The alert settings. Valid value: `INHERITED`.
* `use_https` - (Optional, ForceNew) Specifies whether to transmit data over HTTPS. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hana Backup Client. It formats as `<vault_id>:<client_id>`.
* `client_id` - The ID of the backup client.
* `instance_id` - The ID of the instance.
* `cluster_id` - The ID of the SAP HANA instance.
* `status` - The status of the Hana Backup Client.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Hana Backup Client.
* `delete` - (Defaults to 5 mins) Used when delete the Hana Backup Client.

## Import

Hybrid Backup Recovery (HBR) Hana Backup Client can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_hana_backup_client.example <vault_id>:<client_id>
```
