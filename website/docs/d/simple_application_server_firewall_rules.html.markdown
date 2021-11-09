---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_firewall_rules"
sidebar_current: "docs-alicloud-datasource-simple-application-server-firewall-rules"
description: |-
  Provides a list of Simple Application Server Firewall Rules to the user.
---

# alicloud\_simple\_application\_server\_firewall\_rules

This data source provides the Simple Application Server Firewall Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_firewall_rules" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "simple_application_server_firewall_rule_id_1" {
  value = data.alicloud_simple_application_server_firewall_rules.ids.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Firewall Rule IDs.
* `instance_id` - (Required, ForceNew) Alibaba Cloud simple application server instance ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of Simple Application Server Firewall Rules. Each element contains the following attributes:
	* `firewall_rule_id` - The ID of the firewall rule.
	* `id` - The ID of the Firewall Rule. The value formats as `<instance_id>:<firewall_rule_id>`.
	* `instance_id` - Alibaba Cloud simple application server instance ID.
	* `port` - The port range of the firewall rule.
	* `remark` - The remarks of the firewall rule.
	* `rule_protocol` - The transport layer protocol. Valid values: `Tcp`, `Udp`, `TcpAndUdp`.