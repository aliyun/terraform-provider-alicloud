---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_probe"
sidebar_current: "docs-alicloud-resource-threat-detection-honeypot-probe"
description: |-
  Provides a Alicloud Threat Detection Honeypot Probe resource.
---

# alicloud_threat_detection_honeypot_probe

Provides a Threat Detection Honeypot Probe resource.

For information about Threat Detection Honeypot Probe and how to use it, see [What is Honeypot Probe](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-createhoneypotprobe).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_honeypot_probe" "default" {
  uuid            = "032b618f-b220-4a0d-bd37-fbdc6ef58b6a"
  probe_type      = "host_probe"
  control_node_id = "a44e1ab3-6945-444c-889d-5bacee7056e8"
  ping            = true
  honeypot_bind_list {
    bind_port_list {
      start_port = 80
      end_port   = 80
    }
    honeypot_id = "ede59ccdb1b7a2e21735d4593a6eb5ed31883af320c5ab63ab33818e94307be9"
  }
  display_name = "apispec"
  arp          = true
}
```

## Argument Reference

The following arguments are supported:
* `arp` - (Optional) ARP spoofing detection.**true**: Enable **false**: Disabled
* `control_node_id` - (Required,ForceNew) The ID of the management node.
* `display_name` - (Required) Probe display name.
* `honeypot_bind_list` - (ForceNew,Optional) Configure the service.See the following `Block HoneypotBindList`.
* `ping` - (Optional) Ping scan detection. Value: **true**: Enable **false**: Disabled
* `probe_type` - (Required,ForceNew) Probe type, support `host_probe` and `vpc_black_hole_probe`.
  * `host_probe`: host probe
  * `vpc_black_hole_probe`: virtual private cloud (VPC) probe
* `proxy_ip` - (ForceNew,Optional) The IP address of the proxy.
* `probe_version` - (ForceNew,Optional) The version of the probe.
* `service_ip_list` - (Computed,Optional) Listen to the IP address list.
* `uuid` - (ForceNew,Optional) Machine uuid, **probe_type** is `host_probe`. This value cannot be empty.
* `vpc_id` - (ForceNew,Optional) The ID of the VPC. **probe_type** is `vpc_black_hole_probe`. This value cannot be empty. 

#### Block HoneypotBindList

The HoneypotBindList supports the following:
* `bind_port_list` - (ForceNew,Optional) List of listening ports.See the following `Block BindPortList`.
* `honeypot_id` - (ForceNew,Optional) Honeypot ID.

#### Block BindPortList

The BindPortList supports the following:
* `bind_port` - (ForceNew,Optional) Whether to bind the port.
* `end_port` - (ForceNew,Optional) End port.
* `fixed` - (ForceNew,Optional) Whether the port is fixed.
* `start_port` - (ForceNew,Optional) Start port.
* `target_port` - (ForceNew,Optional) Destination port.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `honeypot_probe_id` - The first ID of the resource
* `service_ip_list` - Listen to the IP address list.
* `status` - The status of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honeypot Probe.
* `delete` - (Defaults to 5 mins) Used when delete the Honeypot Probe.
* `update` - (Defaults to 5 mins) Used when update the Honeypot Probe.

## Import

Threat Detection Honeypot Probe can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honeypot_probe.example <id>
```