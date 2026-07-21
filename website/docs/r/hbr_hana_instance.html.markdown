---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_instance"
sidebar_current: "docs-alicloud-resource-hbr-hana-instance"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Hana Instance resource.
---

# alicloud_hbr_hana_instance

Provides a Hybrid Backup Recovery (HBR) Hana Instance resource.

For information about Hybrid Backup Recovery (HBR) Hana Instance and how to use it, see [What is Hana Instance](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-hbr-2017-09-08-createhanainstance).

-> **NOTE:** Available since v1.178.0.

-> **NOTE:** The `sid` attribute is required when destroying resources.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_hana_instance&exampleId=530e500c-6119-001b-373a-6eb81c410aca8e4bb25b&activeTab=example&spm=docs.r.hbr_hana_instance.0.530e500c61&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_hbr_vault" "example" {
  vault_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_hbr_hana_instance" "example" {
  alert_setting        = "INHERITED"
  hana_name            = "terraform-example-${random_integer.default.result}"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_hbr_hana_instance&spm=docs.r.hbr_hana_instance.example&intl_lang=EN_US)

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Hana Instance.
* `delete` - (Defaults to 1 mins) Used when delete the Hana Instance.
* `update` - (Defaults to 1 mins) Used when update the Hana Instance.

## Import

Hybrid Backup Recovery (HBR) Hana Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_hana_instance.example <vault_id>:<hana_instance_id>
```