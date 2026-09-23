---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_probes"
sidebar_current: "docs-alicloud-datasource-threat-detection-honeypot-probes"
description: |-
  Provides a list of Threat Detection Honeypot Probe owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_honeypot_probes

This data source provides Threat Detection Honeypot Probe available to the user.[What is Honeypot Probe](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypotprobe)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
variable "name" {
  default = "tf-testAccThreatDetectionHoneypotProbe"
}

resource "alicloud_threat_detection_honeypot_probe" "default" {
  uuid            = "e52c7872-29d1-4aa1-9908-0299abd53606"
  probe_type      = "host_probe"
  control_node_id = "e1397077-4941-4b14-b533-ca2bdebd00a3"
  ping            = true
  honeypot_bind_list {
    bind_port_list {
      start_port = 80
      end_port   = 80
    }
    honeypot_id = "4925bf9784de992ecd017ad051528a03b3927ef814eeff76c2ebb3ab9a84bf05"
  }
  display_name = var.name
  arp          = true
}

data "alicloud_threat_detection_honeypot_probes" "default" {
  ids            = ["${alicloud_threat_detection_honeypot_probe.default.id}"]
  display_name   = var.name
  probe_type     = "host_probe"
  enable_details = true
}

output "alicloud_threat_detection_honeypot_probe_example_id" {
  value = data.alicloud_threat_detection_honeypot_probes.default.probes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `display_name` - (ForceNew, Optional) Probe name
* `probe_type` - (ForceNew, Optional) Probe type
* `ids` - (Optional, ForceNew, Computed) A list of Honeypot Probe IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional, ForceNew) Default to `false`. Set it to `true` can output more details about resource attributes.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by display name.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Honeypot Probe IDs.
* `probes` - A list of Honeypot Probe Entries. Each element contains the following attributes:
    * `id` - The ID of the honeypot probe. Its value is the same as `honeypot_probe_id`.
    * `control_node_id` - The ID of the management node.
    * `display_name` - Probe name.
    * `honeypot_probe_id` - The first ID of the resource
    * `probe_type` - Probe type, support `host_probe` and `vpc_black_hole_probe`.
    * `status` - The status of the resource.
    * `uuid` - Machine uuid. Has a value when the type is `host_probe`.
    * `vpc_id` - The ID of the VPC. Has a value when the type is `vpc_black_hole_probe`.
    * `ping` - Ping scan detection. Value:**true**: Enable **false**: Disabled. Available when `enable_details` is on.
    * `arp` - ARP spoofing detection.-**true**: Enable-**false**: Disabled. Available when `enable_details` is on.
    * `service_ip_list` - Listen to the IP address list. Available when `enable_details` is on.
    * `proxy_ip` - Proxy IP. Available when `enable_details` is on.
    * `bind_port_list` - List of listening ports. Available when `enable_details` is on.
      * `bind_port` - Whether to bind the port.
      * `end_port` - End port.
      * `fixed` - Whether the port is fixed.
      * `start_port` - Start port.
      * `target_port` - Destination port.
