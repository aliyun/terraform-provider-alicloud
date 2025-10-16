---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_traffic_qos_rule"
description: |-
  Provides a Alicloud Express Connect Traffic Qos Rule resource.
---

# alicloud_express_connect_traffic_qos_rule

Provides a Express Connect Traffic Qos Rule resource.

Express Connect Traffic QoS Rule.

For information about Express Connect Traffic Qos Rule and how to use it, see [What is Traffic Qos Rule](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateExpressConnectTrafficQosRule).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos" "createQos" {
  qos_name        = var.name
  qos_description = "terraform-example"
}

resource "alicloud_express_connect_traffic_qos_association" "associateQos" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.1
  qos_id        = alicloud_express_connect_traffic_qos.createQos.id
  instance_type = "PHYSICALCONNECTION"
}

resource "alicloud_express_connect_traffic_qos_queue" "createQosQueue" {
  qos_id            = alicloud_express_connect_traffic_qos.createQos.id
  bandwidth_percent = "60"
  queue_description = "terraform-example"
  queue_name        = var.name
  queue_type        = "Medium"
}


resource "alicloud_express_connect_traffic_qos_rule" "default" {
  rule_description = "terraform-example"
  priority         = "1"
  protocol         = "ALL"
  src_port_range   = "-1/-1"
  dst_ipv6_cidr    = "2001:db8:1234:5678::/64"
  src_ipv6_cidr    = "2001:db8:1234:5678::/64"
  dst_port_range   = "-1/-1"
  remarking_dscp   = "-1"
  queue_id         = alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id
  qos_id           = alicloud_express_connect_traffic_qos.createQos.id
  match_dscp       = "-1"
  rule_name        = var.name
}
```

## Argument Reference

The following arguments are supported:
* `dst_cidr` - (Optional) The traffic of the QoS rule matches the Destination IPv4 network segment.

-> **NOTE:**  If this parameter is not supported, enter `SrcIPv6Cidr` or **DstIPv6Cidr * *.

* `dst_ipv6_cidr` - (Optional) The QoS rule traffic matches the Destination IPv6 network segment.

-> **NOTE:**  If this parameter is not supported, enter `SrcCidr` or **DstCidr * *.

* `dst_port_range` - (Optional, Computed) QoS rule traffic matches the destination port number range. Value range: `0` to `65535`. If not, the value is - 1. Currently, only a single port number is supported, and the start and end of the port number must be the same. The corresponding destination port number range is fixed for different protocol types. The values are as follows:
  - `ALL`:-1/-1, not editable.
  - **ICMP(IPv4)**:-1/-1, non-editable.
  - **ICMPv6(IPv6)**:-1/-1, non-editable.
  - `TCP`:-1/-1, editable.
  - `UDP`:-1/-1, editable.
  - `GRE`:-1/-1, not editable.
  - `SSH`:22/22, not editable.
  - `Telnet`:23/23, not editable.
  - `HTTP`:80/80, non-editable.
  - `HTTPS`:443/443, which cannot be edited.
  - **MS SQL**:1443/1443, which cannot be edited.
  - `Oracle`:1521/1521, non-editable.
  - `MySql`:3306/3306, non-editable.
  - `RDP`:3389/3389, non-editable.
  - `PostgreSQL`:5432/5432, non-editable.
  - `Redis`:6379/6379, non-editable.
* `match_dscp` - (Optional, Computed, Int) The DSCP value of the traffic matched by the QoS rule. Value range: `0` to `63`. If not, the value is - 1.
* `priority` - (Required, Int) QoS rule priority. Value range: `1` to `9000`. The larger the number, the higher the priority. The priority of a QoS rule cannot be repeated in the same QoS policy.
* `protocol` - (Required) QoS rule protocol type, value:
  - `ALL`
  - **ICMP(IPv4)**
  - **ICMPv6(IPv6)* *
  - `TCP`
  - `UDP`
  - `GRE`
  - `SSH`
  - `Telnet`
  - `HTTP`
  - `HTTPS`
  - **MS SQL**
  - `Oracle`
  - `MySql`
  - `RDP`
  - `PostgreSQL`
  - `Redis`
* `qos_id` - (Required, ForceNew) The QoS policy ID.
* `queue_id` - (Required, ForceNew) The QoS queue ID.
* `remarking_dscp` - (Optional, Computed, Int) Modify The DSCP value in the flow. Value range: `0` to `63`. If the value is not modified, the value is - 1.
* `rule_description` - (Optional) The description of the QoS rule.
The length is 0 to 256 characters and cannot start with 'http:// 'or 'https.
* `rule_name` - (Optional) The name of the QoS rule.
The length is 0 to 128 characters and cannot start with 'http:// 'or 'https.
* `src_cidr` - (Optional) The QoS rule traffic matches the source IPv4 CIDR block.

-> **NOTE:**  If this parameter is not supported, enter `SrcIPv6Cidr` or **DstIPv6Cidr * *.

* `src_ipv6_cidr` - (Optional) The QoS rule traffic matches the source IPv6 network segment.

-> **NOTE:**  If this parameter is not supported, enter `SrcCidr` or **DstCidr * *.

* `src_port_range` - (Optional, Computed) The source port number of the QoS rule traffic matching. The value range is `0` to `65535`. If the traffic does not match, the value is - 1. Currently, only a single port number is supported, and the start and end of the port number must be the same.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<qos_id>:<queue_id>:<rule_id>`.
* `rule_id` - The ID of the QoS rule.
* `status` - The status of the QoS rule. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Qos Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Qos Rule.
* `update` - (Defaults to 5 mins) Used when update the Traffic Qos Rule.

## Import

Express Connect Traffic Qos Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_traffic_qos_rule.example <qos_id>:<queue_id>:<rule_id>
```