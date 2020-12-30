---
subcategory: "Ga"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_listener"
sidebar_current: "docs-alicloud-resource-ga-listener"
description: |-
  Provides a Alicloud Ga Listener resource.
---

# alicloud\_ga\_listener

Provides a Ga Listener resource.

For information about Ga Listener and how to use it, see [What is Listener](https://help.aliyun.com/document_detail/153253.html).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_listener" "example" {
  accelerator_id = "alicloud_ga_accelerator.example.id"
  port_ranges {
    from_port = 60
    to_port   = 70
  }
}

```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required) The accelerator id.
* `certificates` - (Optional) The certificates of the listener.
* `client_affinity` - (Optional) The clientAffinity of the listener. Default value is `NONE`. Valid values:
    `NONE`: client affinity is not maintained, that is, connection requests from the same client cannot always be directed to the same terminal node.
    `SOURCE_IP`: maintain client affinity. When a client accesses a stateful application, all requests from the same client can be directed to the same terminal node, regardless of the source port and protocol.
* `description` - (Optional) The description of the listener.
* `name` - (Optional) The name of the listener. The length of the name is 2-128 characters. It starts with uppercase and lowercase letters or Chinese characters. It can contain numbers and underscores and dashes.
* `port_ranges` - (Required) The portRanges of the listener.
* `protocol` - (Optional) Type of network transport protocol monitored. Default value is `TCP`. Valid values: `TCP`, `UDP`.
* `proxy_protocol` - (Optional) The proxy protocol of the listener.

#### Block port_ranges

The port_ranges supports the following: 

* `from_port` - (Required) The initial listening port used to receive requests and forward them to terminal nodes.
* `to_port` - (Required) The end listening port used to receive requests and forward them to terminal nodes.

#### Block certificates

The certificates supports the following: 

* `id` - (Optional) The id of the certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener.
* `status` - The status of the listener.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Listener.
* `delete` - (Defaults to 6 mins) Used when update the Listener.
* `update` - (Defaults to 6 mins) Used when terminating the Listener.

## Import

Ga Listener can be imported using the id, e.g.

```
$ terraform import alicloud_ga_listener.example <id>
```