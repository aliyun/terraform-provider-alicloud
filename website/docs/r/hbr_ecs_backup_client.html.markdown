---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_client"
sidebar_current: "docs-alicloud-resource-hbr-ecs-backup-client"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Ecs Backup Client resource.
---

# alicloud\_hbr\_ecs\_backup\_client

Provides a Hybrid Backup Recovery (HBR) Ecs Backup Client resource.

For information about Hybrid Backup Recovery (HBR) Ecs Backup Client and how to use it, see [What is Ecs Backup Client](https://www.alibabacloud.com/help/doc-detail/186570.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_ecs_backup_client&exampleId=c40439fa-ab48-48c5-df88-bfcf39dc0a95bc01aab6&activeTab=example&spm=docs.r.hbr_ecs_backup_client.0.c40439faab&intl_lang=EN_US" target="_blank">
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
  name_regex = "^ubuntu_18.*64"
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

resource "alicloud_hbr_ecs_backup_client" "example" {
  instance_id        = alicloud_instance.example.id
  use_https          = false
  data_network_type  = "VPC"
  max_cpu_core       = 2
  max_worker         = 4
  data_proxy_setting = "USE_CONTROL_PROXY"
  proxy_host         = "192.168.11.101"
  proxy_port         = 80
  proxy_user         = "user"
  proxy_password     = "password"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of ECS instance.
* `use_https` - (Optional) Indicates whether to use the HTTPS protocol. Valid values: `true`, `false`.
* `status` - (Optional, Computed) Status of client. Valid values: `ACTIVATED`, `STOPPED`. You can start or stop the client by specifying the status.
* `data_network_type` - (Optional) The data plane access point type. Valid values: `CLASSIC`, `PUBLIC`, `VPC`. **NOTE:** The value of `CLASSIC` has been deprecated in v1.161.0+.
* `data_proxy_setting` - (Optional, Computed) The data plane proxy settings. Valid values: `CUSTOM`, `DISABLE`, `USE_CONTROL_PROXY`.
* `max_cpu_core` - (Optional) The number of CPU cores used by a single backup task, 0 means no restrictions.
* `max_worker` - (Optional) The number of concurrent jobs for a single backup task, 0 means no restrictions.
* `proxy_host` - (Optional) Custom data plane proxy server host address.
* `proxy_port` - (Optional) Custom data plane proxy server host port.
* `proxy_user` - (Optional) The username of custom data plane proxy server.
* `proxy_password` - (Optional) The password of custom data plane proxy server.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ecs Backup Client.

## Notice

-> **Note:** Please read the following precautions carefully before deleting a client:
1. You cannot delete active clients that have received heartbeat packets within one hour.
2. You can make the client inactive by change the status of client to `STOPPED`.
3. The resources bound to the client will be deleted in cascade, including:
    - Backup plan
    - Backup task (Running in the background)
    - Snapshot


## Import

Hybrid Backup Recovery (HBR) Ecs Backup Client can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_ecs_backup_client.example <id>
```
