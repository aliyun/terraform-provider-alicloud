---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_firewall_rule"
sidebar_current: "docs-alicloud-resource-simple-application-server-firewall-rule"
description: |-
  Provides a Alicloud Simple Application Server Firewall Rule resource.
---

# alicloud\_simple\_application\_server\_firewall\_rule

Provides a Simple Application Server Firewall Rule resource.

For information about Simple Application Server Firewall Rule and how to use it, see [What is Firewall Rule](https://www.alibabacloud.com/help/doc-detail/190449.htm).

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_instances" "default" {}

data "alicloud_simple_application_server_images" "default" {}

data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server_instance" "default" {
  count          = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = "tf-testaccswas-firewallrule"
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

resource "alicloud_simple_application_server_firewall_rule" "default" {
  instance_id   = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? data.alicloud_simple_application_server_instances.default.ids.0 : alicloud_simple_application_server_instance.default.0.id
  rule_protocol = "Tcp"
  port          = "9999"
  remark        = "example_value"
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Alibaba Cloud simple application server instance ID.
* `port` - (Required, ForceNew) The port range. Valid values of port numbers: `1` to `65535`. Specify a port range in the format of `<start port number>/<end port number>`. Example: `1024/1055`, which indicates the port range of `1024` through `1055`.
* `remark` - (Optional, ForceNew) The remarks of the firewall rule.
* `rule_protocol` - (Required, ForceNew) The transport layer protocol. Valid values: `Tcp`, `Udp`, `TcpAndUdp`.

## Attributes Reference

The following attributes are exported:

* `firewall_rule_id` - The ID of the firewall rule.
* `id` - The resource ID of Firewall Rule. The value formats as `<instance_id>:<firewall_rule_id>`.

## Import

Simple Application Server Firewall Rule can be imported using the id, e.g.

```
$ terraform import alicloud_simple_application_server_firewall_rule.example <instance_id>:<firewall_rule_id>
```