---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_listener_acl_attachment"
sidebar_current: "docs-alicloud-resource-slb-listener-acl-attachment"
description: |-
  Provides a Alicloud Classic Load Balancer (SLB) Listener ACL Attachment resource.
---

# alicloud_slb_listener_acl_attachment

Provides a Classic Load Balancer (SLB) Listener ACL Attachment resource. It binds an existing access control policy group (ACL) to an existing SLB listener.

For information about Classic Load Balancer (SLB) and how to use it, see [What is Classic Load Balancer (SLB)](https://www.alibabacloud.com/help/en/slb/classic-load-balancer/product-overview/what-is-slb).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** This resource manages only the access control settings of the listener. It can be used with listeners that are not managed by Terraform, for example the listeners created by an ACK cluster. Use the `id` element of `alicloud_slb_listeners` (or the `id` of `alicloud_slb_listener`) as `listener_id`.

-> **NOTE:** An SLB listener expresses a single set of ACLs. This resource exclusively manages the ACL binding of its listener and overwrites the ACLs already bound to that listener. Do not manage the same listener through both this resource and the `acl_status`/`acl_type`/`acl_id`/`acl_ids` arguments of `alicloud_slb_listener`. When the listener itself is still created by `alicloud_slb_listener`, add `ignore_changes = [acl_status, acl_type, acl_id, acl_ids]` to that listener so that only this resource drives the ACL settings:

```terraform
resource "alicloud_slb_listener" "default" {
  # ... other listener arguments ...

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}
```

-> **NOTE:** Enabling access control takes effect immediately and may affect traffic. When a white list is enabled with no entries, all requests are forwarded; when a black list is enabled with no entries, all requests are denied. Make sure the policy group is configured before binding it.

-> **NOTE:** Deleting this resource only disables access control (`AclStatus=off`) and never deletes or rebuilds the listener. SLB has no dedicated de-association API, so the server may keep the previous `AclId` after the binding is removed.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = var.name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  load_balancer_spec   = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
}

resource "alicloud_slb_acl" "default" {
  name       = "${var.name}-acl"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.10.0/24"
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.default.id
  acl_type    = "white"
}
```

### Attach an ACL to a listener that is not managed by Terraform

An ACK-created apiServer SLB listener can be located through the `alicloud_slb_listeners` data source and bound directly:

```terraform
data "alicloud_slb_listeners" "apiserver" {
  load_balancer_id = "lb-xxxxxxxxxxxx"
}

resource "alicloud_slb_listener_acl_attachment" "apiserver" {
  listener_id = data.alicloud_slb_listeners.apiserver.slb_listeners[0].id
  acl_id      = alicloud_slb_acl.apiserver.id
  acl_type    = "white"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The ID of the SLB listener, in the format of `<load_balancer_id>:<protocol>:<listener_port>`. The protocol is one of `http`, `https`, `tcp` and `udp`. It is the same as the `id` of `alicloud_slb_listener`.
* `acl_id` - (Required) The ID of the access control policy group bound to the listener.
* `acl_type` - (Required) The access control type. Valid values: `white` and `black`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the listener ACL attachment. It formats as `<load_balancer_id>:<protocol>:<listener_port>`, which is the ID of the listener whose ACL binding is managed.
* `acl_status` - The access control status of the listener.

## Import

Classic Load Balancer (SLB) Listener ACL Attachment can be imported using the listener ID, optionally followed by the ACL ID, e.g.

```shell
$ terraform import alicloud_slb_listener_acl_attachment.example <load_balancer_id>:<protocol>:<listener_port>
```

or

```shell
$ terraform import alicloud_slb_listener_acl_attachment.example <load_balancer_id>:<protocol>:<listener_port>:<acl_id>
```
