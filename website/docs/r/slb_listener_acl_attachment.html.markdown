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

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** This resource manages only the access control settings of the listener. It can be used with listeners that are not managed by Terraform, for example the listeners created by an ACK cluster. Use the `id` element of `alicloud_slb_listeners` (or the `id` of `alicloud_slb_listener`) as `listener_id`.

-> **NOTE:** An SLB listener expresses a single set of ACLs. This resource exclusively manages the ACL binding of its listener and overwrites the ACLs already bound to that listener. Do not manage the same listener through both this resource and the `acl_status`/`acl_type`/`acl_id`/`acl_ids` arguments of `alicloud_slb_listener`. When the listener itself is still created by `alicloud_slb_listener`, add `ignore_changes = [acl_status, acl_type, acl_ids]` to that listener so that only this resource drives the ACL settings.

-> **NOTE:** Enabling access control takes effect immediately and may affect traffic. A white list or a black list with no entries forwards all requests. Make sure the policy group is configured before binding it.

-> **NOTE:** Deleting this resource only disables access control (`AclStatus=off`) and never deletes or rebuilds the listener. SLB has no dedicated de-association API, so the server may keep the previous `AclId` after the binding is removed.

-> **NOTE:** When the listener is bound to several policy groups, use the four segment import form `<load_balancer_id>:<protocol>:<listener_port>:<acl_id>` to select the managed group; importing with only the listener ID then fails with a diagnostic instead of guessing one. Deleting this resource still disables access control for the whole listener, not only for the imported policy group.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-exampleSlbListenerAcl"
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

  lifecycle {
    # The attachment resource owns the access control settings of this listener.
    # `acl_id` is deprecated, so only the current arguments are ignored.
    ignore_changes = [acl_status, acl_type, acl_ids]

    precondition {
      condition     = alicloud_slb_acl.default.id != ""
      error_message = "The policy group must exist before the listener so that destroying the listener releases the ACL association first."
    }
  }
}

resource "alicloud_slb_acl" "default" {
  name       = "${var.name}-acl"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "default" {
  acl_id = alicloud_slb_acl.default.id
  entries {
    entry   = "10.10.10.0/24"
    comment = "example entry"
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.default.id
  acl_type    = "white"
}
```

### Locate an existing listener by protocol and frontend port

An ACK-created apiServer SLB listener can be located through the `alicloud_slb_listeners` data source and bound directly, without importing the listener into Terraform. A load balancer can expose several listeners, so the query must be filtered by the target protocol and frontend port and must return exactly one listener.

The example below creates its own load balancer, two listeners and a policy group so it can run on its own; the second listener proves that the filters select the apiServer listener and not just the first entry of the data source. In an ACK cluster, drop the `alicloud_slb_load_balancer`, `alicloud_slb_listener` and `alicloud_slb_acl` resources, set `load_balancer_id` on the data source to the ACK load balancer ID, and keep the protocol and frontend port filters, the `one(...)` lookup and the precondition.

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-exampleSlbListenerAclAck"
}

resource "alicloud_slb_load_balancer" "apiserver" {
  load_balancer_name   = var.name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  load_balancer_spec   = "slb.s1.small"
}

resource "alicloud_slb_listener" "apiserver" {
  load_balancer_id = alicloud_slb_load_balancer.apiserver.id
  backend_port     = 6443
  frontend_port    = 6443
  protocol         = "tcp"
  bandwidth        = 10
  health_check     = "on"

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_ids]
  }
}

resource "alicloud_slb_listener" "other" {
  load_balancer_id = alicloud_slb_load_balancer.apiserver.id
  backend_port     = 80
  frontend_port    = 80
  protocol         = "http"
  bandwidth        = 10
}

resource "alicloud_slb_acl" "apiserver" {
  name       = "${var.name}-acl"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "apiserver" {
  acl_id = alicloud_slb_acl.apiserver.id
  entries {
    entry = "10.10.10.0/24"
  }
}

data "alicloud_slb_listeners" "apiserver" {
  load_balancer_id = alicloud_slb_load_balancer.apiserver.id
  # Referring to the listener attributes makes the query wait until the
  # listeners exist, so the result does not depend on creation order.
  protocol      = alicloud_slb_listener.apiserver.protocol
  frontend_port = alicloud_slb_listener.apiserver.frontend_port
}

resource "alicloud_slb_listener_acl_attachment" "apiserver" {
  # The filters above identify the listener by protocol and frontend port; a
  # load balancer exposes at most one listener per protocol and port. `one`
  # keeps the example from silently using the first of several matches.
  listener_id = one([for listener in data.alicloud_slb_listeners.apiserver.slb_listeners : listener.id])
  acl_id      = alicloud_slb_acl.apiserver.id
  acl_type    = "white"

  lifecycle {
    precondition {
      condition     = length(data.alicloud_slb_listeners.apiserver.slb_listeners) == 1
      error_message = "The protocol and frontend_port filters must match exactly one listener; check the input variables before binding the policy group."
    }
  }
}
```

A load balancer exposes at most one listener per protocol and frontend port, so the filtered query must return exactly one element. The example enforces that requirement with a precondition and `one(...)`, so a zero or ambiguous match fails during plan instead of binding the policy group to an unintended listener.

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew, String) The ID of the SLB listener, in the format of `<load_balancer_id>:<protocol>:<listener_port>`. The protocol is one of `http`, `https`, `tcp` and `udp`, and the port ranges from `1` to `65535`. It is the same as the `id` of `alicloud_slb_listener`.
* `acl_id` - (Required, String) The ID of the access control policy group bound to the listener. Exactly one policy group is managed by this resource; a comma separated value is rejected during validation.
* `acl_type` - (Required, String) The access control type. Valid values: `white` and `black`.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The resource ID in terraform of the listener ACL attachment. It formats as `<load_balancer_id>:<protocol>:<listener_port>`, which is the ID of the listener whose ACL binding is managed.
* `acl_status` - (String) The access control status of the listener, as reported by the last successful read.
* `acl_pending` - (Bool) Whether the last access control write has not been confirmed by the listener yet. It is `true` after a write whose result could not be verified: the planned `acl_id` and `acl_type` are kept in state, a refresh does not remove the resource while the listener still reports the previous group or a disabled ACL, and a destroy only releases the binding that is attributed to this resource. It becomes `false` once the listener reports the planned binding.

  -> **NOTE:** When `acl_pending` is `true`, the configuration and the state agree, so the next `terraform apply` has nothing to change. Run `terraform apply -replace=alicloud_slb_listener_acl_attachment.<name>` (or `terraform taint`) to write the binding again. Destroying the resource in this state only disables access control when the listener actually reports the planned policy group, so a binding that was never applied does not release a group managed by something else.

## Import

Classic Load Balancer (SLB) Listener ACL Attachment can be imported using the listener ID, optionally followed by the ACL ID, e.g.

```shell
$ terraform import alicloud_slb_listener_acl_attachment.example <load_balancer_id>:<protocol>:<listener_port>
```

or

```shell
$ terraform import alicloud_slb_listener_acl_attachment.example <load_balancer_id>:<protocol>:<listener_port>:<acl_id>
```
