---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_node"
sidebar_current: "docs-alicloud-resource-threat_detection-honeypot-node"
description: |-
  Provides a Alicloud Threat Detection Honeypot Node resource.
---

# alicloud_threat_detection_honeypot_node

Provides a Threat Detection Honeypot Node resource.

For information about Threat Detection Honeypot Node and how to use it, see [What is Honeypot Node](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypotnode).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name                    = var.name
  available_probe_num          = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}
```

## Argument Reference

The following arguments are supported:
* `allow_honeypot_access_internet` - (Optional, ForceNew) Whether to allow honeypot access to the external network. Value:-**true**: Allow-**false**: Disabled
* `available_probe_num` - (Required) Number of probes available.
* `node_name` - (Required) Management node name.
* `security_group_probe_ip_list` - (Optional) Release the collection of network segments.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Honeypot Node.
* `delete` - (Defaults to 5 mins) Used when delete the Honeypot Node.
* `update` - (Defaults to 5 mins) Used when update the Honeypot Node.

## Import

Threat Detection Honeypot Node can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honeypot_node.example <id>
```