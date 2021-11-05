---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_slb_attachment"
sidebar_current: "docs-alicloud-resource-edas-slb-attachment"
description: |-
  Binds SLB to an EDAS application.
---

# alicloud\_edas\_slb\_attachment

Binds SLB to an EDAS application.

-> **NOTE:** Available in 1.82.0+

## Example Usage

Basic Usage

```
resource "alicloud_edas_slb_attachment" "default" {
  app_id           = var.app_id
  slb_id           = var.slb_id
  slb_ip           = var.slb_ip
  type             = var.type
  listener_port    = var.listener_port
  vserver_group_id = var.vserver_group_id
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application to which you want to bind an SLB instance.
* `slb_id` - (Required, ForceNew) The ID of the SLB instance that is going to be bound.
* `slb_ip` - (Required, ForceNew) The IP address that is allocated to the bound SLB instance.
* `type` - (Required, ForceNew) The type of the bound SLB instance.
* `listener_port` - (Optional, ForceNew) The listening port for the bound SLB instance.
* `vserver_group_id` - (Optional, ForceNew) The ID of the virtual server (VServer) group associated with the intranet SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<slb_id>`.
* `slb_status` - Running Status of SLB instance. Inactive：The instance is stopped, and listener will not monitor and forward traffic. Active：The instance is running. After the instance is created, the default state is active. Locked：The instance is locked, the instance has been owed or locked by Alibaba Cloud. Expired: The instance has expired.
* `vswitch_id` - VPC related vswitch ID.


