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

For information about Hybrid Backup Recovery (HBR) Hana Backup Client and how to use it, see [What is Hana Backup Client](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-hbr-2017-09-08-createclients).

-> **NOTE:** Available in v1.198.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_hana_backup_client&exampleId=319f318f-bbc7-3db8-f74a-5826f2f4afb15d7bebca&activeTab=example&spm=docs.r.hbr_hana_backup_client.0.319f318fbb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}



resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types.0.id
  availability_zone    = data.alicloud_zones.example.zones.0.id
  security_groups      = [alicloud_security_group.example.id]
  instance_name        = "terraform-example"
  internet_charge_type = "PayByBandwidth"
  vswitch_id           = alicloud_vswitch.example.id
}


data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

resource "alicloud_hbr_vault" "example" {
  vault_name = "terraform-example"
}

resource "alicloud_hbr_hana_instance" "example" {
  alert_setting        = "INHERITED"
  hana_name            = "terraform-example"
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

resource "alicloud_hbr_hana_backup_client" "default" {
  vault_id      = alicloud_hbr_vault.example.id
  client_info   = "[ { \"instanceId\": \"${alicloud_instance.example.id}\", \"clusterId\": \"${alicloud_hbr_hana_instance.example.hana_instance_id}\", \"sourceTypes\": [ \"HANA\" ]  }]"
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
