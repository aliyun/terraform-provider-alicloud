---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listener_acl_attachment"
sidebar_current: "docs-alicloud-resource-alb-listener-acl-attachment"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Listener Acl Attachment resource.
---

# alicloud\_alb\_listener\_acl\_attachment"

Provides a Application Load Balancer (ALB) Listener Acl Attachment resource.

For information about Application Load Balancer (ALB) Listener Acl Attachment and how to use it, see [What is Listener Acl Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/associateaclswithlistener).

-> **NOTE:** Available in v1.163.0+.

-> **NOTE:** You can associate at most three ACLs with a listener.

-> **NOTE:** You can only configure either a whitelist or a blacklist for listener, not at the same time.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = "example_value"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  acl_entries {
    description = "description"
    entry       = "10.0.0.0/24"
  }
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}

data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = "example_value"
  load_balancer_edition  = "Standard"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default_1.ids[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default_2.ids[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.vpcs.0.id
  server_group_name = "example_value"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
  tags = {
    Created = "TF"
  }
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default.id
  listener_protocol    = "HTTP"
  listener_port        = 80
  listener_description = "example_value"
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
  }
}

resource "alicloud_alb_listener_acl_attachment" "default" {
  acl_id      = alicloud_alb_acl.default.id
  listener_id = alicloud_alb_listener.default.id
  acl_type    = "White"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `listener_id` - (Required, ForceNew) The ID of the ALB listener.
* `acl_type` - (Required, ForceNew) The type of the ACL. Valid values: 
  - White: a whitelist. Only requests from the IP addresses or CIDR blocks in the ACL are forwarded. The whitelist applies to scenarios in which you want to allow only specific IP addresses to access an application. Risks may arise if you specify an ACL as a whitelist. After a whitelist is configured, only IP addresses in the whitelist can access the Application Load Balancer (ALB) listener. If you enable a whitelist but the whitelist does not contain an IP address, the listener forwards all requests. 
  - Black: a blacklist. All requests from the IP addresses or CIDR blocks in the ACL are blocked. The blacklist applies to scenarios in which you want to block access from specific IP addresses to an application. If you enable a blacklist but the blacklist does not contain an IP address, the listener forwards all requests.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Listener Acl Attachment. The value formats as `<listener_id>:<acl_id>`.
* `status` - The status of the Listener Acl Attachment.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Listener Acl Attachment.
* `delete` - (Defaults to 2 mins) Used when delete the Listener Acl Attachment.

## Import

Application Load Balancer (ALB) Listener Acl Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_alb_listener_acl_attachment.example <listener_id>:<acl_id>
```