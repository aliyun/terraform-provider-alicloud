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

```terraform
data "alicloud_instances" "default" {
  name_regex = "ecs_instance_name"
  status     = "Running"
}

resource "alicloud_hbr_ecs_backup_client" "example" {
  instance_id        = data.alicloud_instances.default.instances.0.id
  use_https          = false
  data_network_type  = "PUBLIC"
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

```
$ terraform import alicloud_hbr_ecs_backup_client.example <id>
```
