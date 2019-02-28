---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_attachment"
sidebar_current: "docs-alicloud-resource-slb-attachment"
description: |-
  Provides an Application Load Banlancer Attachment resource.
---

# alicloud\_slb\_attachment

Add a group of backend servers (ECS instance) to the Server Load Balancer or remove them from it.

## Example Usage

```
# Create a new load balancer attachment for classic
resource "alicloud_slb" "default" {
  # Other parameters...
}

resource "alicloud_instance" "default" {
  # Other parameters...
}

resource "alicloud_slb_attachment" "default" {
  load_balancer_id    = "${alicloud_slb.default.id}"
  instance_ids = ["${alicloud_instance.default.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the load balancer.
* `instance_ids` - (Required) A list of instance ids to added backend server in the SLB.
* `weight` - (Optional) Weight of the instances. Valid value range: [0-100]. Default to 100.
* `slb_id` - (Deprecated) It has been deprecated from provider version 1.6.0. New field 'load_balancer_id' replaces it.
* `instances` - (Deprecated) It has been deprecated from provider version 1.6.0. New field 'instance_ids' replaces it.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource.
* `load_balancer_id` - ID of the load balancer.
* `instance_ids` - A list of instance ids that have been added in the SLB.
* `weight` - (Optional) Weight of the instances.
* `backend_servers` - The backend servers of the load balancer.

## Import

Load balancer attachment can be imported using the id or load balancer id, e.g.

```
$ terraform import alicloud_slb_attachment.example lb-abc123456
```
