---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_port"
sidebar_current: "docs-alicloud-resource-ddoscoo-port"
description: |-
  Provides a Alicloud Anti-DDoS Pro Port resource.
---

# alicloud\_ddoscoo\_port

Provides a Anti-DDoS Pro Port resource.

For information about Anti-DDoS Pro Port and how to use it, see [What is Port](https://www.alibabacloud.com/help/en/doc-detail/157482.htm).

-> **NOTE:** Available in v1.123.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ddoscoo_instance" "example" {
  name              = "yourDdoscooInstanceName"
  bandwidth         = "30"
  base_bandwidth    = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
}

resource "alicloud_ddoscoo_port" "example" {
  instance_id       = alicloud_ddoscoo_instance.example.id
  frontend_port     = "7001"
  frontend_protocol = "tcp"
  real_servers      = ["1.1.1.1", "2.2.2.2"]
}

```

## Argument Reference

The following arguments are supported:

* `backend_port` - (Optional, ForceNew) The port of the origin server. Valid values: [1~65535].
* `frontend_port` - (Required, ForceNew) The forwarding port. Valid values: [1~65535].
* `instance_id` - (Required, ForceNew) The ID of Ddoscoo instance.
* `frontend_protocol` - (Required, ForceNew) The forwarding protocol. Valid values `tcp` and `udp`.
* `real_servers` - (Required) List of source IP addresses.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Port. The value is formatted `<instance_id>:<frontend_port>:<frontend_protocol>`.

## Import

Anti-DDoS Pro Port can be imported using the id, e.g.

```
$ terraform import alicloud_ddoscoo_port.example <instance_id>:<frontend_port>:<frontend_protocol>
```
