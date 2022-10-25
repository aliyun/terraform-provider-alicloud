---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_instance"
sidebar_current: "docs-alicloud-resource-hbr-hana-instance"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Hana Instance resource.
---

# alicloud\_hbr\_hana\_instance

Provides a Hybrid Backup Recovery (HBR) Hana Instance resource.

For information about Hybrid Backup Recovery (HBR) Hana Instance and how to use it, see [What is Hana Instance](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-doc-hbr-2017-09-08-api-doc-createhanainstance).

-> **NOTE:** Available in v1.178.0+.

-> **NOTE:** The `sid` attribute is required when destroying resources.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbr_vault" "example" {
  vault_name = var.name
}
data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

resource "alicloud_hbr_hana_instance" "example" {
  alert_setting        = "INHERITED"
  hana_name            = var.name
  host                 = "1.1.1.1"
  instance_number      = 1
  password             = "YouPassword123"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.example.groups.0.id
  sid                  = "HXE"
  use_ssl              = false
  user_name            = "admin"
  validate_certificate = false
  vault_id             = alicloud_hbr_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `vault_id` - (Required, ForceNew) The ID of the backup vault.
* `alert_setting` - (Optional) The alert settings. Valid value: `INHERITED`, which indicates that the backup client sends alert notifications in the same way as the backup vault.
* `ecs_instance_ids` - (Optional) The IDs of ECS instances that host the SAP HANA instance to be registered. HBR installs backup clients on the specified ECS instances.
* `hana_name` - (Optional) The name of the SAP HANA instance.
* `host` - (Optional) The private or internal IP address of the host where the primary node of the SAP HANA instance resides.
* `instance_number` - (Optional) The instance number of the SAP HANA system.
* `password` - (Optional, ForceNew) The password that is used to connect with the SAP HANA database.
* `resource_group_id` - (Optional) The ID of the resource group.
* `sid` - (Optional, ForceNew) The security identifier (SID) of the SAP HANA database.
* `use_ssl` - (Optional) Specifies whether to connect with the SAP HANA database over Secure Sockets Layer (SSL).
* `user_name` - (Optional) The username of the SYSTEMDB database.
* `validate_certificate` - (Optional) Specifies whether to verify the SSL certificate of the SAP HANA database.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Hana Instance. The value formats as `<vault_id>:<hana_instance_id>`.
* `hana_instance_id` - The id of the Hana Instance.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Hana Instance.
* `delete` - (Defaults to 1 mins) Used when delete the Hana Instance.
* `update` - (Defaults to 1 mins) Used when update the Hana Instance.

## Import

Hybrid Backup Recovery (HBR) Hana Instance can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_hana_instance.example <vault_id>:<hana_instance_id>
```