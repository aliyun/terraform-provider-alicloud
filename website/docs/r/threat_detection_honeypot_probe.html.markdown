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

For information about Threat Detection Honeypot Probe and how to use it, see [What is Honeypot Probe](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypotprobe).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_honeypot_probe&exampleId=42e6e58c-9d76-ac2b-21ff-3da80fd33b3b329fe6dc&activeTab=example&spm=docs.r.threat_detection_honeypot_probe.0.42e6e58c9d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Honeypot Probe.
* `delete` - (Defaults to 5 mins) Used when delete the Honeypot Probe.
* `update` - (Defaults to 5 mins) Used when update the Honeypot Probe.

## Import

Threat Detection Honeypot Probe can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_honeypot_probe.example <id>
```