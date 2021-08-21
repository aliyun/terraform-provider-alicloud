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
  instance_id = data.alicloud_instances.default.instances.0.id
  use_https =          false
  data_network_type =  "PUBLIC"
  max_cpu_core =       2
  max_worker   =       4
  data_proxy_setting = "USE_CONTROL_PROXY"
  proxy_host =         "192.168.11.101"
  proxy_port =         80
  proxy_user =         "user"
  proxy_password =     "password"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ECS Instance Id.
* `use_https` - (Optional) Indicates Whether to Use the Https Transport Data Plane Data. Valid values: `true`, `false`.
* `status` - (Optional, Computed) Status of client. Valid values: `ACTIVATED`, `STOPPED`.
* `data_network_type` - (Optional) The Data Plane Data Access Point Type. Valid values: `CLASSIC`, `PUBLIC`, `VPC`.
* `data_proxy_setting` - (Optional, Computed) The Data Plane Proxy Settings. Valid values: `CUSTOM`, `DISABLE`, `USE_CONTROL_PROXY`.
* `max_cpu_core` - (Optional) A Single Backup Task Uses for Example, Instances Can Be Grouped According to CPU Core Count, 0 Means No Restrictions.
* `max_worker` - (Optional) A Single Backup Task Parallel Work, the Number of 0 Means No Restrictions.
* `proxy_host` - (Optional) Custom Data Plane Proxy Server Host Address.
* `proxy_port` - (Optional) Custom Data Plane Proxy Server Host Port.
* `proxy_user` - (Optional) Custom Data Plane Proxy Server Username.
* `proxy_password` - (Optional) Custom Data Plane Proxy Password.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ecs Backup Client.

## Notice

> Please read the following precautions carefully before deleting a client:

1. You cannot delete active clients that have received heartbeat packets within one hour.
2. You can make the client inactive by change the status of client to `ACTIVATED`.
3. The resources bound to the client will be deleted in cascade, including:
- Backup plan
- Backup task (Running in the background)
- Snapshot

## Import

Hybrid Backup Recovery (HBR) Ecs Backup Client can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_ecs_backup_client.example <id>
```