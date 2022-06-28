---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_qos_policy"
sidebar_current: "docs-alicloud-resource-sag-qos-policy"
description: |-
  Provides a Sag Qos Policy resource.
---

# alicloud\_sag\_qos\_policy

Provides a Sag qos policy resource. 
You need to create a QoS policy to set priorities, rate limits, and quintuple rules for different messages.

For information about Sag Qos Policy and how to use it, see [What is Qos Policy](https://www.alibabacloud.com/help/doc-detail/140065.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_qos" "default" {
  name = "tf-testAccSagQosName"
}
resource "alicloud_sag_qos_policy" "default" {
  qos_id            = alicloud_sag_qos.default.id
  name              = "tf-testSagQosPolicyName"
  description       = "tf-testSagQosPolicyDescription"
  priority          = "1"
  ip_protocol       = "ALL"
  source_cidr       = "192.168.0.0/24"
  source_port_range = "-1/-1"
  dest_cidr         = "10.10.0.0/24"
  dest_port_range   = "-1/-1"
  start_time        = "2019-10-25T16:41:33+0800"
  end_time          = "2019-10-26T16:41:33+0800"
}
```
## Argument Reference

The following arguments are supported:

* `qos_id` - (Required) The instance ID of the QoS policy to which the quintuple rule is created.
* `name` - (Optional) The name of the QoS policy.
* `description` - (Optional) The description of the QoS policy.
* `priority` - (Required) The priority of the quintuple rule. A smaller value indicates a higher priority. If the priorities of two quintuple rules are the same, the rule created earlier is applied first.Value range: 1 to 7.
* `ip_protocol` - (Required) The transport layer protocol.
* `source_cidr` - (Required) The source CIDR block.
* `source_port_range` - (Required) The source port range of the transport layer.
* `dest_cidr` - (Required) The destination CIDR block.
* `dest_port_range` - (Required) The destination port range.
* `start_time` - (Optional) The time when the quintuple rule takes effect.
* `end_time` - (Optional) The expiration time of the quintuple rule. 


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Qos Policy id and formates as `<qos_id>:<qos_policy_id>`.

## Import

The Sag Qos Policy can be imported using the id, e.g.

```
$ terraform import alicloud_sag_qos_policy.example qos-abc123456:qospy-abc123456
```

