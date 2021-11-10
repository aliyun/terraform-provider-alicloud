---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_listener"
sidebar_current: "docs-alicloud-resource-ga-listener"
description: |-
  Provides a Alicloud Global Accelerator (GA) Listener resource.
---

# alicloud\_ga\_listener

Provides a Global Accelerator (GA) Listener resource.

For information about Global Accelerator (GA) Listener and how to use it, see [What is Listener](https://help.aliyun.com/document_detail/153253.html).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "de" {
  bandwidth      = "100"
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "de" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.de.id
}
resource "alicloud_ga_listener" "example" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.de]
  accelerator_id = alicloud_ga_accelerator.example.id
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

-> **NOTE:** This parameter needs to be configured only for monitoring of the HTTPS protocol.
             
* `client_affinity` - (Optional) The clientAffinity of the listener. Default value is `NONE`. Valid values:
    `NONE`: client affinity is not maintained, that is, connection requests from the same client cannot always be directed to the same terminal node.
    `SOURCE_IP`: maintain client affinity. When a client accesses a stateful application, all requests from the same client can be directed to the same terminal node, regardless of the source port and protocol.
* `description` - (Optional) The description of the listener.
* `name` - (Optional) The name of the listener. The length of the name is 2-128 characters. It starts with uppercase and lowercase letters or Chinese characters. It can contain numbers and underscores and dashes.
* `port_ranges` - (Required) The portRanges of the listener.

-> **NOTE:** For HTTP or HTTPS protocol monitoring, only one monitoring port can be configured, that is, the start monitoring port and end monitoring port should be the same. 

* `protocol` - (Optional) Type of network transport protocol monitored. Default value is `TCP`. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`.

-> **NOTE:** At present, the white list of HTTP and HTTPS monitoring protocols is open. If you need to use it, please submit a work order.
             
* `proxy_protocol` - (Optional) The proxy protocol of the listener. Default value is `false`. Valid value:
    `true`: Turn on the keep client source IP function. After it is turned on, the back-end service is supported to view the original IP address of the client. 
    `false`: keep client source IP function is not turned on.

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

* `create` - (Defaults to 10 mins) Used when create the Listener.
* `delete` - (Defaults to 6 mins) Used when update the Listener.
* `update` - (Defaults to 3 mins) Used when terminating the Listener.

## Import

Ga Listener can be imported using the id, e.g.

```
$ terraform import alicloud_ga_listener.example <id>
```
