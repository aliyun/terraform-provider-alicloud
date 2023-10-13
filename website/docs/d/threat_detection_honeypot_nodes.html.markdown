---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_nodes"
sidebar_current: "docs-alicloud-datasource-threat_detection-honeypot-nodes"
description: |-
  Provides a list of Threat Detection Honeypot Node owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_honeypot_nodes

This data source provides Threat Detection Honeypot Node available to the user.[What is Honeypot Node](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createhoneypotnode)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_threat_detection_honeypot_nodes" "default" {
  ids = ["${alicloud_threat_detection_honeypot_node.default.id}"]
}

output "alicloud_threat_detection_honeypot_node_example_id" {
  value = data.alicloud_threat_detection_honeypot_nodes.default.nodes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `node_id` - (ForceNew,Optional) Honeypot management node id.
* `name_regex` - (ForceNew,Optional)  A regex string to filter results by Honeypot Node name.
* `node_name` - (ForceNew,Optional) The name of the management node.
* `ids` - (Optional, ForceNew, Computed) A list of Honeypot Node IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Honeypot Node IDs.
* `names` - A list of Honeypot Node names.
* `nodes` - A list of Honeypot Node Entries. Each element contains the following attributes:
  * `allow_honeypot_access_internet` - Whether to allow honeypot access to the external network. Value:-**true**: Allow-**false**: Disabled
  * `available_probe_num` - Number of probes available.
  * `node_id` - Honeypot management node id.
  * `id` - The ID of the Honeypot management node.
  * `node_name` - Management node name.
  * `page_total` - Total pages.
  * `security_group_probe_ip_list` - Release the collection of network segments.
